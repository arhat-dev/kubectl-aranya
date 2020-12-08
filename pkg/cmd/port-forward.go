package cmd

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"arhat.dev/aranya-proto/aranyagopb"
	"arhat.dev/aranya-proto/aranyagopb/aranyagoconst"
	"arhat.dev/arhat-proto/arhatgopb"
	"arhat.dev/libext/codec"
	"arhat.dev/pkg/nethelper"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	kubeerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/httpstream"
	"k8s.io/klog/v2"

	"arhat.dev/kubectl-aranya/pkg/conf"

	// add network support for nethelper
	_ "arhat.dev/pkg/nethelper/stdnet" // tcp/udp/unix

	// add protobuf codec support
	_ "arhat.dev/libext/codec/gogoprotobuf"
)

func NewPortForwardCmd(appCtx *context.Context, opts *conf.PortForwardOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "port-forward",
		Args:          cobra.ExactArgs(1),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPortForward(*appCtx, args[0])
		},
	}

	flags := cmd.Flags()
	flags.AddFlagSet(opts.Flags())

	err := viper.BindPFlags(flags)
	if err != nil {
		panic(err)
	}

	return cmd
}

// nolint:gocyclo
func runPortForward(appCtx context.Context, podName string) error {
	pbCodec, ok := codec.Get(arhatgopb.CODEC_PROTOBUF)
	if !ok {
		panic("protobuf codec support not built-in")
	}

	config, kubeClient, _, tlsConfig, namespace := getAppOpts(appCtx)

	// validate config options

	opts := config.PortForwardOptions

	listenAddr := opts.LocalAddress
	switch {
	case strings.HasPrefix(opts.LocalNetwork, "unix"):
	case strings.HasPrefix(opts.LocalNetwork, "ip"):
	default:
		// not unix/ip network
		listenAddr = net.JoinHostPort(opts.LocalAddress, strconv.FormatInt(int64(opts.LocalPort), 10))
	}

	rawListener, err := nethelper.Listen(appCtx, nil, opts.LocalNetwork, listenAddr, nil)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	listenerCloser, ok := rawListener.(io.Closer)
	if !ok {
		panic(fmt.Sprintf(
			"invalid %q listener not implementing io.Closer", opts.LocalNetwork,
		))
	}
	defer func() {
		_ = listenerCloser.Close()
	}()

	go func() {
		<-appCtx.Done()

		_ = listenerCloser.Close()
	}()

	// validate listener
	var (
		listenNetwork string
		isPacket      bool
		fwd           forwarder

		appExited             = appCtx.Done()
		reqRemoteConnCh       = make(chan struct{})
		reqCancelRemoteConnCh = make(chan struct{}, 16)

		// MUST not be buffered so we can cancel
		preparedRemoteConnCh = make(chan preparedRemoteConn)
	)

	switch t := rawListener.(type) {
	case net.Listener:
		// stream oriented
		listenNetwork = t.Addr().Network()
		listenAddr = t.Addr().String()
		isPacket = false

		fwd = &streamForwarder{
			appExited:             appExited,
			reqRemoteConnCh:       reqRemoteConnCh,
			reqCancelRemoteConnCh: reqCancelRemoteConnCh,
			preparedRemoteConnCh:  preparedRemoteConnCh,

			pbCodec:  pbCodec,
			listener: t,
		}
	case net.PacketConn:
		// packet oriented
		listenNetwork = t.LocalAddr().Network()
		listenAddr = t.LocalAddr().String()
		isPacket = true

		fwd = &packetForwarder{
			appExited:             appExited,
			reqRemoteConnCh:       reqRemoteConnCh,
			reqCancelRemoteConnCh: reqCancelRemoteConnCh,
			preparedRemoteConnCh:  preparedRemoteConnCh,

			pbCodec:  pbCodec,
			listener: t,

			localEndpoints: make(map[net.Addr]*localPacketEndpoint),
			alive:          make(map[net.Addr]struct{}),
			creating:       make(map[net.Addr]struct{}),
		}
	default:
		return fmt.Errorf("unknown local network listener implementation: %s", reflect.TypeOf(rawListener).String())
	}

	pfOptsBytes, err := json.Marshal(&aranyagoconst.CustomPortForwardOptions{
		Network: opts.RemoteNetwork,
		Address: opts.RemoteAddress,
		Port:    opts.RemotePort,
		Packet:  isPacket,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal port-forward options: %w", err)
	}

	_, err = kubeClient.CoreV1().Pods(namespace).Get(appCtx, podName, metav1.GetOptions{})
	if err != nil {
		if !kubeerrors.IsForbidden(err) {
			// no permission to get pod, port-forward directly
			return err
		}
	}

	pfReqURL := kubeClient.CoreV1().RESTClient().Post().
		Resource("pods").Namespace(namespace).Name(podName).
		SubResource("portforward").URL()

	// prepare new connection for port-forwarding
	go func() {
		reqURL := pfReqURL.String()
		dialer := &net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}

		createConn := func() (_ net.Conn, sid uint64, mtu int64, err error) {
			conn, err := dialer.DialContext(appCtx, "tcp", pfReqURL.Host)
			if err != nil {
				return nil, 0, 0, err
			}
			defer func() {
				if err != nil {
					_ = conn.Close()
				}
			}()

			conn = tls.Client(conn, tlsConfig)

			// nolint:staticcheck
			postReq, err := http.NewRequestWithContext(
				appCtx, http.MethodPost, reqURL,
				bytes.NewReader(pfOptsBytes),
			)
			postReq.ContentLength = -1
			postReq.Close = false

			// Set mandatory Headers
			//   Connection: Upgrade
			// required by kubernetes api-server for port-forward api endpoint
			postReq.Header.Set(httpstream.HeaderConnection, httpstream.HeaderUpgrade)

			// do request
			err = postReq.Write(conn)
			if err != nil {
				return nil, 0, 0, fmt.Errorf("failed to write request: %w", err)
			}

			resp, err := http.ReadResponse(bufio.NewReader(conn), postReq)
			defer func() {
				if err != nil && resp != nil {
					respData, err2 := ioutil.ReadAll(resp.Body)
					_ = resp.Body.Close()
					if err2 != nil && err2 != io.EOF {
						klog.Infoln(err2)
					}

					if len(respData) != 0 {
						klog.Infoln(string(respData))
					}
				}
			}()
			if err != nil {
				return nil, 0, 0, fmt.Errorf("failed to read response: %w", err)
			}

			if resp.StatusCode != http.StatusSwitchingProtocols {
				return nil, 0, 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
			}

			sidStr := resp.Header.Get(aranyagoconst.HeaderSessionID)
			if sidStr == "" {
				return nil, 0, 0, fmt.Errorf("unexpected empty session id")
			}

			sid, err = strconv.ParseUint(sidStr, 10, 64)
			if err != nil {
				return nil, 0, 0, fmt.Errorf("invalid session id %q: %w", sidStr, err)
			}

			mtuStr := resp.Header.Get(aranyagoconst.HeaderMaxPayloadSize)
			if mtuStr == "" {
				return nil, 0, 0, fmt.Errorf("unexpected no max payload size set")
			}

			mtu, err = strconv.ParseInt(mtuStr, 10, 64)
			if err != nil {
				return nil, 0, 0, fmt.Errorf("invalid max payload size %q: %w", mtuStr, err)
			}

			if mtu <= 0 {
				return nil, 0, 0, fmt.Errorf("unexpected max payload size %q", mtuStr)
			}

			startMsg := []byte("\r\nport-forward\r\n")
			var (
				n   int
				sum int
			)
			for sum != len(startMsg) {
				n, err = conn.Write(startMsg)
				if err != nil {
					return nil, 0, 0, fmt.Errorf("failed to mark port-forward start: %w", err)
				}
				sum += n
			}

			return conn, sid, mtu, nil
		}

		for {
			select {
			case <-appExited:
				return
			case <-reqRemoteConnCh:
				// establish a new connection
				go func() {
					for {
						conn, sid, mtu, err2 := createConn()
						if err2 != nil {
							// retry try with delay
							klog.Infof("failed to create new connection to kubernetes: %v\n", err2)

							select {
							case <-reqCancelRemoteConnCh:
								return
							case <-appExited:
								return
							default:
								// TODO: backoff?
								time.Sleep(time.Second)
								continue
							}
						}

						bufSize := 65537
						if mtu < 65537 {
							bufSize = int(mtu)
						}

						select {
						case <-reqCancelRemoteConnCh:
							_ = conn.Close()
							return
						case <-appExited:
							_ = conn.Close()
							return
						case preparedRemoteConnCh <- preparedRemoteConn{
							sid:     sid,
							bufSize: bufSize,
							conn:    conn,
						}:
							return
						}
					}
				}()
			}
		}
	}()

	fmt.Printf("Forwarding %s://%s -> %s://%s:%d@%s\n",
		listenNetwork, listenAddr,
		opts.RemoteNetwork, opts.RemoteAddress,
		opts.RemotePort, podName,
	)

	err = fwd.run()
	if err != nil {
		klog.Infoln(err)
	}

	return nil
}

type preparedRemoteConn struct {
	sid     uint64
	bufSize int
	conn    net.Conn
}

type forwarder interface {
	run() error
}

type streamForwarder struct {
	appExited             <-chan struct{}
	reqRemoteConnCh       chan<- struct{}
	reqCancelRemoteConnCh chan<- struct{}
	preparedRemoteConnCh  <-chan preparedRemoteConn

	pbCodec  codec.Interface
	listener net.Listener
}

func (sf *streamForwarder) run() error {
	// stream oriented, data is ordered
	for {
		conn, err2 := sf.listener.Accept()
		if err2 != nil {
			klog.Info(err2)
			if strings.Contains(err2.Error(), "closed") {
				return nil
			}

			return err2
		}

		// request a new remote connection
		select {
		case <-sf.appExited:
			_ = conn.Close()
			return nil
		case sf.reqRemoteConnCh <- struct{}{}:
		}

		klog.Infoln("Handling connection", conn.RemoteAddr().String())

		go sf.handleRemote(conn)
	}
}

func (sf *streamForwarder) handleRemote(conn net.Conn) {
	finished := make(chan struct{})
	defer func() {
		_ = conn.Close()

		close(finished)
	}()

	var remoteConn preparedRemoteConn
	timer := time.NewTimer(time.Second)
checkLoop:
	for {
		select {
		case remoteConn = <-sf.preparedRemoteConnCh:
			if !timer.Stop() {
				<-timer.C
			}

			break checkLoop
		case <-sf.appExited:
			select {
			case sf.reqCancelRemoteConnCh <- struct{}{}:
			default:
				// do not wait since application exited
			}

			if !timer.Stop() {
				<-timer.C
			}

			return
		case <-timer.C:
			zero := make([]byte, 0)
			_, err2 := conn.Write(zero)
			if err2 != nil {
				select {
				case sf.reqCancelRemoteConnCh <- struct{}{}:
					// cancel connection creation
				case <-sf.appExited:
				}

				klog.Infoln(err2)

				_ = timer.Stop()
				return
			}

			timer.Reset(time.Second)
		}
	}

	go func() {
		// data is ordered, do not need to decode aranya-proto packets, always raw data
		_, err2 := io.Copy(conn, remoteConn.conn)
		if err2 != nil && err2 != io.EOF {
			klog.Infoln(err2)
		}
	}()

	buf := codec.GetBytesBuf(remoteConn.bufSize)
	defer func() {
		codec.PutBytesBuf(&buf)
	}()

	enc := sf.pbCodec.NewEncoder(remoteConn.conn)
	seq := uint64(0)

	defer func() {
		// close remote connection since local read has finished
		_ = remoteConn.conn.Close()
	}()

	for {
		n, err2 := conn.Read(buf)
		if err2 != nil {
			if err2 != io.EOF {
				klog.Infoln(err2)
			}

			if n == 0 {
				return
			}
		}

		data := make([]byte, n)
		_ = copy(data, buf)

		err2 = enc.Encode(&aranyagopb.Cmd{
			Kind:    aranyagopb.CMD_DATA_UPSTREAM,
			Sid:     remoteConn.sid,
			Seq:     seq, // set sequence so aranya won't need to decode and encode
			Payload: data,
		})
		if err2 != nil {
			klog.Infoln(err2)
			return
		}

		seq++
	}
}

type packetForwarder struct {
	appExited             <-chan struct{}
	reqRemoteConnCh       chan<- struct{}
	reqCancelRemoteConnCh chan<- struct{}
	preparedRemoteConnCh  <-chan preparedRemoteConn

	pbCodec  codec.Interface
	listener net.PacketConn

	localEndpoints map[net.Addr]*localPacketEndpoint
	alive          map[net.Addr]struct{}
	creating       map[net.Addr]struct{}

	_working uint32
}

func (pf *packetForwarder) checkEndpointAlive() {
	timer := time.NewTimer(10 * time.Second)
	defer func() {
		if !timer.Stop() {
			<-timer.C
		}
	}()

	for {
		select {
		case <-pf.appExited:
			return
		case <-timer.C:
			pf.doExclusive(func() {
				for raddr, ep := range pf.localEndpoints {
					if _, ok := pf.alive[raddr]; ok {
						delete(pf.alive, raddr)
						continue
					}

					// already marked not alive
					ep.close()
				}
			})
		}
	}
}

func (pf *packetForwarder) run() error {
	// request a remote connection to check packet size limit
	var initialRemoteConn preparedRemoteConn
	select {
	case <-pf.appExited:
		return nil
	case pf.reqRemoteConnCh <- struct{}{}:
	}

	select {
	case <-pf.appExited:
		return nil
	case initialRemoteConn = <-pf.preparedRemoteConnCh:
	}

	go pf.checkEndpointAlive()

	// bufSize is immutable for one EdgeDevice when it's online
	buf := make([]byte, initialRemoteConn.bufSize)
	for {
		n, raddr, err := pf.listener.ReadFrom(buf)
		if err != nil {
			klog.Infoln(err)
			if n == 0 {
				return err
			}
		}

		data := make([]byte, n)
		_ = copy(data, buf)

		var (
			ep *localPacketEndpoint
			ok bool
		)
		pf.doExclusive(func() {
			ep, ok = pf.localEndpoints[raddr]
			if ok {
				pf.alive[raddr] = struct{}{}
				return
			}

			if _, creating := pf.creating[raddr]; creating {
				return
			}

			pf.creating[raddr] = struct{}{}

			go func() {
				// this may block for a while
				newEp := pf.newLocalEndpoint(raddr)
				if newEp == nil {
					return
				}

				newEp.writeToRemote(data)

				pf.doExclusive(func() {
					pf.localEndpoints[raddr] = newEp
					pf.alive[raddr] = struct{}{}
					delete(pf.creating, raddr)
				})

				ep.run()

				// local endpoint exited, mark it not alive
				pf.doExclusive(func() {
					delete(pf.localEndpoints, raddr)
					delete(pf.alive, raddr)
				})
			}()
		})

		if ok {
			ep.writeToRemote(data)
		}
	}
}

func (pf *packetForwarder) doExclusive(f func()) {
	for !atomic.CompareAndSwapUint32(&pf._working, 0, 1) {
		runtime.Gosched()
	}

	f()

	atomic.StoreUint32(&pf._working, 0)
}

func (pf *packetForwarder) newLocalEndpoint(raddr net.Addr) *localPacketEndpoint {
	var remoteConn preparedRemoteConn
	select {
	case <-pf.appExited:
		return nil
	case pf.reqRemoteConnCh <- struct{}{}:
	}

	select {
	case <-pf.appExited:
		return nil
	case remoteConn = <-pf.preparedRemoteConnCh:
	}

	return &localPacketEndpoint{
		raddr:    raddr,
		listener: pf.listener,

		sid:        remoteConn.sid,
		seq:        0,
		remoteConn: remoteConn.conn,
		enc:        pf.pbCodec.NewEncoder(remoteConn.conn),
	}
}

type localPacketEndpoint struct {
	raddr    net.Addr
	listener net.PacketConn

	sid        uint64
	seq        uint64
	remoteConn net.Conn
	enc        codec.Encoder
}

func (pep *localPacketEndpoint) run() {
	br := bufio.NewReader(pep.remoteConn)
	for {
		size, err := binary.ReadUvarint(br)
		if err != nil {
			klog.Infoln(err)

			// close remote conn to ensure local listener won't send data through it
			_ = pep.remoteConn.Close()
			return
		}

		data := make([]byte, size)
		_, err = io.ReadFull(br, data)
		if err != nil {
			klog.Infoln(err)
			return
		}

		_, err = pep.listener.WriteTo(data, pep.raddr)
		if err != nil && err != io.EOF {
			klog.Infoln(err)
			// failed to write to local conn, we can do nothing
			return
		}
	}
}

func (pep *localPacketEndpoint) close() {
	_ = pep.remoteConn.Close()
}

func (pep *localPacketEndpoint) writeToRemote(data []byte) {
	err := pep.enc.Encode(&aranyagopb.Cmd{
		Kind:     aranyagopb.CMD_DATA_UPSTREAM,
		Sid:      pep.sid,
		Seq:      atomic.AddUint64(&pep.seq, 1) - 1,
		Complete: false,
		Payload:  data,
	})
	if err != nil {
		klog.Infoln(err)
		return
	}
}
