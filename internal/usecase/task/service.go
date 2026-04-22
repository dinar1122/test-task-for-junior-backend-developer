package task

import (
	"context"
	"fmt"
	"strings"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type Service struct {
	repo Repository
	now  func() time.Time
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
		now:  func() time.Time { return time.Now().UTC() },
	}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (*taskdomain.Task, error) {
	normalized, err := validateCreateInput(input)
	if err != nil {
		return nil, err
	}

	model := &taskdomain.Task{
		Title:       normalized.Title,
		Description: normalized.Description,
		Status:      normalized.Status,
		Frequency:   normalized.Frequency,
	}
	now := s.now()
	model.CreatedAt = now
	model.UpdatedAt = now

	created, err := s.repo.Create(ctx, model)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	return s.repo.GetByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, id int64, input UpdateInput) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	normalized, err := validateUpdateInput(input)
	if err != nil {
		return nil, err
	}

	model := &taskdomain.Task{
		ID:          id,
		Title:       normalized.Title,
		Description: normalized.Description,
		Status:      normalized.Status,
		Frequency:   normalized.Frequency,
		UpdatedAt:   s.now(),
	}

	updated, err := s.repo.Update(ctx, model)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	return s.repo.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]taskdomain.Task, error) {
	return s.repo.List(ctx)
}

func validateCreateInput(input CreateInput) (CreateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return CreateInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	if input.Status == "" {
		input.Status = taskdomain.StatusNew
	}

	if !input.Status.Valid() {
		return CreateInput{}, fmt.Errorf("%w: invalid status", ErrInvalidInput)
	}

	if err := validateFrequency(input.Frequency); err != nil {
		return CreateInput{}, err
	}

	return input, nil
}

func validateUpdateInput(input UpdateInput) (UpdateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return UpdateInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	if !input.Status.Valid() {
		return UpdateInput{}, fmt.Errorf("%w: invalid status", ErrInvalidInput)
	}

	if err := validateFrequency(input.Frequency); err != nil {
		return UpdateInput{}, err
	}

	return input, nil
}

func validateFrequency(f *taskdomain.Frequency) error {
	if f == nil {
		return nil
	}

	switch f.Type {
	case taskdomain.FrequencyNone:
		return nil

	case taskdomain.FrequencyDaily:
		if f.EveryNDays == nil || *f.EveryNDays <= 0 {
			return fmt.Errorf("%w: every_n_days must be greater than 0", ErrInvalidInput)
		}

	case taskdomain.FrequencyMonthly:
		if len(f.DaysOfMonth) == 0 {
			return fmt.Errorf("%w: days_of_month is required", ErrInvalidInput)
		}

		for _, day := range f.DaysOfMonth {
			if day < 1 || day > 30 {
				return fmt.Errorf("%w: days_of_month values must be between 1 and 30", ErrInvalidInput)
			}
		}

	case taskdomain.FrequencyDates:
		if len(f.Dates) == 0 {
			return fmt.Errorf("%w: dates is required", ErrInvalidInput)
		}

	case taskdomain.FrequencyMonthParity:
		if f.MonthParity == nil {
			return fmt.Errorf("%w: month_parity is required", ErrInvalidInput)
		}

		if *f.MonthParity != taskdomain.MonthParityOdd && *f.MonthParity != taskdomain.MonthParityEven {
			return fmt.Errorf("%w: month_parity must be odd or even", ErrInvalidInput)
		}

	default:
		return fmt.Errorf("%w: invalid frequency type", ErrInvalidInput)
	}

	if f.StartDate != nil && f.EndDate != nil && f.EndDate.Before(*f.StartDate) {
		return fmt.Errorf("%w: end_date must not be before start_date", ErrInvalidInput)
	}

	return nil
}
