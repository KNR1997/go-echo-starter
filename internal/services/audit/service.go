package audit

import (
	"context"
	"fmt"
	"go-echo-starter/internal/domain"
	"go-echo-starter/internal/models"
	"go-echo-starter/internal/utils"
)

type auditRepository interface {
	Create(ctx context.Context, dept *models.AuditLog) error
	GetAuditLogPaginated(
		ctx context.Context,
		pagination domain.Pagination,
		searchConditions []utils.SearchCondition,
		searchJoin string,
		method string,
	) ([]models.AuditLog, int64, error)
}

type Service struct {
	auditRepository auditRepository
}

func NewService(auditRepository auditRepository) *Service {
	return &Service{auditRepository: auditRepository}
}

func (s *Service) Create(ctx context.Context, auditLog *models.AuditLog) error {
	return s.auditRepository.Create(ctx, auditLog)
}

func (s *Service) GetAuditLogPaginated(
	ctx context.Context,
	pagination domain.Pagination,
	searchConditions []utils.SearchCondition,
	searchJoin string,
	method string,
) ([]models.AuditLog, int64, error) {
	auditLogs, total, err := s.auditRepository.GetAuditLogPaginated(
		ctx,
		pagination,
		searchConditions,
		searchJoin,
		method,
	)
	if err != nil {
		return nil, 0, fmt.Errorf(
			"get auditLogs from repository: %w",
			err,
		)
	}

	return auditLogs, total, nil
}
