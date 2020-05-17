package cms

import (
	"encoding/json"
	"github.com/westcoastcode-se/gocms/content"
	"net/http"
)

func getPages(repository content.Repository, ctx *RequestContext) {
	pages := repository.GetAll()
	rw := ctx.Response
	loginResponseJson, err := json.Marshal(pages)
	if err != nil {
		returnErrorResponse(rw, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	_, _ = rw.Write(loginResponseJson)
}
