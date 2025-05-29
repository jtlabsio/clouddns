package models

import (
	"context"
	"errors"
	"strings"

	"google.golang.org/api/dns/v1"
	"google.golang.org/api/option"
)

const recordType string = "A"

type DNS struct {
	s   *Settings
	svc *dns.Service
}

func (d *DNS) Create(ip string) error {
	call := d.svc.ResourceRecordSets.Create(d.s.Google.ProjectID, d.s.Google.ManagedZone, &dns.ResourceRecordSet{
		Name:    d.s.Google.Record,
		Type:    recordType,
		Ttl:     int64(d.s.Google.TTL),
		Rrdatas: []string{ip},
	})

	if call == nil {
		return errors.New("unable to create DNS request")
	}

	if _, err := call.Do(); err != nil {
		return errors.Join(err, errors.New("unable to create DNS record"))
	}

	return nil
}

func (d *DNS) RecordExists() (bool, error) {
	call := d.svc.ResourceRecordSets.Get(d.s.Google.ProjectID, d.s.Google.ManagedZone, d.s.Google.Record, recordType)
	if call == nil {
		return false, errors.New("unable to create DNS request")
	}

	if _, err := call.Do(); err != nil {
		if strings.Contains(err.Error(), "Error 404") {
			return false, nil
		}

		return false, errors.Join(err, errors.New("unable to lookup DNS record"))
	}

	return true, nil
}

func (d *DNS) Update(ip string) error {
	call := d.svc.ResourceRecordSets.Patch(d.s.Google.ProjectID, d.s.Google.ManagedZone, d.s.Google.Record, recordType, &dns.ResourceRecordSet{
		Name:    d.s.Google.Record,
		Type:    recordType,
		Ttl:     int64(d.s.Google.TTL),
		Rrdatas: []string{ip},
	})

	if call == nil {
		return errors.New("unable to create DNS request")
	}

	if _, err := call.Do(); err != nil {
		return errors.Join(err, errors.New("unable to update DNS record"))
	}

	return nil
}

func NewDNS(s *Settings) (*DNS, error) {
	svc, err := dns.NewService(context.Background(), option.WithCredentialsFile(s.Google.CredentialsPath))
	if err != nil {
		return nil, err
	}

	return &DNS{s: s, svc: svc}, nil
}
