package enqueue

import (
	"github.com/hibiken/asynq"

	"shared/infrastructure/queue"
	"shared/infrastructure/queue/payloads"
)

func EnqueueGenerateReport(
	client *asynq.Client,
	payload payloads.GenerateReportPayload,
) error {
	return EnqueueTask(
		client,
		queue.TaskGenerateReport,
		payload,
		TaskEnqueueOptions{
			QueueName:      "reports",
			MaxRetries:     2,
			TimeoutSeconds: 5000,
		},
	)
}
