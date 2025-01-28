package validate

import (
	"chat-bots-api/domain"
	"errors"
	"fmt"
	"slices"
	"strings"
)

func UserID(id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid owner id provided %v", id)
	}
	return nil
}

func UserEmail(email string) error {
	if len(email) == 0 {
		return fmt.Errorf("empty email provided")
	}
	if !strings.Contains(email, "@") {
		return fmt.Errorf("invalid email format")
	}
	if !strings.Contains(email, ".") {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

func Plan(plan string) error {
	validPlans := []string{
		domain.FreePlan,
		domain.BusinessPlan,
		domain.EnterprisePlan,
	}
	if !slices.Contains(validPlans, plan) {
		return fmt.Errorf("invalid user plan provided")
	}
	return nil
}

func SaveUser(email, password, plan string) error {
	if err := Plan(plan); err != nil {
		return err
	}
	if err := UserEmail(email); err != nil {
		return err
	}
	if password == "" {
		return errors.New("empty password provided")
	}
	if len(password) < 8 {
		return errors.New("password is too short")
	}
	return nil
}
