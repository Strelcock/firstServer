package user

import "firstServer/pkg/db"

type UserRepository struct {
	Database *db.Db
}

func NewUserRepository(database *db.Db) *UserRepository {
	return &UserRepository{
		Database: database,
	}
}

func (ur *UserRepository) Create(user *User) (*User, error) {
	result := ur.Database.DB.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (ur *UserRepository) FindByEmail(email string) (*User, error) {
	var foundUser User
	result := ur.Database.DB.First(&foundUser, "email = ?", email)
	if result.Error != nil {
		return nil, result.Error
	}
	return &foundUser, nil
}
