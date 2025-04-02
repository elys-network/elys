package types

import (
	"fmt"
	"strings"
)

const (
	// WeekEpochID defines the identifier for weekly epochs
	WeekEpochID = "week"
	// DayEpochID defines the identifier for daily epochs
	DayEpochID = "day"
	// TwelveHoursEpochID defines the identifier for twelve hours
	TwelveHoursEpochID = "twelvehours"
	// SixHoursEpochID defines the identifier for six hours
	SixHoursEpochID = "sixhours"
	// FourHoursEpochID defines the identifier for four hours
	FourHoursEpochID = "fourhours"
	// TwoHoursEpochID defines the identifier for two hours
	TwoHoursEpochID = "twohours"
	// HourEpochID defines the identifier for hourly epochs
	HourEpochID = "hour"
	// HalfHourEpochID defines the identifier for half hour
	HalfHourEpochID = "halfhour"
	// BandEpochID defines the identifier for band epochs
	BandEpochID = "band_epoch"
	// TenSecondsEpochID defines the identifier for 10 seconds
	TenSecondsEpochID  = "tenseconds"
	TenDaysEpochID     = "ten_days"
	FiveMinutesEpochID = "five_minutes"
	EightHoursEpochID  = "eight_hours"
)

// ValidateEpochIdentifierInterface performs a stateless
// validation of the epoch ID interface.
func ValidateEpochIdentifierInterface(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if err := ValidateEpochIdentifierString(v); err != nil {
		return err
	}

	return nil
}

// ValidateEpochIdentifierInterface performs a stateless
// validation of the epoch ID.
func ValidateEpochIdentifierString(s string) error {
	s = strings.TrimSpace(s)
	if s == "" {
		return fmt.Errorf("blank epoch identifier: %s", s)
	}
	return nil
}
