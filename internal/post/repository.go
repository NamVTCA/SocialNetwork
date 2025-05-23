package post

import (
    "context"
    "errors"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "socialnetwork/models"
)

type PostRepository struct {
    collection *mongo.Collection
}

func NewPostRepository(db *mongo.Database) PostRepository {
    return PostRepository{
        collection: db.Collection("posts"),
    }
}

// Tạo post mới
func (r *PostRepository) Create(ctx context.Context, post *models.Post) error {
    now := time.Now()
    post.ID = primitive.NewObjectID()
    post.CreatedAt = now
    post.UpdatedAt = now

    _, err := r.collection.InsertOne(ctx, post)
    return err
}

// Lấy post theo ID
func (r *PostRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Post, error) {
    var post models.Post
    err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&post)
    if err == mongo.ErrNoDocuments {
        return nil, errors.New("post không tồn tại")
    }
    return &post, err
}

// Cập nhật post (chỉ sửa content và media, image_url)
func (r *PostRepository) Update(ctx context.Context, id primitive.ObjectID, updateData bson.M) error {
    updateData["updated_at"] = time.Now()
    update := bson.M{"$set": updateData}
    res, err := r.collection.UpdateByID(ctx, id, update)
    if err != nil {
        return err
    }
    if res.MatchedCount == 0 {
        return errors.New("post không tồn tại")
    }
    return nil
}

// Xóa post
func (r *PostRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
    res, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
    if err != nil {
        return err
    }
    if res.DeletedCount == 0 {
        return errors.New("post không tồn tại")
    }
    return nil
}

// Lấy danh sách post có phân trang, sắp xếp theo created_at giảm dần
func (r *PostRepository) List(ctx context.Context, page, limit int64) ([]models.Post, error) {
    opts := options.Find()
    opts.SetSort(bson.M{"created_at": -1})
    opts.SetSkip((page - 1) * limit)
    opts.SetLimit(limit)

    cursor, err := r.collection.Find(ctx, bson.M{}, opts)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var posts []models.Post
    for cursor.Next(ctx) {
        var post models.Post
        if err := cursor.Decode(&post); err != nil {
            return nil, err
        }
        posts = append(posts, post)
    }
    return posts, nil
}

func (r *PostRepository) FindByOwnerAndVisibility(ctx context.Context, ownerID primitive.ObjectID, visibility string) ([]models.Post, error) {
    filter := bson.M{
        "user_id":    ownerID,
        "visibility": visibility,
    }
    cursor, err := r.collection.Find(ctx, filter)
    if err != nil {
        return nil, err
    }
    var posts []models.Post
    if err = cursor.All(ctx, &posts); err != nil {
        return nil, err
    }
    return posts, nil
}
