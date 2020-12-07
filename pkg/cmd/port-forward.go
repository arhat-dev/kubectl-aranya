package cmd

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"reflect"
	"strconv"
	"strings"
	"time"

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

type portForwardOptions struct {
	Network string `json:"network"`
	Address string `json:"address"`
	Port    int32  `json:"port"`
}

// nolint:gocyclo
func runPortForward(appCtx context.Context, podName string) error {
	config, kubeClient, _, tlsConfig, namespace := getAppOpts(appCtx)

	// validate config options

	opts := config.PortForwardOptions

	pfOpts := &portForwardOptions{
		Network: opts.RemoteNetwork,
		Address: opts.RemoteAddress,
		Port:    opts.RemotePort,
	}
	pfOptsBytes, err := json.Marshal(pfOpts)
	if err != nil {
		return err
	}

	listenAddr := opts.LocalAddress
	if !strings.HasPrefix(opts.LocalNetwork, "unix") {
		// not unix listen addr, add port any way
		listenAddr = net.JoinHostPort(opts.LocalAddress, strconv.FormatInt(int64(opts.LocalPort), 10))
	}

	rawListener, err := nethelper.Listen(appCtx, nil, opts.LocalNetwork, listenAddr, nil)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	listenerCloser, ok := rawListener.(io.Closer)
	if !ok {
		return fmt.Errorf("invalid network listen not implementing io.Closer")
	}
	defer func() {
		_ = listenerCloser.Close()
	}()

	go func() {
		<-appCtx.Done()

		_ = listenerCloser.Close()
	}()

	// validate listener
	var listenNetwork string
	switch t := rawListener.(type) {
	case net.Listener:
		// stream oriented
		listenNetwork = t.Addr().Network()
		listenAddr = t.Addr().String()
	case net.Conn:
		// packet oriented
		listenNetwork = t.LocalAddr().Network()
		listenAddr = t.LocalAddr().String()
	default:
		return fmt.Errorf("unknown local network listener implementation: %s", reflect.TypeOf(rawListener).String())
	}

	_, err = kubeClient.CoreV1().Pods(namespace).Get(appCtx, podName, metav1.GetOptions{})
	if err != nil {
		if !kubeerrors.IsForbidden(err) {
			// no permission for get pod, port-forward directly
			return err
		}
	}

	pfReqURL := kubeClient.CoreV1().RESTClient().Post().
		Resource("pods").
		Namespace(namespace).
		Name(podName).
		SubResource("portforward").URL()

	// MUST not be buffered
	preparedRemoteConnCh := make(chan net.Conn)
	reqRemoteConnCh := make(chan struct{})
	reqCancelRemoteConnCh := make(chan struct{}, 16)
	appExited := appCtx.Done()

	// prepare new connection for port-forwarding
	go func() {
		reqURL := pfReqURL.String()
		dialer := &net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}

		createConn := func() (_ net.Conn, err error) {
			conn, err := dialer.DialContext(appCtx, "tcp", pfReqURL.Host)
			if err != nil {
				return nil, err
			}
			defer func() {
				if err != nil {
					_ = conn.Close()
				}
			}()

			conn = tls.Client(conn, tlsConfig)

			// nolint:staticcheck
			httpClient := httputil.NewClientConn(conn, nil)
			postReq, err := http.NewRequestWithContext(appCtx, http.MethodPost, reqURL, bytes.NewReader(pfOptsBytes))
			if err != nil {
				return nil, err
			}

			postReq.Header.Set(httpstream.HeaderConnection, httpstream.HeaderUpgrade)

			resp, err := httpClient.Do(postReq)
			if err != nil {
				return nil, err
			}

			if resp.StatusCode != http.StatusSwitchingProtocols {
				_ = resp.Body.Close()
				return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
			}

			// expect no data remain unread
			_, _ = httpClient.Hijack()

			return conn, nil
		}

		for {
			select {
			case <-appExited:
				return
			case <-reqRemoteConnCh:
				// establish a new connection
				go func() {
					for {
						conn, err2 := createConn()
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

						select {
						case <-reqCancelRemoteConnCh:
							_ = conn.Close()
							return
						case <-appExited:
							_ = conn.Close()
							return
						case preparedRemoteConnCh <- conn:
							return
						}
					}
				}()
			}
		}
	}()

	fmt.Printf("forwarding %s://%s -> %s://%s:%d@%s\n",
		listenNetwork, listenAddr,
		opts.RemoteNetwork, opts.RemoteAddress, opts.RemotePort, podName,
	)

	switch l := rawListener.(type) {
	case net.Listener:
		// stream oriented
		for {
			conn, err2 := l.Accept()
			if err2 != nil {
				klog.Info(err2)
				if strings.Contains(err2.Error(), "closed") {
					return nil
				}

				return err2
			}

			// request a new connection
			select {
			case <-appExited:
				_ = conn.Close()
				return nil
			case reqRemoteConnCh <- struct{}{}:
			}

			klog.Infoln("handling new connection", conn.RemoteAddr().String())

			go func() {
				finished := make(chan struct{})
				defer func() {
					_ = conn.Close()

					close(finished)
				}()

				var remoteConn net.Conn
				timer := time.NewTimer(time.Second)
			checkLoop:
				for {
					select {
					case remoteConn = <-preparedRemoteConnCh:
						if !timer.Stop() {
							<-timer.C
						}

						break checkLoop
					case <-appExited:
						select {
						case reqCancelRemoteConnCh <- struct{}{}:
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
							case reqCancelRemoteConnCh <- struct{}{}:
								// cancel connection creation
							case <-appExited:
							}

							klog.Infoln(err2)

							_ = timer.Stop()
							return
						}

						timer.Reset(time.Second)
					}
				}

				go func() {
					_, err2 := io.Copy(conn, remoteConn)
					if err2 != nil && err2 != io.EOF {
						klog.Infoln(err2)
					}
				}()

				_, err2 := io.Copy(remoteConn, conn)
				if err2 != nil && err2 != io.EOF {
					klog.Infoln(err2)
				}
			}()
		}
	case net.Conn:
		// TODO: use chunked data transmission
		// packet oriented connection
		for {
			// request a new connection
			select {
			case <-appExited:
				return nil
			case reqRemoteConnCh <- struct{}{}:
			}

			var remoteConn net.Conn
			select {
			case <-appExited:
				return nil
			case remoteConn = <-preparedRemoteConnCh:
			}

			go func() {
				buf := make([]byte, 65537)
				_, err2 := io.CopyBuffer(l, remoteConn, buf)
				if err2 != nil && err2 != io.EOF {
					klog.Infoln(err2)
				}
			}()

			buf := make([]byte, 65537)
			_, err = io.CopyBuffer(remoteConn, l, buf)
			if err != nil && err != io.EOF {
				klog.Infoln(err)
			}
		}
	default:
		return fmt.Errorf("unknown local network listener implementation: %s", reflect.TypeOf(rawListener).String())
	}
}
