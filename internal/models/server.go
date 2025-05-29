package models

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
)

type Server struct {
	Err chan error         // channel for errors
	Req chan *http.Request // channel for incoming requests
	sts *Status
}

func (srv *Server) handle(w http.ResponseWriter, r *http.Request) {
	// notify request was received
	srv.Req <- r

	// GET /
	// GET /status
	if r.Method == http.MethodGet && slices.Contains([]string{"/", "/status"}, r.URL.Path) {
		// get current IP
		ip := GetPublicIP(srv.sts.Settings.Updater.PublicIPEndpoint)
		srv.sts.Current.IP = ip.IP
		srv.sts.Current.Error = ip.Error

		// note the error if any occurred during the IP lookup
		if ip.Error != nil {
			srv.Err <- ip.Error
		}

		// write updated status to the HTTP response
		srv.sts.writeHTTPResponse(w)
		return
	}

	// 404 any other requests
	err := map[string]any{}
	err["status"] = http.StatusNotFound
	err["message"] = fmt.Sprintf("path not found: %s %s", r.Method, r.URL.Path)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(err)
}

func (srv *Server) Start() error {
	// bind the handler to the root path
	http.HandleFunc("/", srv.handle)

	// start the server on address and port specified in settings
	if err := http.ListenAndServe(srv.sts.Settings.Server.Address, nil); err != nil {
		return err
	}

	return nil
}

func NewServer(sts *Status) *Server {
	return &Server{
		Err: make(chan error),
		Req: make(chan *http.Request),
		sts: sts,
	}
}
