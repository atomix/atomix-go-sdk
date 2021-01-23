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

package value

import (
	"context"
	valueapi "github.com/atomix/api/go/atomix/primitive/value"
	"github.com/atomix/go-framework/pkg/atomix/errors"
	"github.com/atomix/go-framework/pkg/atomix/meta"
	"github.com/atomix/go-framework/pkg/atomix/time"
	"sync"
)

func init() {
	registerService(func(protocol Protocol) Service {
		clock := time.NewLogicalClock()
		return &valueService{
			protocol: protocol,
			clock:    clock,
			value: &valueapi.Value{
				ObjectMeta: meta.NewTimestamped(clock.Get()).AsTombstone().Proto(),
			},
		}
	})
}

type valueService struct {
	protocol Protocol
	clock    time.Clock
	value    *valueapi.Value
	streams  []chan<- valueapi.EventsResponse
	mu       sync.RWMutex
}

func (s *valueService) Protocol() Protocol {
	return s.protocol
}

func (s *valueService) Read(ctx context.Context) (*valueapi.Value, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.value, nil
}

func (s *valueService) Update(ctx context.Context, value *valueapi.Value) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if meta.FromProto(value.ObjectMeta).After(meta.FromProto(s.value.ObjectMeta)) {
		s.value = value
	}
	return nil
}

func (s *valueService) Set(ctx context.Context, input *valueapi.SetRequest) (*valueapi.SetResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.value != nil && meta.FromProto(s.value.ObjectMeta).After(meta.FromProto(input.Value.ObjectMeta)) {
		return &valueapi.SetResponse{
			Value: *s.value,
		}, nil
	}

	err := checkPreconditions(s.value, input.Preconditions)
	if err != nil {
		return nil, err
	}

	s.value = &input.Value

	s.notify(valueapi.Event{
		Type:  valueapi.Event_UPDATE,
		Value: input.Value,
	})

	err = s.Protocol().Broadcast(ctx, &input.Value)
	if err != nil {
		return nil, err
	}
	return &valueapi.SetResponse{
		Value: input.Value,
	}, nil
}

func (s *valueService) Get(ctx context.Context, input *valueapi.GetRequest) (*valueapi.GetResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var value valueapi.Value
	if s.value != nil {
		value = *s.value
	}
	return &valueapi.GetResponse{
		Value: value,
	}, nil
}

func (s *valueService) Events(ctx context.Context, request *valueapi.EventsRequest, ch chan<- valueapi.EventsResponse) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.streams = append(s.streams, ch)
	ch <- valueapi.EventsResponse{}
	return nil
}

func (s *valueService) notify(event valueapi.Event) {
	for _, ch := range s.streams {
		ch <- valueapi.EventsResponse{
			Event: event,
		}
	}
}

func checkPreconditions(value *valueapi.Value, preconditions []valueapi.Precondition) error {
	for _, precondition := range preconditions {
		switch p := precondition.Precondition.(type) {
		case *valueapi.Precondition_Metadata:
			if value == nil {
				return errors.NewConflict("metadata precondition failed")
			}
			if !meta.FromProto(value.ObjectMeta).Equal(meta.FromProto(*p.Metadata)) {
				return errors.NewConflict("metadata mismatch")
			}
		}
	}
	return nil
}
