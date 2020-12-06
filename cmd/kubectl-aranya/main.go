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

package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"math/rand"
	"os"
	"strings"
	"time"

	"arhat.dev/kubectl-aranya/pkg/cmd"
	"arhat.dev/kubectl-aranya/pkg/version"

	// add auth plugins for cloud providers
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	cobra.OnInitialize(viper.AutomaticEnv)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	rootCmd := cmd.NewRootCmd()
	rootCmd.AddCommand(version.NewVersionCmd())

	err := rootCmd.Execute()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to run kubectl-aranya %v: %v\n", os.Args, err)
		os.Exit(1)
	}
}
