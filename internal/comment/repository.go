package comment

import (
    "context"
    "time"
    "errors"
    "socialnetwork/models"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
    Create(ctx context.Context, c *models.Comment) (*models.Comment, error)
    GetByID(ctx context.Context, id primitive.ObjectID) (*models.Comment, error)
    Update(ctx context.Context, id, userID primitive.ObjectID, data map[string]interface{}) error
    Delete(ctx context.Context, id, userID primitive.ObjectID) error
    ListByPost(ctx context.Context, postID primitive.ObjectID) ([]*models.Comment, error)
    ToggleLike(ctx context.Context, commentID, userID primitive.ObjectID) error
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

func (r *commentRepo) Update(ctx context.Context, id, userID primitive.ObjectID, data map[string]interface{}) error {
    data["updated_at"] = time.Now()
    res, err := r.collection.UpdateOne(
        ctx,
        bson.M{"_id": id, "user_id": userID}, // Check đúng người
        bson.M{"$set": data},
    )
    if err != nil {
        return err
    }
    if res.MatchedCount == 0 {
        return errors.New("unauthorized or comment not found")
    }
    return nil
}


func (r *commentRepo) Delete(ctx context.Context, id, userID primitive.ObjectID) error {
    res, err := r.collection.DeleteOne(ctx, bson.M{"_id": id, "user_id": userID})
    if err != nil {
        return err
    }
    if res.DeletedCount == 0 {
        return errors.New("unauthorized or comment not found")
    }
    return nil
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

func (r *commentRepo) ToggleLike(ctx context.Context, commentID, userID primitive.ObjectID) error {
    var comment models.Comment
    err := r.collection.FindOne(ctx, bson.M{"_id": commentID}).Decode(&comment)
    if err != nil {
        return err
    }

    // Kiểm tra user đã like chưa
    liked := false
    for _, uid := range comment.Likes {
        if uid == userID {
            liked = true
            break
        }
    }

    var update bson.M
    if liked {
        update = bson.M{"$pull": bson.M{"likes": userID}}
    } else {
        update = bson.M{"$addToSet": bson.M{"likes": userID}}
    }

    _, err = r.collection.UpdateOne(ctx, bson.M{"_id": commentID}, update)
    return err
}


