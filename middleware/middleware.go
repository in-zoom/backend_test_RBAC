package middleware

import (
	"backend_test_RBAC/db"
	"backend_test_RBAC/handlers"
	"backend_test_RBAC/model"
	"errors"
	"github.com/alexedwards/scs/v2"
	"github.com/casbin/casbin/v2"
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
					handlers.ResponseError(w, http.StatusInternalServerError, errors.New("ошибка"))
				} else if !*exists {
					handlers.ResponseError(w, http.StatusForbidden, errors.New("в доступе отказано"))
					return
				}
			}

			res, err := e.Enforce(role, r.URL.Path, r.Method)
			if err != nil {
				handlers.ResponseError(w, http.StatusInternalServerError, errors.New("ошибка"))
				return
			}
			if res {
				next.ServeHTTP(w, r)
			} else {
				handlers.ResponseError(w, http.StatusForbidden, errors.New("в доступе отказано"))
				return
			}
		}

		return http.HandlerFunc(fn)
	}
}
