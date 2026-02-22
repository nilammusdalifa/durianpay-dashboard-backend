package usecase_test

import (
	"testing"

	"github.com/durianpay/fullstack-boilerplate/internal/entity"
	"github.com/durianpay/fullstack-boilerplate/internal/module/payment/usecase"
)

// mockPaymentRepo implements repository.PaymentRepository
type mockPaymentRepo struct {
	payments []entity.Payment
	err      error
}

func (m *mockPaymentRepo) ListPayments(filter entity.PaymentFilter) ([]entity.Payment, error) {
	if m.err != nil {
		return nil, m.err
	}
	if filter.Status == nil {
		return m.payments, nil
	}
	var result []entity.Payment
	for _, p := range m.payments {
		if p.Status == *filter.Status {
			result = append(result, p)
		}
	}
	return result, nil
}

func TestListPayments_NoFilter(t *testing.T) {
	payments := []entity.Payment{
		{ID: "1", Merchant: "A", Status: "completed", Amount: "100.00"},
		{ID: "2", Merchant: "B", Status: "failed", Amount: "200.00"},
	}
	repo := &mockPaymentRepo{payments: payments}
	uc := usecase.NewPaymentUsecase(repo)

	result, err := uc.ListPayments(entity.PaymentFilter{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 payments, got %d", len(result))
	}
}

func TestListPayments_FilterByStatus(t *testing.T) {
	payments := []entity.Payment{
		{ID: "1", Merchant: "A", Status: "completed", Amount: "100.00"},
		{ID: "2", Merchant: "B", Status: "failed", Amount: "200.00"},
		{ID: "3", Merchant: "C", Status: "completed", Amount: "300.00"},
	}
	repo := &mockPaymentRepo{payments: payments}
	uc := usecase.NewPaymentUsecase(repo)

	status := "completed"
	result, err := uc.ListPayments(entity.PaymentFilter{Status: &status})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 completed payments, got %d", len(result))
	}
}

func TestListPayments_InvalidStatus(t *testing.T) {
	repo := &mockPaymentRepo{}
	uc := usecase.NewPaymentUsecase(repo)

	status := "invalid_status"
	_, err := uc.ListPayments(entity.PaymentFilter{Status: &status})
	if err == nil {
		t.Fatal("expected error for invalid status, got nil")
	}
}

func TestListPayments_ValidStatuses(t *testing.T) {
	repo := &mockPaymentRepo{}
	uc := usecase.NewPaymentUsecase(repo)

	for _, s := range []string{"completed", "processing", "failed"} {
		sc := s
		_, err := uc.ListPayments(entity.PaymentFilter{Status: &sc})
		if err != nil {
			t.Errorf("expected no error for status %q, got %v", s, err)
		}
	}
}
