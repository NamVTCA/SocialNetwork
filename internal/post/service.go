package post


import (
    "context"
    "socialnetwork/models"


    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type PostService interface {
    CreatePost(ctx context.Context, post *models.Post) (*models.Post, error)
    GetPostByID(ctx context.Context, id primitive.ObjectID) (*models.Post, error)
    UpdatePost(ctx context.Context, id primitive.ObjectID, updateData bson.M) error
    DeletePost(ctx context.Context, id primitive.ObjectID) error
    ListPosts(ctx context.Context, page, limit int64) ([]models.Post, error)
}

type postService struct {
    repo *PostRepository
}

func NewPostService(repo *PostRepository) PostService {
    return &postService{repo: repo}
}

func (s *postService) CreatePost(ctx context.Context, post *models.Post) (*models.Post, error) {
    err := s.repo.Create(ctx, post)
    return post, err
}

func (s *postService) GetPostByID(ctx context.Context, id primitive.ObjectID) (*models.Post, error) {
    return s.repo.GetByID(ctx, id)
}

func (s *postService) UpdatePost(ctx context.Context, id primitive.ObjectID, updateData bson.M) error {
    return s.repo.Update(ctx, id, updateData)
}

func (s *postService) DeletePost(ctx context.Context, id primitive.ObjectID) error {
    return s.repo.Delete(ctx, id)
}

func (s *postService) ListPosts(ctx context.Context, page, limit int64) ([]models.Post, error) {
    return s.repo.List(ctx, page, limit)
}