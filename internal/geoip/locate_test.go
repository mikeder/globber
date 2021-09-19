package geoip_test

import (
	"testing"

	"github.com/mikeder/globber/internal/geoip"
)

func TestLocate(t *testing.T) {
	loc, err := geoip.NewGeoIPLocator()
	checkErrFail(t, err)

	record := loc.Lookup("100.15.156.135")
	if record.CountryAlpha2 == ""{
		t.Log("country empty")
		t.Fail()
	}
}

func checkErrFail(t *testing.T, err error) {
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}
