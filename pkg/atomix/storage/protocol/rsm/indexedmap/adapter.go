// Code generated by atomix-go-framework. DO NOT EDIT.
package indexedmap

import (
	"github.com/atomix/atomix-go-framework/pkg/atomix/errors"
	"github.com/atomix/atomix-go-framework/pkg/atomix/logging"
	"github.com/atomix/atomix-go-framework/pkg/atomix/storage/protocol/rsm"
	"github.com/gogo/protobuf/proto"
	"io"
)

var log = logging.GetLogger("atomix", "indexedmap", "service")

const Type = "IndexedMap"

const (
	sizeOp       = "Size"
	putOp        = "Put"
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
	case 2:
		p, err := newPutProposal(command)
		if err != nil {
			err = errors.NewInternal(err.Error())
			log.Error(err)
			command.Output(nil, err)
			return
		}

		log.Debugf("Proposal PutProposal %s", p)
		response, err := s.rsm.Put(p)
		if err != nil {
			log.Warnf("Proposal PutProposal %s failed: %v", p, err)
			command.Output(nil, err)
		} else {
			output, err := proto.Marshal(response)
			if err != nil {
				err = errors.NewInternal(err.Error())
				log.Errorf("Proposal PutProposal %s failed: %v", p, err)
				command.Output(nil, err)
			} else {
				log.Debugf("Proposal PutProposal %s complete: %+v", p, response)
				command.Output(output, nil)
			}
		}
	case 8:
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
				log.Debugf("Proposal RemoveProposal %s complete: %+v", p, response)
				command.Output(output, nil)
			}
		}
	case 9:
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
				log.Debugf("Proposal ClearProposal %s complete: %+v", p, response)
				command.Output(output, nil)
			}
		}
	case 10:
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
				log.Debugf("Querying SizeQuery %s complete: %+v", q, response)
				query.Output(output, nil)
			}
		}
	case 3:
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
				log.Debugf("Querying GetQuery %s complete: %+v", q, response)
				query.Output(output, nil)
			}
		}
	case 4:
		q, err := newFirstEntryQuery(query)
		if err != nil {
			err = errors.NewInternal(err.Error())
			log.Error(err)
			query.Output(nil, err)
			return
		}

		log.Debugf("Querying FirstEntryQuery %s", q)
		response, err := s.rsm.FirstEntry(q)
		if err != nil {
			log.Warnf("Querying FirstEntryQuery %s failed: %v", q, err)
			query.Output(nil, err)
		} else {
			output, err := proto.Marshal(response)
			if err != nil {
				err = errors.NewInternal(err.Error())
				log.Errorf("Querying FirstEntryQuery %s failed: %v", q, err)
				query.Output(nil, err)
			} else {
				log.Debugf("Querying FirstEntryQuery %s complete: %+v", q, response)
				query.Output(output, nil)
			}
		}
	case 5:
		q, err := newLastEntryQuery(query)
		if err != nil {
			err = errors.NewInternal(err.Error())
			log.Error(err)
			query.Output(nil, err)
			return
		}

		log.Debugf("Querying LastEntryQuery %s", q)
		response, err := s.rsm.LastEntry(q)
		if err != nil {
			log.Warnf("Querying LastEntryQuery %s failed: %v", q, err)
			query.Output(nil, err)
		} else {
			output, err := proto.Marshal(response)
			if err != nil {
				err = errors.NewInternal(err.Error())
				log.Errorf("Querying LastEntryQuery %s failed: %v", q, err)
				query.Output(nil, err)
			} else {
				log.Debugf("Querying LastEntryQuery %s complete: %+v", q, response)
				query.Output(output, nil)
			}
		}
	case 6:
		q, err := newPrevEntryQuery(query)
		if err != nil {
			err = errors.NewInternal(err.Error())
			log.Error(err)
			query.Output(nil, err)
			return
		}

		log.Debugf("Querying PrevEntryQuery %s", q)
		response, err := s.rsm.PrevEntry(q)
		if err != nil {
			log.Warnf("Querying PrevEntryQuery %s failed: %v", q, err)
			query.Output(nil, err)
		} else {
			output, err := proto.Marshal(response)
			if err != nil {
				err = errors.NewInternal(err.Error())
				log.Errorf("Querying PrevEntryQuery %s failed: %v", q, err)
				query.Output(nil, err)
			} else {
				log.Debugf("Querying PrevEntryQuery %s complete: %+v", q, response)
				query.Output(output, nil)
			}
		}
	case 7:
		q, err := newNextEntryQuery(query)
		if err != nil {
			err = errors.NewInternal(err.Error())
			log.Error(err)
			query.Output(nil, err)
			return
		}

		log.Debugf("Querying NextEntryQuery %s", q)
		response, err := s.rsm.NextEntry(q)
		if err != nil {
			log.Warnf("Querying NextEntryQuery %s failed: %v", q, err)
			query.Output(nil, err)
		} else {
			output, err := proto.Marshal(response)
			if err != nil {
				err = errors.NewInternal(err.Error())
				log.Errorf("Querying NextEntryQuery %s failed: %v", q, err)
				query.Output(nil, err)
			} else {
				log.Debugf("Querying NextEntryQuery %s complete: %+v", q, response)
				query.Output(output, nil)
			}
		}
	case 11:
		q, err := newEntriesQuery(query)
		if err != nil {
			err = errors.NewInternal(err.Error())
			log.Error(err)
			query.Output(nil, err)
			return
		}

		log.Debugf("Querying EntriesQuery %s", q)
		s.rsm.Entries(q)
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
