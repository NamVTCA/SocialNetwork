package comment

import (
    "context"
    "time"

    "socialnetwork/models"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
    Create(ctx context.Context, c *models.Comment) (*models.Comment, error)
    GetByID(ctx context.Context, id primitive.ObjectID) (*models.Comment, error)
    Update(ctx context.Context, id primitive.ObjectID, data map[string]interface{}) error
    Delete(ctx context.Context, id primitive.ObjectID) error
    ListByPost(ctx context.Context, postID primitive.ObjectID) ([]*models.Comment, error)
}

type commentRepo struct {
    collection *mongo.Collection
}

func NewCommentRepository(db *mongo.Database) Repository {
    return &commentRepo{
        collection: db.Collection("comments"),
    }
}

func (r *commentRepo) Create(ctx context.Context, c *models.Comment) (*models.Comment, error) {
    c.ID = primitive.NewObjectID()
    c.CreatedAt = time.Now()
    c.UpdatedAt = c.CreatedAt

    _, err := r.collection.InsertOne(ctx, c)
    return c, err
}

func (r *commentRepo) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Comment, error) {
    var comment models.Comment
    err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&comment)
    return &comment, err
}

func (r *commentRepo) Update(ctx context.Context, id primitive.ObjectID, data map[string]interface{}) error {
    data["updated_at"] = time.Now()
    _, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": data})
    return err
}

func (r *commentRepo) Delete(ctx context.Context, id primitive.ObjectID) error {
    _, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
    return err
}

func (r *commentRepo) ListByPost(ctx context.Context, postID primitive.ObjectID) ([]*models.Comment, error) {
    cursor, err := r.collection.Find(ctx, bson.M{"post_id": postID})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var comments []*models.Comment
    for cursor.Next(ctx) {
        var comment models.Comment
        if err := cursor.Decode(&comment); err != nil {
            return nil, err
        }
        comments = append(comments, &comment)
    }
    return comments, nil
}
