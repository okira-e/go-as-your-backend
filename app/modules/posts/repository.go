package posts

import (
	"context"
	"errors"

	"github.com/okira-e/go-as-your-backend/app/models"
	"github.com/okira-e/go-as-your-backend/app/spec"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) spec.Repository[models.Post] {
	return &Repository{db: db}
}

func (self *Repository) Create(ctx context.Context, entity *models.Post) (*models.Post, error) {
	if entity == nil {
		return nil, errors.New("entity cannot be nil")
	}

	err := self.db.WithContext(ctx).Create(entity).Error
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (self *Repository) FindByID(ctx context.Context, id string) (*models.Post, error) {
	var entity models.Post

	err := self.db.WithContext(ctx).First(&entity, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("entity not found")
		}

		return nil, err
	}

	return &entity, nil
}

func (self *Repository) FindAll(ctx context.Context, queryOptions *spec.QueryOptions, filter *spec.Filter) ([]models.Post, error) {
	var entities []models.Post

	tx := self.db.WithContext(ctx)
	tx = spec.ApplyPagination(tx, queryOptions)
	tx, err := spec.ApplyFilters(tx, filter, models.Post{})
	if err != nil {
		return entities, err
	}

	err = tx.Find(&entities).Error
	if err != nil {
		return nil, err
	}

	return entities, nil
}

func (self *Repository) Count(ctx context.Context, filter *spec.Filter) (int64, error) {
	var count int64

	tx := self.db.WithContext(ctx)
	tx, err := spec.ApplyFilters(tx, filter, models.Post{})
	if err != nil {
		return 0, err
	}

	if err := tx.Model(&models.Post{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (self *Repository) Update(ctx context.Context, entity *models.Post) error {
	if entity == nil {
		return errors.New("entity cannot be nil")
	}

	result := self.db.WithContext(ctx).Model(entity).Updates(entity)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("entity not found")
	}

	return nil
}

func (self *Repository) Delete(ctx context.Context, id string) error {
	result := self.db.WithContext(ctx).Delete(&models.Post{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("entity not found")
	}

	return nil
}

func (self *Repository) Exists(ctx context.Context, id string) (bool, error) {
	var count int64

	err := self.db.WithContext(ctx).Model(&models.Post{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
