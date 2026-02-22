package usecase

import (
	"github.com/durianpay/fullstack-boilerplate/internal/entity"
	"github.com/durianpay/fullstack-boilerplate/internal/module/payment/repository"
)

type PaymentUsecase interface {
	ListPayments(filter entity.PaymentFilter) ([]entity.Payment, error)
}

type Payment struct {
	repo repository.PaymentRepository
}

func NewPaymentUsecase(repo repository.PaymentRepository) *Payment {
	return &Payment{repo: repo}
}

func (p *Payment) ListPayments(filter entity.PaymentFilter) ([]entity.Payment, error) {
	// Validate status if provided
	if filter.Status != nil && *filter.Status != "" {
		valid := map[string]bool{"completed": true, "processing": true, "failed": true}
		if !valid[*filter.Status] {
			return nil, entity.ErrorBadRequest("invalid status, must be one of: completed, processing, failed")
		}
	}
	return p.repo.ListPayments(filter)
}
