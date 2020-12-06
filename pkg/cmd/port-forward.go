package cmd

import (
	"arhat.dev/kubectl-aranya/pkg/conf"
	"arhat.dev/pkg/nethelper"
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	kubeerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	"net"
	"net/http"
	"net/http/httputil"
	"reflect"
	"strconv"
	"time"
)

func NewPortForwardCmd(appCtx *context.Context) *cobra.Command {
	opts := new(conf.PortForwardOptions)

	cmd := &cobra.Command{
		Use:           "kubectl-aranya",
		Args:          cobra.ExactArgs(1),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPortForward(*appCtx, args[0])
		},
	}

	flags := cmd.Flags()
	flags.AddFlagSet(opts.Flags())

	return cmd
}

type portForwardOptions struct {
	Network string `json:"network"`
	Address string `json:"host"`
	Port    int32  `json:"port"`
}

func runPortForward(appCtx context.Context, podName string) error {
	kubeClient, kubeConfig, namespace, config := getAppOpts(appCtx)
	tlsConfig, err := rest.TLSConfigFor(kubeConfig)
	if err != nil {
		return err
	}

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
	if opts.LocalPort > 0 {
		listenAddr = net.JoinHostPort(opts.LocalAddress, strconv.FormatInt(int64(opts.LocalPort), 10))
	}

	rawListener, err := nethelper.Listen(appCtx, nil, opts.LocalNetwork, listenAddr, nil)
	if err != nil {
		return err
	}

	listenerCloser, ok := rawListener.(io.Closer)
	if !ok {
		return fmt.Errorf("invalid network listen not implementing io.Closer")
	}

	defer func() {
		_ = listenerCloser.Close()
	}()

	switch rawListener.(type) {
	case net.Listener:
		// stream oriented
	case net.Conn:
		// packet oriented
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

	pfReqURL := kubeClient.RESTClient().Post().
		Resource("pods").
		Namespace(namespace).
		Name(podName).
		SubResource("portforward").URL()

	preparedRemoteConnCh := make(chan net.Conn)
	reqRemoteConnCh := make(chan struct{})
	appExited := appCtx.Done()

	// prepare new connection for port-forwarding
	go func() {
		reqUrl := pfReqURL.String()
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

			httpClient := httputil.NewClientConn(conn, nil)

			postReq, err := http.NewRequestWithContext(appCtx, http.MethodPost, reqUrl, bytes.NewReader(pfOptsBytes))
			resp, err := httpClient.Do(postReq)
			if err != nil {
				return nil, err
			}

			if resp.StatusCode != http.StatusSwitchingProtocols {
				return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
			}

			// expect no data remain unread
			_, _ = httpClient.Hijack()

			return conn, nil
		}

		for {
		routine:
			select {
			case <-appExited:
				return
			case <-reqRemoteConnCh:
				// establish a new connection
				for i := 0; i < 5; i++ {
					conn, err2 := createConn()
					if err2 != nil {
						// retry try with delay
						time.Sleep(time.Second)
						continue
					}

					select {
					case <-appExited:
						_ = conn.Close()
						return
					case preparedRemoteConnCh <- conn:
						goto routine
					}
				}
			}
		}
	}()

	switch l := rawListener.(type) {
	case net.Listener:
		// stream oriented
		for {
			conn, err := l.Accept()
			if err != nil {
				return fmt.Errorf("accept: %w", err)
			}

			// request a new connection
			select {
			case <-appExited:
				_ = conn.Close()
				return nil
			case reqRemoteConnCh <- struct{}{}:
			}

			go func() {
				finished := make(chan struct{})
				defer func() {
					_ = conn.Close()

					close(finished)
				}()

				var remoteConn net.Conn
				select {
				case remoteConn = <-preparedRemoteConnCh:
				case <-appExited:
					return
				}

				// close on app exit, will it do anything?
				go func() {
					select {
					case <-finished:
					case <-appExited:
						_ = conn.Close()
					}
				}()

				_, err2 := io.Copy(remoteConn, conn)
				if err2 != nil && err2 != io.EOF {
					// TODO: log error with message
					klog.V(2).Infoln(err2)
				}
			}()
		}
	case net.Conn:
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

			_, err = io.Copy(l, remoteConn)
			if err != nil && err != io.EOF {
				// TODO: log error with message
				klog.V(2).Infoln(err)
			}
		}
	default:
		return fmt.Errorf("unknown local network listener implementation: %s", reflect.TypeOf(rawListener).String())
	}
}
