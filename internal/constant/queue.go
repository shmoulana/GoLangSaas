package constant

const (
	// Type of job EMAIL, LONGJOB (consider if job will be take some minutes)
	TypeQueueEmail   = "EMAIL"
	TypeQueueLongJob = "LONGJOB"
	TypeQueueOthers  = "OTHERS"

	// two type of queue status PENDING and ONQUEUE, it just using as a tagging, so worker can check job is on queue or pending
	QueueStatusPending = "PENDING"
	QueueStatusOnQueue = "ONQUEUE"
	// QueueStatusOnProcess = "ONPROCESS"
)
