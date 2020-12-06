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
	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type Config struct {
	kubeConfigFlags *genericclioptions.ConfigFlags

	PortForwardOptions PortForwardOptions
}

func (c *Config) Flags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("kubectl-aranya", pflag.ContinueOnError)

	c.kubeConfigFlags = genericclioptions.NewConfigFlags(false)
	c.kubeConfigFlags.AddFlags(fs)

	return fs
}
