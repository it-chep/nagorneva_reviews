package upload_doctor_photo

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/doctors/dto"
	"github.com/nagorneva/nagorneva_reviews/internal/pkg/domain_error"
)

type doctorStorage struct {
	getErr, setErr error
	photo          string
}

func (s *doctorStorage) GetDoctor(context.Context, int64) (*dto.Doctor, error) {
	return &dto.Doctor{ID: 7}, s.getErr
}
func (s *doctorStorage) SetPhoto(_ context.Context, id int64, photo string) (*dto.Doctor, error) {
	s.photo = photo
	return &dto.Doctor{ID: id, Photo: photo}, s.setErr
}

type imageStorage struct {
	err       error
	calls     int
	body, key string
}

func (s *imageStorage) PutDoctorPhoto(_ context.Context, key, filename string, body io.Reader) (string, error) {
	s.calls++
	s.key = key
	data, err := io.ReadAll(body)
	if err != nil {
		return "", err
	}
	s.body = string(data)
	return "https://example.test/photo.png", s.err
}

func TestUploadFailureBoundaries(t *testing.T) {
	failed := errors.New("storage failed")
	for _, tc := range []struct {
		name                              string
		getErr, putErr, setErr, errorWant error
		size                              int
		filename, contentType             string
		uploads                           int
	}{
		{name: "empty", errorWant: ErrInvalidPhoto},
		{name: "too large", size: MaxPhotoSize + 1, errorWant: ErrInvalidPhoto},
		{name: "invalid metadata", size: 1, filename: "file.txt", contentType: "text/plain", errorWant: ErrInvalidPhotoMeta},
		{name: "missing doctor", size: 1, filename: "photo.png", contentType: "image/png", getErr: domain_error.ErrNotFound, errorWant: domain_error.ErrNotFound},
		{name: "S3 failure", size: 1, filename: "photo.png", contentType: "image/png", putErr: failed, errorWant: failed, uploads: 1},
		{name: "database failure", size: 1, filename: "photo.png", contentType: "image/png", setErr: failed, errorWant: failed, uploads: 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			db := &doctorStorage{getErr: tc.getErr, setErr: tc.setErr}
			images := &imageStorage{err: tc.putErr}
			_, err := New(db, images).Execute(context.Background(), 7, make([]byte, tc.size), tc.filename, tc.contentType)
			if !errors.Is(err, tc.errorWant) {
				t.Fatalf("error=%v want=%v", err, tc.errorWant)
			}
			if images.calls != tc.uploads {
				t.Fatalf("uploads=%d", images.calls)
			}
			if tc.putErr != nil && db.photo != "" {
				t.Fatal("database updated after failed upload")
			}
		})
	}
	if _, err := New(&doctorStorage{}, nil).Execute(context.Background(), 7, []byte("image"), "photo.png", "image/png"); !errors.Is(err, ErrS3NotConfigured) {
		t.Fatal(err)
	}
}
func TestUploadPersistsReturnedURL(t *testing.T) {
	db := &doctorStorage{}
	images := &imageStorage{}
	doctor, err := New(db, images).Execute(context.Background(), 7, []byte("image"), "photo.png", "image/png")
	if err != nil {
		t.Fatal(err)
	}
	if doctor.Photo != "https://example.test/photo.png" || db.photo != doctor.Photo || images.body != "image" || !strings.HasPrefix(images.key, "images/user_7_") {
		t.Fatalf("doctor=%+v images=%+v", doctor, images)
	}
}
