package login

import (
	"context"
	"errors"
	"testing"

	"github.com/nagorneva/nagorneva_reviews/internal/pkg/domain_error"
	"github.com/nagorneva/nagorneva_reviews/internal/pkg/password"
)

type storage struct {
	hash string
	err  error
}

func (s storage) FindByEmail(context.Context, string) (int64, string, error) { return 7, s.hash, s.err }

type issuer struct{ issued int64 }

func (i *issuer) Issue(id int64) (string, error) { i.issued = id; return "token", nil }
func TestLoginErrorsAndIssuance(t *testing.T) {
	hash, err := password.Hash("password")
	if err != nil {
		t.Fatal(err)
	}
	unavailable := errors.New("database unavailable")
	for _, tc := range []struct {
		name, raw      string
		dbErr, wantErr error
	}{
		{name: "correct", raw: "password"},
		{name: "wrong password", raw: "wrong", wantErr: ErrInvalidCredentials},
		{name: "unknown user", raw: "password", dbErr: domain_error.ErrNotFound, wantErr: ErrInvalidCredentials},
		{name: "database failure", raw: "password", dbErr: unavailable, wantErr: unavailable},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tokens := &issuer{}
			token, err := New(storage{hash: hash, err: tc.dbErr}, tokens).Execute(context.Background(), "admin@example.test", tc.raw)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("error=%v", err)
			}
			if tc.wantErr != nil && (tokens.issued != 0 || token != "") {
				t.Fatal("issued token on failure")
			}
			if tc.wantErr == nil && (token != "token" || tokens.issued != 7) {
				t.Fatal("wrong token subject")
			}
		})
	}
}
