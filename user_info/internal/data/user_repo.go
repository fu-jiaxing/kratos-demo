package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

func NewUserRepo(data *Data, logger log.Logger) UserRepo {
	return &userRepo{data: data, log: log.NewHelper(logger)}
}

type User struct {
	ID      int64
	Name    string
	Email   string
	Address string
	Tel     string

	Visits int64
}

func (u User) TableName() string {
	return "user_info"
}

type UserRepo interface {
	FindByID(ctx context.Context, id int64) (*User, error)
	Update(ctx context.Context, u *User) (*User, error)
}
type userRepo struct {
	data *Data
	log  *log.Helper
}

func (r *userRepo) FindByID(ctx context.Context, id int64) (*User, error) {
	var user User
	res := r.data.db.Select("name", "email", "address", "tel").Where("id=?", id).First(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	return &user, nil
}

func (r *userRepo) Update(ctx context.Context, user *User) (*User, error) {
	return nil, nil
}
