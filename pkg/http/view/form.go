package view

import (
	"net/http"
	"net/url"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/monoculum/formam/v3"
)

func DecodeRequestForm(r *http.Request, dst any) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	return DecodeForm(r.PostForm, dst)
}

func DecodeForm(form url.Values, dst any) error {
	// We are instantiating a new formam.Decoder instance every time we call this function
	// to avoid a known bug that reset IgnoreUnknownKeys to false after a certain number of calls
	return defaultDecoder().Decode(form, dst)
}

func WrapValidationError(err error) error {
	if err == nil {
		return make(validation.Errors)
	}

	return err
}

func defaultDecoder() *formam.Decoder {
	return formam.
		NewDecoder(&formam.DecoderOptions{
			IgnoreUnknownKeys: true,
		})
}
