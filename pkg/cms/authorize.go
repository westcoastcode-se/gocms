package cms

import (
	"encoding/json"
	"github.com/westcoastcode-se/gocms/pkg/log"
	"github.com/westcoastcode-se/gocms/pkg/security"
	"net/http"
)

type LoginRequest struct {
	Username string
	Password string
}

type LoginResponse struct {
	Name  string
	Token string
}

func login(loginService security.LoginService, tokenizer security.Tokenizer, ctx *RequestContext) {
	decoder := json.NewDecoder(ctx.Request.Body)
	var body LoginRequest
	err := decoder.Decode(&body)
	defer ctx.Request.Body.Close()
	rw := ctx.Response
	if err != nil {
		log.Warnf(ctx.Request.Context(), "Could not parse request: %e", err.Error())
		returnErrorResponse(rw, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	log.Infof(ctx.Request.Context(), "Trying to logging in user: %s", body.Username)
	user, err := loginService.Login(body.Username, body.Password)
	if err != nil {
		log.Warnf(ctx.Request.Context(), "Invalid username or password for user: %e", err.Error())
		returnErrorResponse(rw, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	result, err := tokenizer.UserToToken(user)
	if err != nil {
		log.Warnf(ctx.Request.Context(), "Could not create token for user: %e", err.Error())
		returnErrorResponse(rw, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	loginResponse := &LoginResponse{Name: user.Name, Token: result}
	loginResponseJson, err := json.Marshal(loginResponse)
	if err != nil {
		returnErrorResponse(rw, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	log.Infof(ctx.Request.Context(), "User: %s has successfully logged in", body.Username)
	rw.Header().Set("Content-Type", "application/json")
	http.SetCookie(rw, &http.Cookie{Name: security.SessionKey, Value: result, Path: "/", MaxAge: 60 * 60})
	rw.WriteHeader(http.StatusOK)
	_, _ = rw.Write(loginResponseJson)
}

func logout(ctx *RequestContext) {
	http.SetCookie(ctx.Response, &http.Cookie{Name: security.SessionKey, Value: "", Path: "/", MaxAge: -1})
	http.Redirect(ctx.Response, ctx.Request, "/login?logout=true", http.StatusTemporaryRedirect)
}
