package adminv1

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/upload_doctor_photo"
	"github.com/nagorneva/nagorneva_reviews/internal/module/admin/action/validation"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func rpcError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, validation.ErrPasswordTooShort) || errors.Is(err, upload_doctor_photo.ErrInvalidPhoto) || errors.Is(err, upload_doctor_photo.ErrInvalidPhotoMeta) {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	if errors.Is(err, upload_doctor_photo.ErrS3NotConfigured) {
		return status.Error(codes.FailedPrecondition, err.Error())
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return status.Error(codes.NotFound, "not found")
	}
	return status.Error(codes.Internal, "internal server error")
}

func isClientActionError(err error) bool {
	return errors.Is(err, upload_doctor_photo.ErrInvalidPhoto) || errors.Is(err, upload_doctor_photo.ErrInvalidPhotoMeta)
}

func httpActionStatus(err error) int {
	if errors.Is(err, upload_doctor_photo.ErrS3NotConfigured) {
		return http.StatusServiceUnavailable
	}
	if isClientActionError(err) {
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}
