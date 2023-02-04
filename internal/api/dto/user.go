package dto

import (
	"errors"

	"github.com/driif/echo-go-starter/internal/server/net/runtime/valid"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id,omitempty" bson:"id"`
	Email     string    `json:"email,omitempty" bson:"email"`
	LastName  string    `json:"last_name,omitempty" bson:"last_name"`
	FirstName string    `json:"first_name,omitempty" bson:"first_name"`
	ShortName string    `json:"short_name,omitempty" bson:"short_name"`
	Active    bool      `json:"active,omitempty" bson:"active"`
	Avatar    string    `json:"avatar,omitempty" bson:"avatar"`
	UpdatedAt string    `json:"updated_at,omitempty" bson:"updated_at"`
	CreatedAt string    `json:"created_at,omitempty" bson:"created_at"`
}

func (u *User) Validate() error {

	err := valid.Email(u.Email)
	if err != nil {
		return err
	}

	if u.LastName == "" {
		return errors.New("last name is required")
	}

	return nil
}
