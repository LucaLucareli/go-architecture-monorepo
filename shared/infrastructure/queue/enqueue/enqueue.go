package enqueue

import (
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
)

type TaskEnqueueOptions struct {
	QueueName      string
	MaxRetries     int
	TimeoutSeconds int
	DelaySeconds   int
	TaskID         string
}

const DefaultQueueName = "default"

func EnqueueTask[T any](
	asynqClient *asynq.Client,
	taskType string,
	payload T,
	opts TaskEnqueueOptions,
) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	maxRetries := opts.MaxRetries
	if maxRetries == 0 {
		maxRetries = 3
	}

	asynqOptions := []asynq.Option{}

	queueName := opts.QueueName
	if queueName == "" {
		queueName = DefaultQueueName
	}
	asynqOptions = append(asynqOptions, asynq.Queue(queueName))

	asynqOptions = append(asynqOptions, asynq.MaxRetry(maxRetries))

	if opts.TimeoutSeconds > 0 {
		asynqOptions = append(
			asynqOptions,
			asynq.Timeout(time.Duration(opts.TimeoutSeconds)*time.Second),
		)
	}

	if opts.DelaySeconds > 0 {
		asynqOptions = append(
			asynqOptions,
			asynq.ProcessIn(time.Duration(opts.DelaySeconds)*time.Second),
		)
	}

	if opts.TaskID != "" {
		asynqOptions = append(
			asynqOptions,
			asynq.TaskID(opts.TaskID),
		)
	}

	task := asynq.NewTask(taskType, payloadBytes, asynqOptions...)

	_, err = asynqClient.Enqueue(task)
	return err
}
