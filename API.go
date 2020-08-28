package main

import (
	"backend_test_RBAC/db"
	"backend_test_RBAC/handlers"
	"backend_test_RBAC/middleware"
	"backend_test_RBAC/model"
	"github.com/alexedwards/scs/v2"
	"github.com/casbin/casbin/v2"
	"github.com/julienschmidt/httprouter"
	"github.com/subosito/gotenv"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	authEnforcer, err := casbin.NewEnforcer("./auth_model.conf", "./policy.csv")
	if err != nil {
		log.Fatal(err)
	}
	sessionManager := scs.New()
	sessionManager.Lifetime = 10 * time.Minute

	initialize()

	db, err := db.OpenDB()
	if err != nil {
		log.Fatal(err)
	}

	router := httprouter.New()
	router.POST("/registration", handlers.RegistrationUser(db))
	router.POST("/login", handlers.Login(model.User{}, *sessionManager, db))
	router.POST("/logout", handlers.Logout(*sessionManager))
	router.GET("/user/foo", handlers.Foo())
	router.GET("/user/bar", handlers.Bar())
	router.GET("/admin/sigma", handlers.Sigma())
	port := os.Getenv("PORT")
	log.Fatal(http.ListenAndServe(":"+port, sessionManager.LoadAndSave(middleware.Authorizer(authEnforcer, model.User{}, sessionManager, db)(router))))
}

func initialize() {
	if err := gotenv.Load(); err != nil {
		log.Print("Файл .env не найден")
	}
}
