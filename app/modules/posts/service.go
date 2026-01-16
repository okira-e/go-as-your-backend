package posts

import (
	"context"

	"github.com/okira-e/go-as-your-backend/app/models"
	"github.com/okira-e/go-as-your-backend/app/spec"
)

type Service struct {
	repository spec.Repository[models.Post]
}

func NewService(repository spec.Repository[models.Post]) *Service {
	return &Service{repository: repository}
}

func (self *Service) FindAll(ctx context.Context, queryOptions *spec.QueryOptions, filter *spec.Filter) ([]models.Post, error) {
	entities, err := self.repository.FindAll(ctx, queryOptions, filter)
	if err != nil {
		return entities, err
	}

	return entities, nil
}

func (self *Service) GetPublished(ctx context.Context) ([]models.Post, error) {
	queryOptions := spec.QueryOptions{
		Limit: 10,
	}

	filter := spec.Filter{
		Where: spec.WhereClause{
			And: []spec.WhereCondition{
				{
					Column:   "published",
					Operator: "=",
					Value:    true,
				},
			},
		},
	}

	entities, err := self.repository.FindAll(ctx, &queryOptions, &filter)
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

func (self *Service) Create(ctx context.Context, entityDto *models.CreatePostDto, userId string) (*models.PostDto, error) {
	entity := entityDto.FromDto(userId)

	entity, err := self.repository.Create(ctx, entity)
	if err != nil {
		return &models.PostDto{}, err
	}

	return entity.ToDto(), nil
}
