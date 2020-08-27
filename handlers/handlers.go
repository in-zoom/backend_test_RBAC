package handlers

import (
	"backend_test_RBAC/data"
	"backend_test_RBAC/db"
	"backend_test_RBAC/hashing"
	"backend_test_RBAC/model"
	"backend_test_RBAC/validation"
	"backend_test_RBAC/verification"
	"encoding/json"
	"github.com/alexedwards/scs"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func RegistrationUser(db *db.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		addedUser := data.RegisterUser{}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		err := json.NewDecoder(r.Body).Decode(&addedUser)
		if err != nil {
			responseError(w, http.StatusBadRequest, err)
			return
		}

		err = db.СreateTable()
		if err != nil {
			responseError(w, http.StatusInternalServerError, err)
			return
		}

		err = db.CreateAdministrator(data.AdminName, data.Role, data.AdminPass)
		if err != nil {
			responseError(w, http.StatusInternalServerError, err)
			return
		}

		resultNameUser, err := validation.ValidateNameUser(addedUser.Name, db)
		if err != nil {
			responseError(w, http.StatusBadRequest, err)
			return
		}

		hashPasswordUser, err := hashing.ValidateAndHashPasswordUser(addedUser.Password)
		if err != nil {
			responseError(w, http.StatusBadRequest, err)
			return
		}

		err = db.AddNewUser(resultNameUser, data.UserRole, hashPasswordUser)
		if err != nil {
			responseError(w, http.StatusInternalServerError, err)
			return
		}
		responceOk(w, http.StatusOK, "Вы успешно зарегистрированы")
	}
}

func Login(users model.User, session scs.SessionManager, db *db.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		auth := data.Auth{}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		err := json.NewDecoder(r.Body).Decode(&auth)
		if err != nil {
			responseError(w, http.StatusBadRequest, err)
			return
		}

		err = verification.VerificationLogin(auth.Username, auth.Password, db)
		if err != nil {
			responseError(w, http.StatusBadRequest, err)
			return
		}

		user, err := users.FindByName(auth.Username, db)
		if err != nil {
			responseError(w, http.StatusBadRequest, err)
			return
		}

		if err := session.RenewToken(r.Context()); err != nil {
			responseError(w, http.StatusInternalServerError, err)
			return
		}

		session.Put(r.Context(), "userID", user.ID)
		session.Put(r.Context(), "role", user.Role)

		responceOk(w, http.StatusOK, "Вы успешно авторизованы")
	}

}

func Logout(session scs.SessionManager) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		if err := session.Destroy(r.Context()); err != nil {
			responseError(w, http.StatusInternalServerError, err)
			return
		}
		responceOk(w, http.StatusOK, "Вы вышли из учетной записи")
	}
}

func Foo() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		responceOk(w, http.StatusOK, "Доступ к ресурсу foo получен")
	}
}

func Bar() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		responceOk(w, http.StatusOK, "Доступ к ресурсу bar получен")
	}
}

func Sigma() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		responceOk(w, http.StatusOK, "Доступ к ресурсу sigma получен")
	}
}

func responceOk(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	m := data.Message{OkMessage: message, Status: code}
	json.NewEncoder(w).Encode(m)
}

func responseError(w http.ResponseWriter, code int, err error) {
	w.WriteHeader(code)
	errMessage := data.ErrMessage{Message: err.Error()}
	json.NewEncoder(w).Encode(errMessage)
}
