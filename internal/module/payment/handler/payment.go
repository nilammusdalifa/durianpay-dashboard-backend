package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/durianpay/fullstack-boilerplate/internal/entity"
	paymentUsecase "github.com/durianpay/fullstack-boilerplate/internal/module/payment/usecase"
	"github.com/durianpay/fullstack-boilerplate/internal/openapigen"
	"github.com/durianpay/fullstack-boilerplate/internal/transport"
	"github.com/golang-jwt/jwt/v5"
)

type PaymentHandler struct {
	paymentUC paymentUsecase.PaymentUsecase
	jwtSecret []byte
}

func NewPaymentHandler(paymentUC paymentUsecase.PaymentUsecase, jwtSecret []byte) *PaymentHandler {
	return &PaymentHandler{paymentUC: paymentUC, jwtSecret: jwtSecret}
}

func (h *PaymentHandler) GetDashboardV1Payments(w http.ResponseWriter, r *http.Request, params openapigen.GetDashboardV1PaymentsParams) {
	// Validate JWT
	if err := h.validateToken(r); err != nil {
		transport.WriteAppError(w, entity.WrapError(err, entity.ErrorCodeUnauthorized, "Unauthenticated: missing or invalid token"))
		return
	}

	filter := entity.PaymentFilter{
		Status: params.Status,
		ID:     params.Id,
		Sort:   params.Sort,
	}

	payments, err := h.paymentUC.ListPayments(filter)
	if err != nil {
		transport.WriteError(w, err)
		return
	}

	// Convert to openapigen type
	var result []openapigen.Payment
	for _, p := range payments {
		pCopy := p
		id := pCopy.ID
		merchant := pCopy.Merchant
		status := pCopy.Status
		amount := pCopy.Amount
		createdAt := pCopy.CreatedAt
		result = append(result, openapigen.Payment{
			Id:        &id,
			Merchant:  &merchant,
			Status:    &status,
			Amount:    &amount,
			CreatedAt: &createdAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(openapigen.PaymentListResponse{Payments: &result})
}

func (h *PaymentHandler) validateToken(r *http.Request) error {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return entity.ErrorUnauthorized("missing token")
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	_, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, entity.ErrorUnauthorized("invalid signing method")
		}
		return h.jwtSecret, nil
	})
	return err
}
