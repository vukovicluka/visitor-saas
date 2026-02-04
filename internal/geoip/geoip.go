package geoip

import (
	"log"
	"net"

	"github.com/oschwald/geoip2-golang"
)

type Resolver struct {
	db *geoip2.Reader
}

func New(path string) *Resolver {
	if path == "" {
		log.Println("GeoIP: no database path configured, country detection disabled")
		return &Resolver{}
	}

	db, err := geoip2.Open(path)
	if err != nil {
		log.Printf("GeoIP: failed to open %s: %v (country detection disabled)", path, err)
		return &Resolver{}
	}

	log.Printf("GeoIP: loaded database from %s", path)
	return &Resolver{db: db}
}

func (r *Resolver) Country(ipStr string) string {
	if r.db == nil {
		return ""
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return ""
	}

	record, err := r.db.Country(ip)
	if err != nil {
		return ""
	}

	return record.Country.IsoCode
}

func (r *Resolver) Close() {
	if r.db != nil {
		r.db.Close()
	}
}