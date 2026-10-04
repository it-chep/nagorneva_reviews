package adminv1

import (
	"time"

	adminpb "github.com/nagorneva/nagorneva_reviews/gen/go/admin"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func user(v map[string]any) *adminpb.User {
	if v == nil {
		return nil
	}
	return &adminpb.User{Id: int64Value(v["id"]), Email: stringValue(v["email"])}
}

func city(v map[string]any) *adminpb.City {
	if v == nil {
		return nil
	}
	return &adminpb.City{Id: int64Value(v["id"]), Name: stringValue(v["name"]), Lat: float64Value(v["lat"]), Lon: float64Value(v["lon"]), DoctorsCount: int64Value(v["doctors_count"])}
}

func specialty(v map[string]any) *adminpb.Specialty {
	if v == nil {
		return nil
	}
	return &adminpb.Specialty{Id: int64Value(v["id"]), Name: stringValue(v["name"]), DoctorsCount: int64Value(v["doctors_count"])}
}

func course(v map[string]any) *adminpb.Course {
	if v == nil {
		return nil
	}
	return &adminpb.Course{Id: int64Value(v["id"]), Name: stringValue(v["name"]), DoctorsCount: int64Value(v["doctors_count"])}
}

func doctor(v map[string]any) *adminpb.Doctor {
	if v == nil {
		return nil
	}
	return &adminpb.Doctor{
		Id:                  int64Value(v["id"]),
		Name:                stringValue(v["name"]),
		FullName:            stringValue(v["full_name"]),
		Photo:               stringValue(v["photo"]),
		PersonalDataConsent: boolValue(v["personal_data_consent"]),
		City: &adminpb.City{
			Id:   int64Value(v["city_id"]),
			Name: stringValue(v["city_name"]),
			Lat:  float64Value(v["city_lat"]),
			Lon:  float64Value(v["city_lon"]),
		},
		Specialty: &adminpb.Specialty{
			Id:   int64Value(v["specialty_id"]),
			Name: stringValue(v["specialty_name"]),
		},
		Lat:          float64Value(v["lat"]),
		Lon:          float64Value(v["lon"]),
		IsActive:     boolValue(v["is_active"]),
		ReviewsCount: int64Value(v["reviews_count"]),
	}
}

func review(v map[string]any) *adminpb.Review {
	if v == nil {
		return nil
	}
	out := &adminpb.Review{Id: int64Value(v["id"]), DoctorId: int64Value(v["doctor_id"]), Rating: int32(int64Value(v["rating"])), Comment: stringValue(v["comment"]), CourseId: int64Value(v["course_id"]), IsActive: boolValue(v["is_active"]), CourseName: stringValue(v["course_name"])}
	if createdAt, ok := v["created_at"].(time.Time); ok {
		out.CreatedAt = timestamppb.New(createdAt)
	}
	return out
}

func users(items []map[string]any) []*adminpb.User {
	out := make([]*adminpb.User, 0, len(items))
	for _, item := range items {
		out = append(out, user(item))
	}
	return out
}

func cities(items []map[string]any) []*adminpb.City {
	out := make([]*adminpb.City, 0, len(items))
	for _, item := range items {
		out = append(out, city(item))
	}
	return out
}

func specialties(items []map[string]any) []*adminpb.Specialty {
	out := make([]*adminpb.Specialty, 0, len(items))
	for _, item := range items {
		out = append(out, specialty(item))
	}
	return out
}

func courses(items []map[string]any) []*adminpb.Course {
	out := make([]*adminpb.Course, 0, len(items))
	for _, item := range items {
		out = append(out, course(item))
	}
	return out
}

func doctors(items []map[string]any) []*adminpb.Doctor {
	out := make([]*adminpb.Doctor, 0, len(items))
	for _, item := range items {
		out = append(out, doctor(item))
	}
	return out
}

func reviews(items []map[string]any) []*adminpb.Review {
	out := make([]*adminpb.Review, 0, len(items))
	for _, item := range items {
		out = append(out, review(item))
	}
	return out
}

func stringValue(value any) string { text, _ := value.(string); return text }
func boolValue(value any) bool     { flag, _ := value.(bool); return flag }

func float64Value(value any) float64 {
	switch value := value.(type) {
	case float64:
		return value
	case float32:
		return float64(value)
	case int64:
		return float64(value)
	case int32:
		return float64(value)
	case int16:
		return float64(value)
	case int:
		return float64(value)
	default:
		return 0
	}
}

func int64Value(value any) int64 {
	switch value := value.(type) {
	case int64:
		return value
	case int32:
		return int64(value)
	case int16:
		return int64(value)
	case int:
		return int64(value)
	case float64:
		return int64(value)
	case float32:
		return int64(value)
	default:
		return 0
	}
}
