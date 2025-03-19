package repository

import (
	"errors"

	"github.com/derkellernerd/dori/core"
	"github.com/derkellernerd/dori/model"
)

type Command struct {
	env *core.Environment
}

func NewCommand(env *core.Environment) *Command {
	return &Command{
		env: env,
	}
}

func (r Command) Migrate() error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	err = db.AutoMigrate(&model.Command{})
	if err != nil {
		return err
	}

	return nil
}

var ErrCommandNotFound = errors.New("Command not found")

func (r Command) CommandFindAll() ([]model.Command, error) {
	var items []model.Command
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

func (r Command) CommandFindById(id uint) (model.Command, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return model.Command{}, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var item model.Command
	result := db.Find(&item, "id = ?", id)
	if result.Error != nil {
		return model.Command{}, result.Error
	}

	if result.RowsAffected == 0 {
		return model.Command{}, ErrCommandNotFound
	}
	return item, result.Error
}

func (r Command) CommandInsert(item *model.Command) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Create(item)
	return result.Error
}

func (r Command) CommandUpdate(item *model.Command) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Updates(item)
	return result.Error
}

func (r Command) CommandDelete(item *model.Command) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Unscoped().Delete(item)
	return result.Error
}

func (r Command) CommandFindByCommand(command string) (model.Command, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return model.Command{}, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var item model.Command
	result := db.Find(&item, "command = ?", command)
	if result.Error != nil {
		return model.Command{}, result.Error
	}

	if result.RowsAffected == 0 {
		return model.Command{}, ErrCommandNotFound
	}
	return item, result.Error
}
