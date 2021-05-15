// Code generated by atomix-go-framework. DO NOT EDIT.
package log

import (
	"context"
	_log "github.com/atomix/atomix-api/go/atomix/primitive/log"
	"github.com/atomix/atomix-go-framework/pkg/atomix/driver/proxy/rsm"
	"github.com/atomix/atomix-go-framework/pkg/atomix/errors"
	"github.com/atomix/atomix-go-framework/pkg/atomix/logging"
	storage "github.com/atomix/atomix-go-framework/pkg/atomix/storage/protocol/rsm"
	streams "github.com/atomix/atomix-go-framework/pkg/atomix/stream"
	"github.com/golang/protobuf/proto"
)

const Type = "Log"

const (
	sizeOp       = "Size"
	appendOp     = "Append"
	getOp        = "Get"
	firstEntryOp = "FirstEntry"
	lastEntryOp  = "LastEntry"
	prevEntryOp  = "PrevEntry"
	nextEntryOp  = "NextEntry"
	removeOp     = "Remove"
	clearOp      = "Clear"
	eventsOp     = "Events"
	entriesOp    = "Entries"
)

// NewProxyServer creates a new ProxyServer
func NewProxyServer(client *rsm.Client) _log.LogServiceServer {
	return &ProxyServer{
		Client: client,
		log:    logging.GetLogger("atomix", "counter"),
	}
}

type ProxyServer struct {
	*rsm.Client
	log logging.Logger
}

func (s *ProxyServer) Size(ctx context.Context, request *_log.SizeRequest) (*_log.SizeResponse, error) {
	s.log.Debugf("Received SizeRequest %+v", request)
	input, err := proto.Marshal(request)
	if err != nil {
		s.log.Errorf("Request SizeRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	partition := s.PartitionBy([]byte(request.Headers.PrimitiveID.String()))

	service := storage.ServiceId{
		Type:      Type,
		Namespace: request.Headers.PrimitiveID.Namespace,
		Name:      request.Headers.PrimitiveID.Name,
	}
	output, err := partition.DoQuery(ctx, service, sizeOp, input)
	if err != nil {
		s.log.Errorf("Request SizeRequest failed: %v", err)
		return nil, errors.Proto(err)
	}

	response := &_log.SizeResponse{}
	err = proto.Unmarshal(output, response)
	if err != nil {
		s.log.Errorf("Request SizeRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	s.log.Debugf("Sending SizeResponse %+v", response)
	return response, nil
}

func (s *ProxyServer) Append(ctx context.Context, request *_log.AppendRequest) (*_log.AppendResponse, error) {
	s.log.Debugf("Received AppendRequest %+v", request)
	input, err := proto.Marshal(request)
	if err != nil {
		s.log.Errorf("Request AppendRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	partition := s.PartitionBy([]byte(request.Headers.PrimitiveID.String()))

	service := storage.ServiceId{
		Type:      Type,
		Namespace: request.Headers.PrimitiveID.Namespace,
		Name:      request.Headers.PrimitiveID.Name,
	}
	output, err := partition.DoCommand(ctx, service, appendOp, input)
	if err != nil {
		s.log.Errorf("Request AppendRequest failed: %v", err)
		return nil, errors.Proto(err)
	}

	response := &_log.AppendResponse{}
	err = proto.Unmarshal(output, response)
	if err != nil {
		s.log.Errorf("Request AppendRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	s.log.Debugf("Sending AppendResponse %+v", response)
	return response, nil
}

func (s *ProxyServer) Get(ctx context.Context, request *_log.GetRequest) (*_log.GetResponse, error) {
	s.log.Debugf("Received GetRequest %+v", request)
	input, err := proto.Marshal(request)
	if err != nil {
		s.log.Errorf("Request GetRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	partition := s.PartitionBy([]byte(request.Headers.PrimitiveID.String()))

	service := storage.ServiceId{
		Type:      Type,
		Namespace: request.Headers.PrimitiveID.Namespace,
		Name:      request.Headers.PrimitiveID.Name,
	}
	output, err := partition.DoQuery(ctx, service, getOp, input)
	if err != nil {
		s.log.Errorf("Request GetRequest failed: %v", err)
		return nil, errors.Proto(err)
	}

	response := &_log.GetResponse{}
	err = proto.Unmarshal(output, response)
	if err != nil {
		s.log.Errorf("Request GetRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	s.log.Debugf("Sending GetResponse %+v", response)
	return response, nil
}

func (s *ProxyServer) FirstEntry(ctx context.Context, request *_log.FirstEntryRequest) (*_log.FirstEntryResponse, error) {
	s.log.Debugf("Received FirstEntryRequest %+v", request)
	input, err := proto.Marshal(request)
	if err != nil {
		s.log.Errorf("Request FirstEntryRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	partition := s.PartitionBy([]byte(request.Headers.PrimitiveID.String()))

	service := storage.ServiceId{
		Type:      Type,
		Namespace: request.Headers.PrimitiveID.Namespace,
		Name:      request.Headers.PrimitiveID.Name,
	}
	output, err := partition.DoQuery(ctx, service, firstEntryOp, input)
	if err != nil {
		s.log.Errorf("Request FirstEntryRequest failed: %v", err)
		return nil, errors.Proto(err)
	}

	response := &_log.FirstEntryResponse{}
	err = proto.Unmarshal(output, response)
	if err != nil {
		s.log.Errorf("Request FirstEntryRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	s.log.Debugf("Sending FirstEntryResponse %+v", response)
	return response, nil
}

func (s *ProxyServer) LastEntry(ctx context.Context, request *_log.LastEntryRequest) (*_log.LastEntryResponse, error) {
	s.log.Debugf("Received LastEntryRequest %+v", request)
	input, err := proto.Marshal(request)
	if err != nil {
		s.log.Errorf("Request LastEntryRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	partition := s.PartitionBy([]byte(request.Headers.PrimitiveID.String()))

	service := storage.ServiceId{
		Type:      Type,
		Namespace: request.Headers.PrimitiveID.Namespace,
		Name:      request.Headers.PrimitiveID.Name,
	}
	output, err := partition.DoQuery(ctx, service, lastEntryOp, input)
	if err != nil {
		s.log.Errorf("Request LastEntryRequest failed: %v", err)
		return nil, errors.Proto(err)
	}

	response := &_log.LastEntryResponse{}
	err = proto.Unmarshal(output, response)
	if err != nil {
		s.log.Errorf("Request LastEntryRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	s.log.Debugf("Sending LastEntryResponse %+v", response)
	return response, nil
}

func (s *ProxyServer) PrevEntry(ctx context.Context, request *_log.PrevEntryRequest) (*_log.PrevEntryResponse, error) {
	s.log.Debugf("Received PrevEntryRequest %+v", request)
	input, err := proto.Marshal(request)
	if err != nil {
		s.log.Errorf("Request PrevEntryRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	partition := s.PartitionBy([]byte(request.Headers.PrimitiveID.String()))

	service := storage.ServiceId{
		Type:      Type,
		Namespace: request.Headers.PrimitiveID.Namespace,
		Name:      request.Headers.PrimitiveID.Name,
	}
	output, err := partition.DoQuery(ctx, service, prevEntryOp, input)
	if err != nil {
		s.log.Errorf("Request PrevEntryRequest failed: %v", err)
		return nil, errors.Proto(err)
	}

	response := &_log.PrevEntryResponse{}
	err = proto.Unmarshal(output, response)
	if err != nil {
		s.log.Errorf("Request PrevEntryRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	s.log.Debugf("Sending PrevEntryResponse %+v", response)
	return response, nil
}

func (s *ProxyServer) NextEntry(ctx context.Context, request *_log.NextEntryRequest) (*_log.NextEntryResponse, error) {
	s.log.Debugf("Received NextEntryRequest %+v", request)
	input, err := proto.Marshal(request)
	if err != nil {
		s.log.Errorf("Request NextEntryRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	partition := s.PartitionBy([]byte(request.Headers.PrimitiveID.String()))

	service := storage.ServiceId{
		Type:      Type,
		Namespace: request.Headers.PrimitiveID.Namespace,
		Name:      request.Headers.PrimitiveID.Name,
	}
	output, err := partition.DoQuery(ctx, service, nextEntryOp, input)
	if err != nil {
		s.log.Errorf("Request NextEntryRequest failed: %v", err)
		return nil, errors.Proto(err)
	}

	response := &_log.NextEntryResponse{}
	err = proto.Unmarshal(output, response)
	if err != nil {
		s.log.Errorf("Request NextEntryRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	s.log.Debugf("Sending NextEntryResponse %+v", response)
	return response, nil
}

func (s *ProxyServer) Remove(ctx context.Context, request *_log.RemoveRequest) (*_log.RemoveResponse, error) {
	s.log.Debugf("Received RemoveRequest %+v", request)
	input, err := proto.Marshal(request)
	if err != nil {
		s.log.Errorf("Request RemoveRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	partition := s.PartitionBy([]byte(request.Headers.PrimitiveID.String()))

	service := storage.ServiceId{
		Type:      Type,
		Namespace: request.Headers.PrimitiveID.Namespace,
		Name:      request.Headers.PrimitiveID.Name,
	}
	output, err := partition.DoCommand(ctx, service, removeOp, input)
	if err != nil {
		s.log.Errorf("Request RemoveRequest failed: %v", err)
		return nil, errors.Proto(err)
	}

	response := &_log.RemoveResponse{}
	err = proto.Unmarshal(output, response)
	if err != nil {
		s.log.Errorf("Request RemoveRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	s.log.Debugf("Sending RemoveResponse %+v", response)
	return response, nil
}

func (s *ProxyServer) Clear(ctx context.Context, request *_log.ClearRequest) (*_log.ClearResponse, error) {
	s.log.Debugf("Received ClearRequest %+v", request)
	input, err := proto.Marshal(request)
	if err != nil {
		s.log.Errorf("Request ClearRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	partition := s.PartitionBy([]byte(request.Headers.PrimitiveID.String()))

	service := storage.ServiceId{
		Type:      Type,
		Namespace: request.Headers.PrimitiveID.Namespace,
		Name:      request.Headers.PrimitiveID.Name,
	}
	output, err := partition.DoCommand(ctx, service, clearOp, input)
	if err != nil {
		s.log.Errorf("Request ClearRequest failed: %v", err)
		return nil, errors.Proto(err)
	}

	response := &_log.ClearResponse{}
	err = proto.Unmarshal(output, response)
	if err != nil {
		s.log.Errorf("Request ClearRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	s.log.Debugf("Sending ClearResponse %+v", response)
	return response, nil
}

func (s *ProxyServer) Events(request *_log.EventsRequest, srv _log.LogService_EventsServer) error {
	s.log.Debugf("Received EventsRequest %+v", request)
	input, err := proto.Marshal(request)
	if err != nil {
		s.log.Errorf("Request EventsRequest failed: %v", err)
		return errors.Proto(err)
	}

	stream := streams.NewBufferedStream()
	partition := s.PartitionBy([]byte(request.Headers.PrimitiveID.String()))

	service := storage.ServiceId{
		Type:      Type,
		Namespace: request.Headers.PrimitiveID.Namespace,
		Name:      request.Headers.PrimitiveID.Name,
	}
	err = partition.DoCommandStream(srv.Context(), service, eventsOp, input, stream)
	if err != nil {
		s.log.Errorf("Request EventsRequest failed: %v", err)
		return errors.Proto(err)
	}

	for {
		result, ok := stream.Receive()
		if !ok {
			break
		}

		if result.Failed() {
			s.log.Errorf("Request EventsRequest failed: %v", result.Error)
			return errors.Proto(result.Error)
		}

		response := &_log.EventsResponse{}
		err = proto.Unmarshal(result.Value.([]byte), response)
		if err != nil {
			s.log.Errorf("Request EventsRequest failed: %v", err)
			return errors.Proto(err)
		}

		s.log.Debugf("Sending EventsResponse %+v", response)
		if err = srv.Send(response); err != nil {
			s.log.Errorf("Response EventsResponse failed: %v", err)
			return errors.Proto(err)
		}
	}
	s.log.Debugf("Finished EventsRequest %+v", request)
	return nil
}

func (s *ProxyServer) Entries(request *_log.EntriesRequest, srv _log.LogService_EntriesServer) error {
	s.log.Debugf("Received EntriesRequest %+v", request)
	input, err := proto.Marshal(request)
	if err != nil {
		s.log.Errorf("Request EntriesRequest failed: %v", err)
		return errors.Proto(err)
	}

	stream := streams.NewBufferedStream()
	partition := s.PartitionBy([]byte(request.Headers.PrimitiveID.String()))

	service := storage.ServiceId{
		Type:      Type,
		Namespace: request.Headers.PrimitiveID.Namespace,
		Name:      request.Headers.PrimitiveID.Name,
	}
	err = partition.DoQueryStream(srv.Context(), service, entriesOp, input, stream)
	if err != nil {
		s.log.Errorf("Request EntriesRequest failed: %v", err)
		return errors.Proto(err)
	}

	for {
		result, ok := stream.Receive()
		if !ok {
			break
		}

		if result.Failed() {
			s.log.Errorf("Request EntriesRequest failed: %v", result.Error)
			return errors.Proto(result.Error)
		}

		response := &_log.EntriesResponse{}
		err = proto.Unmarshal(result.Value.([]byte), response)
		if err != nil {
			s.log.Errorf("Request EntriesRequest failed: %v", err)
			return errors.Proto(err)
		}

		s.log.Debugf("Sending EntriesResponse %+v", response)
		if err = srv.Send(response); err != nil {
			s.log.Errorf("Response EntriesResponse failed: %v", err)
			return errors.Proto(err)
		}
	}
	s.log.Debugf("Finished EntriesRequest %+v", request)
	return nil
}
