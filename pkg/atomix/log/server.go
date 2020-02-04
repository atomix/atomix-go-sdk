// Copyright 2020-present Open Networking Foundation.
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

package log

import (
	"context"

	"github.com/atomix/api/proto/atomix/headers"
	api "github.com/atomix/api/proto/atomix/log"
	"github.com/atomix/go-framework/pkg/atomix/node"
	"github.com/atomix/go-framework/pkg/atomix/server"
	streams "github.com/atomix/go-framework/pkg/atomix/stream"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func init() {
	node.RegisterServer(registerServer)
}

// registerServer registers a log server with the given gRPC server
func registerServer(server *grpc.Server, protocol node.Protocol) {
	api.RegisterLogServiceServer(server, newServer(protocol.Client()))
}

func newServer(client node.Client) api.LogServiceServer {
	return &Server{
		SessionizedServer: &server.SessionizedServer{
			Type:   logType,
			Client: client,
		},
	}
}

// Server is an implementation of LogServiceServer for the log primitive
type Server struct {
	api.LogServiceServer
	*server.SessionizedServer
}

// Create opens a new session
func (s *Server) Create(ctx context.Context, request *api.CreateRequest) (*api.CreateResponse, error) {
	log.Tracef("Received CreateRequest %+v", request)
	header, err := s.OpenSession(ctx, request.Header, request.Timeout)
	if err != nil {
		return nil, err
	}
	response := &api.CreateResponse{
		Header: header,
	}
	log.Tracef("Sending CreateResponse %+v", response)
	return response, nil
}

// KeepAlive keeps an existing session alive
func (s *Server) KeepAlive(ctx context.Context, request *api.KeepAliveRequest) (*api.KeepAliveResponse, error) {
	log.Tracef("Received KeepAliveRequest %+v", request)
	header, err := s.KeepAliveSession(ctx, request.Header)
	if err != nil {
		return nil, err
	}
	response := &api.KeepAliveResponse{
		Header: header,
	}
	log.Tracef("Sending KeepAliveResponse %+v", response)
	return response, nil
}

// Close closes a session
func (s *Server) Close(ctx context.Context, request *api.CloseRequest) (*api.CloseResponse, error) {
	log.Tracef("Received CloseRequest %+v", request)
	if request.Delete {
		header, err := s.Delete(ctx, request.Header)
		if err != nil {
			return nil, err
		}
		response := &api.CloseResponse{
			Header: header,
		}
		log.Tracef("Sending CloseResponse %+v", response)
		return response, nil
	}

	header, err := s.CloseSession(ctx, request.Header)
	if err != nil {
		return nil, err
	}
	response := &api.CloseResponse{
		Header: header,
	}
	log.Tracef("Sending CloseResponse %+v", response)
	return response, nil
}

// Size gets the number of entries in the log
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

// Exists checks whether the log contains an index
func (s *Server) Exists(ctx context.Context, request *api.ExistsRequest) (*api.ExistsResponse, error) {
	log.Tracef("Received ExistsRequest %+v", request)
	in, err := proto.Marshal(&ContainsIndexRequest{
		Index: request.Index,
	})
	if err != nil {
		return nil, err
	}

	out, header, err := s.Query(ctx, opExists, in, request.Header)
	if err != nil {
		return nil, err
	}

	containsResponse := &ContainsIndexResponse{}
	if err = proto.Unmarshal(out, containsResponse); err != nil {
		return nil, err
	}

	response := &api.ExistsResponse{
		Header:        header,
		ContainsIndex: containsResponse.ContainsIndex,
	}
	log.Tracef("Sending ExistsResponse %+v", response)
	return response, nil
}

// Append appends a value to the end of the log
func (s *Server) Append(ctx context.Context, request *api.AppendRequest) (*api.AppendResponse, error) {
	log.Tracef("Received PutRequest %+v", request)
	in, err := proto.Marshal(&AppendRequest{
		Index:   uint64(request.Index),
		Value:   request.Value,
		Version: uint64(request.Version),
		IfEmpty: request.Version == -1,
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
		Header:          header,
		Status:          getResponseStatus(appendResponse.Status),
		Index:           int64(appendResponse.Index),
		Timestamp:       appendResponse.Timestamp,
		PreviousValue:   appendResponse.PreviousValue,
		PreviousVersion: int64(appendResponse.PreviousVersion),
	}
	log.Tracef("Sending PutResponse %+v", response)
	return response, nil
}

// Get gets the value of an index
func (s *Server) Get(ctx context.Context, request *api.GetRequest) (*api.GetResponse, error) {
	log.Tracef("Received GetRequest %+v", request)
	in, err := proto.Marshal(&GetRequest{
		Index: uint64(request.Index),
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
		Header:    header,
		Index:     int64(getResponse.Index),
		Value:     getResponse.Value,
		Version:   int64(getResponse.Version),
		Timestamp: getResponse.Timestamp,
	}
	log.Tracef("Sending GetResponse %+v", response)
	return response, nil
}

// FirstEntry gets the first entry in the log
func (s *Server) FirstEntry(ctx context.Context, request *api.FirstEntryRequest) (*api.FirstEntryResponse, error) {
	log.Tracef("Received FirstEntryRequest %+v", request)
	in, err := proto.Marshal(&FirstEntryRequest{})
	if err != nil {
		return nil, err
	}

	out, header, err := s.Query(ctx, opFirstEntry, in, request.Header)
	if err != nil {
		return nil, err
	}

	firstEntryResponse := &FirstEntryResponse{}
	if err = proto.Unmarshal(out, firstEntryResponse); err != nil {
		return nil, err
	}

	response := &api.FirstEntryResponse{
		Header:    header,
		Index:     int64(firstEntryResponse.Index),
		Value:     firstEntryResponse.Value,
		Version:   int64(firstEntryResponse.Version),
		Timestamp: firstEntryResponse.Timestamp,
	}
	log.Tracef("Sending FirstEntryResponse %+v", response)
	return response, nil
}

// LastEntry gets the last entry in the log
func (s *Server) LastEntry(ctx context.Context, request *api.LastEntryRequest) (*api.LastEntryResponse, error) {
	log.Tracef("Received LastEntryRequest %+v", request)
	in, err := proto.Marshal(&LastEntryRequest{})
	if err != nil {
		return nil, err
	}

	out, header, err := s.Query(ctx, opLastEntry, in, request.Header)
	if err != nil {
		return nil, err
	}

	lastEntryResponse := &LastEntryResponse{}
	if err = proto.Unmarshal(out, lastEntryResponse); err != nil {
		return nil, err
	}

	response := &api.LastEntryResponse{
		Header:    header,
		Index:     int64(lastEntryResponse.Index),
		Value:     lastEntryResponse.Value,
		Version:   int64(lastEntryResponse.Version),
		Timestamp: lastEntryResponse.Timestamp,
	}
	log.Tracef("Sending LastEntryResponse %+v", response)
	return response, nil
}

// PrevEntry gets the previous entry in the log
func (s *Server) PrevEntry(ctx context.Context, request *api.PrevEntryRequest) (*api.PrevEntryResponse, error) {
	log.Tracef("Received PrevEntryRequest %+v", request)
	in, err := proto.Marshal(&PrevEntryRequest{
		Index: uint64(request.Index),
	})
	if err != nil {
		return nil, err
	}

	out, header, err := s.Query(ctx, opPrevEntry, in, request.Header)
	if err != nil {
		return nil, err
	}

	prevEntryResponse := &PrevEntryResponse{}
	if err = proto.Unmarshal(out, prevEntryResponse); err != nil {
		return nil, err
	}

	response := &api.PrevEntryResponse{
		Header:    header,
		Index:     int64(prevEntryResponse.Index),
		Value:     prevEntryResponse.Value,
		Version:   int64(prevEntryResponse.Version),
		Timestamp: prevEntryResponse.Timestamp,
	}
	log.Tracef("Sending PrevEntryResponse %+v", response)
	return response, nil
}

// NextEntry gets the next entry in the log
func (s *Server) NextEntry(ctx context.Context, request *api.NextEntryRequest) (*api.NextEntryResponse, error) {
	log.Tracef("Received NextEntryRequest %+v", request)
	in, err := proto.Marshal(&NextEntryRequest{
		Index: uint64(request.Index),
	})
	if err != nil {
		return nil, err
	}

	out, header, err := s.Query(ctx, opNextEntry, in, request.Header)
	if err != nil {
		return nil, err
	}

	nextEntryResponse := &NextEntryResponse{}
	if err = proto.Unmarshal(out, nextEntryResponse); err != nil {
		return nil, err
	}

	response := &api.NextEntryResponse{
		Header:    header,
		Index:     int64(nextEntryResponse.Index),
		Value:     nextEntryResponse.Value,
		Version:   int64(nextEntryResponse.Version),
		Timestamp: nextEntryResponse.Timestamp,
	}
	log.Tracef("Sending NextEntryResponse %+v", response)
	return response, nil
}

// Remove removes a key from the log
func (s *Server) Remove(ctx context.Context, request *api.RemoveRequest) (*api.RemoveResponse, error) {
	log.Tracef("Received RemoveRequest %+v", request)
	in, err := proto.Marshal(&RemoveRequest{
		Index:   uint64(request.Index),
		Value:   request.Value,
		Version: uint64(request.Version),
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
		Header:          header,
		Status:          getResponseStatus(removeResponse.Status),
		Index:           int64(removeResponse.Index),
		PreviousValue:   removeResponse.PreviousValue,
		PreviousVersion: int64(removeResponse.PreviousVersion),
	}
	log.Tracef("Sending RemoveRequest %+v", response)
	return response, nil
}

// Events listens for log change events
func (s *Server) Events(request *api.EventRequest, srv api.LogService_EventsServer) error {
	log.Tracef("Received EventRequest %+v", request)
	in, err := proto.Marshal(&ListenRequest{
		Replay: request.Replay,
		Index:  uint64(request.Index),
	})
	if err != nil {
		return err
	}

	stream := streams.NewBufferedStream()
	if err := s.CommandStream(srv.Context(), opEvents, in, request.Header, stream); err != nil {
		return err
	}

	for {
		result, ok := stream.Receive()
		if !ok {
			break
		}

		if result.Failed() {
			return result.Error
		}

		response := &ListenResponse{}
		output := result.Value.(server.SessionOutput)
		if err = proto.Unmarshal(output.Value.([]byte), response); err != nil {
			return err
		}

		var eventResponse *api.EventResponse
		switch output.Header.Type {
		case headers.ResponseType_OPEN_STREAM:
			eventResponse = &api.EventResponse{
				Header: output.Header,
			}
		case headers.ResponseType_CLOSE_STREAM:
			eventResponse = &api.EventResponse{
				Header: output.Header,
			}
		default:
			eventResponse = &api.EventResponse{
				Header:    output.Header,
				Type:      getEventType(response.Type),
				Index:     int64(response.Index),
				Value:     response.Value,
				Version:   int64(response.Version),
				Timestamp: response.Timestamp,
			}
		}

		log.Tracef("Sending EventResponse %+v", eventResponse)
		if err = srv.Send(eventResponse); err != nil {
			return err
		}
	}
	log.Tracef("Finished EventRequest %+v", request)
	return nil
}

// Clear removes all keys from the log
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

	serviceResponse := &ClearResponse{}
	if err = proto.Unmarshal(out, serviceResponse); err != nil {
		return nil, err
	}

	response := &api.ClearResponse{
		Header: header,
	}
	log.Tracef("Sending ClearResponse %+v", response)
	return response, nil
}

func getResponseStatus(status UpdateStatus) api.ResponseStatus {
	switch status {
	case UpdateStatus_OK:
		return api.ResponseStatus_OK
	case UpdateStatus_NOOP:
		return api.ResponseStatus_NOOP
	case UpdateStatus_PRECONDITION_FAILED:
		return api.ResponseStatus_PRECONDITION_FAILED
	case UpdateStatus_WRITE_LOCK:
		return api.ResponseStatus_WRITE_LOCK
	}
	return api.ResponseStatus_OK
}

func getEventType(eventType ListenResponse_Type) api.EventResponse_Type {
	switch eventType {
	case ListenResponse_NONE:
		return api.EventResponse_NONE
	case ListenResponse_APPENDED:
		return api.EventResponse_APPENDED
	case ListenResponse_REMOVED:
		return api.EventResponse_REMOVED
	}
	return api.EventResponse_NONE
}
