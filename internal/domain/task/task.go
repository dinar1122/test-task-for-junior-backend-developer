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

func (t *Task) OccursOn(day time.Time) bool {
	if t == nil {
		return false
	}

	if t.Frequency == nil || t.Frequency.Type == "" || t.Frequency.Type == FrequencyNone {
		return true
	}

	day = normalizeDate(day)

	if !isWithinRange(day, t.Frequency.StartDate, t.Frequency.EndDate) {
		return false
	}

	switch t.Frequency.Type {
	case FrequencyDaily:
		if t.Frequency.EveryNDays == nil || *t.Frequency.EveryNDays <= 0 {
			return false
		}

		if t.Frequency.StartDate == nil {
			return true
		}

		start := normalizeDate(*t.Frequency.StartDate)
		diffDays := int(day.Sub(start).Hours() / 24)

		return diffDays >= 0 && diffDays%*t.Frequency.EveryNDays == 0

	case FrequencyMonthly:
		currentDay := day.Day()
		for _, d := range t.Frequency.DaysOfMonth {
			if d == currentDay {
				return true
			}
		}
		return false

	case FrequencyDates:
		for _, d := range t.Frequency.Dates {
			if normalizeDate(d).Equal(day) {
				return true
			}
		}
		return false

	case FrequencyMonthParity:
		if t.Frequency.MonthParity == nil {
			return false
		}

		switch *t.Frequency.MonthParity {
		case MonthParityEven:
			return day.Day()%2 == 0
		case MonthParityOdd:
			return day.Day()%2 != 0
		default:
			return false
		}

	default:
		return false
	}
}

func normalizeDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

func isWithinRange(day time.Time, startDate, endDate *time.Time) bool {
	if startDate != nil {
		start := normalizeDate(*startDate)
		if day.Before(start) {
			return false
		}
	}

	if endDate != nil {
		end := normalizeDate(*endDate)
		if day.After(end) {
			return false
		}
	}

	return true
}
