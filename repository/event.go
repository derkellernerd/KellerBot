package repository

import (
	"errors"

	"github.com/derkellernerd/kellerbot/core"
	"github.com/derkellernerd/kellerbot/model"
)

type Event struct {
	env *core.Environment
}

func NewEvent(env *core.Environment) *Event {
	return &Event{
		env: env,
	}
}

func (r Event) Migrate() error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	err = db.AutoMigrate(&model.Event{})
	if err != nil {
		return err
	}

	return nil
}

var ErrEventNotFound = errors.New("Event not found")

func (r Event) EventFindAll() ([]model.Event, error) {
	var items []model.Event
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return items, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Order("created_at desc").Find(&items)
	if result.Error != nil {
		return items, result.Error
	}
	return items, result.Error
}

func (r Event) EventFindById(id uint) (model.Event, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return model.Event{}, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var item model.Event
	result := db.Find(&item, "id = ?", id)
	if result.Error != nil {
		return model.Event{}, result.Error
	}

	if result.RowsAffected == 0 {
		return model.Event{}, ErrEventNotFound
	}
	return item, result.Error
}

func (r Event) EventInsert(item *model.Event) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Create(item)
	return result.Error
}

func (r Event) EventUpdate(item *model.Event) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Updates(item)
	return result.Error
}

func (r Event) EventDelete(item *model.Event) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Unscoped().Delete(item)
	return result.Error
}
