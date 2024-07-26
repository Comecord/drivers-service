package services

import (
	"context"
	"drivers-service/config"
	"drivers-service/pkg/logging"
	"go.mongodb.org/mongo-driver/mongo"
)

type BaseCollection struct {
	string
}

type BaseService[T any, Tc any, Tu any, Tr any] struct {
	Mongo      *mongo.Database
	Collection *mongo.Collection
	Logger     logging.Logger
	ctx        context.Context
}

func NewBaseService[T any, Tc any, Tu any, Tr any](db *mongo.Database, cfg *config.Config, ctx context.Context, collectionName string) *BaseService[T, Tc, Tu, Tr] {
	baseService := &BaseService[T, Tc, Tu, Tr]{
		Mongo:  db,
		Logger: logging.NewLogger(cfg),
		ctx:    ctx,
	}
	baseService.Collection = db.Collection(collectionName)

	return baseService
}
