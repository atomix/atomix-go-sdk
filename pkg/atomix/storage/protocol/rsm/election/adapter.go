// Code generated by atomix-go-framework. DO NOT EDIT.
package election

import (
	"github.com/atomix/atomix-go-framework/pkg/atomix/errors"
	"github.com/atomix/atomix-go-framework/pkg/atomix/logging"
	"github.com/atomix/atomix-go-framework/pkg/atomix/storage/protocol/rsm"
	"github.com/gogo/protobuf/proto"
	"io"
)

var log = logging.GetLogger("atomix", "election", "service")

const Type = "Election"

const (
	enterOp    = "Enter"
	withdrawOp = "Withdraw"
	anointOp   = "Anoint"
	promoteOp  = "Promote"
	evictOp    = "Evict"
	getTermOp  = "GetTerm"
	eventsOp   = "Events"
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
		p, err := newEnterProposal(command)
		if err != nil {
			err = errors.NewInternal(err.Error())
			log.Error(err)
			command.Output(nil, err)
			return
		}

		log.Debugf("Proposal EnterProposal %s", p)
		response, err := s.rsm.Enter(p)
		if err != nil {
			log.Warnf("Proposal EnterProposal %s failed: %v", p, err)
			command.Output(nil, err)
		} else {
			output, err := proto.Marshal(response)
			if err != nil {
				err = errors.NewInternal(err.Error())
				log.Errorf("Proposal EnterProposal %s failed: %v", p, err)
				command.Output(nil, err)
			} else {
				log.Errorf("Proposal EnterProposal %s complete: %+v", p, response)
				command.Output(output, nil)
			}
		}
	case 2:
		p, err := newWithdrawProposal(command)
		if err != nil {
			err = errors.NewInternal(err.Error())
			log.Error(err)
			command.Output(nil, err)
			return
		}

		log.Debugf("Proposal WithdrawProposal %s", p)
		response, err := s.rsm.Withdraw(p)
		if err != nil {
			log.Warnf("Proposal WithdrawProposal %s failed: %v", p, err)
			command.Output(nil, err)
		} else {
			output, err := proto.Marshal(response)
			if err != nil {
				err = errors.NewInternal(err.Error())
				log.Errorf("Proposal WithdrawProposal %s failed: %v", p, err)
				command.Output(nil, err)
			} else {
				log.Errorf("Proposal WithdrawProposal %s complete: %+v", p, response)
				command.Output(output, nil)
			}
		}
	case 3:
		p, err := newAnointProposal(command)
		if err != nil {
			err = errors.NewInternal(err.Error())
			log.Error(err)
			command.Output(nil, err)
			return
		}

		log.Debugf("Proposal AnointProposal %s", p)
		response, err := s.rsm.Anoint(p)
		if err != nil {
			log.Warnf("Proposal AnointProposal %s failed: %v", p, err)
			command.Output(nil, err)
		} else {
			output, err := proto.Marshal(response)
			if err != nil {
				err = errors.NewInternal(err.Error())
				log.Errorf("Proposal AnointProposal %s failed: %v", p, err)
				command.Output(nil, err)
			} else {
				log.Errorf("Proposal AnointProposal %s complete: %+v", p, response)
				command.Output(output, nil)
			}
		}
	case 4:
		p, err := newPromoteProposal(command)
		if err != nil {
			err = errors.NewInternal(err.Error())
			log.Error(err)
			command.Output(nil, err)
			return
		}

		log.Debugf("Proposal PromoteProposal %s", p)
		response, err := s.rsm.Promote(p)
		if err != nil {
			log.Warnf("Proposal PromoteProposal %s failed: %v", p, err)
			command.Output(nil, err)
		} else {
			output, err := proto.Marshal(response)
			if err != nil {
				err = errors.NewInternal(err.Error())
				log.Errorf("Proposal PromoteProposal %s failed: %v", p, err)
				command.Output(nil, err)
			} else {
				log.Errorf("Proposal PromoteProposal %s complete: %+v", p, response)
				command.Output(output, nil)
			}
		}
	case 5:
		p, err := newEvictProposal(command)
		if err != nil {
			err = errors.NewInternal(err.Error())
			log.Error(err)
			command.Output(nil, err)
			return
		}

		log.Debugf("Proposal EvictProposal %s", p)
		response, err := s.rsm.Evict(p)
		if err != nil {
			log.Warnf("Proposal EvictProposal %s failed: %v", p, err)
			command.Output(nil, err)
		} else {
			output, err := proto.Marshal(response)
			if err != nil {
				err = errors.NewInternal(err.Error())
				log.Errorf("Proposal EvictProposal %s failed: %v", p, err)
				command.Output(nil, err)
			} else {
				log.Errorf("Proposal EvictProposal %s complete: %+v", p, response)
				command.Output(output, nil)
			}
		}
	case 7:
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
	case 6:
		q, err := newGetTermQuery(query)
		if err != nil {
			err = errors.NewInternal(err.Error())
			log.Error(err)
			query.Output(nil, err)
			return
		}

		log.Debugf("Querying GetTermQuery %s", q)
		response, err := s.rsm.GetTerm(q)
		if err != nil {
			log.Warnf("Querying GetTermQuery %s failed: %v", q, err)
			query.Output(nil, err)
		} else {
			output, err := proto.Marshal(response)
			if err != nil {
				err = errors.NewInternal(err.Error())
				log.Errorf("Querying GetTermQuery %s failed: %v", q, err)
				query.Output(nil, err)
			} else {
				log.Errorf("Querying GetTermQuery %s complete: %+v", q, response)
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
