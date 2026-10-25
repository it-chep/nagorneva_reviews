package upload_doctor_photo

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"io"

	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/doctors/dto"
)

const MaxPhotoSize = 10 << 20

var (
	ErrS3NotConfigured  = errors.New("S3 is not configured")
	ErrInvalidPhoto     = errors.New("photo must be between 1 byte and 10 MiB")
	ErrInvalidPhotoMeta = errors.New("filename and image content type are required")
)

type Storage interface {
	GetDoctor(context.Context, int64) (*dto.Doctor, error)
	SetPhoto(context.Context, int64, string) (*dto.Doctor, error)
}
type ImageStorage interface {
	PutDoctorPhoto(context.Context, string, string, io.Reader) (string, error)
}
type Action struct {
	storage Storage
	s3      ImageStorage
}

func New(storage Storage, images ImageStorage) *Action { return &Action{storage: storage, s3: images} }

func (a *Action) Execute(ctx context.Context, doctorID int64, photo []byte, filename, contentType string) (*dto.Doctor, error) {
	if a.s3 == nil {
		return nil, ErrS3NotConfigured
	}
	if len(photo) == 0 || len(photo) > MaxPhotoSize {
		return nil, ErrInvalidPhoto
	}
	if filename == "" || !strings.HasPrefix(contentType, "image/") {
		return nil, ErrInvalidPhotoMeta
	}
	if _, err := a.storage.GetDoctor(ctx, doctorID); err != nil {
		return nil, err
	}
	key := fmt.Sprintf("images/user_%d_%d%s", doctorID, time.Now().UnixNano(), filepath.Ext(filename))
	photoURL, err := a.s3.PutDoctorPhoto(ctx, key, filename, bytes.NewReader(photo))
	if err != nil {
		return nil, err
	}
	return a.storage.SetPhoto(ctx, doctorID, photoURL)
}
