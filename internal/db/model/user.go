package model

import (
	"github.com/driif/echo-go-starter/internal/api/dto"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `bson:"_id"`
	Email     string    `bson:"email"`
	LastName  string    `bson:"last_name"`
	FirstName string    `bson:"first_name"`
	ShortName string    `bson:"short_name"`
	Active    bool      `bson:"active"`
	Avatar    string    `bson:"avatar"`
	UpdatedAt string    `bson:"updated_at"`
	CreatedAt string    `bson:"created_at"`
}

func (u *User) FromDTO(dto *dto.User) {
	u.ID = dto.ID
	u.Email = dto.Email
	u.LastName = dto.LastName
	u.FirstName = dto.FirstName
	u.ShortName = dto.ShortName
	u.Active = dto.Active
	u.Avatar = dto.Avatar
	u.UpdatedAt = dto.UpdatedAt
	u.CreatedAt = dto.CreatedAt
}

func (u *User) ToDTO() *dto.User {
	return &dto.User{
		ID:        u.ID,
		Email:     u.Email,
		LastName:  u.LastName,
		FirstName: u.FirstName,
		ShortName: u.ShortName,
		Active:    u.Active,
		Avatar:    u.Avatar,
		UpdatedAt: u.UpdatedAt,
		CreatedAt: u.CreatedAt,
	}
}

func (u User) GetCollectionName() string {
	return "users"
}
