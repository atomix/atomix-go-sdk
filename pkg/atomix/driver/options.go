// Copyright 2019-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package driver

const (
	defaultID   = "atomix-driver"
	defaultHost = ""
	defaultPort = 5252
)

type driverOptions struct {
	driverID string
	nodeID   string
	host     string
	port     int
}

func applyOptions(opts ...Option) driverOptions {
	options := driverOptions{
		driverID: defaultID,
		host:     defaultHost,
		port:     defaultPort,
	}
	for _, opt := range opts {
		opt(&options)
	}
	return options
}

// Option is a driver option
type Option func(opts *driverOptions)

func WithDriverID(id string) Option {
	return func(opts *driverOptions) {
		opts.driverID = id
	}
}

func WithNodeID(id string) Option {
	return func(opts *driverOptions) {
		opts.nodeID = id
	}
}

func WithHost(host string) Option {
	return func(opts *driverOptions) {
		opts.host = host
	}
}

func WithPort(port int) Option {
	return func(opts *driverOptions) {
		opts.port = port
	}
}
