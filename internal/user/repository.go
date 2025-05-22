package user

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"socialnetwork/models"
	"time"
)

type Repository interface {
	Create(ctx context.Context, user *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id string) (*models.User, error)
	GetAllUsers(ctx context.Context) ([]*models.User, error)
	FindAll(ctx context.Context) ([]*models.User, error)
	UpdateByID(ctx context.Context, id string, update bson.M) error
	SendFriendRequest(ctx context.Context, fromID, toID primitive.ObjectID) error
	AcceptFriendRequest(ctx context.Context, userID, requesterID primitive.ObjectID) error
	BlockUser(ctx context.Context, userID, targetID primitive.ObjectID) error
	ToggleHideProfile(ctx context.Context, userID primitive.ObjectID, hide bool) error
	FriendRequestExists(ctx context.Context, fromID, toID primitive.ObjectID) (bool, error)
	CancelFriendRequest(ctx context.Context, fromID, toID primitive.ObjectID) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*models.User, error)
	IncrementFollowerCount(ctx context.Context, userID primitive.ObjectID) error
	DecrementFollowerCount(ctx context.Context, userID primitive.ObjectID) error
	IncrementFollowingCount(ctx context.Context, userID primitive.ObjectID) error
	DecrementFollowingCount(ctx context.Context, userID primitive.ObjectID) error
}

func (r *repository) FindByID(ctx context.Context, id string) (*models.User, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user models.User
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

type repository struct {
	collection *mongo.Collection
	db         *mongo.Database
}

func NewRepository(db *mongo.Database) Repository {
	return &repository{

		collection: db.Collection("users"),
		db:         db,
	}
}

func (r *repository) Create(ctx context.Context, user *models.User) error {
	user.CreatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	return &user, err
}

func (r *repository) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*models.User
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func (r *repository) FindAll(ctx context.Context) ([]*models.User, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*models.User
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *repository) UpdateByID(ctx context.Context, id string, update bson.M) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("Không tìm thấy user với ID: %s", id)
	}
	return nil
}

func (r *repository) SendFriendRequest(ctx context.Context, fromID, toID primitive.ObjectID) error {
	// addToSet trên User[toID].FriendRequests
	_, err := r.collection.UpdateOne(ctx,
		bson.M{"_id": toID},
		bson.M{"$addToSet": bson.M{"friendRequests": fromID}},
	)
	return err
}

func (r *repository) AcceptFriendRequest(ctx context.Context, userID, requesterID primitive.ObjectID) error {
	// Xóa lời mời kết bạn
	_, err := r.db.Collection("friend_requests").DeleteOne(ctx, bson.M{
		"from": requesterID,
		"to":   userID,
	})
	if err != nil {
		return err
	}

	// Thêm vào danh sách bạn bè (cả 2 chiều)
	usersCol := r.db.Collection("users")

	// Thêm requester vào danh sách bạn của user
	_, err = usersCol.UpdateOne(ctx, bson.M{"_id": userID}, bson.M{
		"$addToSet": bson.M{"friends": requesterID},
	})
	if err != nil {
		return err
	}

	// Thêm user vào danh sách bạn của requester
	_, err = usersCol.UpdateOne(ctx, bson.M{"_id": requesterID}, bson.M{
		"$addToSet": bson.M{"friends": userID},
	})
	return err
}

func (r *repository) BlockUser(ctx context.Context, userID, targetID primitive.ObjectID) error {
	// addToSet vào blockedUsers
	_, err := r.collection.UpdateOne(ctx,
		bson.M{"_id": userID},
		bson.M{"$addToSet": bson.M{"blockedUsers": targetID}},
	)
	return err
}

func (r *repository) ToggleHideProfile(ctx context.Context, userID primitive.ObjectID, hide bool) error {
	_, err := r.collection.UpdateOne(ctx,
		bson.M{"_id": userID},
		bson.M{"$set": bson.M{"hideProfile": hide}},
	)
	return err
}

func (r *repository) FriendRequestExists(ctx context.Context, fromID, toID primitive.ObjectID) (bool, error) {
	filter := bson.M{"from": fromID, "to": toID}
	count, err := r.db.Collection("friend_requests").CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *repository) CancelFriendRequest(ctx context.Context, fromID, toID primitive.ObjectID) error {
	_, err := r.db.Collection("friend_requests").DeleteOne(ctx, bson.M{
		"from": fromID,
		"to":   toID,
	})
	return err
}

func (r *repository) IncrementFollowerCount(ctx context.Context, userID primitive.ObjectID) error {
	_, err := r.collection.UpdateOne(ctx,
		bson.M{"_id": userID},
		bson.M{"$inc": bson.M{"followerCount": 1}},
	)
	return err
}
func (r *repository) DecrementFollowerCount(ctx context.Context, userID primitive.ObjectID) error {
	_, err := r.collection.UpdateOne(ctx,
		bson.M{"_id": userID},
		bson.M{"$inc": bson.M{"followerCount": -1}},
	)
	return err
}
func (r *repository) IncrementFollowingCount(ctx context.Context, userID primitive.ObjectID) error {
	_, err := r.collection.UpdateOne(ctx,
		bson.M{"_id": userID},
		bson.M{"$inc": bson.M{"followingCount": 1}},
	)
	return err
}
func (r *repository) DecrementFollowingCount(ctx context.Context, userID primitive.ObjectID) error {
	_, err := r.collection.UpdateOne(ctx,
		bson.M{"_id": userID},
		bson.M{"$inc": bson.M{"followingCount": -1}},
	)
	return err
}
func (r *repository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}