package httperrors

import (
	"net/http"
)

func Code(err error) int {
	switch err.Error() {
	case IDNegative,
		ParamIDInvalid,
		NoField,
		NameNull,
		PhoneInvalid,
		BloodInvalid:
		return http.StatusBadRequest
	case PatientFailed,
		CreateFailed,
		UpdateFailed,
		DeleteFailed:
		return http.StatusInternalServerError
	case PatientNotFound:
		return http.StatusNotFound
	default:
		return http.StatusBadGateway
	}
}
