package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/durianpay/fullstack-boilerplate/internal/api"
	"github.com/durianpay/fullstack-boilerplate/internal/config"
	"github.com/durianpay/fullstack-boilerplate/internal/entity"
	ah "github.com/durianpay/fullstack-boilerplate/internal/module/auth/handler"
	ar "github.com/durianpay/fullstack-boilerplate/internal/module/auth/repository"
	au "github.com/durianpay/fullstack-boilerplate/internal/module/auth/usecase"
	ph "github.com/durianpay/fullstack-boilerplate/internal/module/payment/handler"
	pr "github.com/durianpay/fullstack-boilerplate/internal/module/payment/repository"
	pu "github.com/durianpay/fullstack-boilerplate/internal/module/payment/usecase"
	srv "github.com/durianpay/fullstack-boilerplate/internal/service/http"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	_ = godotenv.Load()

	userRepo := ar.NewInMemoryUserRepo()
	seedUsers(userRepo)

	jwtExpired, err := time.ParseDuration(config.JwtExpired)
	if err != nil {
		log.Fatal("invalid JWT_EXPIRED:", err)
	}
	authUC := au.NewAuthUsecase(userRepo, config.JwtSecret, jwtExpired)
	authH := ah.NewAuthHandler(authUC)

	paymentRepo := pr.NewInMemoryPaymentRepo()
	paymentRepo.Seed(seedPayments())

	paymentUC := pu.NewPaymentUsecase(paymentRepo)
	paymentH := ph.NewPaymentHandler(paymentUC, config.JwtSecret)

	// --- Wire & serve ---
	apiHandler := &api.APIHandler{
		Auth:    authH,
		Payment: paymentH,
	}

	server := srv.NewServer(apiHandler, config.OpenapiYamlLocation)
	log.Printf("🚀 Server starting on %s (in-memory store)", config.HttpAddress)
	server.Start(config.HttpAddress)
}

func seedUsers(repo *ar.InMemoryUserRepo) {
	users := []struct {
		email    string
		password string
		role     string
	}{
		{"cs@test.com", "password", "cs"},
		{"operation@test.com", "password", "operation"},
	}
	for _, u := range users {
		hash, err := bcrypt.GenerateFromPassword([]byte(u.password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal("bcrypt error:", err)
		}
		repo.AddUser(&entity.User{
			ID:           u.email,
			Email:        u.email,
			PasswordHash: string(hash),
			Role:         u.role,
		})
	}
	log.Println("✓ seeded 2 users")
}

func seedPayments() []entity.Payment {
	merchants := []string{
		"Tokopedia", "Shopee", "Bukalapak", "Lazada", "Blibli",
		"GoJek", "Grab", "OVO", "Dana", "LinkAja",
		"Traveloka", "Tiket.com", "Pegipegi", "RedDoorz", "OYO",
	}
	statuses := []string{"completed", "processing", "failed"}
	weights := []int{6, 2, 2}

	rng := rand.New(rand.NewSource(42))
	pickStatus := func() string {
		n := rng.Intn(10)
		cum := 0
		for i, w := range weights {
			cum += w
			if n < cum {
				return statuses[i]
			}
		}
		return statuses[0]
	}

	now := time.Now()
	payments := make([]entity.Payment, 50)
	for i := range payments {
		payments[i] = entity.Payment{
			ID:        fmt.Sprintf("PAY-%05d", i+1),
			Merchant:  merchants[rng.Intn(len(merchants))],
			Status:    pickStatus(),
			Amount:    fmt.Sprintf("%.2f", float64(rng.Intn(9900000)+100000)/100.0),
			CreatedAt: now.Add(-time.Duration(rng.Intn(30*24)) * time.Hour),
		}
	}
	log.Printf("✓ seeded %d payments", len(payments))
	return payments
}
