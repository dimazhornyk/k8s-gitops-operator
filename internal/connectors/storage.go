package connectors

import "diploma/internal/common"

type storage struct {
}

func (s storage) GetLastEvent() common.RepositoryEvent {
	return common.RepositoryEvent{}
}

func (s storage) SaveEvent(event common.RepositoryEvent) error {
	return nil
}
