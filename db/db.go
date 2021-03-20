package db

type EntityHandle interface{}
type DbHandle interface{}
type ObjectId interface{}
type Row interface{}
type Schema interface{}

type Match map[string]interface{}

var Dbh DbHandle
var DbInstance DB

type DB interface {
	Connect(string, string) DbHandle
	Create(string) EntityHandle
	GetEntity(string) EntityHandle
	Insert(EntityHandle, Row) error
	Remove(entity EntityHandle, id ObjectId) error
	Get(EntityHandle, func(fn func(interface{}) error) (interface{}, error), ObjectId) (Row, error)
	GetMatched(EntityHandle, func(fn func(interface{}) error) (interface{}, error), Match) ([]Row, error)
}
