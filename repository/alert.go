package repository

import (
	"errors"

	"github.com/derkellernerd/dori/core"
	"github.com/derkellernerd/dori/model"
)

type Alert struct {
	env *core.Environment
}

var ErrAlertNotFound = errors.New("Alert not found")

func NewAlert(env *core.Environment) *Alert {
	return &Alert{
		env: env,
	}
}

func (r Alert) Migrate() error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	err = db.AutoMigrate(&model.Alert{})
	if err != nil {
		return err
	}

	return nil
}

func (r Alert) AlertFindAll() ([]model.Alert, error) {
	var items []model.Alert
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

func (r Alert) AlertFindById(id uint) (model.Alert, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return model.Alert{}, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var item model.Alert
	result := db.Find(&item, "id = ?", id)
	if result.Error != nil {
		return model.Alert{}, result.Error
	}

	if result.RowsAffected == 0 {
		return model.Alert{}, ErrAlertNotFound
	}
	return item, result.Error
}

func (r Alert) AlertInsert(item *model.Alert) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Create(item)
	return result.Error
}

func (r Alert) AlertUpdate(item *model.Alert) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Updates(item)
	return result.Error
}

func (r Alert) AlertDelete(item *model.Alert) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Delete(item)
	return result.Error
}

func (r Alert) AlertFindByName(name string) (model.Alert, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return model.Alert{}, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var item model.Alert
	result := db.Find(&item, "name = ?", name)
	if result.Error != nil {
		return model.Alert{}, result.Error
	}

	if result.RowsAffected == 0 {
		return model.Alert{}, ErrAlertNotFound
	}
	return item, result.Error
}
