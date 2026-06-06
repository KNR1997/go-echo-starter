package dept

import (
	"context"
	"fmt"
	"go-echo-starter/internal/domain"
	"go-echo-starter/internal/models"
)

type deptRepository interface {
	GetById(ctx context.Context, id uint) (models.Department, error)
	GetDepartments(ctx context.Context) ([]models.Department, error)
	GetDepartmentPaginated(ctx context.Context, pagination domain.Pagination) ([]models.Department, int64, error)
	Create(ctx context.Context, dept *models.Department) error
	Update(ctx context.Context, dept *models.Department) error
	Delete(ctx context.Context, post *models.Department) error
	ExistsByName(ctx context.Context, name string) (bool, error)
}

type Service struct {
	deptRepository deptRepository
}

func NewService(deptRepository deptRepository) *Service {
	return &Service{deptRepository: deptRepository}
}

func (s *Service) GetDepartments(ctx context.Context) ([]models.Department, error) {
	depts, err := s.deptRepository.GetDepartments(ctx)
	if err != nil {
		return nil, fmt.Errorf("get departments from repository: %w", err)
	}

	return depts, nil
}

func (s *Service) GetDepartmentPaginated(
	ctx context.Context,
	pagination domain.Pagination,
) ([]models.Department, int64, error) {

	depts, total, err := s.deptRepository.GetDepartmentPaginated(
		ctx,
		pagination,
	)
	if err != nil {
		return nil, 0, fmt.Errorf(
			"get departments from repository: %w",
			err,
		)
	}

	return depts, total, nil
}

func (s *Service) Create(ctx context.Context, dept *models.Department) error {
	exists, err := s.deptRepository.ExistsByName(ctx, dept.Name)
	if err != nil {
		return fmt.Errorf("check department exists: %w", err)
	}

	if exists {
		return fmt.Errorf("department name already exists")
	}

	if err := s.deptRepository.Create(ctx, dept); err != nil {
		return fmt.Errorf("create department in repository: %w", err)
	}

	return nil
}

func (s *Service) Update(ctx context.Context, request domain.UpdateDepartmentRequest) (*models.Department, error) {
	dept, err := s.deptRepository.GetById(ctx, request.DeptID)
	if err != nil {
		return nil, fmt.Errorf("get stored department from repository: %w", err)
	}

	if dept.Name != request.Name {
		exists, err := s.deptRepository.ExistsByName(ctx, request.Name)
		if err != nil {
			return nil, fmt.Errorf("check department exists: %w", err)
		}

		if exists {
			return nil, fmt.Errorf("department name already exists")
		}
	}

	dept.Name = request.Name
	dept.Desc = request.Desc

	if err := s.deptRepository.Update(ctx, &dept); err != nil {
		return nil, fmt.Errorf("update department in repository: %w", err)
	}

	return &dept, nil
}

func (s *Service) Delete(ctx context.Context, request domain.DeleteDepartmentRequest) error {
	dept, err := s.deptRepository.GetById(ctx, request.DeptID)
	if err != nil {
		return fmt.Errorf("get stored department from repository: %w", err)
	}

	if err := s.deptRepository.Delete(ctx, &dept); err != nil {
		return fmt.Errorf("delete department in repository: %w", err)
	}

	return nil
}
