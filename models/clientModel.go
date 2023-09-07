package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Client struct {
	Id          primitive.ObjectID `bson:"_id"`
	CompanyName *string            `json:"name,omitempty" validate:"required,min=2,max=16"`
	RedirectUri *string            `json:"redirectUri,omitempty" validate:"required,url"`
	UserId      primitive.ObjectID `bson:"_id"`
	CreatedAt   time.Time          `json:"createdTime"`
	UpdatedAt   time.Time          `json:"updatedTime"`
}

type CreateClientRequest struct {
	Id          primitive.ObjectID `bson:"_id"`
	CompanyName *string            `form:"name"`
	RedirectUri *string            `form:"redirectUri"`
}
