package database

import "github.com/txfs19260817/url-shortener/model"

var DB Database

// Database is an interface contains a set of CRUD methods to be implemented.
// Thus, we can call these methods in a database agnostic manner.
type Database interface {
	CreateUrl(*model.Url) error
	ReadUrl(string) (*model.Url, error)
}
