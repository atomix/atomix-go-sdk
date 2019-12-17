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

package server

import (
	"context"
	"errors"
	"github.com/atomix/atomix-api/proto/atomix/headers"
	"github.com/atomix/atomix-go-node/pkg/atomix/node"
	"github.com/atomix/atomix-go-node/pkg/atomix/service"
	streams "github.com/atomix/atomix-go-node/pkg/atomix/stream"
	"github.com/golang/protobuf/proto"
	"time"
)

// SessionizedServer is a base server for servers that support sessions
type SessionizedServer struct {
	Client node.Client
	Type   string
}

// Write sends a write to the service
func (s *SessionizedServer) Write(ctx context.Context, request []byte, header *headers.RequestHeader) ([]byte, error) {
	serviceRequest := &service.ServiceRequest{
		Id: &service.ServiceId{
			Type:      s.Type,
			Name:      header.Name.Name,
			Namespace: header.Name.Namespace,
		},
		Request: &service.ServiceRequest_Command{
			Command: request,
		},
	}

	bytes, err := proto.Marshal(serviceRequest)
	if err != nil {
		return nil, err
	}

	// Create a write channel
	ch := make(chan streams.Result)

	// Write the request
	if err := s.Client.Write(ctx, bytes, streams.NewChannelStream(ch)); err != nil {
		return nil, err
	}

	// Wait for the result
	result, ok := <-ch
	if !ok {
		return nil, errors.New("write channel closed")
	}

	// If the result failed, return the error
	if result.Failed() {
		return nil, result.Error
	}

	// Decode and return the response
	serviceResponse := &service.ServiceResponse{}
	err = proto.Unmarshal(result.Value, serviceResponse)
	if err != nil {
		return nil, err
	}
	return serviceResponse.GetCommand(), nil
}

// WriteStream sends a streaming write to the service
func (s *SessionizedServer) WriteStream(request []byte, header *headers.RequestHeader, ch chan<- streams.Result) error {
	serviceRequest := &service.ServiceRequest{
		Id: &service.ServiceId{
			Type:      s.Type,
			Name:      header.Name.Name,
			Namespace: header.Name.Namespace,
		},
		Request: &service.ServiceRequest_Command{
			Command: request,
		},
	}

	bytes, err := proto.Marshal(serviceRequest)
	if err != nil {
		return err
	}

	streamCh := make(chan streams.Result)
	if err := s.Client.Write(context.TODO(), bytes, streams.NewChannelStream(streamCh)); err != nil {
		return err
	}

	// Create a goroutine to convert the results into raw form
	go func() {
		defer close(ch)
		for result := range streamCh {
			if result.Failed() {
				ch <- result
			} else {
				serviceResponse := &service.ServiceResponse{}
				err := proto.Unmarshal(result.Value, serviceResponse)
				if err != nil {
					ch <- streams.Result{
						Error: err,
					}
				} else {
					ch <- streams.Result{
						Value: serviceResponse.GetCommand(),
					}
				}
			}
		}
	}()

	return nil
}

// Read sends a read to the service
func (s *SessionizedServer) Read(ctx context.Context, request []byte, header *headers.RequestHeader) ([]byte, error) {
	serviceRequest := &service.ServiceRequest{
		Id: &service.ServiceId{
			Type:      s.Type,
			Name:      header.Name.Name,
			Namespace: header.Name.Namespace,
		},
		Request: &service.ServiceRequest_Query{
			Query: request,
		},
	}

	bytes, err := proto.Marshal(serviceRequest)
	if err != nil {
		return nil, err
	}

	// Create a read channel
	ch := make(chan streams.Result)

	// Read the request
	if err := s.Client.Read(ctx, bytes, streams.NewChannelStream(ch)); err != nil {
		return nil, err
	}

	// Wait for the result
	result, ok := <-ch
	if !ok {
		return nil, errors.New("write channel closed")
	}

	// If the result failed, return the error
	if result.Failed() {
		return nil, result.Error
	}

	serviceResponse := &service.ServiceResponse{}
	err = proto.Unmarshal(result.Value, serviceResponse)
	if err != nil {
		return nil, err
	}
	return serviceResponse.GetQuery(), nil
}

// ReadStream sends a streaming read to the service
func (s *SessionizedServer) ReadStream(request []byte, header *headers.RequestHeader, ch chan<- streams.Result) error {
	serviceRequest := &service.ServiceRequest{
		Id: &service.ServiceId{
			Type:      s.Type,
			Name:      header.Name.Name,
			Namespace: header.Name.Namespace,
		},
		Request: &service.ServiceRequest_Query{
			Query: request,
		},
	}

	bytes, err := proto.Marshal(serviceRequest)
	if err != nil {
		return err
	}

	streamCh := make(chan streams.Result)
	if err := s.Client.Read(context.TODO(), bytes, streams.NewChannelStream(streamCh)); err != nil {
		return err
	}

	// Create a goroutine to convert the results into raw form
	go func() {
		defer close(ch)
		for result := range streamCh {
			if result.Failed() {
				ch <- result
			} else {
				serviceResponse := &service.ServiceResponse{}
				err := proto.Unmarshal(result.Value, serviceResponse)
				if err != nil {
					ch <- streams.Result{
						Error: err,
					}
				} else {
					ch <- streams.Result{
						Value: serviceResponse.GetQuery(),
					}
				}
			}
		}
	}()

	return nil
}

// Command submits a command to the service
func (s *SessionizedServer) Command(ctx context.Context, name string, input []byte, header *headers.RequestHeader) ([]byte, *headers.ResponseHeader, error) {
	sessionRequest := &service.SessionRequest{
		Request: &service.SessionRequest_Command{
			Command: &service.SessionCommandRequest{
				Context: &service.SessionCommandContext{
					SessionID:      header.SessionID,
					SequenceNumber: header.RequestID,
				},
				Name:  name,
				Input: input,
			},
		},
	}

	bytes, err := proto.Marshal(sessionRequest)
	if err != nil {
		return nil, nil, err
	}

	bytes, err = s.Write(ctx, bytes, header)
	if err != nil {
		return nil, nil, err
	}

	sessionResponse := &service.SessionResponse{}
	err = proto.Unmarshal(bytes, sessionResponse)
	if err != nil {
		return nil, nil, err
	}

	commandResponse := sessionResponse.GetCommand()
	responseHeader := &headers.ResponseHeader{
		SessionID:  header.SessionID,
		StreamID:   commandResponse.Context.StreamID,
		ResponseID: commandResponse.Context.Sequence,
		Index:      commandResponse.Context.Index,
	}
	return commandResponse.Output, responseHeader, nil
}

// CommandStream submits a streaming command to the service
func (s *SessionizedServer) CommandStream(name string, input []byte, header *headers.RequestHeader, ch chan<- SessionOutput) error {
	sessionRequest := &service.SessionRequest{
		Request: &service.SessionRequest_Command{
			Command: &service.SessionCommandRequest{
				Context: &service.SessionCommandContext{
					SessionID:      header.SessionID,
					SequenceNumber: header.RequestID,
				},
				Name:  name,
				Input: input,
			},
		},
	}

	bytes, err := proto.Marshal(sessionRequest)
	if err != nil {
		return err
	}

	resultCh := make(chan streams.Result)
	if err = s.WriteStream(bytes, header, resultCh); err != nil {
		return err
	}

	go func() {
		defer close(ch)
		for result := range resultCh {
			if result.Failed() {
				ch <- SessionOutput{
					Result: result,
				}
			} else {
				sessionResponse := &service.SessionResponse{}
				err = proto.Unmarshal(result.Value, sessionResponse)
				if err != nil {
					ch <- SessionOutput{
						Result: streams.Result{
							Error: err,
						},
					}
				} else {
					commandResponse := sessionResponse.GetCommand()
					responseHeader := &headers.ResponseHeader{
						SessionID:  header.SessionID,
						StreamID:   commandResponse.Context.StreamID,
						ResponseID: commandResponse.Context.Sequence,
						Index:      commandResponse.Context.Index,
					}
					ch <- SessionOutput{
						Header: responseHeader,
						Result: streams.Result{
							Value: commandResponse.Output,
						},
					}
				}
			}
		}
	}()

	return nil
}

// Query submits a query to the service
func (s *SessionizedServer) Query(ctx context.Context, name string, input []byte, header *headers.RequestHeader) ([]byte, *headers.ResponseHeader, error) {
	sessionRequest := &service.SessionRequest{
		Request: &service.SessionRequest_Query{
			Query: &service.SessionQueryRequest{
				Context: &service.SessionQueryContext{
					SessionID:          header.SessionID,
					LastIndex:          header.Index,
					LastSequenceNumber: header.RequestID,
				},
				Name:  name,
				Input: input,
			},
		},
	}

	bytes, err := proto.Marshal(sessionRequest)
	if err != nil {
		return nil, nil, err
	}

	bytes, err = s.Read(ctx, bytes, header)
	if err != nil {
		return nil, nil, err
	}

	sessionResponse := &service.SessionResponse{}
	err = proto.Unmarshal(bytes, sessionResponse)
	if err != nil {
		return nil, nil, err
	}

	queryResponse := sessionResponse.GetQuery()
	responseHeader := &headers.ResponseHeader{
		SessionID: header.SessionID,
		Index:     queryResponse.Context.Index,
	}
	return queryResponse.Output, responseHeader, nil
}

// QueryStream submits a streaming query to the service
func (s *SessionizedServer) QueryStream(name string, input []byte, header *headers.RequestHeader, ch chan<- SessionOutput) error {
	sessionRequest := &service.SessionRequest{
		Request: &service.SessionRequest_Query{
			Query: &service.SessionQueryRequest{
				Context: &service.SessionQueryContext{
					SessionID:          header.SessionID,
					LastIndex:          header.Index,
					LastSequenceNumber: header.RequestID,
				},
				Name:  name,
				Input: input,
			},
		},
	}

	bytes, err := proto.Marshal(sessionRequest)
	if err != nil {
		return err
	}

	resultCh := make(chan streams.Result)
	if err = s.ReadStream(bytes, header, resultCh); err != nil {
		return err
	}

	go func() {
		defer close(ch)
		for result := range resultCh {
			if result.Failed() {
				ch <- SessionOutput{
					Result: result,
				}
			} else {
				sessionResponse := &service.SessionResponse{}
				err = proto.Unmarshal(result.Value, sessionResponse)
				if err != nil {
					ch <- SessionOutput{
						Result: streams.Result{
							Error: err,
						},
					}
				} else {
					queryResponse := sessionResponse.GetQuery()
					responseHeader := &headers.ResponseHeader{
						SessionID: header.SessionID,
						Index:     queryResponse.Context.Index,
					}
					ch <- SessionOutput{
						Header: responseHeader,
						Result: streams.Result{
							Value: queryResponse.Output,
						},
					}
				}
			}
		}
	}()

	return nil
}

// OpenSession opens a new session
func (s *SessionizedServer) OpenSession(ctx context.Context, header *headers.RequestHeader, timeout *time.Duration) (uint64, error) {
	sessionRequest := &service.SessionRequest{
		Request: &service.SessionRequest_OpenSession{
			OpenSession: &service.OpenSessionRequest{
				Timeout: timeout,
			},
		},
	}

	bytes, err := proto.Marshal(sessionRequest)
	if err != nil {
		return 0, err
	}

	bytes, err = s.Write(ctx, bytes, header)
	if err != nil {
		return 0, err
	}

	sessionResponse := &service.SessionResponse{}
	err = proto.Unmarshal(bytes, sessionResponse)
	if err != nil {
		return 0, err
	}

	return sessionResponse.GetOpenSession().SessionID, nil
}

// KeepAliveSession keeps a session alive
func (s *SessionizedServer) KeepAliveSession(ctx context.Context, header *headers.RequestHeader) error {
	streams := make(map[uint64]uint64)
	for _, stream := range header.Streams {
		streams[stream.StreamID] = stream.ResponseID
	}

	sessionRequest := &service.SessionRequest{
		Request: &service.SessionRequest_KeepAlive{
			KeepAlive: &service.KeepAliveRequest{
				SessionID:       header.SessionID,
				CommandSequence: header.RequestID,
				Streams:         streams,
			},
		},
	}

	bytes, err := proto.Marshal(sessionRequest)
	if err != nil {
		return err
	}

	bytes, err = s.Write(ctx, bytes, header)
	if err != nil {
		return err
	}

	sessionResponse := &service.SessionResponse{}
	return proto.Unmarshal(bytes, sessionResponse)
}

// CloseSession closes a session
func (s *SessionizedServer) CloseSession(ctx context.Context, header *headers.RequestHeader) error {
	sessionRequest := &service.SessionRequest{
		Request: &service.SessionRequest_CloseSession{
			CloseSession: &service.CloseSessionRequest{
				SessionID: header.SessionID,
			},
		},
	}

	bytes, err := proto.Marshal(sessionRequest)
	if err != nil {
		return err
	}

	bytes, err = s.Write(ctx, bytes, header)
	if err != nil {
		return err
	}

	sessionResponse := &service.SessionResponse{}
	return proto.Unmarshal(bytes, sessionResponse)
}

// Delete deletes the service
func (s *SessionizedServer) Delete(ctx context.Context, header *headers.RequestHeader) error {
	serviceRequest := &service.ServiceRequest{
		Id: &service.ServiceId{
			Type:      s.Type,
			Name:      header.Name.Name,
			Namespace: header.Name.Namespace,
		},
		Request: &service.ServiceRequest_Delete{
			Delete: &service.DeleteRequest{},
		},
	}

	bytes, err := proto.Marshal(serviceRequest)
	if err != nil {
		return err
	}

	// Create a write channel
	ch := make(chan streams.Result)

	// Write the request
	if err := s.Client.Write(ctx, bytes, streams.NewChannelStream(ch)); err != nil {
		return err
	}

	// Wait for the result
	result, ok := <-ch
	if !ok {
		return errors.New("write channel closed")
	}

	// If the result failed, return the error
	if result.Failed() {
		return result.Error
	}

	// Decode and return the response
	serviceResponse := &service.ServiceResponse{}
	err = proto.Unmarshal(result.Value, serviceResponse)
	if err != nil {
		return err
	}
	return nil
}

// SessionOutput is a result for session-supporting servers containing session header information
type SessionOutput struct {
	streams.Result
	Header *headers.ResponseHeader
}
