package cms

import (
	"encoding/json"
	"github.com/westcoastcode-se/gocms/pkg/content"
	"github.com/westcoastcode-se/gocms/pkg/log"
	"github.com/westcoastcode-se/gocms/pkg/security"
	"net/http"
)

type CheckoutRequest struct {
	Commit string
}

type CheckoutResponse struct {
	Commit string
}

func checkout(controller content.Controller, ctx *RequestContext) {
	user := ctx.User
	rw := ctx.Response
	r := ctx.Request
	if !user.IsLoggedIn() || !user.HasRole(security.Admin) {
		returnForbidden(rw)
		return
	}

	if r.Method == http.MethodPost {
		var body CheckoutRequest
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&body)
		defer ctx.Request.Body.Close()
		if err != nil {
			log.Warnf(ctx.Request.Context(), "Could not checkout content: %e", err)
			returnErrorResponse(rw, http.StatusInternalServerError, "Could not checkout new content")
			return
		}

		err = controller.Update(ctx.Request.Context(), body.Commit)
		if err != nil {
			log.Warnf(ctx.Request.Context(), "Could not pull content from remove server: %e", err)
			returnErrorResponse(rw, http.StatusInternalServerError, "Could not checkout: "+body.Commit)
			return
		}

		returnSuccess(rw, &CheckoutResponse{body.Commit})
		return
	}

	returnNotFound(rw)
}
