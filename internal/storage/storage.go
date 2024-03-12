package storage

import (
	"database/sql"
	"github.com/ajugalushkin/url-shortener-version2/internal/model"
	_ "github.com/mattn/go-sqlite3"
)

const file string = "internal/db/store.db"
const create string = `CREATE TABLE IF NOT EXISTS urls (
	key	TEXT NOT NULL PRIMARY key,
	url	TEXT NOT NULL 
)`

type Storage struct {
	db *sql.DB
}

func NewStorage() (*Storage, error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(create); err != nil {
		return nil, err
	}
	return &Storage{
		db: db,
	}, nil
}

func (c *Storage) Insert(urlData model.URLData) (int, error) {
	res, err := c.db.Exec("insert into urls (key,url) values ($1,$2)", urlData.Key, urlData.URL)
	if err != nil {
		return 0, err
	}

	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return 0, err
	}
	return int(id), nil
}

func (c *Storage) Retrieve(id string) (model.URLData, error) {
	row := c.db.QueryRow("select * from urls where key = $1", id)

	urlData := model.URLData{}
	var err error
	if err = row.Scan(&urlData.Key, &urlData.URL); err == sql.ErrNoRows {
		return model.URLData{}, err
	}
	return urlData, err
}

func (c *Storage) RetrieveByURL(url string) (model.URLData, error) {
	row := c.db.QueryRow("select * from urls where url = $1", url)

	urlData := model.URLData{}
	var err error
	if err = row.Scan(&urlData.Key, &urlData.URL); err == sql.ErrNoRows {
		return model.URLData{}, err
	}
	return urlData, err
}
