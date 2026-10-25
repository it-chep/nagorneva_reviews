//go:build integration

package e2e

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	citydal "github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/cities/dal"
	citydto "github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/cities/dto"
	coursedal "github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/courses/dal"
	coursedto "github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/courses/dto"
	doctordal "github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/doctors/dal"
	doctordto "github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/doctors/dto"
	reviewdal "github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/reviews/dal"
	reviewdto "github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/reviews/dto"
	specialtydal "github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/specialties/dal"
	specialtydto "github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/specialties/dto"
	userdal "github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/users/dal"
	userdto "github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/users/dto"
	filterdal "github.com/nagorneva/nagorneva_reviews/internal/modules/reviews/action/filter_doctors/dal"
	filterdto "github.com/nagorneva/nagorneva_reviews/internal/modules/reviews/action/filter_doctors/dto"
	publicreviewdal "github.com/nagorneva/nagorneva_reviews/internal/modules/reviews/action/filter_reviews/dal"
	mapdal "github.com/nagorneva/nagorneva_reviews/internal/modules/reviews/action/get_doctors_on_map/dal"
	"github.com/nagorneva/nagorneva_reviews/internal/pkg/domain_error"
)

func testDatabase(t *testing.T) *pgxpool.Pool {
	t.Helper()
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Fatal("TEST_DATABASE_URL must point to a disposable test database")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	admin, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	schema := fmt.Sprintf("test_reviews_%d", time.Now().UnixNano())
	if _, err := admin.Exec(ctx, "CREATE SCHEMA "+schema); err != nil {
		admin.Close()
		t.Fatal(err)
	}
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if _, err := admin.Exec(ctx, "DROP SCHEMA "+schema+" CASCADE"); err != nil {
			t.Errorf("drop test schema: %v", err)
		}
		admin.Close()
	})
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		t.Fatal(err)
	}
	cfg.MaxConns = 1 // Queries must release rows before acquiring another connection.
	cfg.ConnConfig.RuntimeParams["search_path"] = schema
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	files, err := filepath.Glob("../migrations/*.sql")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no migrations")
	}
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		up := strings.Split(string(data), "-- +goose Down")[0]
		if _, err := pool.Exec(ctx, up); err != nil {
			t.Fatalf("migration %s: %v", file, err)
		}
	}
	return pool
}
func ptr[T any](v T) *T { return &v }

func TestTypedRepositoriesAndPublicQueries(t *testing.T) {
	db := testDatabase(t)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	cities := citydal.New(db)
	specialties := specialtydal.New(db)
	courses := coursedal.New(db)
	doctors := doctordal.New(db)
	reviews := reviewdal.New(db)
	users := userdal.New(db)
	city, err := cities.CreateCity(ctx, citydto.CreateInput{Name: "Moscow", Lat: 55, Lon: 37})
	if err != nil {
		t.Fatal(err)
	}
	specialty, err := specialties.CreateSpecialty(ctx, specialtydto.CreateInput{Name: "Therapist"})
	if err != nil {
		t.Fatal(err)
	}
	course, err := courses.CreateCourse(ctx, coursedto.CreateInput{Name: "Course", SiteLink: "https://example.test"})
	if err != nil {
		t.Fatal(err)
	}
	doctor, err := doctors.CreateDoctor(ctx, doctordto.CreateInput{Name: "Public name", FullName: "Full name", CityID: city.ID, SpecialtyID: specialty.ID, Lat: 55, Lon: 37, IsActive: true, PersonalDataConsent: true})
	if err != nil {
		t.Fatal(err)
	}
	if doctor.CityName != "Moscow" || doctor.SpecialtyName != "Therapist" {
		t.Fatalf("joined doctor: %+v", doctor)
	}
	doctor, err = doctors.SetPhoto(ctx, doctor.ID, "https://example.test/photo.jpg")
	if err != nil {
		t.Fatal(err)
	}
	doctor, err = doctors.UpdateDoctor(ctx, doctor.ID, doctordto.UpdateInput{Name: ptr(""), Lat: ptr(0.0), IsActive: ptr(false), PersonalDataConsent: ptr(false)})
	if err != nil {
		t.Fatal(err)
	}
	if doctor.Name != "" || doctor.Lat != 0 || doctor.IsActive || doctor.PersonalDataConsent || doctor.FullName != "Full name" || doctor.Photo == "" {
		t.Fatalf("PATCH lost zero/absent fields: %+v", doctor)
	}
	doctor, err = doctors.UpdateDoctor(ctx, doctor.ID, doctordto.UpdateInput{IsActive: ptr(true)})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(ctx, "INSERT INTO doctor_courses(doctor_id,course_id) VALUES($1,$2)", doctor.ID, course.ID); err != nil {
		t.Fatal(err)
	}
	review, err := reviews.CreateReview(ctx, reviewdto.CreateInput{DoctorID: doctor.ID, Rating: 5, Comment: "Review", CourseID: course.ID, IsActive: true})
	if err != nil {
		t.Fatal(err)
	}
	if review.CreatedAt.IsZero() {
		t.Fatal("review timestamp missing")
	}
	if _, err := reviews.CreateReview(ctx, reviewdto.CreateInput{DoctorID: doctor.ID, Rating: 10, CourseID: course.ID}); !errors.Is(err, domain_error.ErrInvalid) {
		t.Fatalf("check constraint: %v", err)
	}
	if _, err := doctors.GetDoctor(ctx, -1); !errors.Is(err, domain_error.ErrNotFound) {
		t.Fatalf("missing doctor: %v", err)
	}
	if err := doctors.DeleteDoctor(ctx, -1); !errors.Is(err, domain_error.ErrNotFound) {
		t.Fatalf("missing delete: %v", err)
	}
	if err := cities.DeleteCity(ctx, city.ID); !errors.Is(err, domain_error.ErrReference) {
		t.Fatalf("FK restriction: %v", err)
	}
	if _, err := cities.CreateCity(ctx, citydto.CreateInput{Name: "Moscow"}); !errors.Is(err, domain_error.ErrConflict) {
		t.Fatalf("unique constraint: %v", err)
	}
	filtered, err := filterdal.New(db).Filter(ctx, []int64{city.ID}, []int64{specialty.ID}, []int64{course.ID}, ptr(true), nil, filterdto.ReviewsSortDescending)
	if err != nil {
		t.Fatal(err)
	}
	if len(filtered) != 1 || filtered[0].Photo != "" || filtered[0].ReviewsCount != 1 || len(filtered[0].CompletedCourses) != 1 {
		t.Fatalf("public filtering: %+v", filtered)
	}
	mapped, err := mapdal.New(db).Get(ctx, 0, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(mapped) != 1 || mapped[0].FullName != "" || mapped[0].Photo != "" || len(mapped[0].Reviews) != 1 {
		t.Fatalf("public map: %+v", mapped)
	}
	publicReviews, err := publicreviewdal.New(db).Filter(ctx, &course.ID, nil, nil, 0)
	if err != nil || len(publicReviews) != 1 {
		t.Fatalf("public reviews: %+v %v", publicReviews, err)
	}
	adminReviews, err := reviews.GetDoctorReviews(ctx, doctor.ID)
	if err != nil || len(adminReviews) != 1 || adminReviews[0].CourseName != "Course" {
		t.Fatalf("doctor reviews: %+v %v", adminReviews, err)
	}
	listedCities, err := cities.ListCities(ctx)
	if err != nil || len(listedCities) != 1 || listedCities[0].DoctorsCount != 1 {
		t.Fatalf("city counts: %+v %v", listedCities, err)
	}
	listedSpecialties, err := specialties.ListSpecialties(ctx)
	if err != nil || listedSpecialties[0].DoctorsCount != 1 {
		t.Fatalf("specialty counts: %+v %v", listedSpecialties, err)
	}
	listedCourses, err := courses.ListCourses(ctx)
	if err != nil || listedCourses[0].DoctorsCount != 1 {
		t.Fatalf("course counts: %+v %v", listedCourses, err)
	}
	if got, err := doctors.ListDoctors(ctx); err != nil || len(got) != 1 || got[0].ReviewsCount != 1 {
		t.Fatalf("doctors: %+v %v", got, err)
	}
	if got, err := reviews.ListReviews(ctx); err != nil || len(got) != 1 {
		t.Fatalf("reviews: %+v %v", got, err)
	}
	city, err = cities.UpdateCity(ctx, city.ID, citydto.UpdateInput{Lat: ptr(0.0)})
	if err != nil || city.Lat != 0 || city.Name != "Moscow" {
		t.Fatalf("city patch: %+v %v", city, err)
	}
	course, err = courses.UpdateCourse(ctx, course.ID, coursedto.UpdateInput{SiteLink: ptr("")})
	if err != nil || course.SiteLink != "" || course.Name != "Course" {
		t.Fatalf("course patch: %+v %v", course, err)
	}
	specialty, err = specialties.UpdateSpecialty(ctx, specialty.ID, specialtydto.UpdateInput{Name: ptr("New")})
	if err != nil || specialty.Name != "New" {
		t.Fatalf("specialty patch: %+v %v", specialty, err)
	}
	review, err = reviews.UpdateReview(ctx, review.ID, reviewdto.UpdateInput{Comment: ptr(""), IsActive: ptr(false)})
	if err != nil || review.IsActive || review.Comment != "" || review.Rating != 5 {
		t.Fatalf("review patch: %+v %v", review, err)
	}
	user, err := users.CreateUser(ctx, userdto.CreateInput{Email: "test@example.test", Password: "hash"})
	if err != nil {
		t.Fatal(err)
	}
	user, err = users.UpdateUser(ctx, user.ID, userdto.UpdateInput{Email: ptr("new@example.test")})
	if err != nil {
		t.Fatal(err)
	}
	var hash string
	if err := db.QueryRow(ctx, "SELECT password FROM users WHERE id=$1", user.ID).Scan(&hash); err != nil || hash != "hash" {
		t.Fatalf("password overwritten: %q %v", hash, err)
	}
	if err := users.CreateInitialUser(ctx, user.Email, "replacement"); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRow(ctx, "SELECT password FROM users WHERE id=$1", user.ID).Scan(&hash); err != nil || hash != "hash" {
		t.Fatal("bootstrap overwrote user")
	}
	if got, err := users.GetUser(ctx, user.ID); err != nil || got.Email != user.Email {
		t.Fatalf("user: %+v %v", got, err)
	}
	if _, err := users.ListUsers(ctx); err != nil {
		t.Fatal(err)
	}
	if err := users.DeleteUser(ctx, user.ID); err != nil {
		t.Fatal(err)
	}
	if err := doctors.DeleteDoctor(ctx, doctor.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := reviews.GetReview(ctx, review.ID); !errors.Is(err, domain_error.ErrNotFound) {
		t.Fatalf("cascade: %v", err)
	}
	if err := courses.DeleteCourse(ctx, course.ID); err != nil {
		t.Fatal(err)
	}
	if err := cities.DeleteCity(ctx, city.ID); err != nil {
		t.Fatal(err)
	}
	if err := specialties.DeleteSpecialty(ctx, specialty.ID); err != nil {
		t.Fatal(err)
	}
}
