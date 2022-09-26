package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/shmoulana/Redios/internal/repository"
	"github.com/shmoulana/Redios/pkg/database"
	"github.com/shmoulana/Redios/pkg/dto"
)

type TenantService struct {
	TenantRepository repository.TenantRepository
}

func (s TenantService) CreateTenant(ctx context.Context, payload dto.TenantRequestV1) error {
	t := &database.Tenant{
		ID:          uuid.New().String(),
		Name:        payload.Name,
		DisplayName: payload.Name,
	}

	lastInsertedId, err := s.TenantRepository.CreateTenant(ctx, *t, payload.SeparateDb)
	if err != nil {
		return nil
	}

	fmt.Printf(*lastInsertedId)

	return nil
}

func (s TenantService) UpdateTenant(ctx context.Context, id string, payload dto.TenantRequestV1) (*string, error) {
	tenant, err := s.TenantRepository.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	tenant.Name = payload.Name

	_, err = s.TenantRepository.UpdateTenant(ctx, *tenant)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func (s TenantService) GetTenants(ctx context.Context) ([]database.Tenant, error) {
	tenants, err := s.TenantRepository.Find(ctx)
	if err != nil {
		return nil, err
	}

	return tenants, nil
}

func (s TenantService) GetTenantById(ctx context.Context, id string) (*database.Tenant, error) {
	tenant, err := s.TenantRepository.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	return tenant, err
}

func (s TenantService) DeleteTenant(ctx context.Context, id string) (*string, error) {
	tenant, err := s.TenantRepository.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	_, err = s.TenantRepository.DeleteTenant(ctx, *tenant)
	if err != nil {
		return nil, err
	}

	return &id, nil
}
