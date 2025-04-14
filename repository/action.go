package repository

import (
	"errors"

	"github.com/derkellernerd/kellerbot/core"
	"github.com/derkellernerd/kellerbot/model"
)

type Action struct {
	env *core.Environment
}

func NewAction(env *core.Environment) *Action {
	return &Action{
		env: env,
	}
}

func (r Action) Migrate() error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	err = db.AutoMigrate(&model.Action{})
	if err != nil {
		return err
	}

	return nil
}

var ErrActionNotFound = errors.New("Action not found")

func (r Action) ActionFindAll() ([]model.Action, error) {
	var items []model.Action
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

func (r Action) ActionFindById(id uint) (model.Action, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return model.Action{}, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var item model.Action
	result := db.Find(&item, "id = ?", id)
	if result.Error != nil {
		return model.Action{}, result.Error
	}

	if result.RowsAffected == 0 {
		return model.Action{}, ErrActionNotFound
	}
	return item, result.Error
}

func (r Action) ActionInsert(item *model.Action) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Create(item)
	return result.Error
}

func (r Action) ActionUpdate(item *model.Action) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Updates(item)
	return result.Error
}

func (r Action) ActionDelete(item *model.Action) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Delete(item)
	return result.Error
}

func (r Action) ActionFindByActionName(actionName string) (model.Action, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return model.Action{}, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var item model.Action
	result := db.Find(&item, "action_name = ?", actionName)
	if result.Error != nil {
		return model.Action{}, result.Error
	}

	if result.RowsAffected == 0 {
		return model.Action{}, ErrActionNotFound
	}
	return item, result.Error
}
