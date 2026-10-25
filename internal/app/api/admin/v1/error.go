package adminv1

import (
	"errors"
	"net/http"

	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/doctors/action/upload_doctor_photo"
	"github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/users/service/password"
	"github.com/nagorneva/nagorneva_reviews/internal/pkg/domain_error"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func rpcError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, password.ErrPasswordTooShort) || errors.Is(err, upload_doctor_photo.ErrInvalidPhoto) || errors.Is(err, upload_doctor_photo.ErrInvalidPhotoMeta) {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	if errors.Is(err, upload_doctor_photo.ErrS3NotConfigured) {
		return status.Error(codes.FailedPrecondition, err.Error())
	}
	if errors.Is(err, domain_error.ErrNotFound) {
		return status.Error(codes.NotFound, "not found")
	}
	if errors.Is(err, domain_error.ErrConflict) {
		return status.Error(codes.AlreadyExists, err.Error())
	}
	if errors.Is(err, domain_error.ErrReference) {
		return status.Error(codes.FailedPrecondition, err.Error())
	}
	if errors.Is(err, domain_error.ErrInvalid) {
		return status.Error(codes.InvalidArgument, err.Error())
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
	if errors.Is(err, domain_error.ErrNotFound) {
		return http.StatusNotFound
	}
	if isClientActionError(err) {
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}
