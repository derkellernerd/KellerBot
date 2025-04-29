package repository

import (
	"errors"

	"github.com/derkellernerd/kellerbot/core"
	"github.com/derkellernerd/kellerbot/model"
)

type User struct {
	env *core.Environment
}

func NewUser(env *core.Environment) *User {
	return &User{
		env: env,
	}
}

func (r User) Migrate() error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	err = db.AutoMigrate(&model.User{})
	if err != nil {
		return err
	}

	return nil
}

var ErrUserNotFound = errors.New("User not found")

func (r User) UserFindAll() ([]model.User, error) {
	var items []model.User
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

func (r User) UserFindById(id uint) (model.User, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return model.User{}, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var item model.User
	result := db.Find(&item, "id = ?", id)
	if result.Error != nil {
		return model.User{}, result.Error
	}

	if result.RowsAffected == 0 {
		return model.User{}, ErrUserNotFound
	}
	return item, result.Error
}

func (r User) UserInsert(item *model.User) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Create(item)
	return result.Error
}

func (r User) UserUpdate(item *model.User) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Updates(item)
	return result.Error
}

func (r User) UserDelete(item *model.User) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Unscoped().Delete(item)
	return result.Error
}

func (r User) UserFindByUsername(username string) (model.User, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return model.User{}, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var item model.User
	result := db.Find(&item, "username = ?", username)
	if result.Error != nil {
		return model.User{}, result.Error
	}

	if result.RowsAffected == 0 {
		return model.User{}, ErrUserNotFound
	}
	return item, result.Error
}
