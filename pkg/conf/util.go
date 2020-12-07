/*
Copyright 2020 The arhat.dev Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package conf

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"

	"arhat.dev/kubectl-aranya/pkg/constant"
)

func ReadConfig(config *Config) (context.Context, error) {
	cmdFactory := cmdutil.NewFactory(config.kubeConfigFlags)

	namespace, _, err := cmdFactory.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return nil, err
	}

	kubeConfig, err := config.kubeConfigFlags.ToRESTConfig()
	if err != nil {
		return nil, err
	}

	tlsConfig, err := rest.TLSConfigFor(kubeConfig)
	if err != nil {
		return nil, err
	}

	kubeClient, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, err
	}

	if tlsConfig != nil {
		if tlsConfig.ServerName == "" && !tlsConfig.InsecureSkipVerify {
			tlsConfig.ServerName, _, _ = net.SplitHostPort(kubeClient.RESTClient().Get().URL().Host)
		}
	}

	appCtx := context.WithValue(context.Background(), constant.ContextKeyConfig, config)
	appCtx = context.WithValue(appCtx, constant.ContextKeyKubeConfig, kubeConfig)
	appCtx = context.WithValue(appCtx, constant.ContextKeyKubeClient, kubeClient)
	appCtx = context.WithValue(appCtx, constant.ContextKeyTLSConfig, tlsConfig)
	appCtx = context.WithValue(appCtx, constant.ContextKeyNamespace, namespace)

	appCtx, exit := context.WithCancel(appCtx)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		exitCount := 0
		for sig := range sigCh {
			switch sig {
			case os.Interrupt, syscall.SIGTERM:
				exitCount++
				if exitCount == 1 {
					exit()
				} else {
					os.Exit(1)
				}
				//case syscall.SIGHUP:
				//	// force reload
			}
		}
	}()

	return appCtx, nil
}
