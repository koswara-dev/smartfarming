package config

import (
	"log"
	"smartfarming/model"

	"golang.org/x/crypto/bcrypt"
)

func SeedDatabase() {
	users := []struct {
		Name     string
		Email    string
		Password string
		Role     string
	}{
		{
			Name:     "Admin SmartFarming",
			Email:    "admin@smartfarming.id",
			Password: "Indonesia",
			Role:     "admin",
		},
		{
			Name:     "Operator SmartFarming",
			Email:    "operator@smartfarming.id",
			Password: "Indonesia",
			Role:     "operator",
		},
		{
			Name:     "User SmartFarming",
			Email:    "user@smartfarming.id",
			Password: "Indonesia",
			Role:     "user",
		},
	}

	for _, u := range users {
		var count int64
		err := DB.Model(&model.User{}).Where("email = ?", u.Email).Count(&count).Error
		if err != nil {
			log.Printf("Failed to check if user %s exists: %v", u.Email, err)
			continue
		}

		if count == 0 {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
			if err != nil {
				log.Printf("Failed to hash password for %s: %v", u.Email, err)
				continue
			}

			user := model.User{
				Name:     u.Name,
				Email:    u.Email,
				Password: string(hashedPassword),
				Role:     u.Role,
			}

			// We write directly using the DB handle since it runs at server start,
			// audit fields created_by, updated_by will be nil for seeded users.
			err = DB.Create(&user).Error
			if err != nil {
				log.Printf("Failed to seed user %s: %v", u.Email, err)
			} else {
				log.Printf("Successfully seeded user: %s (%s)", u.Email, u.Role)
			}
		}
	}
}
