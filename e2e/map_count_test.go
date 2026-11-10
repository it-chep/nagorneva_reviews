//go:build integration

package e2e

import (
	"context"
	"testing"

	reviewspb "github.com/nagorneva/nagorneva_reviews/gen/go/reviews"
	reviewsapi "github.com/nagorneva/nagorneva_reviews/internal/app/api/reviews/v1"
	reviewsmodule "github.com/nagorneva/nagorneva_reviews/internal/modules/reviews"
	mapdal "github.com/nagorneva/nagorneva_reviews/internal/modules/reviews/action/get_doctors_on_map/dal"
)

func TestMapTotalExcludesHiddenAndNoCityDoctors(t *testing.T) {
	db := testDatabase(t)
	ctx := context.Background()
	_, err := db.Exec(ctx, `
INSERT INTO cities(id,name,lat,lon) VALUES(1,'Known city',55,37),(2,'No city',0,0);
INSERT INTO specialties(id,name) VALUES(1,'Specialty');
INSERT INTO doctors(id,name,full_name,city_id,specialty_id,lat,lon,is_active) VALUES
(1,'Nearby','Nearby',1,1,55,37,true),
(2,'Far away','Far away',1,1,60,30,true),
(3,'No city','No city',2,1,55,37,true),
(4,'Hidden','Hidden',1,1,55,37,false);`)
	if err != nil {
		t.Fatal(err)
	}
	count, err := mapdal.New(db).CountAll(ctx)
	if err != nil || count != 2 {
		t.Fatalf("count=%d err=%v", count, err)
	}
	api := reviewsapi.New(reviewsmodule.New(db))
	for _, test := range []struct {
		name    string
		request *reviewspb.GetDoctorsOnMapRequest
		wantIDs []int64
	}{
		{"nearby", &reviewspb.GetDoctorsOnMapRequest{Lat: 55, Lon: 37, RadiusKm: 10}, []int64{1}},
		{"empty viewport", &reviewspb.GetDoctorsOnMapRequest{Lat: -40, Lon: -70, RadiusKm: 10}, nil},
		{"all", &reviewspb.GetDoctorsOnMapRequest{}, []int64{1, 2}},
		{"cached all", &reviewspb.GetDoctorsOnMapRequest{}, []int64{1, 2}},
	} {
		t.Run(test.name, func(t *testing.T) {
			response, err := api.GetDoctorsOnMap(ctx, test.request)
			if err != nil {
				t.Fatal(err)
			}
			if response.GetAllDoctorsCount() != 2 || len(response.GetDoctors()) != len(test.wantIDs) {
				t.Fatalf("unexpected map response: %v", response)
			}
			for i, id := range test.wantIDs {
				if response.Doctors[i].Id != id {
					t.Fatalf("doctor[%d]=%d, want %d", i, response.Doctors[i].Id, id)
				}
			}
		})
	}
}
