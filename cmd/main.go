package main

import (
	"strings"
	"time"

	"go.jtlabs.io/clouddns/internal/models"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.jtlabs.io/settings"
)

func loadSettings() (*models.Settings, error) {
	s := &models.Settings{}

	// configure settings options
	opts := settings.Options().
		SetArgsMap(map[string]string{
			"--logging-level": "Logging.Level",
		}).
		SetBasePath("./settings/settings.yaml").
		SetEnvOverride("ENV", "GO_ENV").
		SetEnvSearchPaths("./settings").
		SetVarsMap(map[string]string{
			"GOOGLE_CREDENTIALS_PATH":    "Google.CredentialsPath",
			"GOOGLE_MANAGED_ZONE":        "Google.ManagedZone",
			"GOOGLE_PROJECT_ID":          "Google.ProjectID",
			"GOOGLE_RECORD":              "Google.Record",
			"GOOGLE_TTL":                 "Google.TTL",
			"LOGGING_LEVEL":              "Logging.Level",
			"SERVER_ADDRESS":             "Server.Address",
			"UPDATER_INTERVAL":           "Updater.Interval",
			"UPDATER_PUBLIC_IP_ENDPOINT": "Updater.PublicIPEndpoint",
		})

	// read settings from file and environment variables
	if err := settings.Gather(opts, s); err != nil {
		return nil, err
	}

	// set global log level
	zerolog.SetGlobalLevel(s.GlobalLogLevel())

	return s, nil
}

func updateIP(sts *models.Status, dns *models.DNS) {
	n := time.Now()

	log.Trace().
		Str("endpoint", sts.Settings.Updater.PublicIPEndpoint).
		Msg("Retrieving public IP")

	ip := models.GetPublicIP(sts.Settings.Updater.PublicIPEndpoint)
	if ip.Error != nil {
		log.Error().Err(ip.Error).Msg("Failed to get public IP")
		sts.UpdateLatest(n, ip.IP, ip.Error)
		return
	}

	log.Trace().
		Str("ip", ip.IP).
		Str("zone", sts.Settings.Google.ManagedZone).
		Msg("Preparing to update IP for zone")

	exists, err := dns.RecordExists()
	if err != nil {
		log.Error().Err(err).Msg("Failed to check if DNS record exists")
		sts.UpdateLatest(n, ip.IP, err)
		return
	}

	// check to see if a new record should be created
	if !exists {
		log.Trace().
			Str("projectID", sts.Settings.Google.ProjectID).
			Str("record", sts.Settings.Google.Record).
			Str("zone", sts.Settings.Google.ManagedZone).
			Msg("DNS record does not exist, creating it")

		if err := dns.Create(ip.IP); err != nil {
			log.Error().Err(err).Msg("Failed to create DNS record")
			sts.UpdateLatest(n, ip.IP, err)
			return
		}

		log.Debug().
			Str("projectID", sts.Settings.Google.ProjectID).
			Str("record", sts.Settings.Google.Record).
			Str("zone", sts.Settings.Google.ManagedZone).
			Str("ip", ip.IP).
			Int("ttl", sts.Settings.Google.TTL).
			Msg("Successfully created DNS record")

		sts.UpdateLatest(n, ip.IP)
		return
	}

	// update the existing record...
	if err := dns.Update(ip.IP); err != nil {
		log.Error().Err(err).Msg("Failed to update DNS record")
		sts.UpdateLatest(n, ip.IP, err)
		return
	}

	log.Debug().
		Str("projectID", sts.Settings.Google.ProjectID).
		Str("record", sts.Settings.Google.Record).
		Str("zone", sts.Settings.Google.ManagedZone).
		Str("ip", ip.IP).
		Int("ttl", sts.Settings.Google.TTL).
		Msg("Successfully updated DNS record")

	sts.UpdateLatest(n, ip.IP)
}

func main() {
	// load settings from file and environment variables
	s, err := loadSettings()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load settings")
	}

	log.Debug().Interface("settings", s).Msg("Loaded settings")

	// create a DNS client for managing DNS records
	dns, err := models.NewDNS(s)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create DNS client")
	}

	// create a status for tracking runtime details
	sts := models.NewStatus(s)

	// create a new server for reporting runtime details
	srv := models.NewServer(sts)

	// asynchronously call the updater on a schedule
	stp := make(chan bool)
	go func() {
		for {
			updateIP(sts, dns)
			select {
			case <-time.After(s.Updater.Interval):
			case <-stp:
				return
			}
		}
	}()

	// setup request logging
	go func() {
		for req := range srv.Req {
			// ignore favicon.ico requests
			if strings.Contains(req.URL.Path, "/favicon.ico") {
				continue
			}

			log.Debug().
				Str("method", req.Method).
				Str("url", req.URL.Path).
				Msg("Received request")
		}
	}()

	// setup error logging
	go func() {
		for err := range srv.Err {
			log.Error().
				Err(err).
				Send()
		}
	}()

	// start the server
	log.Trace().
		Str("address", s.Server.Address).
		Msg("Starting server")

	if err := srv.Start(); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
