package update_user

import (
	"context"
	"errors"
	"testing"

	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/users/dto"
	validation "github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/users/service/password"
	"github.com/nagorneva/nagorneva_reviews/internal/pkg/password"
)

type storage struct {
	input dto.UpdateInput
	calls int
}

func (s *storage) UpdateUser(_ context.Context, id int64, input dto.UpdateInput) (*dto.User, error) {
	s.input = input
	s.calls++
	return &dto.User{ID: id}, nil
}
func TestPasswordPatch(t *testing.T) {
	s := &storage{}
	action := New(s)
	email := "changed@example.test"
	if _, err := action.Execute(context.Background(), 1, &email, nil); err != nil {
		t.Fatal(err)
	}
	if s.input.Password != nil || s.input.Email != &email {
		t.Fatal("omitted password must not be changed")
	}
	short := "short"
	if _, err := action.Execute(context.Background(), 1, nil, &short); !errors.Is(err, validation.ErrPasswordTooShort) {
		t.Fatal(err)
	}
	if s.calls != 1 {
		t.Fatal("invalid password reached repository")
	}
	raw := "long-password"
	if _, err := action.Execute(context.Background(), 1, nil, &raw); err != nil {
		t.Fatal(err)
	}
	if s.input.Password == nil || *s.input.Password == raw || password.Check(*s.input.Password, raw) != nil {
		t.Fatal("password was not bcrypt-hashed")
	}
}
