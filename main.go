package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/bgentry/que-go"
)

type printNameArgs struct {
	Name string
}

func main() {
	qc, err := NewQueClient("postgresql://localhost/quetest")
	if err != nil {
		log.Fatal(err)
	}
	wm := que.WorkMap{
		"PrintName":          PrintNameJob,
		"SelectItem":         SelectItem,
		"UpdateItem":         UpdateItem,
		"UpdateMultipleItem": UpdateMultipleItem,
	}
	log.Println("create worker pool")
	workers := que.NewWorkerPool(qc, wm, 2)
	workers.Interval = 2 * time.Second
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

	sj := &que.Job{
		Type:  "SelectItem",
		RunAt: time.Now().UTC().Add(2 * time.Second),
	}
	if err := qc.Enqueue(sj); err != nil {
		log.Fatal(err)
	}

	args, err = json.Marshal(UpdateItemArgs{ID: 1})
	uj := &que.Job{
		Type: "UpdateItem",
		Args: args,
	}
	if err := qc.Enqueue(uj); err != nil {
		log.Fatal(err)
	}

	args, err = json.Marshal(UpdateItemArgs{ID: 1})
	muj := &que.Job{
		Type: "UpdateMultipleItem",
		Args: args,
	}
	if err := qc.Enqueue(muj); err != nil {
		log.Fatal(err)
	}

	log.Println("waiting for jobs to be completed")
	time.Sleep(time.Second * 10)
	log.Println("done")
}
