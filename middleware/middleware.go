package middleware

import (
	"backend_test_RBAC/db"
	"backend_test_RBAC/model"
	"errors"
	"github.com/alexedwards/scs"
	"github.com/casbin/casbin"
	"log"
	"net/http"
)

func Authorizer(e *casbin.Enforcer, users model.User, session *scs.SessionManager, db *db.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			role := session.GetString(r.Context(), "role")

			if role == "" {
				role = "anonymous"
			}

			if role == "user" {
				uid := session.GetInt(r.Context(), "userID")

				exists, err := users.Exists(uid, db)
				if err != nil {
					writeError(http.StatusInternalServerError, "ERROR", w, err)
				} else if !*exists {
					writeError(http.StatusForbidden, "ЗАПРЕЩЕНО", w, errors.New("Пользователь не существует"))
					return
				}
			}

			res, err := e.Enforce(role, r.URL.Path, r.Method)
			if err != nil {
				writeError(http.StatusInternalServerError, "ERROR", w, err)
				return
			}
			if res {
				next.ServeHTTP(w, r)
			} else {
				writeError(http.StatusForbidden, "ЗАПРЕЩЕНО", w, errors.New("Несанкционированный доступ"))
				return
			}
		}

		return http.HandlerFunc(fn)
	}
}

func writeError(status int, message string, w http.ResponseWriter, err error) {
	log.Print("ERROR: ", err.Error())
	w.WriteHeader(status)
	w.Write([]byte(message))
}
