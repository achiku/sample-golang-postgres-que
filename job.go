package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"time"

	"github.com/bgentry/que-go"
	"github.com/gocraft/dbr"
	"github.com/gocraft/dbr/dialect"
)

// PrintNameArgs print name args
type PrintNameArgs struct {
	Name string
}

// HTTPGetResponse response struct
type HTTPGetResponse struct {
	Args struct {
	} `json:"args"`
	Headers struct {
		Accept                  string `json:"Accept"`
		AcceptEncoding          string `json:"Accept-Encoding"`
		AcceptLanguage          string `json:"Accept-Language"`
		Cookie                  string `json:"Cookie"`
		Host                    string `json:"Host"`
		Referer                 string `json:"Referer"`
		UpgradeInsecureRequests string `json:"Upgrade-Insecure-Requests"`
		UserAgent               string `json:"User-Agent"`
	} `json:"headers"`
	Origin string `json:"origin"`
	URL    string `json:"url"`
}

// PrintNameJob print name
func PrintNameJob(j *que.Job) error {
	var args PrintNameArgs
	if err := json.Unmarshal(j.Args, &args); err != nil {
		return err
	}
	log.Printf("Hello, %s!\n", args.Name)
	return nil
}

// SelectItem select item
func SelectItem(j *que.Job) error {
	var n int32
	var name string
	var t time.Time
	conn := j.Conn()
	err := conn.QueryRow(`
		SELECT
		  id
		  ,name
		  ,updated_at
		FROM item
		LIMIT 1
	`).Scan(&n, &name, &t)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("id: %d name: %s updated_at: %s", n, name, t)
	return nil
}

// UpdateItemArgs updte item struct
type UpdateItemArgs struct {
	ID int32
}

// UpdateItem update item
func UpdateItem(j *que.Job) error {
	conn := j.Conn()
	tx, err := conn.Begin()
	if err != nil {
		log.Println(err)
		return err
	}
	defer tx.Rollback()

	var args UpdateItemArgs
	if err := json.Unmarshal(j.Args, &args); err != nil {
		return err
	}
	log.Printf("args id => %d", args.ID)
	res, err := tx.Exec(`
	UPDATE item
	SET updated_at = current_timestamp
	WHERE id = $1
	`, args.ID)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("rows affected: %d", res.RowsAffected())
	tx.Commit()
	return nil
}

// IncrementItemPrice increment item price by 100
func IncrementItemPrice(j *que.Job) error {
	conn := j.Conn()
	tx, err := conn.Begin()
	if err != nil {
		log.Println(err)
		return err
	}
	defer tx.Rollback()

	var args UpdateItemArgs
	if err := json.Unmarshal(j.Args, &args); err != nil {
		return err
	}
	log.Printf("args id => %d", args.ID)
	res, err := tx.Exec(`
	UPDATE item_attribute
	SET price = (
		select price from item_attribute where item_id = $1
	) + 100
	WHERE id = $1
	`, args.ID)
	if err != nil {
		log.Println(err)
		return err
	}
	tx.Commit()
	log.Printf("[increment] updated => %d", res.RowsAffected())
	return nil
}

// UpdateMultipleItem update multiple items
func UpdateMultipleItem(j *que.Job) error {
	conn := j.Conn()
	tx, err := conn.Begin()
	if err != nil {
		log.Println(err)
		return err
	}
	defer tx.Rollback()

	var args UpdateItemArgs
	if err := json.Unmarshal(j.Args, &args); err != nil {
		return err
	}
	log.Printf("args id => %d", args.ID)
	res, err := tx.Exec(`
	UPDATE item
	SET updated_at = current_timestamp
	WHERE id = $1
	`, args.ID)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("[item] rows affected: %d", res.RowsAffected())
	res, err = tx.Exec(`
	UPDATE item_attribute
	SET updated_at = current_timestamp
	WHERE item_id = $1
	`, args.ID)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("[item_attribute] rows affected: %d", res.RowsAffected())
	tx.Commit()
	return nil
}

// HTTPGetRequestItem http request item
func HTTPGetRequestItem(j *que.Job) error {
	log.Println("HTTPRequestItem started")
	conn := j.Conn()
	tx, err := conn.Begin()
	if err != nil {
		log.Println(err)
		return err
	}
	defer tx.Rollback()

	var args UpdateItemArgs
	if err := json.Unmarshal(j.Args, &args); err != nil {
		return err
	}

	log.Printf("args id => %d", args.ID)
	res, err := tx.Exec(`
	UPDATE item
	SET updated_at = current_timestamp
	WHERE id = $1
	`, args.ID)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("[item] rows affected: %d", res.RowsAffected())

	getURL := "https://httpbin.org/get"
	log.Printf("[client] target %s", getURL)
	c := NewHTTPClient()
	resp, err := c.Get(getURL)
	if err != nil {
		log.Println(err)
		return err
	}

	var result HTTPGetResponse
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return err
	}
	if err := json.Unmarshal(b, &result); err != nil {
		log.Println(err)
		return err
	}
	log.Printf("resp: %+v", result)

	return nil
}

// FailingJob job with simple failure
func FailingJob(j *que.Job) error {
	conn := j.Conn()
	var cnt int64
	err := conn.QueryRow(`
		SELECT count(*) FROM item
	`).Scan(&cnt)
	if err != nil {
		log.Println(err)
		return err
	}
	if cnt != 2 {
		return errors.New("count doesn't match")
	}
	log.Printf("[FailingJob] count: %d", cnt)
	return nil
}

func toSQL(stmt *dbr.SelectStmt) (dbr.Buffer, error) {
	buf := dbr.NewBuffer()
	if err := stmt.Build(dialect.PostgreSQL, buf); err != nil {
		log.Println(err)
		return buf, err
	}
	return buf, nil
}

// DbrQueryBuilderJob using dbr as query builder
func DbrQueryBuilderJob(j *que.Job) error {
	stmt := dbr.Select("id", "name").
		From("item").
		OrderDesc("id")
	buf, err := toSQL(stmt)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(buf.String())
	return nil
}
