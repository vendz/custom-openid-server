package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id           primitive.ObjectID `bson:"_id"`
	Name         *string            `json:"Name,omitempty" validate:"required,min=2,max=16"`
	Email        *string            `json:"email,omitempty" validate:"required,email"`
	Password     *string            `json:"-" validate:"required,min=8,max=64"`
	PhoneNumber  *string            `json:"phoneNumber,omitempty" validate:"required,e164,min=10,max=13"`
	Token        *string            `json:"token,omitempty"`
	AccessTokens map[string]string  `json:"accessTokens,omitempty"`
	Bio          *string            `json:"bio,omitempty"`
	IsClient     bool               `json:"isClient,omitempty"`
	CreatedAt    time.Time          `json:"createdTime"`
	UpdatedAt    time.Time          `json:"updatedTime"`
}

type CreateUserRequest struct {
	Id          primitive.ObjectID `bson:"_id"`
	Name        *string            `form:"Name"`
	Email       *string            `form:"email"`
	Password    *string            `form:"password"`
	PhoneNumber *string            `form:"phoneNumber"`
	Bio         *string            `form:"bio,omitempty"`
}

type LoginUserRequest struct {
	Email    *string `json:"email" validate:"required,email"`
	Password *string `json:"password" validate:"required,min=8,max=64"`
}
