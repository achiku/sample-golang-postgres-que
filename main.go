package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/bgentry/que-go"
	"github.com/jackc/pgx"
)

type printNameArgs struct {
	Name string
}

func main() {
	log.Println("start que example")
	printName := func(j *que.Job) error {
		var args printNameArgs
		if err := json.Unmarshal(j.Args, &args); err != nil {
			return err
		}
		fmt.Printf("Hello %s!\n", args.Name)
		return nil
	}

	pgxcfg, err := pgx.ParseURI("postgresql://localhost/pgtest")
	if err != nil {
		log.Fatal(err)
	}

	pgxpool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:   pgxcfg,
		AfterConnect: que.PrepareStatements,
	})
	defer pgxpool.Close()
	if err != nil {
		log.Fatal(err)
	}

	qc := que.NewClient(pgxpool)
	wm := que.WorkMap{
		"PrintName": printName,
	}
	workers := que.NewWorkerPool(qc, wm, 2)
	log.Println("start workers")
	go workers.Start()

	args, err := json.Marshal(printNameArgs{Name: "achiku"})
	if err != nil {
		log.Fatal(err)
	}

	j := &que.Job{
		Type: "PrintName",
		Args: args,
	}
	log.Println("enqueue the first PrintName job")
	if err := qc.Enqueue(j); err != nil {
		log.Fatal(err)
	}

	j = &que.Job{
		Type:  "PrintName",
		RunAt: time.Now().UTC().Add(30 * time.Second),
		Args:  args,
	}
	log.Println("enqueue the delayed PrintName job")
	if err := qc.Enqueue(j); err != nil {
		log.Fatal(err)
	}

	time.Sleep(35 * time.Second)
	workers.Shutdown()
}
