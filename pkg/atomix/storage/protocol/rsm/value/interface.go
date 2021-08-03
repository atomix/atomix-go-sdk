// Code generated by atomix-go-framework. DO NOT EDIT.
package value

import (
	"fmt"
	value "github.com/atomix/atomix-api/go/atomix/primitive/value"
	errors "github.com/atomix/atomix-go-framework/pkg/atomix/errors"
	rsm "github.com/atomix/atomix-go-framework/pkg/atomix/storage/protocol/rsm"
	util "github.com/atomix/atomix-go-framework/pkg/atomix/util"
	proto "github.com/golang/protobuf/proto"
	"io"
)

type Service interface {
	ServiceContext
	Backup(SnapshotWriter) error
	Restore(SnapshotReader) error
	// Set sets the value
	Set(SetProposal) (*value.SetResponse, error)
	// Get gets the value
	Get(GetQuery) (*value.GetResponse, error)
	// Events listens for value change events
	Events(EventsProposal)
}

type ServiceContext interface {
	Scheduler() rsm.Scheduler
	Sessions() Sessions
	Proposals() Proposals
}

func newServiceContext(service rsm.ServiceContext) ServiceContext {
	return &serviceContext{
		scheduler: service.Scheduler(),
		sessions:  newSessions(service.Sessions()),
		proposals: newProposals(service.Commands()),
	}
}

type serviceContext struct {
	scheduler rsm.Scheduler
	sessions  Sessions
	proposals Proposals
}

func (s *serviceContext) Scheduler() rsm.Scheduler {
	return s.scheduler
}

func (s *serviceContext) Sessions() Sessions {
	return s.sessions
}

func (s *serviceContext) Proposals() Proposals {
	return s.proposals
}

var _ ServiceContext = &serviceContext{}

type SnapshotWriter interface {
	WriteState(*ValueState) error
}

func newSnapshotWriter(writer io.Writer) SnapshotWriter {
	return &serviceSnapshotWriter{
		writer: writer,
	}
}

type serviceSnapshotWriter struct {
	writer io.Writer
}

func (w *serviceSnapshotWriter) WriteState(state *ValueState) error {
	bytes, err := proto.Marshal(state)
	if err != nil {
		return err
	}
	err = util.WriteBytes(w.writer, bytes)
	if err != nil {
		return err
	}
	return err
}

var _ SnapshotWriter = &serviceSnapshotWriter{}

type SnapshotReader interface {
	ReadState() (*ValueState, error)
}

func newSnapshotReader(reader io.Reader) SnapshotReader {
	return &serviceSnapshotReader{
		reader: reader,
	}
}

type serviceSnapshotReader struct {
	reader io.Reader
}

func (r *serviceSnapshotReader) ReadState() (*ValueState, error) {
	bytes, err := util.ReadBytes(r.reader)
	if err != nil {
		return nil, err
	}
	state := &ValueState{}
	err = proto.Unmarshal(bytes, state)
	if err != nil {
		return nil, err
	}
	return state, nil
}

var _ SnapshotReader = &serviceSnapshotReader{}

type Sessions interface {
	Get(SessionID) (Session, bool)
	List() []Session
}

func newSessions(sessions rsm.Sessions) Sessions {
	return &serviceSessions{
		sessions: sessions,
	}
}

type serviceSessions struct {
	sessions rsm.Sessions
}

func (s *serviceSessions) Get(id SessionID) (Session, bool) {
	session, ok := s.sessions.Get(rsm.SessionID(id))
	if !ok {
		return nil, false
	}
	return newSession(session), true
}

func (s *serviceSessions) List() []Session {
	serviceSessions := s.sessions.List()
	sessions := make([]Session, len(serviceSessions))
	for i, serviceSession := range serviceSessions {
		sessions[i] = newSession(serviceSession)
	}
	return sessions
}

var _ Sessions = &serviceSessions{}

type SessionID uint64

type SessionState int

const (
	SessionClosed SessionState = iota
	SessionOpen
)

type Watcher interface {
	Cancel()
}

func newWatcher(watcher rsm.Watcher) Watcher {
	return &serviceWatcher{
		watcher: watcher,
	}
}

type serviceWatcher struct {
	watcher rsm.Watcher
}

func (s *serviceWatcher) Cancel() {
	s.watcher.Cancel()
}

var _ Watcher = &serviceWatcher{}

type Session interface {
	ID() SessionID
	State() SessionState
	Watch(func(SessionState)) Watcher
	Proposals() Proposals
}

func newSession(session rsm.Session) Session {
	return &serviceSession{
		session:   session,
		proposals: newProposals(session.Commands()),
	}
}

type serviceSession struct {
	session   rsm.Session
	proposals Proposals
}

func (s *serviceSession) ID() SessionID {
	return SessionID(s.session.ID())
}

func (s *serviceSession) Proposals() Proposals {
	return s.proposals
}

func (s *serviceSession) State() SessionState {
	return SessionState(s.session.State())
}

func (s *serviceSession) Watch(f func(SessionState)) Watcher {
	return newWatcher(s.session.Watch(func(state rsm.SessionState) {
		f(SessionState(state))
	}))
}

var _ Session = &serviceSession{}

type Proposals interface {
	Set() SetProposals
	Events() EventsProposals
}

func newProposals(commands rsm.Commands) Proposals {
	return &serviceProposals{
		setProposals:    newSetProposals(commands),
		eventsProposals: newEventsProposals(commands),
	}
}

type serviceProposals struct {
	setProposals    SetProposals
	eventsProposals EventsProposals
}

func (s *serviceProposals) Set() SetProposals {
	return s.setProposals
}
func (s *serviceProposals) Events() EventsProposals {
	return s.eventsProposals
}

var _ Proposals = &serviceProposals{}

type ProposalID uint64

type ProposalState int

const (
	ProposalComplete ProposalState = iota
	ProposalOpen
)

type Proposal interface {
	fmt.Stringer
	ID() ProposalID
	Session() Session
	State() ProposalState
	Watch(func(ProposalState)) Watcher
}

func newProposal(command rsm.Command) Proposal {
	return &serviceProposal{
		command: command,
	}
}

type serviceProposal struct {
	command rsm.Command
}

func (p *serviceProposal) ID() ProposalID {
	return ProposalID(p.command.ID())
}

func (p *serviceProposal) Session() Session {
	return newSession(p.command.Session())
}

func (p *serviceProposal) State() ProposalState {
	return ProposalState(p.command.State())
}

func (p *serviceProposal) Watch(f func(ProposalState)) Watcher {
	return newWatcher(p.command.Watch(func(state rsm.CommandState) {
		f(ProposalState(state))
	}))
}

func (p *serviceProposal) String() string {
	return fmt.Sprintf("ProposalID: %d, SessionID: %d", p.ID(), p.Session().ID())
}

var _ Proposal = &serviceProposal{}

type Query interface {
	fmt.Stringer
	Session() Session
}

func newQuery(query rsm.Query) Query {
	return &serviceQuery{
		query: query,
	}
}

type serviceQuery struct {
	query rsm.Query
}

func (p *serviceQuery) Session() Session {
	return newSession(p.query.Session())
}

func (p *serviceQuery) String() string {
	return fmt.Sprintf("SessionID: %d", p.Session().ID())
}

var _ Query = &serviceQuery{}

type SetProposals interface {
	Get(ProposalID) (SetProposal, bool)
	List() []SetProposal
}

func newSetProposals(commands rsm.Commands) SetProposals {
	return &setProposals{
		commands: commands,
	}
}

type setProposals struct {
	commands rsm.Commands
}

func (p *setProposals) Get(id ProposalID) (SetProposal, bool) {
	command, ok := p.commands.Get(rsm.CommandID(id))
	if !ok {
		return nil, false
	}
	proposal, err := newSetProposal(command)
	if err != nil {
		log.Error(err)
		return nil, false
	}
	return proposal, true
}

func (p *setProposals) List() []SetProposal {
	commands := p.commands.List(rsm.OperationID(1))
	proposals := make([]SetProposal, len(commands))
	for i, command := range commands {
		proposal, err := newSetProposal(command)
		if err != nil {
			log.Error(err)
		} else {
			proposals[i] = proposal
		}
	}
	return proposals
}

var _ SetProposals = &setProposals{}

type SetProposal interface {
	Proposal
	Request() *value.SetRequest
}

func newSetProposal(command rsm.Command) (SetProposal, error) {
	request := &value.SetRequest{}
	if err := proto.Unmarshal(command.Input(), request); err != nil {
		return nil, err
	}
	return &setProposal{
		Proposal: newProposal(command),
		command:  command,
		request:  request,
	}, nil
}

type setProposal struct {
	Proposal
	command rsm.Command
	request *value.SetRequest
}

func (p *setProposal) Request() *value.SetRequest {
	return p.request
}

func (p *setProposal) String() string {
	return fmt.Sprintf("ProposalID=%d, SessionID=%d", p.ID(), p.Session().ID())
}

var _ SetProposal = &setProposal{}

type GetQuery interface {
	Query
	Request() *value.GetRequest
}

func newGetQuery(query rsm.Query) (GetQuery, error) {
	request := &value.GetRequest{}
	if err := proto.Unmarshal(query.Input(), request); err != nil {
		return nil, err
	}
	return &getQuery{
		Query:   newQuery(query),
		query:   query,
		request: request,
	}, nil
}

type getQuery struct {
	Query
	query   rsm.Query
	request *value.GetRequest
}

func (p *getQuery) Request() *value.GetRequest {
	return p.request
}

func (p *getQuery) String() string {
	return fmt.Sprintf("SessionID=%d", p.Session().ID())
}

var _ GetQuery = &getQuery{}

type EventsProposals interface {
	Get(ProposalID) (EventsProposal, bool)
	List() []EventsProposal
}

func newEventsProposals(commands rsm.Commands) EventsProposals {
	return &eventsProposals{
		commands: commands,
	}
}

type eventsProposals struct {
	commands rsm.Commands
}

func (p *eventsProposals) Get(id ProposalID) (EventsProposal, bool) {
	command, ok := p.commands.Get(rsm.CommandID(id))
	if !ok {
		return nil, false
	}
	proposal, err := newEventsProposal(command)
	if err != nil {
		log.Error(err)
		return nil, false
	}
	return proposal, true
}

func (p *eventsProposals) List() []EventsProposal {
	commands := p.commands.List(rsm.OperationID(3))
	proposals := make([]EventsProposal, len(commands))
	for i, command := range commands {
		proposal, err := newEventsProposal(command)
		if err != nil {
			log.Error(err)
		} else {
			proposals[i] = proposal
		}
	}
	return proposals
}

var _ EventsProposals = &eventsProposals{}

type EventsProposal interface {
	Proposal
	Request() *value.EventsRequest
	Notify(*value.EventsResponse)
	Close()
}

func newEventsProposal(command rsm.Command) (EventsProposal, error) {
	request := &value.EventsRequest{}
	if err := proto.Unmarshal(command.Input(), request); err != nil {
		return nil, err
	}
	return &eventsProposal{
		Proposal: newProposal(command),
		command:  command,
		request:  request,
	}, nil
}

type eventsProposal struct {
	Proposal
	command rsm.Command
	request *value.EventsRequest
	closed  bool
}

func (p *eventsProposal) Request() *value.EventsRequest {
	return p.request
}

func (p *eventsProposal) Notify(response *value.EventsResponse) {
	if p.closed {
		return
	}
	log.Debugf("Notifying EventsProposal %s: %s", p, response)
	output, err := proto.Marshal(response)
	if err != nil {
		err = errors.NewInternal(err.Error())
		log.Errorf("Notifying EventsProposal %s failed: %v", p, err)
		p.command.Output(nil, err)
		p.command.Close()
		p.closed = true
	} else {
		log.Debugf("Notifying EventsProposal %s: %s", p, response)
		p.command.Output(output, nil)
	}
}

func (p *eventsProposal) Close() {
	p.command.Close()
	p.closed = true
}

func (p *eventsProposal) String() string {
	return fmt.Sprintf("ProposalID=%d, SessionID=%d", p.ID(), p.Session().ID())
}

var _ EventsProposal = &eventsProposal{}
