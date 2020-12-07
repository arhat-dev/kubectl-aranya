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
	"flag"

	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/klog/v2"
)

type Config struct {
	kubeConfigFlags *genericclioptions.ConfigFlags

	PortForwardOptions PortForwardOptions
}

func (c *Config) Flags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("kubectl-aranya", pflag.ContinueOnError)

	c.kubeConfigFlags = genericclioptions.NewConfigFlags(false)
	c.kubeConfigFlags.AddFlags(fs)

	goFlags := flag.NewFlagSet("klog", flag.ContinueOnError)
	klog.InitFlags(goFlags)
	fs.AddGoFlagSet(goFlags)

	return fs
}
