// Original code with Dispatcher
package main

import (
	"context"
	_ "expvar"
	"fmt"
	_ "net/http/pprof"
	"time"

	"github.com/shmoulana/Redios/configs"
	"github.com/shmoulana/Redios/internal"
	"github.com/shmoulana/Redios/internal/service/worker"
)

func main() {
	// var max

	configs.Init()
	conf := configs.Get()

	// Create the job queue.
	// Max Queue Size Default - 100
	jobQueue := make(chan worker.Job, 100)

	// Start the dispatcher.
	// Max Workers Default - 3
	dispatcher := worker.NewDispatcher(jobQueue, 3)
	dispatcher.Run()

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

		time.Sleep(10 * time.Second)
	}
}
