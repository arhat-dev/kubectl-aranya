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

package cmd

import (
	"arhat.dev/kubectl-aranya/pkg/conf"
	"arhat.dev/kubectl-aranya/pkg/constant"
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func NewRootCmd() *cobra.Command {
	var (
		appCtx context.Context
		config = new(conf.Config)
	)

	rootCmd := &cobra.Command{
		Use:           "kubectl-aranya",
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Use == "version" {
				return nil
			}

			var err error
			appCtx, err = conf.ReadConfig(config)
			if err != nil {
				return err
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("please run sub commands")
		},
	}

	// initialize config fields, MUST be called before any other access to config
	appConfigFlags := config.Flags()

	// add sub commands
	for _, c := range []*cobra.Command{
		NewPortForwardCmd(&appCtx),
	} {
		subCmd := c
		subCmd.Flags().AddFlagSet(appConfigFlags)

		rootCmd.AddCommand(subCmd)
	}

	appFlags := rootCmd.Flags()
	appFlags.AddFlagSet(appConfigFlags)
	err := viper.BindPFlags(appFlags)
	if err != nil {
		panic(err)
	}

	return rootCmd
}

func getAppOpts(appCtx context.Context) (_ *kubernetes.Clientset, _ *rest.Config, namespace string, _ *conf.Config) {
	return appCtx.Value(constant.ContextKeyKubeClient).(*kubernetes.Clientset),
		appCtx.Value(constant.ContextKeyKubeConfig).(*rest.Config),
		appCtx.Value(constant.ContextKeyNamespace).(string),
		appCtx.Value(constant.ContextKeyConfig).(*conf.Config)
}
