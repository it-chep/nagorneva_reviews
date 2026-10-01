package crud

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
)

type DAL struct{ db *pgxpool.Pool }

func (d DAL) create(ctx context.Context, r string, v map[string]any) (map[string]any, error) {
	if err := validate(r, v); err != nil {
		return nil, err
	}
	v = normalized(r, v)
	sets, args := setClause(v, 1)
	cols := make([]string, len(sets))
	for i := range sets {
		cols[i] = strings.Split(sets[i], "=")[0]
	}
	q := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING %s", r, strings.Join(cols, ","), placeholders(len(args)), fieldList(r))
	return d.one(ctx, q, args...)
}
func (d DAL) get(ctx context.Context, r string, id int64) (map[string]any, error) {
	if _, ok := columns[r]; !ok {
		return nil, fmt.Errorf("unknown resource")
	}
	return d.one(ctx, fmt.Sprintf("SELECT %s FROM %s WHERE id=$1", fieldList(r), r), id)
}
func (d DAL) list(ctx context.Context, r string) ([]map[string]any, error) {
	if _, ok := columns[r]; !ok {
		return nil, fmt.Errorf("unknown resource")
	}
	rows, err := d.db.Query(ctx, fmt.Sprintf("SELECT %s FROM %s ORDER BY id", fieldList(r), r))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []map[string]any{}
	for rows.Next() {
		v, err := rowMap(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}
func (d DAL) update(ctx context.Context, r string, id int64, v map[string]any) (map[string]any, error) {
	if err := validate(r, v); err != nil {
		return nil, err
	}
	if len(v) == 0 {
		return d.get(ctx, r, id)
	}
	v = normalized(r, v)
	sets, args := setClause(v, 1)
	args = append(args, id)
	if r == "doctors" {
		sets = append(sets, "updated_at=now()")
	}
	q := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d RETURNING %s", r, strings.Join(sets, ","), len(args), fieldList(r))
	return d.one(ctx, q, args...)
}
func (d DAL) delete(ctx context.Context, r string, id int64) error {
	if _, ok := columns[r]; !ok {
		return fmt.Errorf("unknown resource")
	}
	tag, err := d.db.Exec(ctx, fmt.Sprintf("DELETE FROM %s WHERE id=$1", r), id)
	if err == nil && tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return err
}
func (d DAL) doctorReviews(ctx context.Context, id int64) ([]map[string]any, error) {
	rows, err := d.db.Query(ctx, `SELECT r.id,r.doctor_id,r.rating,r.comment,r.course_id,r.created_at,r.is_active,c.name AS course_name FROM reviews r JOIN courses c ON c.id=r.course_id WHERE r.doctor_id=$1 ORDER BY r.created_at DESC`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []map[string]any{}
	for rows.Next() {
		v, e := rowMap(rows)
		if e != nil {
			return nil, e
		}
		out = append(out, v)
	}
	return out, rows.Err()
}
func (d DAL) one(ctx context.Context, q string, args ...any) (map[string]any, error) {
	rows, err := d.db.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, pgx.ErrNoRows
	}
	return rowMap(rows)
}
func rowMap(rows pgx.Rows) (map[string]any, error) {
	vals, err := rows.Values()
	if err != nil {
		return nil, err
	}
	fds := rows.FieldDescriptions()
	out := make(map[string]any, len(vals))
	for i, v := range vals {
		out[string(fds[i].Name)] = v
	}
	return out, nil
}
func placeholders(n int) string {
	p := make([]string, n)
	for i := range p {
		p[i] = fmt.Sprintf("$%d", i+1)
	}
	return strings.Join(p, ",")
}
