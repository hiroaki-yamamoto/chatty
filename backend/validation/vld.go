package validation

import (
	"context"
	"errors"
	"log"
	"net"

	"github.com/hiroaki-yamamoto/recaptcha"
	"google.golang.org/grpc/peer"
	"gopkg.in/go-playground/validator.v9"
)

// New creates a new form validation with the provided context.
// the provided context is used to obtain the remote address for recaptcha.
var New = func(
	reqCtx context.Context,
	recapSecret string,
) (vld *validator.Validate, err error) {
	if peer, ok := peer.FromContext(reqCtx); ok {
		var remoteIP string
		remoteIP, _, err = net.SplitHostPort(peer.Addr.String())
		if err != nil {
			return
		}
		recap := recaptcha.New(recapSecret)
		vld = validator.New()
		vld.RegisterValidation("recap", func(fl validator.FieldLevel) bool {
			res, err := recap.Check(remoteIP, fl.Field().String())
			if err != nil {
				log.Print(err)
				return false
			}
			return res.Success
		})
		return
	}
	err = errors.New("reqCtx doesn't seem to contain the remote peer data")
	return
}
