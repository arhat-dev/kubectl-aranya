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

// nolint:maligned
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

	fs.StringVar(&opts.RemoteNetwork, "remote-network", "tcp", "set network protocol to be forwarded")
	fs.StringVar(&opts.RemoteAddress, "remote-address", "localhost", "set target address to be forwarded")
	fs.Int32Var(&opts.RemotePort, "remote-port", 0, "set port for ip network")

	fs.StringVar(&opts.LocalNetwork, "local-network", "tcp", "set local network protocol")
	fs.StringVar(&opts.LocalAddress, "local-address", "localhost", "set local listen address")
	fs.Int32Var(&opts.LocalPort, "local-port", 0, "set local listen port")

	return fs
}
