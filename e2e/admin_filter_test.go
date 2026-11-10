//go:build integration

package e2e

import (
	"context"
	"testing"
	"time"

	adminpb "github.com/nagorneva/nagorneva_reviews/gen/go/admin"
	reviewspb "github.com/nagorneva/nagorneva_reviews/gen/go/reviews"
	adminapi "github.com/nagorneva/nagorneva_reviews/internal/app/api/admin/v1"
	reviewsapi "github.com/nagorneva/nagorneva_reviews/internal/app/api/reviews/v1"
	adminmodule "github.com/nagorneva/nagorneva_reviews/internal/modules/admin"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities"
	reviewsmodule "github.com/nagorneva/nagorneva_reviews/internal/modules/reviews"
)

func TestAdminAndPublicDoctorVisibility(t *testing.T) {
	db := testDatabase(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := db.Exec(ctx, `
        INSERT INTO cities(id,name,lat,lon) VALUES(1,'City',0,0),(2,'Не указано',0,0);
        INSERT INTO specialties(id,name) VALUES(1,'Specialty'),(2,'Other');
        INSERT INTO courses(id,name,site_link) VALUES(1,'Course','https://example.test');
        INSERT INTO doctors(id,name,full_name,city_id,specialty_id,lat,lon,is_active,personal_data_consent,photo) VALUES
        (1,'Public','Public',1,1,0,0,true,false,'photo-1'),
        (2,'No city','No city',2,1,0,0,true,true,'photo-2'),
        (3,'Inactive','Inactive',1,2,0,0,false,true,'photo-3');
        INSERT INTO doctor_courses(doctor_id,course_id) VALUES(1,1),(2,1);
        INSERT INTO reviews(doctor_id,rating,comment,course_id,is_active) VALUES
        (1,5,'Visible',1,true),(1,5,'Hidden',1,false),(2,5,'No city',1,true),(3,5,'Inactive',1,true);`)
	if err != nil {
		t.Fatal(err)
	}
	admin := adminapi.New(&adminmodule.Module{Actions: entities.New(db, nil, nil)})
	public := reviewsapi.New(reviewsmodule.New(db))
	all, err := admin.FilterDoctors(ctx, &adminpb.FilterDoctorsRequest{})
	if err != nil || all.GetDoctorCount() != 3 || len(all.GetDoctors()) != 3 {
		t.Fatalf("admin all: %v %v", all, err)
	}
	if all.Doctors[0].Photo != "photo-1" || all.Doctors[0].ReviewsCount != 2 {
		t.Fatalf("admin data: %v", all.Doctors[0])
	}
	if all.Doctors[1].City.Id != 2 || all.Doctors[2].IsActive {
		t.Fatalf("admin visibility: %v", all)
	}
	filtered, err := admin.FilterDoctors(ctx, &adminpb.FilterDoctorsRequest{CityIds: []int64{2}, SpecialtyIds: []int64{1}, CourseIds: []int64{1}, IsActive: ptr(true)})
	if err != nil || filtered.GetDoctorCount() != 1 || len(filtered.GetDoctors()) != 1 || filtered.Doctors[0].Id != 2 {
		t.Fatalf("city2 filter: %v %v", filtered, err)
	}
	if len(filtered.Doctors[0].CompletedCourses) != 1 || filtered.Doctors[0].CompletedCourses[0].SiteLink != "https://example.test" {
		t.Fatalf("courses: %v", filtered)
	}
	inactive, err := admin.FilterDoctors(ctx, &adminpb.FilterDoctorsRequest{IsActive: ptr(false)})
	if err != nil || inactive.GetDoctorCount() != 1 || inactive.Doctors[0].Id != 3 {
		t.Fatalf("inactive filter: %v %v", inactive, err)
	}
	noConsent, err := admin.FilterDoctors(ctx, &adminpb.FilterDoctorsRequest{PersonalDataConsent: ptr(false), ReviewsSort: adminpb.ReviewsSort_REVIEWS_SORT_DESC})
	if err != nil || noConsent.GetDoctorCount() != 1 || noConsent.Doctors[0].Id != 1 {
		t.Fatalf("consent filter: %v %v", noConsent, err)
	}
	result, err := public.FilterDoctors(ctx, &reviewspb.FilterDoctorsRequest{})
	if err != nil || result.GetDoctorCount() != 1 || len(result.GetDoctors()) != 1 || result.Doctors[0].Id != 1 {
		t.Fatalf("public filter: %v %v", result, err)
	}
	if result.Doctors[0].Photo != "" || result.Doctors[0].ReviewsCount != 1 {
		t.Fatalf("public data: %v", result)
	}
	hidden, err := public.FilterDoctors(ctx, &reviewspb.FilterDoctorsRequest{CityIds: []int64{2}})
	if err != nil || hidden.GetDoctorCount() != 0 {
		t.Fatalf("public city2: %v %v", hidden, err)
	}
	publicReviews, err := public.FilterReviews(ctx, &reviewspb.FilterReviewsRequest{})
	if err != nil || len(publicReviews.GetReviews()) != 1 || publicReviews.Reviews[0].DoctorId != 1 {
		t.Fatalf("public reviews: %v %v", publicReviews, err)
	}
}
