// Code generated by atomix-go-framework. DO NOT EDIT.
package election

import (
	election "github.com/atomix/atomix-api/go/atomix/primitive/election"
	"github.com/atomix/atomix-go-framework/pkg/atomix/errors"
	"github.com/atomix/atomix-go-framework/pkg/atomix/logging"
	"github.com/atomix/atomix-go-framework/pkg/atomix/storage/protocol/rsm"
	"github.com/golang/protobuf/proto"
	"io"
)

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
	newServiceFunc = func(scheduler rsm.Scheduler, context rsm.ServiceContext) rsm.Service {
		service := &ServiceAdaptor{
			Service: rsm.NewService(scheduler, context),
			rsm:     rsmf(newServiceContext(scheduler)),
			log:     logging.GetLogger("atomix", "election", "service"),
		}
		service.init()
		return service
	}
}

type NewServiceFunc func(ServiceContext) Service

// RegisterService registers the election primitive service on the given node
func RegisterService(node *rsm.Node) {
	node.RegisterService(Type, newServiceFunc)
}

type ServiceAdaptor struct {
	rsm.Service
	rsm Service
	log logging.Logger
}

func (s *ServiceAdaptor) init() {
	s.RegisterUnaryOperation(enterOp, s.enter)
	s.RegisterUnaryOperation(withdrawOp, s.withdraw)
	s.RegisterUnaryOperation(anointOp, s.anoint)
	s.RegisterUnaryOperation(promoteOp, s.promote)
	s.RegisterUnaryOperation(evictOp, s.evict)
	s.RegisterUnaryOperation(getTermOp, s.getTerm)
	s.RegisterStreamOperation(eventsOp, s.events)
}
func (s *ServiceAdaptor) SessionOpen(rsmSession rsm.Session) {
	s.rsm.Sessions().open(newSession(rsmSession))
}

func (s *ServiceAdaptor) SessionExpired(session rsm.Session) {
	s.rsm.Sessions().expire(SessionID(session.ID()))
}

func (s *ServiceAdaptor) SessionClosed(session rsm.Session) {
	s.rsm.Sessions().close(SessionID(session.ID()))
}
func (s *ServiceAdaptor) Backup(writer io.Writer) error {
	err := s.rsm.Backup(newSnapshotWriter(writer))
	if err != nil {
		s.log.Error(err)
		return err
	}
	return nil
}

func (s *ServiceAdaptor) Restore(reader io.Reader) error {
	err := s.rsm.Restore(newSnapshotReader(reader))
	if err != nil {
		s.log.Error(err)
		return err
	}
	return nil
}
func (s *ServiceAdaptor) enter(input []byte, rsmSession rsm.Session) ([]byte, error) {
	request := &election.EnterRequest{}
	err := proto.Unmarshal(input, request)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	session, ok := s.rsm.Sessions().Get(SessionID(rsmSession.ID()))
	if !ok {
		err := errors.NewConflict("session %d not found", rsmSession.ID())
		s.log.Error(err)
		return nil, err
	}

	var response *election.EnterResponse
	proposal := newEnterProposal(ProposalID(s.Index()), session, request, response)

	s.rsm.Proposals().Enter().register(proposal)
	session.Proposals().Enter().register(proposal)

	defer func() {
		session.Proposals().Enter().unregister(proposal.ID())
		s.rsm.Proposals().Enter().unregister(proposal.ID())
	}()

	err = s.rsm.Enter(proposal)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	output, err := proto.Marshal(response)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}
	return output, nil
}
func (s *ServiceAdaptor) withdraw(input []byte, rsmSession rsm.Session) ([]byte, error) {
	request := &election.WithdrawRequest{}
	err := proto.Unmarshal(input, request)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	session, ok := s.rsm.Sessions().Get(SessionID(rsmSession.ID()))
	if !ok {
		err := errors.NewConflict("session %d not found", rsmSession.ID())
		s.log.Error(err)
		return nil, err
	}

	var response *election.WithdrawResponse
	proposal := newWithdrawProposal(ProposalID(s.Index()), session, request, response)

	s.rsm.Proposals().Withdraw().register(proposal)
	session.Proposals().Withdraw().register(proposal)

	defer func() {
		session.Proposals().Withdraw().unregister(proposal.ID())
		s.rsm.Proposals().Withdraw().unregister(proposal.ID())
	}()

	err = s.rsm.Withdraw(proposal)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	output, err := proto.Marshal(response)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}
	return output, nil
}
func (s *ServiceAdaptor) anoint(input []byte, rsmSession rsm.Session) ([]byte, error) {
	request := &election.AnointRequest{}
	err := proto.Unmarshal(input, request)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	session, ok := s.rsm.Sessions().Get(SessionID(rsmSession.ID()))
	if !ok {
		err := errors.NewConflict("session %d not found", rsmSession.ID())
		s.log.Error(err)
		return nil, err
	}

	var response *election.AnointResponse
	proposal := newAnointProposal(ProposalID(s.Index()), session, request, response)

	s.rsm.Proposals().Anoint().register(proposal)
	session.Proposals().Anoint().register(proposal)

	defer func() {
		session.Proposals().Anoint().unregister(proposal.ID())
		s.rsm.Proposals().Anoint().unregister(proposal.ID())
	}()

	err = s.rsm.Anoint(proposal)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	output, err := proto.Marshal(response)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}
	return output, nil
}
func (s *ServiceAdaptor) promote(input []byte, rsmSession rsm.Session) ([]byte, error) {
	request := &election.PromoteRequest{}
	err := proto.Unmarshal(input, request)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	session, ok := s.rsm.Sessions().Get(SessionID(rsmSession.ID()))
	if !ok {
		err := errors.NewConflict("session %d not found", rsmSession.ID())
		s.log.Error(err)
		return nil, err
	}

	var response *election.PromoteResponse
	proposal := newPromoteProposal(ProposalID(s.Index()), session, request, response)

	s.rsm.Proposals().Promote().register(proposal)
	session.Proposals().Promote().register(proposal)

	defer func() {
		session.Proposals().Promote().unregister(proposal.ID())
		s.rsm.Proposals().Promote().unregister(proposal.ID())
	}()

	err = s.rsm.Promote(proposal)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	output, err := proto.Marshal(response)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}
	return output, nil
}
func (s *ServiceAdaptor) evict(input []byte, rsmSession rsm.Session) ([]byte, error) {
	request := &election.EvictRequest{}
	err := proto.Unmarshal(input, request)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	session, ok := s.rsm.Sessions().Get(SessionID(rsmSession.ID()))
	if !ok {
		err := errors.NewConflict("session %d not found", rsmSession.ID())
		s.log.Error(err)
		return nil, err
	}

	var response *election.EvictResponse
	proposal := newEvictProposal(ProposalID(s.Index()), session, request, response)

	s.rsm.Proposals().Evict().register(proposal)
	session.Proposals().Evict().register(proposal)

	defer func() {
		session.Proposals().Evict().unregister(proposal.ID())
		s.rsm.Proposals().Evict().unregister(proposal.ID())
	}()

	err = s.rsm.Evict(proposal)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	output, err := proto.Marshal(response)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}
	return output, nil
}
func (s *ServiceAdaptor) getTerm(input []byte, rsmSession rsm.Session) ([]byte, error) {
	request := &election.GetTermRequest{}
	err := proto.Unmarshal(input, request)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	session, ok := s.rsm.Sessions().Get(SessionID(rsmSession.ID()))
	if !ok {
		err := errors.NewConflict("session %d not found", rsmSession.ID())
		s.log.Error(err)
		return nil, err
	}

	var response *election.GetTermResponse
	proposal := newGetTermProposal(ProposalID(s.Index()), session, request, response)

	s.rsm.Proposals().GetTerm().register(proposal)
	session.Proposals().GetTerm().register(proposal)

	defer func() {
		session.Proposals().GetTerm().unregister(proposal.ID())
		s.rsm.Proposals().GetTerm().unregister(proposal.ID())
	}()

	err = s.rsm.GetTerm(proposal)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	output, err := proto.Marshal(response)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}
	return output, nil
}
func (s *ServiceAdaptor) events(input []byte, rsmSession rsm.Session, stream rsm.Stream) (rsm.StreamCloser, error) {
	request := &election.EventsRequest{}
	err := proto.Unmarshal(input, request)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	session, ok := s.rsm.Sessions().Get(SessionID(rsmSession.ID()))
	if !ok {
		err := errors.NewConflict("session %d not found", rsmSession.ID())
		s.log.Error(err)
		return nil, err
	}

	proposal := newEventsProposal(ProposalID(stream.ID()), session, request, stream)

	s.rsm.Proposals().Events().register(proposal)
	session.Proposals().Events().register(proposal)

	err = s.rsm.Events(proposal)
	if err != nil {
		s.log.Error(err)
		return nil, err
	}
	return func() {
		session.Proposals().Events().unregister(proposal.ID())
		s.rsm.Proposals().Events().unregister(proposal.ID())
	}, nil
}

var _ rsm.Service = &ServiceAdaptor{}
