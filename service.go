package main

import (
	"fmt"

	"github.com/thomaspoignant/go-feature-flag/ffuser"
)

type Payment struct {
	OrderID string  `json:"order_id"`
	Amount  float64 `json:"amount"`
	Notes   string  `json:"notes"`
	Version string  `json:"version,omitempty"`
}

type service struct {
	fg FFallger
}

type Servicer interface {
	Pay(payload Payment) error

	Get(userID string, orderID string) (*Payment, error)

	Generate() (map[string]string, error)
}

type FFallger interface {
	BoolVariation(flagKey string, user ffuser.User, defaultValue bool) (bool, error)
	StringVariation(flagKey string, user ffuser.User, defaultValue string) (string, error)
}

func NewService(fg FFallger) Servicer {
	return &service{
		fg: fg,
	}
}

func (s *service) Pay(payload Payment) error {
	return nil
}

func (s *service) Get(userID string, orderID string) (*Payment, error) {
	payment := &Payment{}

	fgStatus, err := s.fg.BoolVariation(flagKeyFeaturePayment, ffuser.NewUser(userID), false)
	if err != nil {
		return nil, fmt.Errorf("feature flagging error: %w", err)
	}

	if fgStatus {
		payment.Version = "new version with new calculation"
	} else {
		payment.Version = "old version"
	}

	return payment, nil
}

func (s *service) Generate() (map[string]string, error) {
	generate := make(map[string]string, generateCount)

	for i := 0; i < generateCount; i++ {
		key := fmt.Sprintf("user%d", i)

		result, err := s.fg.StringVariation(flagKeyGenerate, ffuser.NewUser(key), "default-arg")
		if err != nil {
			return nil, fmt.Errorf("something wrong the flag: %w", err)
		}

		generate[key] = result
	}

	return generate, nil
}
