package web

import (
	"slices"
	"strconv"

	"github.com/BYT0723/go-tools/monitor"
)

func CodeEqual(code int) monitor.AlertRule[Statistics] {
	return func(s *Statistics) (*monitor.Alert, bool) {
		if s.Code == code {
			return monitor.NewAlert(
				monitor.SeverityError,
				monitor.SourceAlertRule,
				"status code: "+strconv.Itoa(code),
				s,
			), true
		}
		return nil, false
	}
}

func CodeNotEqual(code int) monitor.AlertRule[Statistics] {
	return func(s *Statistics) (*monitor.Alert, bool) {
		if s.Code != code {
			return monitor.NewAlert(
				monitor.SeverityError,
				monitor.SourceAlertRule,
				"status code: "+strconv.Itoa(code),
				s,
			), true
		}
		return nil, false
	}
}

func HeaderContains(key, value string) monitor.AlertRule[Statistics] {
	return func(s *Statistics) (*monitor.Alert, bool) {
		if slices.Contains(s.Header[key], value) {
			return monitor.NewAlert(
				monitor.SeverityError,
				monitor.SourceAlertRule,
				"header: "+key+"="+value,
				s,
			), true
		}
		return nil, false
	}
}

func HeaderNotContains(key, value string) monitor.AlertRule[Statistics] {
	return func(s *Statistics) (*monitor.Alert, bool) {
		if !slices.Contains(s.Header[key], value) {
			return monitor.NewAlert(
				monitor.SeverityError,
				monitor.SourceAlertRule,
				"header: "+key+"="+value,
				s,
			), true
		}
		return nil, false
	}
}
