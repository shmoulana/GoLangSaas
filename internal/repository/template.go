package repository

import (
	"context"

	"github.com/shmoulana/Redios/internal/model"
	"github.com/shmoulana/Redios/pkg/database"
)

type TemplateRespository struct {
	databaseRepo database.DatabaseRepo
}

func (r TemplateRespository) Insert(ctx context.Context, template model.Template) (*int, error) {
	db, err := r.databaseRepo.GetDB(ctx)
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&template)
	if err != nil {
		return nil, err
	}

	result := db.Model(&template).Create(&template)
	if result.Error != nil {
		return nil, result.Error
	}

	return &template.ID, nil
}

func (r TemplateRespository) FindById(ctx context.Context, id int) (*model.Template, error) {
	var template model.Template
	db, err := r.databaseRepo.GetDB(ctx)
	if err != nil {
		return nil, err
	}

	result := db.Where("id = ?", id).Find(&template)
	if result.Error != nil {
		return nil, result.Error
	}

	result.Scan(&template)

	return &template, nil
}

func NewTemplateRepostory(db database.DatabaseRepo) TemplateRespository {
	return TemplateRespository{
		databaseRepo: db,
	}
}
