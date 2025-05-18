package comment

import (
    "context"

    "socialnetwork/models"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
    Create(ctx context.Context, c *models.Comment) (*models.Comment, error)
    GetByID(ctx context.Context, id primitive.ObjectID) (*models.Comment, error)
    Update(ctx context.Context, id, userID primitive.ObjectID, data map[string]interface{}) error
    Delete(ctx context.Context, id, userID primitive.ObjectID) error
    ListByPost(ctx context.Context, postID primitive.ObjectID) ([]*models.Comment, error)
	ToggleLike(ctx context.Context, commentID, userID primitive.ObjectID) error
}

type commentService struct {
    repo Repository
}

func NewCommentService(repo Repository) Service {
    return &commentService{repo: repo}
}

func (s *commentService) Create(ctx context.Context, c *models.Comment) (*models.Comment, error) {
    return s.repo.Create(ctx, c)
}

func (s *commentService) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Comment, error) {
    return s.repo.GetByID(ctx, id)
}

func (s *commentService) Update(ctx context.Context, id, userID primitive.ObjectID, data map[string]interface{}) error {
    c, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return err
    }
    if c.UserID != userID {
        return ErrUnauthorized
    }
    return s.repo.Update(ctx, id, userID, data)
}

func (s *commentService) Delete(ctx context.Context, id, userID primitive.ObjectID) error {
    c, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return err
    }
    if c.UserID != userID {
        return ErrUnauthorized
    }
    return s.repo.Delete(ctx, id, userID)
}

func (s *commentService) ListByPost(ctx context.Context, postID primitive.ObjectID) ([]*models.Comment, error) {
    return s.repo.ListByPost(ctx, postID)
}


func (s *commentService) ToggleLike(ctx context.Context, commentID, userID primitive.ObjectID) error {
    return s.repo.ToggleLike(ctx, commentID, userID)
}