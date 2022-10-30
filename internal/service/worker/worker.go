package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/shmoulana/Redios/internal/service/email"
	"github.com/shmoulana/Redios/internal/service/queue"
)

// JobFunc holds the attributes needed to perform unit of work.
type JobFunc interface {
	// Name  string
	// Delay time.Duration
	Run(key string)
}

type Job struct {
	Key     string
	JobFunc JobFunc
}

// NewWorker creates takes a numeric id and a channel w/ worker pool.
func NewWorker(id int, workerPool chan chan Job) Worker {
	return Worker{
		id:         id,
		jobQueue:   make(chan Job),
		workerPool: workerPool,
		quitChan:   make(chan bool),
	}
}

type Worker struct {
	id         int
	jobQueue   chan Job
	workerPool chan chan Job
	quitChan   chan bool
}

func (w Worker) start() {
	go func() {
		for {
			// Add my jobQueue to the worker pool.
			w.workerPool <- w.jobQueue

			select {
			case job := <-w.jobQueue:
				// Dispatcher has added a job to my jobQueue.
				// fmt.Printf("worker%d: started %s, blocking for %f seconds\n", w.id, job.Name, job.Delay.Seconds())
				job.JobFunc.Run(job.Key)
				time.Sleep(1)
				// fmt.Printf("worker%d: completed %s!\n", w.id, job.Name)
			case <-w.quitChan:
				// We have been asked to stop.
				fmt.Printf("worker%d stopping\n", w.id)
				return
			}
		}
	}()
}

func (w Worker) stop() {
	go func() {
		w.quitChan <- true
	}()
}

// NewDispatcher creates, and returns a new Dispatcher object.
func NewDispatcher(jobQueue chan Job, maxWorkers int) *Dispatcher {
	workerPool := make(chan chan Job, maxWorkers)

	return &Dispatcher{
		jobQueue:   jobQueue,
		maxWorkers: maxWorkers,
		workerPool: workerPool,
	}
}

type Dispatcher struct {
	workerPool chan chan Job
	maxWorkers int
	jobQueue   chan Job
}

func (d *Dispatcher) Run() {
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(i+1, d.workerPool)
		worker.start()
	}

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-d.jobQueue:
			go func() {
				// fmt.Printf("fetching workerJobQueue for: %s\n", job.Name)
				workerJobQueue := <-d.workerPool
				// fmt.Printf("adding %s to workerJobQueue\n", job.Name)
				workerJobQueue <- job
			}()
		}
	}
}

type EmailJob struct {
	EmailService email.EmailService
	QueueService queue.QueueService
}

// type EmailJobValue struct {
// }

func (j *EmailJob) Run(key string) {
	fmt.Println("Sending Test email")
	ctx := context.Background()

	data, err := j.QueueService.GetQueue(ctx, key)
	if err != nil {
		fmt.Println(err)
		return
	}

	var value email.EmailJobValue

	err = json.Unmarshal([]byte(data.Value), &value)
	if err != nil {
		fmt.Print(err)
	}

	err = j.EmailService.SendNow(ctx, value.To, value.Msg)
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Println("Done sending test email")

	err = j.delete(ctx, key)
	if err != nil {
		fmt.Printf("Err: %+v", err)
		return
	}

	return
}

func (j *EmailJob) delete(ctx context.Context, key string) error {
	return j.QueueService.DeleteQueue(ctx, key)
}
