package repository

import (
	"errors"

	"github.com/derkellernerd/dori/core"
	"github.com/derkellernerd/dori/model"
)

type TwitchEvent struct {
	env *core.Environment
}

func NewTwitchEvent(env *core.Environment) *TwitchEvent {
	return &TwitchEvent{
		env: env,
	}
}

func (r TwitchEvent) Migrate() error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	err = db.AutoMigrate(&model.TwitchEvent{})
	if err != nil {
		return err
	}

	return nil
}

var ErrTwitchEventNotFound = errors.New("TwitchEvent not found")

func (r TwitchEvent) TwitchEventFindAll() ([]model.TwitchEvent, error) {
	var items []model.TwitchEvent
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return items, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Find(&items)
	if result.Error != nil {
		return items, result.Error
	}
	return items, result.Error
}

func (r TwitchEvent) TwitchEventFindById(id uint) (model.TwitchEvent, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return model.TwitchEvent{}, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var item model.TwitchEvent
	result := db.Find(&item, "id = ?", id)
	if result.Error != nil {
		return model.TwitchEvent{}, result.Error
	}

	if result.RowsAffected == 0 {
		return model.TwitchEvent{}, ErrTwitchEventNotFound
	}
	return item, result.Error
}

func (r TwitchEvent) TwitchEventInsert(item *model.TwitchEvent) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Create(item)
	return result.Error
}

func (r TwitchEvent) TwitchEventUpdate(item *model.TwitchEvent) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Updates(item)
	return result.Error
}

func (r TwitchEvent) TwitchEventDelete(item *model.TwitchEvent) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Unscoped().Delete(item)
	return result.Error
}

func (r TwitchEvent) TwitchEventFindByTwitchEventSubscripton(twitchEventSubscription string) (model.TwitchEvent, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return model.TwitchEvent{}, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var item model.TwitchEvent
	result := db.Find(&item, "twitch_event_subscription = ?", twitchEventSubscription)
	if result.Error != nil {
		return model.TwitchEvent{}, result.Error
	}

	if result.RowsAffected == 0 {
		return model.TwitchEvent{}, ErrTwitchEventNotFound
	}
	return item, result.Error
}
