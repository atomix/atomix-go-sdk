// Code generated by atomix-go-framework. DO NOT EDIT.
package indexedmap

import (
	indexedmap "github.com/atomix/atomix-api/go/atomix/primitive/indexedmap"
	errors "github.com/atomix/atomix-go-framework/pkg/atomix/errors"
	rsm "github.com/atomix/atomix-go-framework/pkg/atomix/storage/protocol/rsm"
	util "github.com/atomix/atomix-go-framework/pkg/atomix/util"
	proto "github.com/golang/protobuf/proto"
	uuid "github.com/google/uuid"
	"io"
)

type Service interface {
	ServiceContext
	Backup(SnapshotWriter) error
	Restore(SnapshotReader) error
	// Size returns the size of the map
	Size(SizeProposal) error
	// Put puts an entry into the map
	Put(PutProposal) error
	// Get gets the entry for a key
	Get(GetProposal) error
	// FirstEntry gets the first entry in the map
	FirstEntry(FirstEntryProposal) error
	// LastEntry gets the last entry in the map
	LastEntry(LastEntryProposal) error
	// PrevEntry gets the previous entry in the map
	PrevEntry(PrevEntryProposal) error
	// NextEntry gets the next entry in the map
	NextEntry(NextEntryProposal) error
	// Remove removes an entry from the map
	Remove(RemoveProposal) error
	// Clear removes all entries from the map
	Clear(ClearProposal) error
	// Events listens for change events
	Events(EventsProposal) error
	// Entries lists all entries in the map
	Entries(EntriesProposal) error
}

type ServiceContext interface {
	Scheduler() rsm.Scheduler
	Sessions() Sessions
	Proposals() Proposals
}

func newServiceContext(scheduler rsm.Scheduler) ServiceContext {
	return &serviceContext{
		scheduler: scheduler,
		sessions:  newSessions(),
		proposals: newProposals(),
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
	WriteState(*IndexedMapState) error
}

func newSnapshotWriter(writer io.Writer) SnapshotWriter {
	return &serviceSnapshotWriter{
		writer: writer,
	}
}

type serviceSnapshotWriter struct {
	writer io.Writer
}

func (w *serviceSnapshotWriter) WriteState(state *IndexedMapState) error {
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
	ReadState() (*IndexedMapState, error)
}

func newSnapshotReader(reader io.Reader) SnapshotReader {
	return &serviceSnapshotReader{
		reader: reader,
	}
}

type serviceSnapshotReader struct {
	reader io.Reader
}

func (r *serviceSnapshotReader) ReadState() (*IndexedMapState, error) {
	bytes, err := util.ReadBytes(r.reader)
	if err != nil {
		return nil, err
	}
	state := &IndexedMapState{}
	err = proto.Unmarshal(bytes, state)
	if err != nil {
		return nil, err
	}
	return state, nil
}

var _ SnapshotReader = &serviceSnapshotReader{}

type Sessions interface {
	open(Session)
	expire(SessionID)
	close(SessionID)
	Get(SessionID) (Session, bool)
	List() []Session
}

func newSessions() Sessions {
	return &serviceSessions{
		sessions: make(map[SessionID]Session),
	}
}

type serviceSessions struct {
	sessions map[SessionID]Session
}

func (s *serviceSessions) open(session Session) {
	s.sessions[session.ID()] = session
	session.setState(SessionOpen)
}

func (s *serviceSessions) expire(sessionID SessionID) {
	session, ok := s.sessions[sessionID]
	if ok {
		session.setState(SessionClosed)
		delete(s.sessions, sessionID)
	}
}

func (s *serviceSessions) close(sessionID SessionID) {
	session, ok := s.sessions[sessionID]
	if ok {
		session.setState(SessionClosed)
		delete(s.sessions, sessionID)
	}
}

func (s *serviceSessions) Get(id SessionID) (Session, bool) {
	session, ok := s.sessions[id]
	return session, ok
}

func (s *serviceSessions) List() []Session {
	sessions := make([]Session, 0, len(s.sessions))
	for _, session := range s.sessions {
		sessions = append(sessions, session)
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

func newWatcher(f func()) Watcher {
	return &serviceWatcher{
		f: f,
	}
}

type serviceWatcher struct {
	f func()
}

func (s *serviceWatcher) Cancel() {
	s.f()
}

var _ Watcher = &serviceWatcher{}

type Session interface {
	ID() SessionID
	State() SessionState
	setState(SessionState)
	Watch(func(SessionState)) Watcher
	Proposals() Proposals
}

func newSession(session rsm.Session) Session {
	return &serviceSession{
		session:   session,
		proposals: newProposals(),
		watchers:  make(map[string]func(SessionState)),
	}
}

type serviceSession struct {
	session   rsm.Session
	proposals Proposals
	state     SessionState
	watchers  map[string]func(SessionState)
}

func (s *serviceSession) ID() SessionID {
	return SessionID(s.session.ID())
}

func (s *serviceSession) Proposals() Proposals {
	return s.proposals
}

func (s *serviceSession) State() SessionState {
	return s.state
}

func (s *serviceSession) setState(state SessionState) {
	if state != s.state {
		s.state = state
		for _, watcher := range s.watchers {
			watcher(state)
		}
	}
}

func (s *serviceSession) Watch(f func(SessionState)) Watcher {
	id := uuid.New().String()
	s.watchers[id] = f
	return newWatcher(func() {
		delete(s.watchers, id)
	})
}

var _ Session = &serviceSession{}

type Proposals interface {
	Size() SizeProposals
	Put() PutProposals
	Get() GetProposals
	FirstEntry() FirstEntryProposals
	LastEntry() LastEntryProposals
	PrevEntry() PrevEntryProposals
	NextEntry() NextEntryProposals
	Remove() RemoveProposals
	Clear() ClearProposals
	Events() EventsProposals
	Entries() EntriesProposals
}

func newProposals() Proposals {
	return &serviceProposals{
		sizeProposals:       newSizeProposals(),
		putProposals:        newPutProposals(),
		getProposals:        newGetProposals(),
		firstEntryProposals: newFirstEntryProposals(),
		lastEntryProposals:  newLastEntryProposals(),
		prevEntryProposals:  newPrevEntryProposals(),
		nextEntryProposals:  newNextEntryProposals(),
		removeProposals:     newRemoveProposals(),
		clearProposals:      newClearProposals(),
		eventsProposals:     newEventsProposals(),
		entriesProposals:    newEntriesProposals(),
	}
}

type serviceProposals struct {
	sizeProposals       SizeProposals
	putProposals        PutProposals
	getProposals        GetProposals
	firstEntryProposals FirstEntryProposals
	lastEntryProposals  LastEntryProposals
	prevEntryProposals  PrevEntryProposals
	nextEntryProposals  NextEntryProposals
	removeProposals     RemoveProposals
	clearProposals      ClearProposals
	eventsProposals     EventsProposals
	entriesProposals    EntriesProposals
}

func (s *serviceProposals) Size() SizeProposals {
	return s.sizeProposals
}
func (s *serviceProposals) Put() PutProposals {
	return s.putProposals
}
func (s *serviceProposals) Get() GetProposals {
	return s.getProposals
}
func (s *serviceProposals) FirstEntry() FirstEntryProposals {
	return s.firstEntryProposals
}
func (s *serviceProposals) LastEntry() LastEntryProposals {
	return s.lastEntryProposals
}
func (s *serviceProposals) PrevEntry() PrevEntryProposals {
	return s.prevEntryProposals
}
func (s *serviceProposals) NextEntry() NextEntryProposals {
	return s.nextEntryProposals
}
func (s *serviceProposals) Remove() RemoveProposals {
	return s.removeProposals
}
func (s *serviceProposals) Clear() ClearProposals {
	return s.clearProposals
}
func (s *serviceProposals) Events() EventsProposals {
	return s.eventsProposals
}
func (s *serviceProposals) Entries() EntriesProposals {
	return s.entriesProposals
}

var _ Proposals = &serviceProposals{}

type ProposalID uint64

type Proposal interface {
	ID() ProposalID
	Session() Session
}

func newProposal(id ProposalID, session Session) Proposal {
	return &serviceProposal{
		id:      id,
		session: session,
	}
}

type serviceProposal struct {
	id      ProposalID
	session Session
}

func (p *serviceProposal) ID() ProposalID {
	return p.id
}

func (p *serviceProposal) Session() Session {
	return p.session
}

var _ Proposal = &serviceProposal{}

type SizeProposals interface {
	register(SizeProposal)
	unregister(ProposalID)
	Get(ProposalID) (SizeProposal, bool)
	List() []SizeProposal
}

func newSizeProposals() SizeProposals {
	return &sizeProposals{
		proposals: make(map[ProposalID]SizeProposal),
	}
}

type sizeProposals struct {
	proposals map[ProposalID]SizeProposal
}

func (p *sizeProposals) register(proposal SizeProposal) {
	p.proposals[proposal.ID()] = proposal
}

func (p *sizeProposals) unregister(id ProposalID) {
	delete(p.proposals, id)
}

func (p *sizeProposals) Get(id ProposalID) (SizeProposal, bool) {
	proposal, ok := p.proposals[id]
	return proposal, ok
}

func (p *sizeProposals) List() []SizeProposal {
	proposals := make([]SizeProposal, 0, len(p.proposals))
	for _, proposal := range p.proposals {
		proposals = append(proposals, proposal)
	}
	return proposals
}

var _ SizeProposals = &sizeProposals{}

type SizeProposal interface {
	Proposal
	Request() *indexedmap.SizeRequest
	Reply(*indexedmap.SizeResponse) error
}

func newSizeProposal(id ProposalID, session Session, request *indexedmap.SizeRequest, response *indexedmap.SizeResponse) SizeProposal {
	return &sizeProposal{
		Proposal: newProposal(id, session),
		request:  request,
		response: response,
	}
}

type sizeProposal struct {
	Proposal
	request  *indexedmap.SizeRequest
	response *indexedmap.SizeResponse
}

func (p *sizeProposal) Request() *indexedmap.SizeRequest {
	return p.request
}

func (p *sizeProposal) Reply(reply *indexedmap.SizeResponse) error {
	if p.response != nil {
		return errors.NewConflict("reply already sent")
	}
	p.response = reply
	return nil
}

var _ SizeProposal = &sizeProposal{}

type PutProposals interface {
	register(PutProposal)
	unregister(ProposalID)
	Get(ProposalID) (PutProposal, bool)
	List() []PutProposal
}

func newPutProposals() PutProposals {
	return &putProposals{
		proposals: make(map[ProposalID]PutProposal),
	}
}

type putProposals struct {
	proposals map[ProposalID]PutProposal
}

func (p *putProposals) register(proposal PutProposal) {
	p.proposals[proposal.ID()] = proposal
}

func (p *putProposals) unregister(id ProposalID) {
	delete(p.proposals, id)
}

func (p *putProposals) Get(id ProposalID) (PutProposal, bool) {
	proposal, ok := p.proposals[id]
	return proposal, ok
}

func (p *putProposals) List() []PutProposal {
	proposals := make([]PutProposal, 0, len(p.proposals))
	for _, proposal := range p.proposals {
		proposals = append(proposals, proposal)
	}
	return proposals
}

var _ PutProposals = &putProposals{}

type PutProposal interface {
	Proposal
	Request() *indexedmap.PutRequest
	Reply(*indexedmap.PutResponse) error
}

func newPutProposal(id ProposalID, session Session, request *indexedmap.PutRequest, response *indexedmap.PutResponse) PutProposal {
	return &putProposal{
		Proposal: newProposal(id, session),
		request:  request,
		response: response,
	}
}

type putProposal struct {
	Proposal
	request  *indexedmap.PutRequest
	response *indexedmap.PutResponse
}

func (p *putProposal) Request() *indexedmap.PutRequest {
	return p.request
}

func (p *putProposal) Reply(reply *indexedmap.PutResponse) error {
	if p.response != nil {
		return errors.NewConflict("reply already sent")
	}
	p.response = reply
	return nil
}

var _ PutProposal = &putProposal{}

type GetProposals interface {
	register(GetProposal)
	unregister(ProposalID)
	Get(ProposalID) (GetProposal, bool)
	List() []GetProposal
}

func newGetProposals() GetProposals {
	return &getProposals{
		proposals: make(map[ProposalID]GetProposal),
	}
}

type getProposals struct {
	proposals map[ProposalID]GetProposal
}

func (p *getProposals) register(proposal GetProposal) {
	p.proposals[proposal.ID()] = proposal
}

func (p *getProposals) unregister(id ProposalID) {
	delete(p.proposals, id)
}

func (p *getProposals) Get(id ProposalID) (GetProposal, bool) {
	proposal, ok := p.proposals[id]
	return proposal, ok
}

func (p *getProposals) List() []GetProposal {
	proposals := make([]GetProposal, 0, len(p.proposals))
	for _, proposal := range p.proposals {
		proposals = append(proposals, proposal)
	}
	return proposals
}

var _ GetProposals = &getProposals{}

type GetProposal interface {
	Proposal
	Request() *indexedmap.GetRequest
	Reply(*indexedmap.GetResponse) error
}

func newGetProposal(id ProposalID, session Session, request *indexedmap.GetRequest, response *indexedmap.GetResponse) GetProposal {
	return &getProposal{
		Proposal: newProposal(id, session),
		request:  request,
		response: response,
	}
}

type getProposal struct {
	Proposal
	request  *indexedmap.GetRequest
	response *indexedmap.GetResponse
}

func (p *getProposal) Request() *indexedmap.GetRequest {
	return p.request
}

func (p *getProposal) Reply(reply *indexedmap.GetResponse) error {
	if p.response != nil {
		return errors.NewConflict("reply already sent")
	}
	p.response = reply
	return nil
}

var _ GetProposal = &getProposal{}

type FirstEntryProposals interface {
	register(FirstEntryProposal)
	unregister(ProposalID)
	Get(ProposalID) (FirstEntryProposal, bool)
	List() []FirstEntryProposal
}

func newFirstEntryProposals() FirstEntryProposals {
	return &firstEntryProposals{
		proposals: make(map[ProposalID]FirstEntryProposal),
	}
}

type firstEntryProposals struct {
	proposals map[ProposalID]FirstEntryProposal
}

func (p *firstEntryProposals) register(proposal FirstEntryProposal) {
	p.proposals[proposal.ID()] = proposal
}

func (p *firstEntryProposals) unregister(id ProposalID) {
	delete(p.proposals, id)
}

func (p *firstEntryProposals) Get(id ProposalID) (FirstEntryProposal, bool) {
	proposal, ok := p.proposals[id]
	return proposal, ok
}

func (p *firstEntryProposals) List() []FirstEntryProposal {
	proposals := make([]FirstEntryProposal, 0, len(p.proposals))
	for _, proposal := range p.proposals {
		proposals = append(proposals, proposal)
	}
	return proposals
}

var _ FirstEntryProposals = &firstEntryProposals{}

type FirstEntryProposal interface {
	Proposal
	Request() *indexedmap.FirstEntryRequest
	Reply(*indexedmap.FirstEntryResponse) error
}

func newFirstEntryProposal(id ProposalID, session Session, request *indexedmap.FirstEntryRequest, response *indexedmap.FirstEntryResponse) FirstEntryProposal {
	return &firstEntryProposal{
		Proposal: newProposal(id, session),
		request:  request,
		response: response,
	}
}

type firstEntryProposal struct {
	Proposal
	request  *indexedmap.FirstEntryRequest
	response *indexedmap.FirstEntryResponse
}

func (p *firstEntryProposal) Request() *indexedmap.FirstEntryRequest {
	return p.request
}

func (p *firstEntryProposal) Reply(reply *indexedmap.FirstEntryResponse) error {
	if p.response != nil {
		return errors.NewConflict("reply already sent")
	}
	p.response = reply
	return nil
}

var _ FirstEntryProposal = &firstEntryProposal{}

type LastEntryProposals interface {
	register(LastEntryProposal)
	unregister(ProposalID)
	Get(ProposalID) (LastEntryProposal, bool)
	List() []LastEntryProposal
}

func newLastEntryProposals() LastEntryProposals {
	return &lastEntryProposals{
		proposals: make(map[ProposalID]LastEntryProposal),
	}
}

type lastEntryProposals struct {
	proposals map[ProposalID]LastEntryProposal
}

func (p *lastEntryProposals) register(proposal LastEntryProposal) {
	p.proposals[proposal.ID()] = proposal
}

func (p *lastEntryProposals) unregister(id ProposalID) {
	delete(p.proposals, id)
}

func (p *lastEntryProposals) Get(id ProposalID) (LastEntryProposal, bool) {
	proposal, ok := p.proposals[id]
	return proposal, ok
}

func (p *lastEntryProposals) List() []LastEntryProposal {
	proposals := make([]LastEntryProposal, 0, len(p.proposals))
	for _, proposal := range p.proposals {
		proposals = append(proposals, proposal)
	}
	return proposals
}

var _ LastEntryProposals = &lastEntryProposals{}

type LastEntryProposal interface {
	Proposal
	Request() *indexedmap.LastEntryRequest
	Reply(*indexedmap.LastEntryResponse) error
}

func newLastEntryProposal(id ProposalID, session Session, request *indexedmap.LastEntryRequest, response *indexedmap.LastEntryResponse) LastEntryProposal {
	return &lastEntryProposal{
		Proposal: newProposal(id, session),
		request:  request,
		response: response,
	}
}

type lastEntryProposal struct {
	Proposal
	request  *indexedmap.LastEntryRequest
	response *indexedmap.LastEntryResponse
}

func (p *lastEntryProposal) Request() *indexedmap.LastEntryRequest {
	return p.request
}

func (p *lastEntryProposal) Reply(reply *indexedmap.LastEntryResponse) error {
	if p.response != nil {
		return errors.NewConflict("reply already sent")
	}
	p.response = reply
	return nil
}

var _ LastEntryProposal = &lastEntryProposal{}

type PrevEntryProposals interface {
	register(PrevEntryProposal)
	unregister(ProposalID)
	Get(ProposalID) (PrevEntryProposal, bool)
	List() []PrevEntryProposal
}

func newPrevEntryProposals() PrevEntryProposals {
	return &prevEntryProposals{
		proposals: make(map[ProposalID]PrevEntryProposal),
	}
}

type prevEntryProposals struct {
	proposals map[ProposalID]PrevEntryProposal
}

func (p *prevEntryProposals) register(proposal PrevEntryProposal) {
	p.proposals[proposal.ID()] = proposal
}

func (p *prevEntryProposals) unregister(id ProposalID) {
	delete(p.proposals, id)
}

func (p *prevEntryProposals) Get(id ProposalID) (PrevEntryProposal, bool) {
	proposal, ok := p.proposals[id]
	return proposal, ok
}

func (p *prevEntryProposals) List() []PrevEntryProposal {
	proposals := make([]PrevEntryProposal, 0, len(p.proposals))
	for _, proposal := range p.proposals {
		proposals = append(proposals, proposal)
	}
	return proposals
}

var _ PrevEntryProposals = &prevEntryProposals{}

type PrevEntryProposal interface {
	Proposal
	Request() *indexedmap.PrevEntryRequest
	Reply(*indexedmap.PrevEntryResponse) error
}

func newPrevEntryProposal(id ProposalID, session Session, request *indexedmap.PrevEntryRequest, response *indexedmap.PrevEntryResponse) PrevEntryProposal {
	return &prevEntryProposal{
		Proposal: newProposal(id, session),
		request:  request,
		response: response,
	}
}

type prevEntryProposal struct {
	Proposal
	request  *indexedmap.PrevEntryRequest
	response *indexedmap.PrevEntryResponse
}

func (p *prevEntryProposal) Request() *indexedmap.PrevEntryRequest {
	return p.request
}

func (p *prevEntryProposal) Reply(reply *indexedmap.PrevEntryResponse) error {
	if p.response != nil {
		return errors.NewConflict("reply already sent")
	}
	p.response = reply
	return nil
}

var _ PrevEntryProposal = &prevEntryProposal{}

type NextEntryProposals interface {
	register(NextEntryProposal)
	unregister(ProposalID)
	Get(ProposalID) (NextEntryProposal, bool)
	List() []NextEntryProposal
}

func newNextEntryProposals() NextEntryProposals {
	return &nextEntryProposals{
		proposals: make(map[ProposalID]NextEntryProposal),
	}
}

type nextEntryProposals struct {
	proposals map[ProposalID]NextEntryProposal
}

func (p *nextEntryProposals) register(proposal NextEntryProposal) {
	p.proposals[proposal.ID()] = proposal
}

func (p *nextEntryProposals) unregister(id ProposalID) {
	delete(p.proposals, id)
}

func (p *nextEntryProposals) Get(id ProposalID) (NextEntryProposal, bool) {
	proposal, ok := p.proposals[id]
	return proposal, ok
}

func (p *nextEntryProposals) List() []NextEntryProposal {
	proposals := make([]NextEntryProposal, 0, len(p.proposals))
	for _, proposal := range p.proposals {
		proposals = append(proposals, proposal)
	}
	return proposals
}

var _ NextEntryProposals = &nextEntryProposals{}

type NextEntryProposal interface {
	Proposal
	Request() *indexedmap.NextEntryRequest
	Reply(*indexedmap.NextEntryResponse) error
}

func newNextEntryProposal(id ProposalID, session Session, request *indexedmap.NextEntryRequest, response *indexedmap.NextEntryResponse) NextEntryProposal {
	return &nextEntryProposal{
		Proposal: newProposal(id, session),
		request:  request,
		response: response,
	}
}

type nextEntryProposal struct {
	Proposal
	request  *indexedmap.NextEntryRequest
	response *indexedmap.NextEntryResponse
}

func (p *nextEntryProposal) Request() *indexedmap.NextEntryRequest {
	return p.request
}

func (p *nextEntryProposal) Reply(reply *indexedmap.NextEntryResponse) error {
	if p.response != nil {
		return errors.NewConflict("reply already sent")
	}
	p.response = reply
	return nil
}

var _ NextEntryProposal = &nextEntryProposal{}

type RemoveProposals interface {
	register(RemoveProposal)
	unregister(ProposalID)
	Get(ProposalID) (RemoveProposal, bool)
	List() []RemoveProposal
}

func newRemoveProposals() RemoveProposals {
	return &removeProposals{
		proposals: make(map[ProposalID]RemoveProposal),
	}
}

type removeProposals struct {
	proposals map[ProposalID]RemoveProposal
}

func (p *removeProposals) register(proposal RemoveProposal) {
	p.proposals[proposal.ID()] = proposal
}

func (p *removeProposals) unregister(id ProposalID) {
	delete(p.proposals, id)
}

func (p *removeProposals) Get(id ProposalID) (RemoveProposal, bool) {
	proposal, ok := p.proposals[id]
	return proposal, ok
}

func (p *removeProposals) List() []RemoveProposal {
	proposals := make([]RemoveProposal, 0, len(p.proposals))
	for _, proposal := range p.proposals {
		proposals = append(proposals, proposal)
	}
	return proposals
}

var _ RemoveProposals = &removeProposals{}

type RemoveProposal interface {
	Proposal
	Request() *indexedmap.RemoveRequest
	Reply(*indexedmap.RemoveResponse) error
}

func newRemoveProposal(id ProposalID, session Session, request *indexedmap.RemoveRequest, response *indexedmap.RemoveResponse) RemoveProposal {
	return &removeProposal{
		Proposal: newProposal(id, session),
		request:  request,
		response: response,
	}
}

type removeProposal struct {
	Proposal
	request  *indexedmap.RemoveRequest
	response *indexedmap.RemoveResponse
}

func (p *removeProposal) Request() *indexedmap.RemoveRequest {
	return p.request
}

func (p *removeProposal) Reply(reply *indexedmap.RemoveResponse) error {
	if p.response != nil {
		return errors.NewConflict("reply already sent")
	}
	p.response = reply
	return nil
}

var _ RemoveProposal = &removeProposal{}

type ClearProposals interface {
	register(ClearProposal)
	unregister(ProposalID)
	Get(ProposalID) (ClearProposal, bool)
	List() []ClearProposal
}

func newClearProposals() ClearProposals {
	return &clearProposals{
		proposals: make(map[ProposalID]ClearProposal),
	}
}

type clearProposals struct {
	proposals map[ProposalID]ClearProposal
}

func (p *clearProposals) register(proposal ClearProposal) {
	p.proposals[proposal.ID()] = proposal
}

func (p *clearProposals) unregister(id ProposalID) {
	delete(p.proposals, id)
}

func (p *clearProposals) Get(id ProposalID) (ClearProposal, bool) {
	proposal, ok := p.proposals[id]
	return proposal, ok
}

func (p *clearProposals) List() []ClearProposal {
	proposals := make([]ClearProposal, 0, len(p.proposals))
	for _, proposal := range p.proposals {
		proposals = append(proposals, proposal)
	}
	return proposals
}

var _ ClearProposals = &clearProposals{}

type ClearProposal interface {
	Proposal
	Request() *indexedmap.ClearRequest
	Reply(*indexedmap.ClearResponse) error
}

func newClearProposal(id ProposalID, session Session, request *indexedmap.ClearRequest, response *indexedmap.ClearResponse) ClearProposal {
	return &clearProposal{
		Proposal: newProposal(id, session),
		request:  request,
		response: response,
	}
}

type clearProposal struct {
	Proposal
	request  *indexedmap.ClearRequest
	response *indexedmap.ClearResponse
}

func (p *clearProposal) Request() *indexedmap.ClearRequest {
	return p.request
}

func (p *clearProposal) Reply(reply *indexedmap.ClearResponse) error {
	if p.response != nil {
		return errors.NewConflict("reply already sent")
	}
	p.response = reply
	return nil
}

var _ ClearProposal = &clearProposal{}

type EventsProposals interface {
	register(EventsProposal)
	unregister(ProposalID)
	Get(ProposalID) (EventsProposal, bool)
	List() []EventsProposal
}

func newEventsProposals() EventsProposals {
	return &eventsProposals{
		proposals: make(map[ProposalID]EventsProposal),
	}
}

type eventsProposals struct {
	proposals map[ProposalID]EventsProposal
}

func (p *eventsProposals) register(proposal EventsProposal) {
	p.proposals[proposal.ID()] = proposal
}

func (p *eventsProposals) unregister(id ProposalID) {
	delete(p.proposals, id)
}

func (p *eventsProposals) Get(id ProposalID) (EventsProposal, bool) {
	proposal, ok := p.proposals[id]
	return proposal, ok
}

func (p *eventsProposals) List() []EventsProposal {
	proposals := make([]EventsProposal, 0, len(p.proposals))
	for _, proposal := range p.proposals {
		proposals = append(proposals, proposal)
	}
	return proposals
}

var _ EventsProposals = &eventsProposals{}

type EventsProposal interface {
	Proposal
	Request() *indexedmap.EventsRequest
	Notify(*indexedmap.EventsResponse) error
	Close() error
}

func newEventsProposal(id ProposalID, session Session, request *indexedmap.EventsRequest, stream rsm.Stream) EventsProposal {
	return &eventsProposal{
		Proposal: newProposal(id, session),
		request:  request,
		stream:   stream,
	}
}

type eventsProposal struct {
	Proposal
	request *indexedmap.EventsRequest
	stream  rsm.Stream
}

func (p *eventsProposal) Request() *indexedmap.EventsRequest {
	return p.request
}

func (p *eventsProposal) Notify(notification *indexedmap.EventsResponse) error {
	bytes, err := proto.Marshal(notification)
	if err != nil {
		return err
	}
	p.stream.Value(bytes)
	return nil
}

func (p *eventsProposal) Close() error {
	p.stream.Close()
	return nil
}

var _ EventsProposal = &eventsProposal{}

type EntriesProposals interface {
	register(EntriesProposal)
	unregister(ProposalID)
	Get(ProposalID) (EntriesProposal, bool)
	List() []EntriesProposal
}

func newEntriesProposals() EntriesProposals {
	return &entriesProposals{
		proposals: make(map[ProposalID]EntriesProposal),
	}
}

type entriesProposals struct {
	proposals map[ProposalID]EntriesProposal
}

func (p *entriesProposals) register(proposal EntriesProposal) {
	p.proposals[proposal.ID()] = proposal
}

func (p *entriesProposals) unregister(id ProposalID) {
	delete(p.proposals, id)
}

func (p *entriesProposals) Get(id ProposalID) (EntriesProposal, bool) {
	proposal, ok := p.proposals[id]
	return proposal, ok
}

func (p *entriesProposals) List() []EntriesProposal {
	proposals := make([]EntriesProposal, 0, len(p.proposals))
	for _, proposal := range p.proposals {
		proposals = append(proposals, proposal)
	}
	return proposals
}

var _ EntriesProposals = &entriesProposals{}

type EntriesProposal interface {
	Proposal
	Request() *indexedmap.EntriesRequest
	Notify(*indexedmap.EntriesResponse) error
	Close() error
}

func newEntriesProposal(id ProposalID, session Session, request *indexedmap.EntriesRequest, stream rsm.Stream) EntriesProposal {
	return &entriesProposal{
		Proposal: newProposal(id, session),
		request:  request,
		stream:   stream,
	}
}

type entriesProposal struct {
	Proposal
	request *indexedmap.EntriesRequest
	stream  rsm.Stream
}

func (p *entriesProposal) Request() *indexedmap.EntriesRequest {
	return p.request
}

func (p *entriesProposal) Notify(notification *indexedmap.EntriesResponse) error {
	bytes, err := proto.Marshal(notification)
	if err != nil {
		return err
	}
	p.stream.Value(bytes)
	return nil
}

func (p *entriesProposal) Close() error {
	p.stream.Close()
	return nil
}

var _ EntriesProposal = &entriesProposal{}
