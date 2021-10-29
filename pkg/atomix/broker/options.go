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

package broker

const (
	defaultPort = 5678
)

type brokerOptions struct {
	namespace string
	name      string
	node      string
	port      int
}

func applyOptions(opts ...Option) brokerOptions {
	options := brokerOptions{
		port: defaultPort,
	}
	for _, opt := range opts {
		opt(&options)
	}
	return options
}

// Option is a broker option
type Option func(opts *brokerOptions)

// WithPort sets the broker port
func WithPort(port int) Option {
	return func(opts *brokerOptions) {
		opts.port = port
	}
}
