// Code generated by atomix-go-framework. DO NOT EDIT.
package value

import (
	"context"
	value "github.com/atomix/atomix-api/go/atomix/primitive/value"
	"github.com/atomix/atomix-go-framework/pkg/atomix/driver/proxy/rsm"
	"github.com/atomix/atomix-go-framework/pkg/atomix/errors"
	"github.com/atomix/atomix-go-framework/pkg/atomix/logging"
	storage "github.com/atomix/atomix-go-framework/pkg/atomix/storage/protocol/rsm"
	streams "github.com/atomix/atomix-go-framework/pkg/atomix/stream"
	"github.com/golang/protobuf/proto"
)

const Type = "Value"

const (
	setOp    = "Set"
	getOp    = "Get"
	eventsOp = "Events"
)

// NewProxyServer creates a new ProxyServer
func NewProxyServer(client *rsm.Client) value.ValueServiceServer {
	return &ProxyServer{
		Client: client,
		log:    logging.GetLogger("atomix", "proxy", "value"),
	}
}

type ProxyServer struct {
	*rsm.Client
	log logging.Logger
}

func (s *ProxyServer) Set(ctx context.Context, request *value.SetRequest) (*value.SetResponse, error) {
	s.log.Debugf("Received SetRequest %+v", request)
	input, err := proto.Marshal(request)
	if err != nil {
		s.log.Errorf("Request SetRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	clusterKey := request.Headers.ClusterKey
	if clusterKey == "" {
		clusterKey = request.Headers.PrimitiveID.String()
	}
	partition := s.PartitionBy([]byte(clusterKey))

	service := storage.ServiceId{
		Type:    Type,
		Cluster: request.Headers.ClusterKey,
		Name:    request.Headers.PrimitiveID.Name,
	}
	output, err := partition.DoCommand(ctx, service, setOp, input)
	if err != nil {
		s.log.Errorf("Request SetRequest failed: %v", err)
		return nil, errors.Proto(err)
	}

	response := &value.SetResponse{}
	err = proto.Unmarshal(output, response)
	if err != nil {
		s.log.Errorf("Request SetRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	s.log.Debugf("Sending SetResponse %+v", response)
	return response, nil
}

func (s *ProxyServer) Get(ctx context.Context, request *value.GetRequest) (*value.GetResponse, error) {
	s.log.Debugf("Received GetRequest %+v", request)
	input, err := proto.Marshal(request)
	if err != nil {
		s.log.Errorf("Request GetRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	clusterKey := request.Headers.ClusterKey
	if clusterKey == "" {
		clusterKey = request.Headers.PrimitiveID.String()
	}
	partition := s.PartitionBy([]byte(clusterKey))

	service := storage.ServiceId{
		Type:    Type,
		Cluster: request.Headers.ClusterKey,
		Name:    request.Headers.PrimitiveID.Name,
	}
	output, err := partition.DoQuery(ctx, service, getOp, input)
	if err != nil {
		s.log.Errorf("Request GetRequest failed: %v", err)
		return nil, errors.Proto(err)
	}

	response := &value.GetResponse{}
	err = proto.Unmarshal(output, response)
	if err != nil {
		s.log.Errorf("Request GetRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	s.log.Debugf("Sending GetResponse %+v", response)
	return response, nil
}

func (s *ProxyServer) Events(request *value.EventsRequest, srv value.ValueService_EventsServer) error {
	s.log.Debugf("Received EventsRequest %+v", request)
	input, err := proto.Marshal(request)
	if err != nil {
		s.log.Errorf("Request EventsRequest failed: %v", err)
		return errors.Proto(err)
	}

	ch := make(chan streams.Result)
	stream := streams.NewChannelStream(ch)
	clusterKey := request.Headers.ClusterKey
	if clusterKey == "" {
		clusterKey = request.Headers.PrimitiveID.String()
	}
	partition := s.PartitionBy([]byte(clusterKey))

	service := storage.ServiceId{
		Type:    Type,
		Cluster: request.Headers.ClusterKey,
		Name:    request.Headers.PrimitiveID.Name,
	}
	err = partition.DoCommandStream(srv.Context(), service, eventsOp, input, stream)
	if err != nil {
		s.log.Errorf("Request EventsRequest failed: %v", err)
		return errors.Proto(err)
	}

	for result := range ch {
		if result.Failed() {
			if result.Error == context.Canceled {
				break
			}
			s.log.Errorf("Request EventsRequest failed: %v", result.Error)
			return errors.Proto(result.Error)
		}

		response := &value.EventsResponse{}
		err = proto.Unmarshal(result.Value.([]byte), response)
		if err != nil {
			s.log.Errorf("Request EventsRequest failed: %v", err)
			return errors.Proto(err)
		}

		s.log.Debugf("Sending EventsResponse %+v", response)
		if err = srv.Send(response); err != nil {
			s.log.Errorf("Response EventsResponse failed: %v", err)
			return err
		}
	}
	s.log.Debugf("Finished EventsRequest %+v", request)
	return nil
}
