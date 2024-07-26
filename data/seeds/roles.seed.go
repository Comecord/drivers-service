package seeds

import (
	"context"
	"drivers-service/config"
	"drivers-service/data/models"
	"drivers-service/pkg/logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func SeedRoles(db *mongo.Database, ctx context.Context) {
	rolesCollection := db.Collection("roles")
	logger := logging.NewLogger(config.GetConfig())
	count, err := rolesCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	if count == 0 {
		roles := []models.MemberRole{
			{ID: "b9fe0246-4b69-4be6-a305-8e3c4f371382", Name: "admin", Permissions: []string{"create", "read", "update", "delete"}},
			{ID: "207f76c2-39f6-4751-bc85-d145b79a5f15", Name: "member", Permissions: []string{"create", "read", "update", "delete"}},
			{ID: "d0357b22-0324-4710-9552-b8889239498f", Name: "guest", Permissions: []string{"create", "read", "update", "delete"}},
		}

		// Преобразование roles в []interface{}
		var interfaces []interface{}
		for _, role := range roles {
			interfaces = append(interfaces, role)
		}

		_, err := rolesCollection.InsertMany(ctx, interfaces)
		if err != nil {
			logger.Error(logging.MongoDB, logging.Seed, err.Error(), nil)
		}

		logger.Info(logging.MongoDB, logging.Seed, "Roles seeded successfully.", nil)
	} else {
		logger.Info(logging.MongoDB, logging.Seed, "Skipping seeding..", nil)

	}
}
