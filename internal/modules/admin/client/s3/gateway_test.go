package s3

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

type fakeClient struct {
	putInput *awss3.PutObjectInput
	putErr   error
}

func (c *fakeClient) PutObject(_ context.Context, input *awss3.PutObjectInput, _ ...func(*awss3.Options)) (*awss3.PutObjectOutput, error) {
	c.putInput = input
	return &awss3.PutObjectOutput{}, c.putErr
}

func (c *fakeClient) DeleteObject(context.Context, *awss3.DeleteObjectInput, ...func(*awss3.Options)) (*awss3.DeleteObjectOutput, error) {
	return &awss3.DeleteObjectOutput{}, nil
}

func TestPutDoctorPhoto(t *testing.T) {
	client := &fakeClient{}
	gateway := &Gateway{
		client:    client,
		bucket:    "doctor-photos",
		publicURL: "https://storage.yandexcloud.net/doctor-photos",
	}

	photoURL, err := gateway.PutDoctorPhoto(context.Background(), "images/user_7_1.png", "avatar.png", strings.NewReader("image"))
	if err != nil {
		t.Fatalf("PutDoctorPhoto() error = %v", err)
	}
	if want := "https://storage.yandexcloud.net/doctor-photos/images/user_7_1.png"; photoURL != want {
		t.Errorf("PutDoctorPhoto() URL = %q, want %q", photoURL, want)
	}
	if got := aws.ToString(client.putInput.Bucket); got != "doctor-photos" {
		t.Errorf("PutObject bucket = %q", got)
	}
	if got := aws.ToString(client.putInput.Key); got != "images/user_7_1.png" {
		t.Errorf("PutObject key = %q", got)
	}
	if got := aws.ToString(client.putInput.ContentType); got != "image/png" {
		t.Errorf("PutObject content type = %q", got)
	}
	if got := client.putInput.Metadata["origin_name"]; got != "avatar.png" {
		t.Errorf("PutObject origin_name = %q", got)
	}
}

func TestPutDoctorPhotoReturnsUploadError(t *testing.T) {
	gateway := &Gateway{client: &fakeClient{putErr: errors.New("unavailable")}, bucket: "doctor-photos", publicURL: "https://storage.yandexcloud.net/doctor-photos"}

	_, err := gateway.PutDoctorPhoto(context.Background(), "images/user_7_1.jpg", "avatar.jpg", strings.NewReader("image"))
	if err == nil || !strings.Contains(err.Error(), "upload doctor photo") {
		t.Fatalf("PutDoctorPhoto() error = %v, want wrapped upload error", err)
	}
}
