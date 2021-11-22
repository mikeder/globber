package geoip

import (
	"embed"
	"encoding/gob"
	"net"

	"github.com/pkg/errors"
)

//go:embed data.gob
var data embed.FS

type Locator struct {
	db GeoLocDB
}

func NewGeoIPLocator() (*Locator, error) {
	loc := new(Locator)
	if err := loc.load(); err != nil {
		return nil, errors.Wrap(err, "initial load allocations")
	}

	return loc, nil
}

func (l *Locator) Lookup(ip string) GeoIpRecord {
	parsedIP := net.ParseIP(ip)
	for _, alloc := range l.db.Records {
		if alloc.IPNet.Contains(parsedIP) {
			return GeoIpRecord{ // return a copy so the caller can't modify the locator records
				CountryAlpha2: alloc.CountryAlpha2,
				ParsedIP:      parsedIP,
				IPNet:         alloc.IPNet,
			}
		}
	}
	return GeoIpRecord{CountryAlpha2: "A0", ParsedIP: parsedIP}
}

func (l *Locator) load() error {
	f, err := data.Open("data.gob")
	if err != nil {
		return err
	}

	decoder := gob.NewDecoder(f)

	var tmp GeoLocDB
	err = decoder.Decode(&tmp)
	if err != nil {
		return err
	}
	if len(tmp.Records) == 0 {
		return errors.New("decoded gob had 0 records")
	}
	l.db = tmp

	return nil
}
