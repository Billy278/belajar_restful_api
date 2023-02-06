package middleware

import (
	"belajar-golang-api/helper"
	"belajar-golang-api/model/web"
	"net/http"
)

type AuthMiddleware struct {
	Handler http.Handler
}

func NewAuthMiddlewate(handler http.Handler) *AuthMiddleware {
	return &AuthMiddleware{
		Handler: handler,
	}
}

func (middleware *AuthMiddleware) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if "RAHASIA" == request.Header.Get("X-API-key") {
		//ok
		middleware.Handler.ServeHTTP(writer, request)
	} else {
		//error
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusUnauthorized)
		webResponse := web.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: "UnAouthorized",
		}
		helper.WriteToResponseBody(writer, webResponse)
	}
}
