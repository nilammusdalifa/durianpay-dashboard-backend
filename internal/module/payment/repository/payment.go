package repository

import (
	"sort"
	"strings"
	"sync"

	"github.com/durianpay/fullstack-boilerplate/internal/entity"
)

type PaymentRepository interface {
	ListPayments(filter entity.PaymentFilter) ([]entity.Payment, error)
}

type InMemoryPaymentRepo struct {
	mu       sync.RWMutex
	payments []entity.Payment
}

func NewInMemoryPaymentRepo() *InMemoryPaymentRepo {
	return &InMemoryPaymentRepo{}
}

func (r *InMemoryPaymentRepo) Seed(payments []entity.Payment) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.payments = payments
}

func (r *InMemoryPaymentRepo) ListPayments(filter entity.PaymentFilter) ([]entity.Payment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]entity.Payment, 0, len(r.payments))
	for _, p := range r.payments {
		if filter.Status != nil && *filter.Status != "" && p.Status != *filter.Status {
			continue
		}
		if filter.ID != nil && *filter.ID != "" && p.ID != *filter.ID {
			continue
		}
		result = append(result, p)
	}

	// Sorting
	sortField := "created_at"
	desc := true
	if filter.Sort != nil && *filter.Sort != "" {
		s := *filter.Sort
		if strings.HasPrefix(s, "-") {
			desc = true
			sortField = s[1:]
		} else {
			desc = false
			sortField = s
		}
	}

	sort.SliceStable(result, func(i, j int) bool {
		var less bool
		switch sortField {
		case "amount":
			less = result[i].Amount < result[j].Amount
		case "merchant":
			less = result[i].Merchant < result[j].Merchant
		case "status":
			less = result[i].Status < result[j].Status
		default: // created_at
			less = result[i].CreatedAt.Before(result[j].CreatedAt)
		}
		if desc {
			return !less
		}
		return less
	})

	return result, nil
}
