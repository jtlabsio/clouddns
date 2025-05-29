package models

import (
	"time"

	"github.com/rs/zerolog"
)

type Settings struct {
	Google struct {
		CredentialsPath string `yaml:"credentialsPath" json:"credentialsPath"`
		ManagedZone     string `yaml:"managedZone" json:"managedZone"`
		ProjectID       string `yaml:"projectID" json:"projectID"`
		Record          string `yaml:"record" json:"record"`
		TTL             int    `yaml:"ttl" json:"ttl"`
	} `yaml:"google" json:"google"`
	Logging struct {
		Level string `yaml:"level" json:"level"`
	} `yaml:"logging" json:"logging"`
	Server struct {
		Address string `yaml:"address" json:"address"`
	} `yaml:"server" json:"server"`
	Updater struct {
		Interval         time.Duration `yaml:"interval" json:"interval"`
		PublicIPEndpoint string        `yaml:"publicIPEndpoint" json:"publicIPEndpoint"`
	} `yaml:"updater" json:"updater"`
}

func (s *Settings) GlobalLogLevel() zerolog.Level {
	switch s.Logging.Level {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	default:
		return zerolog.InfoLevel
	}
}
