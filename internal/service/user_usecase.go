package service

import (
	"context"
	fs "github.com/core-go/firestore"

	. "go-service/internal/model"
)

func NewUserService(repository fs.Repository) *UserUsecase {
	return &UserUsecase{repository: repository}
}

type UserUsecase struct {
	repository fs.Repository
}

func (s *UserUsecase) Load(ctx context.Context, id string) (*User, error) {
	var user User
	ok, err := s.repository.Get(ctx, id, &user)
	if !ok {
		return nil, err
	} else {
		return &user, err
	}
}
func (s *UserUsecase) Create(ctx context.Context, user *User) (int64, error) {
	return s.repository.Insert(ctx, user)
}
func (s *UserUsecase) Update(ctx context.Context, user *User) (int64, error) {
	return s.repository.Update(ctx, user)
}
func (s *UserUsecase) Patch(ctx context.Context, user map[string]interface{}) (int64, error) {
	return s.repository.Patch(ctx, user)
}
func (s *UserUsecase) Delete(ctx context.Context, id string) (int64, error) {
	return s.repository.Delete(ctx, id)
}
