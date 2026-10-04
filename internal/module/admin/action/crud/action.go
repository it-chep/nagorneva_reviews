// Package crud contains independent admin endpoint actions. Each method maps 1:1
// to an RPC in api/admin/admin.proto; it does not import any other action.
package crud

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Action struct{ dal DAL }

func New(db *pgxpool.Pool) Action { return Action{DAL{db: db}} }
func (a Action) Create(ctx context.Context, resource string, values map[string]any) (map[string]any, error) {
	return a.dal.create(ctx, resource, values)
}
func (a Action) Get(ctx context.Context, resource string, id int64) (map[string]any, error) {
	return a.dal.get(ctx, resource, id)
}
func (a Action) List(ctx context.Context, resource string) ([]map[string]any, error) {
	return a.dal.list(ctx, resource)
}
func (a Action) Update(ctx context.Context, resource string, id int64, values map[string]any) (map[string]any, error) {
	return a.dal.update(ctx, resource, id, values)
}
func (a Action) Delete(ctx context.Context, resource string, id int64) error {
	return a.dal.delete(ctx, resource, id)
}
func (a Action) DoctorReviews(ctx context.Context, doctorID int64) ([]map[string]any, error) {
	return a.dal.doctorReviews(ctx, doctorID)
}
func validate(resource string, values map[string]any) error {
	allowed, ok := columns[resource]
	if !ok {
		return fmt.Errorf("unknown resource %q", resource)
	}
	for key := range values {
		if _, ok := allowed[key]; !ok {
			return fmt.Errorf("field %q cannot be changed", key)
		}
	}
	return nil
}

var columns = map[string]map[string]struct{}{
	"users":       {"email": {}, "password": {}},
	"cities":      {"name": {}, "lat": {}, "lon": {}},
	"specialties": {"name": {}},
	"courses":     {"name": {}, "site_link": {}},
	"doctors":     {"name": {}, "full_name": {}, "photo": {}, "personal_data_consent": {}, "city_id": {}, "specialty_id": {}, "lat": {}, "lon": {}, "is_active": {}},
	"reviews":     {"doctor_id": {}, "rating": {}, "comment": {}, "course_id": {}, "is_active": {}},
}

func setClause(values map[string]any, start int) ([]string, []any) {
	clauses := make([]string, 0, len(values))
	args := make([]any, 0, len(values))
	i := start
	for col, v := range values {
		clauses = append(clauses, fmt.Sprintf("%s=$%d", col, i))
		args = append(args, v)
		i++
	}
	return clauses, args
}
func fields(resource string) []string {
	r := make([]string, 0, len(columns[resource])+1)
	r = append(r, "id")
	for c := range columns[resource] {
		if c != "password" {
			r = append(r, c)
		}
	}
	return r
}
func normalized(_ string, values map[string]any) map[string]any { return values }
func fieldList(resource string) string                          { return strings.Join(fields(resource), ",") }
