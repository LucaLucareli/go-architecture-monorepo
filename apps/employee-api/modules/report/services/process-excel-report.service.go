package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"shared/domain/repositories"
	"shared/infrastructure/queue/payloads"
	"shared/infrastructure/report"

	"github.com/hibiken/asynq"
)

type ProcessExcelReportService struct {
	userRepo repositories.UsersRepository
}

func NewProcessExcelReportService(userRepo repositories.UsersRepository) *ProcessExcelReportService {
	return &ProcessExcelReportService{userRepo: userRepo}
}

func (s *ProcessExcelReportService) Execute(
	ctx context.Context,
	task *asynq.Task,
) error {
	var payload payloads.GenerateReportPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}

	fileName := fmt.Sprintf("report-%d.xlsx", payload.ReportID)
	writer, err := report.NewWriter(payload.ReportType, fileName)
	if err != nil {
		return err
	}
	defer writer.Close()

	writer.WriteHeader([]string{"ID", "Nome", "Documento"})

	usersChan, err := s.userRepo.FindManyToReport(ctx)
	if err != nil {
		return err
	}

	for item := range usersChan {
		if item.Err != nil {
			return item.Err
		}
		writer.WriteRow([]string{
			item.User.ID.String(),
			item.User.Name,
			item.User.Document,
		})
	}

	fmt.Printf("Relatório %s gerado com sucesso para o usuário %s\n", fileName, payload.UserID)

	// In a real scenario, you'd upload the file and delete the local copy
	// For now, let's just keep it or remove it.
	os.Remove(fileName)

	return nil
}
