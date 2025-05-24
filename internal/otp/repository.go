package otp

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"socialnetwork/models"
)

// OTPrepository định nghĩa các method cần thiết
type OTPrepository interface {
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByPhone(ctx context.Context, phone string) (*models.User, error)
	UpdateByID(ctx context.Context, id string, update bson.M) error
}

// otpRepository là struct implement OTPrepository
type otpRepository struct {
	collection *mongo.Collection
}

// NewOTPRepository khởi tạo otpRepository với collection MongoDB
func NewOTPRepository(collection *mongo.Collection) OTPrepository {
	return &otpRepository{collection: collection}
}

// FindByEmail tìm user theo email
func (r *otpRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // Không tìm thấy user
		}
		return nil, err
	}
	return &user, nil
}

// FindByPhone tìm user theo số điện thoại
func (r *otpRepository) FindByPhone(ctx context.Context, phone string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"phone": phone}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // Không tìm thấy user
		}
		return nil, err
	}
	return &user, nil
}

// UpdateByID cập nhật thông tin user theo ID
func (r *otpRepository) UpdateByID(ctx context.Context, id string, update bson.M) error {
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
