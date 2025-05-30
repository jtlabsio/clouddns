package models

import (
	"encoding/json"
	"net/http"
	"time"
)

const version string = "0.1.1"

type Status struct {
	Current struct {
		Error     error     `json:"error,omitempty"`
		IP        string    `json:"ip"`
		Timestamp time.Time `json:"timestamp"`
	} `json:"current"`
	Latest struct {
		Duration  string    `json:"duration"`
		Error     error     `json:"error,omitempty"`
		IP        string    `json:"ip"`
		Timestamp time.Time `json:"timestamp"`
	} `json:"previousRun"`
	Settings *Settings `json:"settings"`
	Version  string
}

func (sts *Status) writeHTTPResponse(w http.ResponseWriter) {
	sts.Current.Timestamp = time.Now()

	w.Header().Set("Content-Type", "application/json")

	if sts.Current.Error == nil {
		w.WriteHeader(http.StatusOK)
	}

	if sts.Current.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(sts)
}

func (sts *Status) UpdateLatest(strt time.Time, ip string, err ...error) {
	sts.Latest.Duration = time.Since(strt).String()

	if len(err) > 0 {
		sts.Latest.Error = err[0]
	}

	if len(err) == 0 {
		sts.Latest.Error = nil
	}

	sts.Latest.IP = ip
	sts.Latest.Timestamp = time.Now()
}

func NewStatus(s *Settings) *Status {
	return &Status{
		Settings: s,
		Version:  version,
	}
}
