package dal

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/nagorneva/nagorneva_reviews/internal/pkg/postgres"
)

// These fakes never connect to a database; unexpected operations fail the test.
type queryDB struct {
	postgres.DB
	query func(context.Context, string, ...any) (pgx.Rows, error)
}

func (d queryDB) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return d.query(ctx, sql, args...)
}

type testRows struct {
	pgx.Rows
	data   [][]any
	index  int
	err    error
	closed bool
}

func (r *testRows) Next() bool { r.index++; return r.index <= len(r.data) }
func (r *testRows) Scan(dest ...any) error {
	for i, value := range r.data[r.index-1] {
		reflect.ValueOf(dest[i]).Elem().Set(reflect.ValueOf(value))
	}
	return nil
}
func (r *testRows) Err() error { return r.err }
func (r *testRows) Close()     { r.closed = true }

func TestGetBatchesReviews(t *testing.T) {
	doctors := &testRows{}
	for _, id := range []int64{1, 2, 3} {
		doctors.data = append(doctors.data, []any{id, "Doctor", "", "", true, int64(1), int64(1), float64(60), float64(60), true, "City"})
	}
	now := time.Now()
	reviews := &testRows{data: [][]any{
		{int64(11), int64(1), int16(5), "new", int64(4), now, true, "Course 4", "https://example.com/4"},
		{int64(12), int64(1), int16(4), "old", int64(5), now.Add(-time.Hour), true, "Course 5", "https://example.com/5"},
		{int64(21), int64(2), int16(3), "other", int64(4), now, true, "Course 4", "https://example.com/4"},
	}}
	calls := 0
	db := queryDB{query: func(_ context.Context, sql string, args ...any) (pgx.Rows, error) {
		calls++
		switch calls {
		case 1:
			return doctors, nil
		case 2:
			if !doctors.closed {
				t.Fatal("doctor rows must close before next query")
			}
			return &testRows{}, nil
		case 3:
			if !strings.Contains(sql, "ANY($1::bigint[])") || !strings.Contains(sql, "AND r.is_active") || !strings.Contains(sql, "r.created_at DESC") {
				t.Fatalf("unexpected reviews query: %s", sql)
			}
			if len(args) != 1 || !reflect.DeepEqual(args[0], []int64{1, 2, 3}) {
				t.Fatalf("unexpected IDs: %v", args)
			}
			return reviews, nil
		default:
			t.Fatalf("unexpected query %d", calls)
			return nil, nil
		}
	}}
	result, err := New(db).Get(context.Background(), 0, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if calls != 3 {
		t.Fatalf("got %d queries, want 3", calls)
	}
	if len(result) != 3 || len(result[0].Reviews) != 2 || len(result[1].Reviews) != 1 || result[2].Reviews == nil || len(result[2].Reviews) != 0 {
		t.Fatalf("incorrect grouping: %+v", result)
	}
	if result[0].Reviews[0].ID != 11 || result[0].Reviews[1].ID != 12 || result[1].Reviews[0].ID != 21 {
		t.Fatal("review order or ownership changed")
	}
	if result[0].Reviews[1].Rating != 4 || result[0].Reviews[1].CourseName != "Course 5" {
		t.Fatal("review fields changed")
	}
	if !reviews.closed {
		t.Fatal("reviews rows not closed")
	}
}

func TestGetEmptySkipsRelatedQueries(t *testing.T) {
	calls := 0
	db := queryDB{query: func(context.Context, string, ...any) (pgx.Rows, error) { calls++; return &testRows{}, nil }}
	result, err := New(db).Get(context.Background(), 0, 0, 0)
	if err != nil || len(result) != 0 || calls != 1 {
		t.Fatalf("result=%v err=%v calls=%d", result, err, calls)
	}
}

func TestReviewsPropagatesErrors(t *testing.T) {
	want := errors.New("query failed")
	for _, rowError := range []bool{false, true} {
		db := queryDB{query: func(context.Context, string, ...any) (pgx.Rows, error) {
			if rowError {
				return &testRows{err: want}, nil
			}
			return nil, want
		}}
		_, err := New(db).reviews(context.Background(), []int64{1, 2})
		if !errors.Is(err, want) {
			t.Fatalf("got %v, want %v", err, want)
		}
	}
}
