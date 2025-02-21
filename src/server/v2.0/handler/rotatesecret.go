package handler

import (
	"context"
	"github.com/go-openapi/runtime/middleware"
	"github.com/goharbor/harbor/src/common/rbac"
	"github.com/goharbor/harbor/src/controller/rotatesecretkey"
	"github.com/goharbor/harbor/src/lib/errors"
	operations "github.com/goharbor/harbor/src/server/v2.0/restapi/operations/rotate_secret_key"
)

func newRotateSecretKeyAPI() *rotateSecretKeyAPI {
	return &rotateSecretKeyAPI{
		controller: rotatesecretkey.Ctl,
	}
}

type rotateSecretKeyAPI struct {
	BaseAPI
	controller rotatesecretkey.Controller
}

func (r *rotateSecretKeyAPI) RotateSecretKey(ctx context.Context, params operations.RotateSecretKeyParams) middleware.Responder {
	if err := r.RequireSystemAccess(ctx, rbac.ActionAll); err != nil {
		return r.SendError(ctx, err)
	}
	req := params.RotateSecretKeyRequest
	if len(req.NewSecretKey) != 16 || len(req.CurrentSecretKey) != 16 {
		return r.SendError(ctx, errors.BadRequestError(nil).WithMessage("Secret key must be 16 characters"))
	}
	if req.NewSecretKey == req.CurrentSecretKey {
		return r.SendError(ctx, errors.BadRequestError(nil).WithMessage("New secret key must be different from the current one"))
	}
	if err := r.controller.Rotate(ctx, req.CurrentSecretKey, req.NewSecretKey, &rotatesecretkey.Option{
		SkipOIDCSecret: req.SkipOIDCSecret,
	}); err != nil {
		return r.SendError(ctx, err)
	}
	return operations.NewRotateSecretKeyOK()
}
