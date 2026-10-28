package get_doctors_on_map

import (
	"context"
	"errors"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/reviews/domain"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type storageFunc func(context.Context, float64, float64, float64) ([]domain.DoctorOnMap, error)

func (f storageFunc) Get(ctx context.Context, a, b, c float64) ([]domain.DoctorOnMap, error) {
	return f(ctx, a, b, c)
}
func TestCacheTTLAndIsolation(t *testing.T) {
	calls := 0
	a := New(storageFunc(func(context.Context, float64, float64, float64) ([]domain.DoctorOnMap, error) {
		calls++
		return []domain.DoctorOnMap{{Reviews: []domain.ReviewWithCourse{{CourseName: "original"}}}}, nil
	}))
	now := time.Now()
	a.cache.now = func() time.Time { return now }
	first, _ := a.Execute(context.Background(), 0, 0, 0)
	first[0].Reviews[0].CourseName = "changed"
	second, _ := a.Execute(context.Background(), 60, 60, 0)
	if calls != 1 || second[0].Reviews[0].CourseName != "original" {
		t.Fatal("cache miss or mutation")
	}
	now = now.Add(59 * time.Second)
	a.Execute(context.Background(), 0, 0, 0)
	if calls != 1 {
		t.Fatal("expired early")
	}
	now = now.Add(time.Second)
	a.Execute(context.Background(), 0, 0, 0)
	if calls != 2 {
		t.Fatal("did not expire")
	}
	a.Execute(context.Background(), 0, 0, 25)
	if calls != 3 {
		t.Fatal("viewport must bypass full-map cache")
	}
}
func TestCacheConcurrentAndCancellation(t *testing.T) {
	var calls atomic.Int32
	started, release := make(chan struct{}), make(chan struct{})
	a := New(storageFunc(func(context.Context, float64, float64, float64) ([]domain.DoctorOnMap, error) {
		calls.Add(1)
		close(started)
		<-release
		return []domain.DoctorOnMap{}, nil
	}))
	var wg sync.WaitGroup
	wg.Add(1)
	go func() { defer wg.Done(); a.Execute(context.Background(), 0, 0, 0) }()
	<-started
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := a.Execute(ctx, 0, 0, 0); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := a.Execute(context.Background(), 0, 0, 0); err != nil {
				t.Error(err)
			}
		}()
	}
	close(release)
	wg.Wait()
	if calls.Load() != 1 {
		t.Fatalf("%d duplicate loads", calls.Load())
	}
}
func TestCacheDoesNotStoreErrors(t *testing.T) {
	calls := 0
	a := New(storageFunc(func(context.Context, float64, float64, float64) ([]domain.DoctorOnMap, error) {
		calls++
		if calls == 1 {
			return nil, errors.New("failed")
		}
		return []domain.DoctorOnMap{}, nil
	}))
	if _, err := a.Execute(context.Background(), 0, 0, 0); err == nil {
		t.Fatal("missing error")
	}
	a.Execute(context.Background(), 0, 0, 0)
	a.Execute(context.Background(), 0, 0, 0)
	if calls != 2 {
		t.Fatalf("calls=%d", calls)
	}
}
