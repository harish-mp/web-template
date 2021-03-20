package db

import (
	"fmt"
	"testing"

	"github.com/harish-mp/web-template/schema"
)

func TestConnectMongoDBCreateUser(t *testing.T) {

	var (
		dbVendor MongoDb
	)
	DbInstance = &dbVendor
	Dbh = DbInstance.Connect("localhost:27017", "testDb")
	DbInstance.Create("User")

}
func TestInsert(t *testing.T) {
	var (
		newUser schema.User
	)
	newUser.Email = "testuser@gmail.com"
	newUser.Pwd = "test123"
	newUser.Name = "TestUser"

	entity := DbInstance.GetEntity("User")
	err := DbInstance.Insert(entity, newUser)
	if err != nil {
		t.Fatalf("err :%s while inserting newUser", err)
	}
}

func TestGetMatchedInvalidEmail(t *testing.T) {

	entity := DbInstance.GetEntity("User")
	rows, err := DbInstance.GetMatched(entity, retrieveUser(), Match{"email": "sample@gmail.com"})
	if err != nil {
		t.Fatalf("err :%s", err)
	} else {
		for _, val := range rows {
			switch v := val.(type) {
			case schema.User:
				gotUser := schema.User(v)
				fmt.Println(gotUser)
			default:
				t.Fatalf("err :%s", err)
			}

		}
	}
}

func TestGetMatchedValidEmail(t *testing.T) {

	entity := DbInstance.GetEntity("User")
	rows, err := DbInstance.GetMatched(entity, retrieveUser(), Match{"email": "testuser@gmail.com"})
	if err != nil {
		t.Fatalf("err :%s", err)
	} else {
		for _, val := range rows {
			switch v := val.(type) {
			case schema.User:
				gotUser := schema.User(v)
				fmt.Println(gotUser)
			default:
				t.Fatalf("err :%s", err)
			}

		}
	}
}

func TestGetRemoveInvalidId(t *testing.T) {
	entity := DbInstance.GetEntity("User")
	err := DbInstance.Remove(entity, "60531e05de9ce631eace3a88")
	if err != nil {
		t.Fatalf("err :%s", err)
	} else {
		fmt.Println("user entry removed")
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
