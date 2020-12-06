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

import "github.com/spf13/pflag"

type PortForwardOptions struct {
	RemoteNetwork string
	RemoteAddress string
	RemotePort    int32

	LocalNetwork string
	LocalAddress string
	LocalPort    int32
}

func (opts *PortForwardOptions) Flags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("port-forward", pflag.ContinueOnError)

	fs.StringVarP(&opts.RemoteNetwork, "remote-network", "N", "tcp", "set network protocol to be forwarded")
	fs.StringVarP(&opts.RemoteAddress, "remote-address", "H", "localhost", "set target address to be forwarded")
	fs.Int32VarP(&opts.RemotePort, "remote-port", "P", 0, "set port for ip network")

	fs.StringVarP(&opts.LocalNetwork, "local-network", "n", "tcp", "set local network protocol")
	fs.StringVarP(&opts.LocalAddress, "local-address", "l", "localhost", "set local listen address")
	fs.Int32VarP(&opts.LocalPort, "local-port", "p", 0, "set local listen port")

	return fs
}
