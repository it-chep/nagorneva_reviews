package reviewsv1

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	reviewspb "github.com/nagorneva/nagorneva_reviews/gen/go/reviews"
)

func TestMapCountJSONIncludesZero(t *testing.T) {
	mux := runtime.NewServeMux()
	_, marshaler := runtime.MarshalerForRequest(mux, httptest.NewRequest("POST", "/api/v1/doctors_on_map", nil))
	for _, count := range []int64{0, 42} {
		data, err := marshaler.Marshal(&reviewspb.GetDoctorsOnMapResponse{AllDoctorsCount: count})
		if err != nil {
			t.Fatal(err)
		}
		var body map[string]json.RawMessage
		if err := json.Unmarshal(data, &body); err != nil {
			t.Fatal(err)
		}
		var actual int64
		// Protobuf JSON represents int64 values as decimal strings.
		if value, ok := body["all_doctors_count"]; !ok {
			t.Fatalf("missing all_doctors_count in %s", data)
		} else {
			var number json.Number
			if err := json.Unmarshal(value, &number); err != nil {
				t.Fatal(err)
			}
			actual, err = number.Int64()
			if err != nil {
				t.Fatal(err)
			}
		}
		if actual != count {
			t.Fatalf("count=%d, want %d: %s", actual, count, data)
		}
	}
}
