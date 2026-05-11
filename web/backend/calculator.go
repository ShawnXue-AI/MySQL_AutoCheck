package main

import (
	"fmt"
	"math"
	"strings"
	"time"
)

const (
	workStartHour   = 8
	workStartMin    = 30
	workEndHour     = 17
	workEndMin      = 30
	hoursPerDay     = 8.0
	coeffWork       = 1.0
	coeffOvertime   = 1.5
	coeffHoliday    = 2.0
)

type CalcResult struct {
	PersonDays    float64  `json:"person_days"`
	WorkHours     float64  `json:"work_hours"`
	OvertimeHours float64  `json:"overtime_hours"`
	HolidayHours  float64  `json:"holiday_hours"`
	Detail        []string `json:"detail"`
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

func isHoliday(t time.Time, holidays map[string]string) bool {
	weekday := t.Weekday()
	if weekday == time.Saturday || weekday == time.Sunday {
		return true
	}
	dateStr := t.Format("2006-01-02")
	_, ok := holidays[dateStr]
	return ok
}

func CalculatePersonDays(startTime, endTime time.Time, holidays map[string]string) *CalcResult {
	result := &CalcResult{
		Detail: make([]string, 0),
	}

	current := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, startTime.Location())
	endDay := time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 0, 0, 0, 0, endTime.Location())

	for !current.After(endDay) {
		dayStart := time.Date(current.Year(), current.Month(), current.Day(), 0, 0, 0, 0, current.Location())
		dayEnd := dayStart.Add(24 * time.Hour)

		effectiveStart := startTime
		if dayStart.After(effectiveStart) {
			effectiveStart = dayStart
		}
		effectiveEnd := endTime
		if dayEnd.Before(effectiveEnd) {
			effectiveEnd = dayEnd
		}

		if !effectiveStart.Before(effectiveEnd) {
			current = current.Add(24 * time.Hour)
			continue
		}

		holiday := isHoliday(current, holidays)
		workStart := time.Date(current.Year(), current.Month(), current.Day(), workStartHour, workStartMin, 0, 0, current.Location())
		workEnd := time.Date(current.Year(), current.Month(), current.Day(), workEndHour, workEndMin, 0, 0, current.Location())

		if holiday {
			hours := effectiveEnd.Sub(effectiveStart).Hours()
			holidayHours := round2(hours)
			result.HolidayHours = round2(result.HolidayHours + holidayHours)
			result.Detail = append(result.Detail, fmt.Sprintf("%s 节假日(系数%.1f): %s ~ %s = %.2f小时",
				current.Format("2006-01-02"), coeffHoliday,
				effectiveStart.Format("15:04"), effectiveEnd.Format("15:04"), holidayHours))
		} else {
			// Working hours: 08:30-17:30, coefficient 1.0
			wStart := effectiveStart
			if workStart.After(wStart) {
				wStart = workStart
			}
			wEnd := effectiveEnd
			if workEnd.Before(wEnd) {
				wEnd = workEnd
			}
			if wStart.Before(wEnd) {
				hours := wEnd.Sub(wStart).Hours()
				result.WorkHours = round2(result.WorkHours + hours)
				result.Detail = append(result.Detail, fmt.Sprintf("%s 工作时间(系数%.1f): %s ~ %s = %.2f小时",
					current.Format("2006-01-02"), coeffWork,
					wStart.Format("15:04"), wEnd.Format("15:04"), round2(hours)))
			}

			// Before working hours: 00:00-08:30, coefficient 1.5
			bStart := effectiveStart
			bEnd := effectiveEnd
			if workStart.Before(bEnd) {
				bEnd = workStart
			}
			if bStart.Before(bEnd) {
				hours := bEnd.Sub(bStart).Hours()
				result.OvertimeHours = round2(result.OvertimeHours + hours)
				result.Detail = append(result.Detail, fmt.Sprintf("%s 非工作时间(系数%.1f): %s ~ %s = %.2f小时",
					current.Format("2006-01-02"), coeffOvertime,
					bStart.Format("15:04"), bEnd.Format("15:04"), round2(hours)))
			}

			// After working hours: 17:30-24:00, coefficient 1.5
			aStart := effectiveStart
			if workEnd.After(aStart) {
				aStart = workEnd
			}
			aEnd := effectiveEnd
			if aStart.Before(aEnd) {
				hours := aEnd.Sub(aStart).Hours()
				result.OvertimeHours = round2(result.OvertimeHours + hours)
				result.Detail = append(result.Detail, fmt.Sprintf("%s 非工作时间(系数%.1f): %s ~ %s = %.2f小时",
					current.Format("2006-01-02"), coeffOvertime,
					aStart.Format("15:04"), aEnd.Format("15:04"), round2(hours)))
			}
		}

		current = current.Add(24 * time.Hour)
	}

	result.WorkHours = round2(result.WorkHours)
	result.OvertimeHours = round2(result.OvertimeHours)
	result.HolidayHours = round2(result.HolidayHours)

	totalWeighted := result.WorkHours*coeffWork + result.OvertimeHours*coeffOvertime + result.HolidayHours*coeffHoliday
	result.PersonDays = round2(totalWeighted / hoursPerDay)

	return result
}

func CalculateAndBuildRecord(customerName, workContent string, startTime, endTime time.Time, holidays map[string]string) *PersonDayRecord {
	calc := CalculatePersonDays(startTime, endTime, holidays)
	return &PersonDayRecord{
		CustomerName:  customerName,
		StartTime:     startTime,
		EndTime:       endTime,
		PersonDays:    calc.PersonDays,
		WorkHours:     calc.WorkHours,
		OvertimeHours: calc.OvertimeHours,
		HolidayHours:  calc.HolidayHours,
		WorkContent:   workContent,
		Detail:        strings.Join(calc.Detail, "\n"),
	}
}
