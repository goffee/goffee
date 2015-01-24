package controllers

import (
	"log"
	"net"
	"net/http"

	"github.com/gophergala/goffee/web/render"
	"github.com/oschwald/geoip2-golang"
)

// IP returnss the remote IP address and other info
func IP(w http.ResponseWriter, req *http.Request) {
	db, err := geoip2.Open("geoip/GeoLite2-Country.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		render.JSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	// If you are using strings that may be invalid, check that ip is not nil
	ip := net.ParseIP(host)

	record, err := db.Country(ip)

	var responseJSON map[string]string

	if err != nil || record.Country.IsoCode == "" {
		responseJSON = map[string]string{
			"IP": ip.String(),
		}
	} else {
		responseJSON = map[string]string{
			"IP":      ip.String(),
			"Country": record.Country.IsoCode,
		}
	}

	render.JSON(w, http.StatusOK, responseJSON)
}
