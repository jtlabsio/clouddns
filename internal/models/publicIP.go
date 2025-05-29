package models

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

type PublicIP struct {
	IP    string `json:"ip"`
	Error error  `json:"error,omitempty"`
}

func GetPublicIP(epnt string) *PublicIP {
	ip := &PublicIP{}

	// create an HTTP request
	res, err := http.Get(epnt)
	if err != nil {
		ip.Error = errors.Join(err, errors.New("unable to fetch public IP address"))
		return ip
	}
	defer res.Body.Close()

	// read the response body
	bdy, err := io.ReadAll(res.Body)
	if err != nil {
		ip.Error = errors.Join(err, errors.New("unable to parse response body"))
		return ip
	}

	// set the IP
	ip.IP = string(bdy)

	// if the response is not 200, return an error
	if res.StatusCode != http.StatusOK {
		ip.Error = fmt.Errorf("unable to fetch public IP address: %s (%s)", res.Status, bdy)
	}

	return ip
}
