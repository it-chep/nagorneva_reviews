package ensure_initial_user

import (
	"context"
	"errors"
	"testing"

	"github.com/nagorneva/nagorneva_reviews/internal/pkg/password"
)

type storage struct {
	exists bool
	err    error
	calls  int
	hash   string
}

func (s *storage) ExistsByEmail(context.Context, string) (bool, error) { return s.exists, s.err }
func (s *storage) CreateInitialUser(_ context.Context, _, hash string) error {
	s.calls++
	s.hash = hash
	return nil
}
func TestBootstrapDoesNotOverwriteExistingUser(t *testing.T) {
	for _, tc := range []struct {
		name, email, raw string
		exists           bool
		err              error
		created          bool
	}{
		{name: "disabled"},
		{name: "existing", email: "admin@example.test", raw: "password", exists: true},
		{name: "database failure", email: "admin@example.test", raw: "password", err: errors.New("database unavailable")},
		{name: "new", email: "admin@example.test", raw: "password", created: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := &storage{exists: tc.exists, err: tc.err}
			err := New(s).Execute(context.Background(), tc.email, tc.raw)
			if !errors.Is(err, tc.err) {
				t.Fatal(err)
			}
			if tc.created {
				if s.calls != 1 || password.Check(s.hash, tc.raw) != nil {
					t.Fatal("bootstrap did not persist bcrypt hash")
				}
			} else if s.calls != 0 {
				t.Fatal("unexpected bootstrap insert")
			}
		})
	}
}
