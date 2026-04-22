package handlers

import (
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type frequencyDTO struct {
	Type        taskdomain.FrequencyType `json:"type"`
	EveryNDays  *int                     `json:"every_n_days,omitempty"`
	DaysOfMonth []int                    `json:"days_of_month,omitempty"`
	Dates       []time.Time              `json:"dates,omitempty"`
	MonthParity *taskdomain.MonthParity  `json:"month_parity,omitempty"`
	StartDate   *time.Time               `json:"start_date,omitempty"`
	EndDate     *time.Time               `json:"end_date,omitempty"`
}

type taskMutationDTO struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      taskdomain.Status `json:"status"`
	Frequency   *frequencyDTO     `json:"frequency,omitempty"`
}

type taskDTO struct {
	ID          int64             `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      taskdomain.Status `json:"status"`
	Frequency   *frequencyDTO     `json:"frequency,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

func newTaskDTO(task *taskdomain.Task) taskDTO {
	return taskDTO{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		Frequency:   newFrequencyDTO(task.Frequency),
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}

func newFrequencyDTO(f *taskdomain.Frequency) *frequencyDTO {
	if f == nil {
		return nil
	}

	return &frequencyDTO{
		Type:        f.Type,
		EveryNDays:  f.EveryNDays,
		DaysOfMonth: f.DaysOfMonth,
		Dates:       f.Dates,
		MonthParity: f.MonthParity,
		StartDate:   f.StartDate,
		EndDate:     f.EndDate,
	}
}

func (d *frequencyDTO) toDomain() *taskdomain.Frequency {
	if d == nil {
		return nil
	}

	return &taskdomain.Frequency{
		Type:        d.Type,
		EveryNDays:  d.EveryNDays,
		DaysOfMonth: d.DaysOfMonth,
		Dates:       d.Dates,
		MonthParity: d.MonthParity,
		StartDate:   d.StartDate,
		EndDate:     d.EndDate,
	}
}
