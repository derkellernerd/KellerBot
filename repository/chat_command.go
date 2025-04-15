package repository

import (
	"errors"

	"github.com/derkellernerd/kellerbot/core"
	"github.com/derkellernerd/kellerbot/model"
)

type ChatCommand struct {
	env *core.Environment
}

func NewChatCommand(env *core.Environment) *ChatCommand {
	return &ChatCommand{
		env: env,
	}
}

func (r ChatCommand) Migrate() error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	err = db.AutoMigrate(&model.ChatCommand{})
	if err != nil {
		return err
	}

	return nil
}

var ErrChatCommandNotFound = errors.New("ChatCommand not found")

func (r ChatCommand) ChatCommandFindAll() ([]model.ChatCommand, error) {
	var items []model.ChatCommand
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

func (r ChatCommand) ChatCommandFindById(id uint) (model.ChatCommand, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return model.ChatCommand{}, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var item model.ChatCommand
	result := db.Find(&item, "id = ?", id)
	if result.Error != nil {
		return model.ChatCommand{}, result.Error
	}

	if result.RowsAffected == 0 {
		return model.ChatCommand{}, ErrChatCommandNotFound
	}
	return item, result.Error
}

func (r ChatCommand) ChatCommandInsert(item *model.ChatCommand) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Create(item)
	return result.Error
}

func (r ChatCommand) ChatCommandUpdate(item *model.ChatCommand) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Updates(item)
	return result.Error
}

func (r ChatCommand) ChatCommandDelete(item *model.ChatCommand) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Unscoped().Delete(item)
	return result.Error
}

func (r ChatCommand) ChatCommandFindByChatCommand(chatCommand string) (model.ChatCommand, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return model.ChatCommand{}, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var item model.ChatCommand
	result := db.Find(&item, "command = ?", chatCommand)
	if result.Error != nil {
		return model.ChatCommand{}, result.Error
	}

	if result.RowsAffected == 0 {
		return model.ChatCommand{}, ErrChatCommandNotFound
	}
	return item, result.Error
}
