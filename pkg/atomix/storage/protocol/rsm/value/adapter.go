// Code generated by atomix-go-framework. DO NOT EDIT.
package value

import (
	"github.com/atomix/atomix-go-framework/pkg/atomix/errors"
	"github.com/atomix/atomix-go-framework/pkg/atomix/logging"
	"github.com/atomix/atomix-go-framework/pkg/atomix/storage/protocol/rsm"
	"github.com/gogo/protobuf/proto"
	"io"
)

var log = logging.GetLogger("atomix", "value", "service")

const Type = "Value"

const (
	setOp    = "Set"
	getOp    = "Get"
	eventsOp = "Events"
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
	case 1:
		p, err := newSetProposal(command)
		if err != nil {
			err = errors.NewInternal(err.Error())
			log.Error(err)
			command.Output(nil, err)
			return
		}

		log.Debugf("Proposal SetProposal %s", p)
		response, err := s.rsm.Set(p)
		if err != nil {
			log.Warnf("Proposal SetProposal %s failed: %v", p, err)
			command.Output(nil, err)
		} else {
			output, err := proto.Marshal(response)
			if err != nil {
				err = errors.NewInternal(err.Error())
				log.Errorf("Proposal SetProposal %s failed: %v", p, err)
				command.Output(nil, err)
			} else {
				log.Errorf("Proposal SetProposal %s complete: %+v", p, response)
				command.Output(output, nil)
			}
		}
	case 3:
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
	case 2:
		q, err := newGetQuery(query)
		if err != nil {
			err = errors.NewInternal(err.Error())
			log.Error(err)
			query.Output(nil, err)
			return
		}

		log.Debugf("Querying GetQuery %s", q)
		response, err := s.rsm.Get(q)
		if err != nil {
			log.Warnf("Querying GetQuery %s failed: %v", q, err)
			query.Output(nil, err)
		} else {
			output, err := proto.Marshal(response)
			if err != nil {
				err = errors.NewInternal(err.Error())
				log.Errorf("Querying GetQuery %s failed: %v", q, err)
				query.Output(nil, err)
			} else {
				log.Errorf("Querying GetQuery %s complete: %+v", q, response)
				query.Output(output, nil)
			}
		}
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
