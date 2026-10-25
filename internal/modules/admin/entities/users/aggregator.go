package users

import (
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/users/action/create_user"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/users/action/delete_user"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/users/action/ensure_initial_user"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/users/action/get_user"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/users/action/list_users"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/users/action/update_user"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/users/dal"
	"github.com/nagorneva/nagorneva_reviews/internal/pkg/postgres"
)

type Aggregator struct {
	CreateUser        *create_user.Action
	GetUser           *get_user.Action
	ListUsers         *list_users.Action
	UpdateUser        *update_user.Action
	DeleteUser        *delete_user.Action
	EnsureInitialUser *ensure_initial_user.Action
}

func New(db postgres.DB) *Aggregator {
	repository := dal.New(db)
	return &Aggregator{
		CreateUser:        create_user.New(repository),
		GetUser:           get_user.New(repository),
		ListUsers:         list_users.New(repository),
		UpdateUser:        update_user.New(repository),
		DeleteUser:        delete_user.New(repository),
		EnsureInitialUser: ensure_initial_user.New(repository),
	}
}
