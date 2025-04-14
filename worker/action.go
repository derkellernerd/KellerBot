package worker

import (
	"github.com/derkellernerd/kellerbot/core"
	"github.com/derkellernerd/kellerbot/model"
	"github.com/derkellernerd/kellerbot/repository"
	"github.com/google/uuid"
)

type Action struct {
	env           *core.Environment
	actionRepo    *repository.Action
	actionChannel map[string]chan model.Action
}

func NewAction(env *core.Environment, actionRepo *repository.Action) *Action {
	return &Action{
		env:           env,
		actionRepo:    actionRepo,
		actionChannel: make(map[string]chan model.Action),
	}
}

func (a *Action) RegisterListener() (string, chan model.Action) {
	id := uuid.NewString()
	a.actionChannel[id] = make(chan model.Action)

	return id, a.actionChannel[id]
}

func (a *Action) UnregisterListener(id string) {
	delete(a.actionChannel, id)
}
