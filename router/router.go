package router

import (
	"net/http"

	"github.com/harish-mp/web-template/controller"
)

func AddRoutes() {
	http.HandleFunc("/login", controller.LoginHandler)
	http.HandleFunc("/adduser", controller.AddUser)
}
