package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/harish-mp/web-template/db"
	"github.com/harish-mp/web-template/params"
	"github.com/harish-mp/web-template/schema"
	"github.com/harish-mp/web-template/token"
)

func extractUserInfo(body io.Reader, user *params.User) error {
	dec := json.NewDecoder(body)
	dec.DisallowUnknownFields()

	err := dec.Decode(user)
	return err
}

func extractLoginInfo(body io.Reader, user *params.Login) error {
	dec := json.NewDecoder(body)
	dec.DisallowUnknownFields()

	err := dec.Decode(user)
	return err
}
func getSecret() string {
	//todo define stronger secret
	return "TOPSECRET"
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var (
		user params.Login
		rows []db.Row
	)
	switch r.Method {
	case "POST":

		err := extractLoginInfo(r.Body, &user)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "invalid login request", http.StatusBadRequest)
		}

		fmt.Printf("%v: %v\n", user.Email, user.Pwd)

		entity := db.DbInstance.GetEntity("User")
		rows, err = db.DbInstance.GetMatched(entity, retrieveUser(), db.Match{"email": user.Email})
		if err != nil {
			fmt.Fprint(w, err)
			break
		}
		if len(rows) > 1 {
			//multiple users with same email id !!!
			fmt.Fprint(w, "multiple users have same emailId!!!")
			break
		} else {
			usr := rows[0].(schema.User)
			//todo check password
			fmt.Println(usr.ID)
			signedToken, _ := token.Create(string("987"), getSecret(), "sha256")
			fmt.Fprint(w, signedToken)
		}

	default:
		http.Error(w, "invalid method", http.StatusBadRequest)
	}
}
func assignUser(apiUser params.User) schema.User {
	return schema.User{Name: apiUser.Name, Email: apiUser.Email}
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	var (
		newUser params.User
		dbUser  schema.User
	)
	switch r.Method {
	case "POST":

		err := extractUserInfo(r.Body, &newUser)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "invalid login request", http.StatusBadRequest)
			break
		}

		userId, err := token.Verify(newUser.Token, getSecret())

		if err != nil {
			fmt.Fprint(w, err)
			break
		} else {
			fmt.Fprintf(w, "userId %s", userId)

			dbUser = assignUser(newUser)
			entity := db.DbInstance.GetEntity("User")
			_ = db.DbInstance.Insert(entity, dbUser)
			fmt.Fprintf(w, "new user add")
		}

	default:
		http.Error(w, "invalid method", http.StatusBadRequest)
	}
}

//similar function to be defined for each schema used in the web app
func retrieveUser() func(func(interface{}) error) (interface{}, error) {
	var usr schema.User
	return func(fn func(interface{}) error) (interface{}, error) {
		err := fn(&usr)
		return usr, err
	}
}
