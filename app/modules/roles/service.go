package roles

import (
	"context"

	"github.com/okira-e/go-as-your-backend/app/models"
	"github.com/okira-e/go-as-your-backend/app/spec"
)

type Service struct {
	repository spec.Repository[models.Role]
}

func NewService(repository spec.Repository[models.Role]) *Service {
	return &Service{repository: repository}
}

func (self *Service) FindAll(ctx context.Context, queryOptions *spec.QueryOptions, filter *spec.Filter) ([]models.Role, error) {
	entities, err := self.repository.FindAll(ctx, queryOptions, filter)
	if err != nil {
		return entities, err
	}

	return entities, nil
}

func (self *Service) GetCount(ctx context.Context, filter *spec.Filter) (int64, error) {
	count, err := self.repository.Count(ctx, filter)
	if err != nil {
		return count, err
	}

	return count, nil
}

func (self *Service) Create(ctx context.Context, entityDto *models.CreateRoleDto) (*models.RoleDto, error) {
	entity := entityDto.FromDto()

	entity, err := self.repository.Create(ctx, entity)
	if err != nil {
		return &models.RoleDto{}, err
	}

	return entity.ToDto(), nil
}
