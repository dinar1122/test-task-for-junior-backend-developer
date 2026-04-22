package task

import "time"

type Status string

const (
	StatusNew        Status = "new"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

type Task struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      Status `json:"status"`

	Frequency *Frequency `json:"frequency,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s Status) Valid() bool {
	switch s {
	case StatusNew, StatusInProgress, StatusDone:
		return true
	default:
		return false
	}
}

type FrequencyType string

const (
	FrequencyNone        FrequencyType = "none"
	FrequencyDaily       FrequencyType = "daily"
	FrequencyMonthly     FrequencyType = "monthly"
	FrequencyDates       FrequencyType = "dates"
	FrequencyMonthParity FrequencyType = "month_parity"
)

type MonthParity string

const (
	MonthParityOdd  MonthParity = "odd"
	MonthParityEven MonthParity = "even"
)

type Frequency struct {
	Type FrequencyType `json:"type"`

	EveryNDays  *int         `json:"every_n_days,omitempty"`
	DaysOfMonth []int        `json:"days_of_month,omitempty"`
	Dates       []time.Time  `json:"dates,omitempty"`
	MonthParity *MonthParity `json:"month_parity,omitempty"`

	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
}
