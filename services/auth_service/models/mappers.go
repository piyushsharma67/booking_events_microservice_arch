package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func MapCreateRequestToDocument(
	req *CreateUserRequest,
) (*UserDocument, error) {

	now := time.Now()

	return &UserDocument{
		ID:           primitive.NewObjectID(),
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: req.PasswordHash,
		Role:         req.Role,
		CreatedAt:    now,
		UpdatedAt:    now,
		IsActive:     true,
	}, nil
}
