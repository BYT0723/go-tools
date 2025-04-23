package web

import (
	"slices"
	"strconv"

	"github.com/BYT0723/go-tools/monitor"
	"github.com/BYT0723/go-tools/transport/httpx"
)

func CodeEqual(code int) monitor.AlertRule[httpx.Response] {
	return func(s *httpx.Response) (*monitor.Alert, bool) {
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

func CodeNotEqual(code int) monitor.AlertRule[httpx.Response] {
	return func(s *httpx.Response) (*monitor.Alert, bool) {
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

func HeaderContains(key, value string) monitor.AlertRule[httpx.Response] {
	return func(s *httpx.Response) (*monitor.Alert, bool) {
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

func HeaderNotContains(key, value string) monitor.AlertRule[httpx.Response] {
	return func(s *httpx.Response) (*monitor.Alert, bool) {
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
