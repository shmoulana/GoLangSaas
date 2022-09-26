package repository

import (
	"context"

	"github.com/shmoulana/Redios/internal/model"
	"github.com/shmoulana/Redios/pkg/database"
)

type UserRepository struct {
	databaseRepo database.DatabaseRepo
}

func (r UserRepository) FindById(ctx context.Context, id int) (*model.User, error) {
	var user model.User
	db, err := r.databaseRepo.GetDB(ctx)
	if err != nil {
		return nil, err
	}

	result := db.Where("id = ?", id).Find(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	result.Scan(&user)

	return &user, nil
}

func (r UserRepository) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	db, err := r.databaseRepo.GetDB(ctx)
	if err != nil {
		return nil, err
	}

	result := db.Where("email = ?", email).Find(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	result.Scan(&user)

	return &user, nil
}

func (r UserRepository) Create(ctx context.Context, user model.User) (*int, error) {
	db, err := r.databaseRepo.GetDB(ctx)
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&user)
	if err != nil {
		return nil, err
	}

	result := db.Model(&user).Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	// logging
	// rowsAffected := result.RowsAffected

	return &user.ID, nil
}

func (r UserRepository) Update(ctx context.Context, user model.User) (*int, error) {
	db, err := r.databaseRepo.GetDB(ctx)
	if err != nil {
		return nil, err
	}

	result := db.Save(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	// logging
	// rowsAffected := result.RowsAffected

	return &user.ID, nil
}

func NewUserRepository(db database.DatabaseRepo) UserRepository {
	return UserRepository{
		databaseRepo: db,
	}
}
