package counter

import (
	"context"
	counter "github.com/atomix/api/go/atomix/primitive/counter"
	"github.com/atomix/go-framework/pkg/atomix/errors"
	"github.com/atomix/go-framework/pkg/atomix/logging"
	"github.com/atomix/go-framework/pkg/atomix/proxy/crdt"
	"google.golang.org/grpc"
)

const Type = "Counter"

// RegisterProxy registers the primitive on the given node
func RegisterProxy(node *crdt.Node) {
	node.RegisterProxy(func(server *grpc.Server, client *crdt.Client) {
		counter.RegisterCounterServiceServer(server, &Proxy{
			Proxy: crdt.NewProxy(client),
			log:   logging.GetLogger("atomix", "counter"),
		})
	})
}

type Proxy struct {
	*crdt.Proxy
	log logging.Logger
}

func (s *Proxy) Set(ctx context.Context, request *counter.SetRequest) (*counter.SetResponse, error) {
	s.log.Debugf("Received SetRequest %+v", request)
	partition, err := s.PartitionFrom(ctx)
	if err != nil {
		return nil, errors.Proto(err)
	}

	conn, err := partition.Connect()
	if err != nil {
		return nil, errors.Proto(err)
	}

	client := counter.NewCounterServiceClient(conn)
	ctx = partition.AddHeaders(ctx)
	response, err := client.Set(ctx, request)
	if err != nil {
		s.log.Errorf("Request SetRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	s.log.Debugf("Sending SetResponse %+v", response)
	return response, nil
}

func (s *Proxy) Get(ctx context.Context, request *counter.GetRequest) (*counter.GetResponse, error) {
	s.log.Debugf("Received GetRequest %+v", request)
	partition, err := s.PartitionFrom(ctx)
	if err != nil {
		return nil, errors.Proto(err)
	}

	conn, err := partition.Connect()
	if err != nil {
		return nil, errors.Proto(err)
	}

	client := counter.NewCounterServiceClient(conn)
	ctx = partition.AddHeaders(ctx)
	response, err := client.Get(ctx, request)
	if err != nil {
		s.log.Errorf("Request GetRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	s.log.Debugf("Sending GetResponse %+v", response)
	return response, nil
}

func (s *Proxy) Increment(ctx context.Context, request *counter.IncrementRequest) (*counter.IncrementResponse, error) {
	s.log.Debugf("Received IncrementRequest %+v", request)
	partition, err := s.PartitionFrom(ctx)
	if err != nil {
		return nil, errors.Proto(err)
	}

	conn, err := partition.Connect()
	if err != nil {
		return nil, errors.Proto(err)
	}

	client := counter.NewCounterServiceClient(conn)
	ctx = partition.AddHeaders(ctx)
	response, err := client.Increment(ctx, request)
	if err != nil {
		s.log.Errorf("Request IncrementRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	s.log.Debugf("Sending IncrementResponse %+v", response)
	return response, nil
}

func (s *Proxy) Decrement(ctx context.Context, request *counter.DecrementRequest) (*counter.DecrementResponse, error) {
	s.log.Debugf("Received DecrementRequest %+v", request)
	partition, err := s.PartitionFrom(ctx)
	if err != nil {
		return nil, errors.Proto(err)
	}

	conn, err := partition.Connect()
	if err != nil {
		return nil, errors.Proto(err)
	}

	client := counter.NewCounterServiceClient(conn)
	ctx = partition.AddHeaders(ctx)
	response, err := client.Decrement(ctx, request)
	if err != nil {
		s.log.Errorf("Request DecrementRequest failed: %v", err)
		return nil, errors.Proto(err)
	}
	s.log.Debugf("Sending DecrementResponse %+v", response)
	return response, nil
}
