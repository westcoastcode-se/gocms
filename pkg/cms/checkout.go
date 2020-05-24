package cms

import (
	"encoding/json"
	"github.com/westcoastcode-se/gocms/pkg/content"
	"github.com/westcoastcode-se/gocms/pkg/security"
	"log"
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
			log.Printf("Could not checkout content. Reason: %e", err)
			returnErrorResponse(rw, http.StatusInternalServerError, "Could not checkout new content")
			return
		}

		err = controller.Update(body.Commit)
		if err != nil {
			log.Printf(`Could not pull content from remove server. Reason: %e`, err)
			returnErrorResponse(rw, http.StatusInternalServerError, "Could not checkout: "+body.Commit)
			return
		}

		returnSuccess(rw, &CheckoutResponse{body.Commit})
		return
	}

	returnNotFound(rw)
}
