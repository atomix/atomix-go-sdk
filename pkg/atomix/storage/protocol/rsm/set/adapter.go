// Code generated by atomix-go-framework. DO NOT EDIT.
package set

import (
	"github.com/atomix/atomix-go-framework/pkg/atomix/errors"
	"github.com/atomix/atomix-go-framework/pkg/atomix/logging"
	"github.com/atomix/atomix-go-framework/pkg/atomix/storage/protocol/rsm"
	"github.com/gogo/protobuf/proto"
	"io"
)

var log = logging.GetLogger("atomix", "set", "service")

const Type = "Set"

const (
	sizeOp     = "Size"
	containsOp = "Contains"
	addOp      = "Add"
	removeOp   = "Remove"
	clearOp    = "Clear"
	eventsOp   = "Events"
	elementsOp = "Elements"
)

var newServiceFunc rsm.NewServiceFunc

func registerServiceFunc(rsmf NewServiceFunc) {
	newServiceFunc = func(context rsm.ServiceContext) rsm.Service {
		return &ServiceAdaptor{
			ServiceContext: context,
			rsm:            rsmf(newServiceContext(context)),
		}
	}
}

type NewServiceFunc func(ServiceContext) Service

// RegisterService registers the election primitive service on the given node
func RegisterService(node *rsm.Node) {
	node.RegisterService(Type, newServiceFunc)
}

type ServiceAdaptor struct {
	rsm.ServiceContext
	rsm Service
}

func (s *ServiceAdaptor) ExecuteCommand(command rsm.Command) {
	switch command.OperationID() {
	case 3:
		p, err := newAddProposal(command)
		if err != nil {
			err = errors.NewInternal(err.Error())
			log.Error(err)
			command.Output(nil, err)
			return
		}

		log.Debugf("Proposal AddProposal %s", p)
		response, err := s.rsm.Add(p)
		if err != nil {
			log.Warnf("Proposal AddProposal %s failed: %v", p, err)
			command.Output(nil, err)
		} else {
			output, err := proto.Marshal(response)
			if err != nil {
				err = errors.NewInternal(err.Error())
				log.Errorf("Proposal AddProposal %s failed: %v", p, err)
				command.Output(nil, err)
			} else {
				log.Errorf("Proposal AddProposal %s complete: %+v", p, response)
				command.Output(output, nil)
			}
		}
	case 4:
		p, err := newRemoveProposal(command)
		if err != nil {
			err = errors.NewInternal(err.Error())
			log.Error(err)
			command.Output(nil, err)
			return
		}

		log.Debugf("Proposal RemoveProposal %s", p)
		response, err := s.rsm.Remove(p)
		if err != nil {
			log.Warnf("Proposal RemoveProposal %s failed: %v", p, err)
			command.Output(nil, err)
		} else {
			output, err := proto.Marshal(response)
			if err != nil {
				err = errors.NewInternal(err.Error())
				log.Errorf("Proposal RemoveProposal %s failed: %v", p, err)
				command.Output(nil, err)
			} else {
				log.Errorf("Proposal RemoveProposal %s complete: %+v", p, response)
				command.Output(output, nil)
			}
		}
	case 5:
		p, err := newClearProposal(command)
		if err != nil {
			err = errors.NewInternal(err.Error())
			log.Error(err)
			command.Output(nil, err)
			return
		}

		log.Debugf("Proposal ClearProposal %s", p)
		response, err := s.rsm.Clear(p)
		if err != nil {
			log.Warnf("Proposal ClearProposal %s failed: %v", p, err)
			command.Output(nil, err)
		} else {
			output, err := proto.Marshal(response)
			if err != nil {
				err = errors.NewInternal(err.Error())
				log.Errorf("Proposal ClearProposal %s failed: %v", p, err)
				command.Output(nil, err)
			} else {
				log.Errorf("Proposal ClearProposal %s complete: %+v", p, response)
				command.Output(output, nil)
			}
		}
	case 6:
		p, err := newEventsProposal(command)
		if err != nil {
			err = errors.NewInternal(err.Error())
			log.Error(err)
			command.Output(nil, err)
			return
		}

		log.Debugf("Proposal EventsProposal %s", p)
		s.rsm.Events(p)
	default:
		err := errors.NewNotSupported("unknown operation %d", command.OperationID())
		log.Warn(err)
		command.Output(nil, err)
	}
}

func (s *ServiceAdaptor) ExecuteQuery(query rsm.Query) {
	switch query.OperationID() {
	case 1:
		q, err := newSizeQuery(query)
		if err != nil {
			err = errors.NewInternal(err.Error())
			log.Error(err)
			query.Output(nil, err)
			return
		}

		log.Debugf("Querying SizeQuery %s", q)
		response, err := s.rsm.Size(q)
		if err != nil {
			log.Warnf("Querying SizeQuery %s failed: %v", q, err)
			query.Output(nil, err)
		} else {
			output, err := proto.Marshal(response)
			if err != nil {
				err = errors.NewInternal(err.Error())
				log.Errorf("Querying SizeQuery %s failed: %v", q, err)
				query.Output(nil, err)
			} else {
				log.Errorf("Querying SizeQuery %s complete: %+v", q, response)
				query.Output(output, nil)
			}
		}
	case 2:
		q, err := newContainsQuery(query)
		if err != nil {
			err = errors.NewInternal(err.Error())
			log.Error(err)
			query.Output(nil, err)
			return
		}

		log.Debugf("Querying ContainsQuery %s", q)
		response, err := s.rsm.Contains(q)
		if err != nil {
			log.Warnf("Querying ContainsQuery %s failed: %v", q, err)
			query.Output(nil, err)
		} else {
			output, err := proto.Marshal(response)
			if err != nil {
				err = errors.NewInternal(err.Error())
				log.Errorf("Querying ContainsQuery %s failed: %v", q, err)
				query.Output(nil, err)
			} else {
				log.Errorf("Querying ContainsQuery %s complete: %+v", q, response)
				query.Output(output, nil)
			}
		}
	case 7:
		q, err := newElementsQuery(query)
		if err != nil {
			err = errors.NewInternal(err.Error())
			log.Error(err)
			query.Output(nil, err)
			return
		}

		log.Debugf("Querying ElementsQuery %s", q)
		s.rsm.Elements(q)
	default:
		err := errors.NewNotSupported("unknown operation %d", query.OperationID())
		log.Warn(err)
		query.Output(nil, err)
	}
}
func (s *ServiceAdaptor) Backup(writer io.Writer) error {
	err := s.rsm.Backup(newSnapshotWriter(writer))
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (s *ServiceAdaptor) Restore(reader io.Reader) error {
	err := s.rsm.Restore(newSnapshotReader(reader))
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

var _ rsm.Service = &ServiceAdaptor{}
