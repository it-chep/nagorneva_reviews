//go:build integration

package e2e

import (
	"context"
	"testing"
	"time"

	adminpb "github.com/nagorneva/nagorneva_reviews/gen/go/admin"
	adminapi "github.com/nagorneva/nagorneva_reviews/internal/app/api/admin/v1"
	adminmodule "github.com/nagorneva/nagorneva_reviews/internal/modules/admin"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestDoctorCoursesLifecycle(t *testing.T) {
	db := testDatabase(t)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_, err := db.Exec(ctx, `
INSERT INTO cities(id,name,lat,lon) VALUES(2,'No city',0,0);
INSERT INTO specialties(id,name) VALUES(1,'Specialty');
INSERT INTO courses(id,name,site_link) VALUES(1,'First','https://example.test/first'),(2,'Second','');
INSERT INTO doctors(id,name,full_name,city_id,specialty_id,lat,lon,is_active) VALUES
(1,'Doctor','Doctor',2,1,0,0,false),(2,'Other','Other',2,1,0,0,false);
INSERT INTO doctor_courses(doctor_id,course_id) VALUES(2,1);
INSERT INTO reviews(doctor_id,rating,comment,course_id,is_active) VALUES(1,5,'Review survives unlinking',1,true);`)
	if err != nil {
		t.Fatal(err)
	}
	service := adminapi.New(&adminmodule.Module{Actions: entities.New(db, nil, nil)})
	result, err := service.GetDoctorCourses(ctx, &adminpb.GetDoctorCoursesRequest{Id: 1})
	if err != nil || result.Courses == nil || len(result.Courses) != 0 {
		t.Fatalf("empty courses: %v %v", result, err)
	}
	for _, courseID := range []int64{2, 1, 1} {
		if _, err := service.AddDoctorCourse(ctx, &adminpb.AddDoctorCourseRequest{Id: 1, CourseId: courseID}); err != nil {
			t.Fatalf("add %d: %v", courseID, err)
		}
	}
	result, err = service.GetDoctorCourses(ctx, &adminpb.GetDoctorCoursesRequest{Id: 1})
	if err != nil || len(result.GetCourses()) != 2 {
		t.Fatalf("courses including duplicate assignment: %v %v", result, err)
	}
	if result.Courses[0].Id != 1 || result.Courses[0].Name != "First" || result.Courses[0].SiteLink != "https://example.test/first" || result.Courses[0].DoctorsCount != 2 || result.Courses[1].Id != 2 {
		t.Fatalf("course metadata/order: %v", result)
	}
	if _, err := service.DeleteDoctorCourse(ctx, &adminpb.DeleteDoctorCourseRequest{Id: 1, CourseId: 1}); err != nil {
		t.Fatal(err)
	}
	result, err = service.GetDoctorCourses(ctx, &adminpb.GetDoctorCoursesRequest{Id: 1})
	if err != nil || len(result.GetCourses()) != 1 || result.Courses[0].Id != 2 {
		t.Fatalf("courses after deleting assignment: %v %v", result, err)
	}
	other, err := service.GetDoctorCourses(ctx, &adminpb.GetDoctorCoursesRequest{Id: 2})
	if err != nil || len(other.GetCourses()) != 1 || other.Courses[0].DoctorsCount != 1 {
		t.Fatalf("other doctor's assignment/count: %v %v", other, err)
	}
	var reviewsCount, coursesCount int
	if err := db.QueryRow(ctx, `SELECT (SELECT count(*) FROM reviews WHERE doctor_id=1 AND course_id=1), (SELECT count(*) FROM courses WHERE id=1)`).Scan(&reviewsCount, &coursesCount); err != nil || reviewsCount != 1 || coursesCount != 1 {
		t.Fatalf("unlink removed a review or course: reviews=%d courses=%d err=%v", reviewsCount, coursesCount, err)
	}
	if _, err := service.DeleteDoctorCourse(ctx, &adminpb.DeleteDoctorCourseRequest{Id: 1, CourseId: 1}); status.Code(err) != codes.NotFound {
		t.Fatalf("missing assignment: %v", err)
	}
	for _, request := range []*adminpb.AddDoctorCourseRequest{{Id: 999, CourseId: 1}, {Id: 1, CourseId: 999}} {
		if _, err := service.AddDoctorCourse(ctx, request); status.Code(err) != codes.NotFound {
			t.Fatalf("missing doctor/course: %v", err)
		}
	}
	if _, err := service.GetDoctorCourses(ctx, &adminpb.GetDoctorCoursesRequest{Id: 999}); status.Code(err) != codes.NotFound {
		t.Fatalf("missing doctor course list: %v", err)
	}
	if _, err := service.GetDoctorCourses(ctx, &adminpb.GetDoctorCoursesRequest{Id: 0}); status.Code(err) != codes.InvalidArgument {
		t.Fatalf("invalid doctor ID: %v", err)
	}
	if _, err := service.AddDoctorCourse(ctx, &adminpb.AddDoctorCourseRequest{Id: 1}); status.Code(err) != codes.InvalidArgument {
		t.Fatalf("missing course ID: %v", err)
	}
	if _, err := service.DeleteDoctorCourse(ctx, &adminpb.DeleteDoctorCourseRequest{Id: 1, CourseId: -1}); status.Code(err) != codes.InvalidArgument {
		t.Fatalf("invalid course ID: %v", err)
	}
}
