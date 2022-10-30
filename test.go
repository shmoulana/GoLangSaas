// Original code with Dispatcher
package main

import (
	"context"
	_ "expvar"
	"flag"
	"fmt"
	_ "net/http/pprof"
	"time"

	"github.com/shmoulana/Redios/configs"
	"github.com/shmoulana/Redios/internal"
	"github.com/shmoulana/Redios/internal/service/worker"
)

func main() {
	var (
		maxWorkers   = flag.Int("max_workers", 5, "The number of workers to start")
		maxQueueSize = flag.Int("max_queue_size", 100, "The size of job queue")
	)

	flag.Parse()

	// Create the job queue.
	jobQueue := make(chan worker.Job, *maxQueueSize)

	// Start the dispatcher.
	dispatcher := worker.NewDispatcher(jobQueue, *maxWorkers)
	dispatcher.Run()

	configs.Init()
	conf := configs.Get()

	factory := internal.Transport{}
	queueService := factory.GetQueueService(*conf)

	// Define all job
	job := &worker.EmailJob{
		EmailService: factory.GetEmailService(*conf),
		QueueService: factory.GetQueueService(*conf),
	}
	ctx := context.Background()

	for {
		fmt.Println("Fetching jobs")

		keys, err := queueService.AllQueues(context.Background())
		if err != nil {
			panic(err)
		}

		fmt.Printf("Found %d jobs", len(keys))
		fmt.Println("")

		for _, key := range keys {
			data, err := queueService.GetQueue(ctx, key)
			if err != nil {
				fmt.Println(err)
				continue
			}

			if data.Status == "ONQUEUE" {
				continue
			}

			if data.TypeQueue == "EMAIL" {
				err := queueService.UpdateStatus(ctx, key, "ONQUEUE")
				if err != nil {
					panic(err)
				}

				job.EmailService.Init()
				jobQueue <- worker.Job{
					Key:     key,
					JobFunc: job,
				}
			}
		}

		// fmt.Println("All jobs done")

		time.Sleep(10 * time.Second)
	}

	// // Start the HTTP handler.
	// http.HandleFunc("/work", func(w http.ResponseWriter, r *http.Request) {
	// 	requestHandler(w, r, jobQueue)
	// })
	// log.Fatal(http.ListenAndServe(":"+*port, nil))
}
