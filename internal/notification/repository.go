package notification

import (
	"context"
	"socialnetwork/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type NotificationRepository interface {
	GetByRecipient(ctx context.Context, userID primitive.ObjectID) ([]models.Notification, error)
	Create(ctx context.Context, notif *models.Notification) error
	MarkAsRead(ctx context.Context, notifID primitive.ObjectID) error
}

type notificationRepository struct {
	Collection *mongo.Collection
}

func NewNotificationRepository(db *mongo.Database) NotificationRepository {
	return &notificationRepository{
		Collection: db.Collection("notifications"),
	}
}

func (r *notificationRepository) GetByRecipient(ctx context.Context, userID primitive.ObjectID) ([]models.Notification, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"recipient": userID}
	cursor, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var notifications []models.Notification
	if err := cursor.All(ctx, &notifications); err != nil {
		return nil, err
	}
	return notifications, nil
}

func (r *notificationRepository) Create(ctx context.Context, notif *models.Notification) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.Collection.InsertOne(ctx, notif)
	return err
}

func (r *notificationRepository) MarkAsRead(ctx context.Context, notifID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": notifID}
	update := bson.M{"$set": bson.M{"isRead": true}}
	_, err := r.Collection.UpdateOne(ctx, filter, update)
	return err
}