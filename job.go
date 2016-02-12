package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/bgentry/que-go"
)

// PrintNameArgs print name args
type PrintNameArgs struct {
	Name string
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
		j.Done()
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
	SET updated_at = now() 
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
	SET updated_at = now()
	WHERE id = $1
	`, args.ID)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("[item] rows affected: %d", res.RowsAffected())
	res, err = tx.Exec(`
	UPDATE item_attribute
	SET updated_at = now()
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
