package get_doctors_on_map

import (
	"context"
	"sync"
	"time"

	"github.com/nagorneva/nagorneva_reviews/internal/modules/reviews/domain"
)

type Storage interface {
	Get(ctx context.Context, lat, lon, radius float64) ([]domain.DoctorOnMap, error)
}
type cacheEntry struct {
	doctors []domain.DoctorOnMap
	expires time.Time
}
type memoryCache struct {
	mu      sync.Mutex
	entries map[string]cacheEntry
	loading chan struct{}
	now     func() time.Time
}
type Action struct {
	dal   Storage
	cache *memoryCache
}

func New(storage Storage) Action {
	return Action{dal: storage, cache: &memoryCache{entries: make(map[string]cacheEntry), now: time.Now}}
}

// Cache the complete map only: arbitrary viewport requests do not grow the cache.
// The pointer keeps synchronization shared when Action is copied by value.
func (a Action) Execute(ctx context.Context, lat, lon, radiusKm float64) ([]domain.DoctorOnMap, error) {
	if radiusKm != 0 {
		return a.dal.Get(ctx, lat, lon, radiusKm)
	}
	const key = "all"
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		a.cache.mu.Lock()
		if entry, ok := a.cache.entries[key]; ok && a.cache.now().Before(entry.expires) {
			result := cloneDoctors(entry.doctors)
			a.cache.mu.Unlock()
			return result, nil
		}
		if done := a.cache.loading; done != nil {
			a.cache.mu.Unlock()
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-done:
				continue
			}
		}
		a.cache.loading = make(chan struct{})
		a.cache.mu.Unlock()
		doctors, err := a.dal.Get(ctx, 0, 0, 0)
		a.cache.mu.Lock()
		if err == nil {
			a.cache.entries[key] = cacheEntry{doctors: cloneDoctors(doctors), expires: a.cache.now().Add(60 * time.Second)}
		}
		close(a.cache.loading)
		a.cache.loading = nil
		a.cache.mu.Unlock()
		return doctors, err
	}
}
func cloneDoctors(source []domain.DoctorOnMap) []domain.DoctorOnMap {
	result := make([]domain.DoctorOnMap, len(source))
	copy(result, source)
	for i := range result {
		if source[i].Reviews != nil {
			result[i].Reviews = append([]domain.ReviewWithCourse{}, source[i].Reviews...)
		}
		if source[i].Courses != nil {
			result[i].Courses = append([]domain.Course{}, source[i].Courses...)
		}
	}
	return result
}
