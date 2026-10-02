package upload_doctor_photo

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/crud"
	s3gateway "github.com/nagorneva/nagorneva_reviews/internal/module/admin/client/s3"
)

const MaxPhotoSize = 10 << 20

var (
	ErrS3NotConfigured  = errors.New("S3 is not configured")
	ErrInvalidPhoto     = errors.New("photo must be between 1 byte and 10 MiB")
	ErrInvalidPhotoMeta = errors.New("filename and image content type are required")
)

type Action struct {
	crud crud.Action
	s3   *s3gateway.Gateway
}

func New(crudAction crud.Action, s3Client *s3gateway.Gateway) *Action {
	return &Action{crud: crudAction, s3: s3Client}
}

func (a *Action) Execute(ctx context.Context, doctorID int64, photo []byte, filename, contentType string) (map[string]any, error) {
	if a.s3 == nil {
		return nil, ErrS3NotConfigured
	}
	if len(photo) == 0 || len(photo) > MaxPhotoSize {
		return nil, ErrInvalidPhoto
	}
	if filename == "" || !strings.HasPrefix(contentType, "image/") {
		return nil, ErrInvalidPhotoMeta
	}
	key := fmt.Sprintf("images/user_%d_%d%s", doctorID, time.Now().UnixNano(), filepath.Ext(filename))
	photoURL, err := a.s3.PutDoctorPhoto(ctx, key, filename, bytes.NewReader(photo))
	if err != nil {
		return nil, err
	}
	return a.crud.Update(ctx, "doctors", doctorID, map[string]any{"photo": photoURL})
}
