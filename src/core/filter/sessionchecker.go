package filter

import (
	"context"
	beegoctx "github.com/astaxie/beego/context"
	"github.com/goharbor/harbor/src/common/utils/log"
	"github.com/goharbor/harbor/src/core/config"
	"net/http"
)

// NoSessionReqKey is the key in the context of a request to mark the request does not carry session
const NoSessionReqKey ContextValueKey = "harbor_no_session_req"

// SessionCheck is a filter to mark the requests that are not carrying a session id, it has to be registered as
// "beego.BeforeStatic" because beego will modify the request after execution of these filters, all requests will
// appear to have a session id cookie.
func SessionCheck(ctx *beegoctx.Context) {
	req := ctx.Request
	_, err := req.Cookie(config.SessionCookieName)
	if err == http.ErrNoCookie {
		ctx.Request = req.WithContext(context.WithValue(req.Context(), NoSessionReqKey, true))
		log.Debug("Mark the request as no-session")
	}
}

// ReqHasNoSession verifies if the request has been marked as "no-session", regardless if beego has modified it
func ReqHasNoSession(req *http.Request) bool {
	r, ok := req.Context().Value(NoSessionReqKey).(bool)
	return ok && r
}
