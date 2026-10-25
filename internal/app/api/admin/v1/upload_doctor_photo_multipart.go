package adminv1

import (
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	uploaddoctorphoto "github.com/nagorneva/nagorneva_reviews/internal/modules/admin/entities/doctors/action/upload_doctor_photo"
	"github.com/nagorneva/nagorneva_reviews/internal/pkg/httpx"
)

// UploadDoctorPhotoMultipart preserves the existing browser multipart endpoint.
func (s *Service) UploadDoctorPhotoMultipart(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid doctor id")
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, uploaddoctorphoto.MaxPhotoSize+(1<<20))
	if err := r.ParseMultipartForm(uploaddoctorphoto.MaxPhotoSize); err != nil {
		httpx.Error(w, http.StatusBadRequest, "photo must be a multipart file up to 10 MiB")
		return
	}
	defer func() { _ = r.MultipartForm.RemoveAll() }()
	file, header, err := r.FormFile("photo")
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "photo field is required")
		return
	}
	defer func() { _ = file.Close() }()
	photo, err := io.ReadAll(io.LimitReader(file, uploaddoctorphoto.MaxPhotoSize+1))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "cannot read photo")
		return
	}
	v, err := s.module.Actions.Doctors.UploadDoctorPhoto.Execute(r.Context(), id, photo, header.Filename, header.Header.Get("Content-Type"))
	if err != nil {
		httpx.Error(w, httpActionStatus(err), err.Error())
		return
	}
	httpx.JSON(w, http.StatusOK, doctor(v))
}
