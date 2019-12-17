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

package list

import (
	"context"
	"github.com/atomix/atomix-api/proto/atomix/headers"
	api "github.com/atomix/atomix-api/proto/atomix/list"
	"github.com/atomix/atomix-go-node/pkg/atomix/node"
	"github.com/atomix/atomix-go-node/pkg/atomix/server"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func init() {
	node.RegisterServer(registerServer)
}

// registerServer registers a list server with the given gRPC server
func registerServer(server *grpc.Server, protocol node.Protocol) {
	api.RegisterListServiceServer(server, newServer(protocol.Client()))
}

func newServer(client node.Client) api.ListServiceServer {
	return &Server{
		SessionizedServer: &server.SessionizedServer{
			Type:   listType,
			Client: client,
		},
	}
}

// Server is an implementation of MapServiceServer for the map primitive
type Server struct {
	*server.SessionizedServer
}

// Create opens a new session
func (s *Server) Create(ctx context.Context, request *api.CreateRequest) (*api.CreateResponse, error) {
	log.Tracef("Received CreateRequest %+v", request)
	session, err := s.OpenSession(ctx, request.Header, request.Timeout)
	if err != nil {
		return nil, err
	}
	response := &api.CreateResponse{
		Header: &headers.ResponseHeader{
			SessionID: session,
			Index:     session,
		},
	}
	log.Tracef("Sending CreateResponse %+v", response)
	return response, nil
}

// KeepAlive keeps an existing session alive
func (s *Server) KeepAlive(ctx context.Context, request *api.KeepAliveRequest) (*api.KeepAliveResponse, error) {
	log.Tracef("Received KeepAliveRequest %+v", request)
	if err := s.KeepAliveSession(ctx, request.Header); err != nil {
		return nil, err
	}
	response := &api.KeepAliveResponse{
		Header: &headers.ResponseHeader{
			SessionID: request.Header.SessionID,
		},
	}
	log.Tracef("Sending KeepAliveResponse %+v", response)
	return response, nil
}

// Close closes a session
func (s *Server) Close(ctx context.Context, request *api.CloseRequest) (*api.CloseResponse, error) {
	log.Tracef("Received CloseRequest %+v", request)
	if request.Delete {
		if err := s.Delete(ctx, request.Header); err != nil {
			return nil, err
		}
	} else {
		if err := s.CloseSession(ctx, request.Header); err != nil {
			return nil, err
		}
	}

	response := &api.CloseResponse{
		Header: &headers.ResponseHeader{
			SessionID: request.Header.SessionID,
		},
	}
	log.Tracef("Sending CloseResponse %+v", response)
	return response, nil
}

// Size gets the number of elements in the list
func (s *Server) Size(ctx context.Context, request *api.SizeRequest) (*api.SizeResponse, error) {
	log.Tracef("Received SizeRequest %+v", request)
	in, err := proto.Marshal(&SizeRequest{})
	if err != nil {
		return nil, err
	}

	out, header, err := s.Query(ctx, opSize, in, request.Header)
	if err != nil {
		return nil, err
	}

	sizeResponse := &SizeResponse{}
	if err = proto.Unmarshal(out, sizeResponse); err != nil {
		return nil, err
	}

	response := &api.SizeResponse{
		Header: header,
		Size_:  sizeResponse.Size_,
	}
	log.Tracef("Sending SizeResponse %+v", response)
	return response, nil
}

// Contains checks whether the list contains a value
func (s *Server) Contains(ctx context.Context, request *api.ContainsRequest) (*api.ContainsResponse, error) {
	log.Tracef("Received ContainsRequest %+v", request)
	in, err := proto.Marshal(&ContainsRequest{
		Value: request.Value,
	})
	if err != nil {
		return nil, err
	}

	out, header, err := s.Query(ctx, opContains, in, request.Header)
	if err != nil {
		return nil, err
	}

	containsResponse := &ContainsResponse{}
	if err = proto.Unmarshal(out, containsResponse); err != nil {
		return nil, err
	}

	response := &api.ContainsResponse{
		Header:   header,
		Contains: containsResponse.Contains,
	}
	log.Tracef("Sending ContainsResponse %+v", response)
	return response, nil
}

// Append adds a value to the end of the list
func (s *Server) Append(ctx context.Context, request *api.AppendRequest) (*api.AppendResponse, error) {
	log.Tracef("Received AppendRequest %+v", request)
	in, err := proto.Marshal(&AppendRequest{
		Value: request.Value,
	})
	if err != nil {
		return nil, err
	}

	out, header, err := s.Command(ctx, opAppend, in, request.Header)
	if err != nil {
		return nil, err
	}

	appendResponse := &AppendResponse{}
	if err = proto.Unmarshal(out, appendResponse); err != nil {
		return nil, err
	}

	response := &api.AppendResponse{
		Header: header,
		Status: getResponseStatus(appendResponse.Status),
	}
	log.Tracef("Sending AppendResponse %+v", response)
	return response, nil
}

// Insert inserts a value at a specific index
func (s *Server) Insert(ctx context.Context, request *api.InsertRequest) (*api.InsertResponse, error) {
	log.Tracef("Received InsertRequest %+v", request)
	in, err := proto.Marshal(&InsertRequest{
		Index: request.Index,
		Value: request.Value,
	})
	if err != nil {
		return nil, err
	}

	out, header, err := s.Command(ctx, opInsert, in, request.Header)
	if err != nil {
		return nil, err
	}

	insertResponse := &InsertResponse{}
	if err = proto.Unmarshal(out, insertResponse); err != nil {
		return nil, err
	}

	response := &api.InsertResponse{
		Header: header,
		Status: getResponseStatus(insertResponse.Status),
	}
	log.Tracef("Sending InsertResponse %+v", response)
	return response, nil
}

// Set sets the value at a specific index
func (s *Server) Set(ctx context.Context, request *api.SetRequest) (*api.SetResponse, error) {
	log.Tracef("Received SetRequest %+v", request)
	in, err := proto.Marshal(&SetRequest{
		Index: request.Index,
		Value: request.Value,
	})
	if err != nil {
		return nil, err
	}

	out, header, err := s.Command(ctx, opSet, in, request.Header)
	if err != nil {
		return nil, err
	}

	setResponse := &SetResponse{}
	if err = proto.Unmarshal(out, setResponse); err != nil {
		return nil, err
	}

	response := &api.SetResponse{
		Header: header,
		Status: getResponseStatus(setResponse.Status),
	}
	log.Tracef("Sending SetResponse %+v", response)
	return response, nil
}

// Get gets the value at a specific index
func (s *Server) Get(ctx context.Context, request *api.GetRequest) (*api.GetResponse, error) {
	log.Tracef("Received GetRequest %+v", request)
	in, err := proto.Marshal(&GetRequest{
		Index: request.Index,
	})
	if err != nil {
		return nil, err
	}

	out, header, err := s.Query(ctx, opGet, in, request.Header)
	if err != nil {
		return nil, err
	}

	getResponse := &GetResponse{}
	if err = proto.Unmarshal(out, getResponse); err != nil {
		return nil, err
	}

	response := &api.GetResponse{
		Header: header,
		Status: getResponseStatus(getResponse.Status),
		Value:  getResponse.Value,
	}
	log.Tracef("Sending GetResponse %+v", response)
	return response, nil
}

// Remove removes an index from the list
func (s *Server) Remove(ctx context.Context, request *api.RemoveRequest) (*api.RemoveResponse, error) {
	log.Tracef("Received RemoveRequest %+v", request)
	in, err := proto.Marshal(&RemoveRequest{
		Index: request.Index,
	})
	if err != nil {
		return nil, err
	}

	out, header, err := s.Command(ctx, opRemove, in, request.Header)
	if err != nil {
		return nil, err
	}

	removeResponse := &RemoveResponse{}
	if err = proto.Unmarshal(out, removeResponse); err != nil {
		return nil, err
	}

	response := &api.RemoveResponse{
		Header: header,
		Status: getResponseStatus(removeResponse.Status),
		Value:  removeResponse.Value,
	}
	log.Tracef("Sending RemoveResponse %+v", response)
	return response, nil
}

// Clear removes all indexes from the list
func (s *Server) Clear(ctx context.Context, request *api.ClearRequest) (*api.ClearResponse, error) {
	log.Tracef("Received ClearRequest %+v", request)
	in, err := proto.Marshal(&ClearRequest{})
	if err != nil {
		return nil, err
	}

	out, header, err := s.Command(ctx, opClear, in, request.Header)
	if err != nil {
		return nil, err
	}

	clearResponse := &ClearResponse{}
	if err = proto.Unmarshal(out, clearResponse); err != nil {
		return nil, err
	}

	response := &api.ClearResponse{
		Header: header,
	}
	log.Tracef("Sending ClearResponse %+v", response)
	return response, nil
}

// Events listens for list change events
func (s *Server) Events(request *api.EventRequest, srv api.ListService_EventsServer) error {
	log.Tracef("Received EventRequest %+v", request)
	in, err := proto.Marshal(&ListenRequest{
		Replay: request.Replay,
	})
	if err != nil {
		return err
	}

	ch := make(chan server.SessionOutput)
	if err := s.CommandStream(opEvents, in, request.Header, ch); err != nil {
		return err
	}

	for result := range ch {
		if result.Failed() {
			return result.Error
		}

		response := &ListenResponse{}
		if err = proto.Unmarshal(result.Value, response); err != nil {
			return err
		}

		eventResponse := &api.EventResponse{
			Header: result.Header,
			Type:   getEventType(response.Type),
			Index:  response.Index,
			Value:  response.Value,
		}
		log.Tracef("Sending EventResponse %+v", response)
		if err = srv.Send(eventResponse); err != nil {
			return err
		}
	}

	log.Tracef("Finished EventRequest %+v", request)
	return nil
}

// Iterate lists all the value in the list
func (s *Server) Iterate(request *api.IterateRequest, srv api.ListService_IterateServer) error {
	log.Tracef("Received IterateRequest %+v", request)
	in, err := proto.Marshal(&IterateRequest{})
	if err != nil {
		return err
	}

	ch := make(chan server.SessionOutput)
	if err := s.QueryStream(opIterate, in, request.Header, ch); err != nil {
		return err
	}

	for result := range ch {
		if result.Failed() {
			return result.Error
		}

		response := &IterateResponse{}
		if err = proto.Unmarshal(result.Value, response); err != nil {
			srv.Context().Done()
		}
		iterateResponse := &api.IterateResponse{
			Header: result.Header,
			Value:  response.Value,
		}
		log.Tracef("Sending IterateResponse %+v", response)
		if err = srv.Send(iterateResponse); err != nil {
			return err
		}
	}

	log.Tracef("Finished IterateRequest %+v", request)
	return nil
}

func getResponseStatus(status ResponseStatus) api.ResponseStatus {
	switch status {
	case ResponseStatus_OK:
		return api.ResponseStatus_OK
	case ResponseStatus_NOOP:
		return api.ResponseStatus_NOOP
	case ResponseStatus_WRITE_LOCK:
		return api.ResponseStatus_WRITE_LOCK
	case ResponseStatus_OUT_OF_BOUNDS:
		return api.ResponseStatus_OUT_OF_BOUNDS
	}
	return api.ResponseStatus_OK
}

func getEventType(eventType ListenResponse_Type) api.EventResponse_Type {
	switch eventType {
	case ListenResponse_NONE:
		return api.EventResponse_NONE
	case ListenResponse_ADDED:
		return api.EventResponse_ADDED
	case ListenResponse_REMOVED:
		return api.EventResponse_REMOVED
	default:
		return api.EventResponse_OPEN
	}
}
