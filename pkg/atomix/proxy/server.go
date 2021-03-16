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

package proxy

import (
	proxyapi "github.com/atomix/api/go/atomix/proxy"
	"github.com/atomix/go-framework/pkg/atomix/cluster"
	"github.com/atomix/go-framework/pkg/atomix/logging"
	"github.com/atomix/go-framework/pkg/atomix/primitives"
	"github.com/atomix/go-framework/pkg/atomix/server"
	"google.golang.org/grpc"
)

var log = logging.GetLogger("atomix", "proxy")

// Node is an interface for proxy nodes
type Node interface {
	server.Node
	Primitives() primitives.Resolver
	PrimitiveTypes() *PrimitiveTypeRegistry
}

// NewNode creates a new proxy server
func NewNode(coordinator *cluster.Replica, cluster cluster.Cluster) Node {
	return &Server{
		Server:         server.NewServer(cluster),
		coordinator:    coordinator,
		primitives:     primitives.NewRegistry(),
		primitiveTypes: NewPrimitiveTypeRegistry(),
	}
}

// Server is a proxy server
type Server struct {
	*server.Server
	coordinator    *cluster.Replica
	primitives     *primitives.Registry
	primitiveTypes *PrimitiveTypeRegistry
}

func (s *Server) Primitives() primitives.Resolver {
	return s.primitives
}

func (s *Server) PrimitiveTypes() *PrimitiveTypeRegistry {
	return s.primitiveTypes
}

// Start starts the node
func (s *Server) Start() error {
	s.RegisterService(func(server *grpc.Server) {
		proxyapi.RegisterProxyConfigServiceServer(server, newConfigServer(s.Cluster))
	})
	return s.Server.Start()
}
