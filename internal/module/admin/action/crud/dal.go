package crud

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DAL struct{ db *pgxpool.Pool }

const doctorSelect = `SELECT
	d.id, d.name, d.full_name, d.photo, d.personal_data_consent,
	c.id AS city_id, c.name AS city_name, c.lat AS city_lat, c.lon AS city_lon,
	s.id AS specialty_id, s.name AS specialty_name,
	d.lat, d.lon, d.is_active,
	(SELECT count(*) FROM reviews review WHERE review.doctor_id = d.id) AS reviews_count
FROM doctors d
JOIN cities c ON c.id = d.city_id
JOIN specialties s ON s.id = d.specialty_id`

const cityListSelect = `SELECT
	city.id, city.name, city.lat, city.lon,
	(SELECT count(*) FROM doctors doctor WHERE doctor.city_id = city.id) AS doctors_count
FROM cities city
ORDER BY city.id`

const specialtyListSelect = `SELECT
	specialty.id, specialty.name,
	(SELECT count(*) FROM doctors doctor WHERE doctor.specialty_id = specialty.id) AS doctors_count
FROM specialties specialty
ORDER BY specialty.id`

const courseListSelect = `SELECT
	course.id, course.name, course.site_link,
	(SELECT count(*) FROM doctor_courses completion WHERE completion.course_id = course.id) AS doctors_count
FROM courses course
ORDER BY course.id`

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
	if r == "doctors" {
		var id int64
		q := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING id", r, strings.Join(cols, ","), placeholders(len(args)))
		if err := d.db.QueryRow(ctx, q, args...).Scan(&id); err != nil {
			return nil, err
		}
		return d.get(ctx, r, id)
	}
	q := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING %s", r, strings.Join(cols, ","), placeholders(len(args)), fieldList(r))
	return d.one(ctx, q, args...)
}
func (d DAL) get(ctx context.Context, r string, id int64) (map[string]any, error) {
	if _, ok := columns[r]; !ok {
		return nil, fmt.Errorf("unknown resource")
	}
	if r == "doctors" {
		return d.one(ctx, doctorSelect+" WHERE d.id=$1", id)
	}
	return d.one(ctx, fmt.Sprintf("SELECT %s FROM %s WHERE id=$1", fieldList(r), r), id)
}
func (d DAL) list(ctx context.Context, r string) ([]map[string]any, error) {
	if _, ok := columns[r]; !ok {
		return nil, fmt.Errorf("unknown resource")
	}
	query := fmt.Sprintf("SELECT %s FROM %s ORDER BY id", fieldList(r), r)
	switch r {
	case "cities":
		query = cityListSelect
	case "specialties":
		query = specialtyListSelect
	case "courses":
		query = courseListSelect
	case "doctors":
		query = doctorSelect + " ORDER BY d.id"
	}
	rows, err := d.db.Query(ctx, query)
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
		var updatedID int64
		q := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d RETURNING id", r, strings.Join(sets, ","), len(args))
		if err := d.db.QueryRow(ctx, q, args...).Scan(&updatedID); err != nil {
			return nil, err
		}
		return d.get(ctx, r, updatedID)
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
