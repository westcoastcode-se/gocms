package cms

import (
	"encoding/json"
	"github.com/westcoastcode-se/gocms/security"
	"log"
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
		log.Print("could not parse request. Reason: " + err.Error())
		returnErrorResponse(rw, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	user, err := loginService.Login(body.Username, body.Password)
	if err != nil {
		log.Print("invalid username or password for user: " + body.Username)
		returnErrorResponse(rw, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	result, err := tokenizer.UserToToken(user)
	if err != nil {
		log.Print("could not create token from user. Reason: " + err.Error())
		returnErrorResponse(rw, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	loginResponse := &LoginResponse{Name: user.Name, Token: result}
	loginResponseJson, err := json.Marshal(loginResponse)
	if err != nil {
		returnErrorResponse(rw, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	http.SetCookie(rw, &http.Cookie{Name: security.SessionKey, Value: result, Path: "/", MaxAge: 60 * 60})
	rw.WriteHeader(http.StatusOK)
	_, _ = rw.Write(loginResponseJson)
}

func logout(ctx *RequestContext) {
	http.SetCookie(ctx.Response, &http.Cookie{Name: security.SessionKey, Value: "", Path: "/", MaxAge: -1})
	http.Redirect(ctx.Response, ctx.Request, "/login?logout=true", http.StatusTemporaryRedirect)
}
