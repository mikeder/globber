package geoip

//go:generate go run gen/gen.go

import (
	"bufio"
	"bytes"
	"context"
	"encoding/csv"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
)

func WriteRecordsGob(outputFile string) error {
	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	generator := NewGenerator()
	records, err := generator.NewRecords(context.Background())
	if err != nil {
		return err
	}

	db := GeoLocDB{
		generated: time.Now(),
		Records:   records,
	}

	enc := gob.NewEncoder(f)
	if err := enc.Encode(db); err != nil {
		return err
	}

	return nil
}

type GeoLocDB struct {
	generated time.Time
	Records   []*GeoIpRecord
}

type GeoIpRecord struct {
	CountryAlpha2 string
	IPNet         net.IPNet
}

type Generator struct {
	httpClient *http.Client
}

func NewGenerator() *Generator {
	g := &Generator{
		httpClient: &http.Client{
			Timeout: time.Second * 5,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				req.URL.Opaque = req.URL.Path
				return nil
			},
		},
	}

	return g
}

// NewRecords builds a fresh slice of allocation records from the public lists, its not fast.
func (g *Generator) NewRecords(ctx context.Context) ([]*GeoIpRecord, error) {
	// arin, ripe, apnic, lacnic, afrinic - ordered for performance
	urls := []string{
		"http://ftp.arin.net/pub/stats/arin/delegated-arin-extended-latest",
		"http://ftp.ripe.net/pub/stats/ripencc/delegated-ripencc-extended-latest",
		"http://ftp.apnic.net/pub/stats/apnic/delegated-apnic-extended-latest",
		"http://ftp.lacnic.net/pub/stats/lacnic/delegated-lacnic-extended-latest",
		"http://ftp.afrinic.net/pub/stats/afrinic/delegated-afrinic-extended-latest",
	}

	// Download and filter only the lines in the files we need (allocated ipv4/ipv6 networks)
	var filtered bytes.Buffer
	var expression, err = regexp.Compile(`^.*(ipv4|ipv6).*(allocated).*$`)
	if err != nil {
		return nil, fmt.Errorf("%s: compile regex", err.Error())
	}

	for _, u := range urls {
		resp, err := g.httpClient.Get(u)
		if err != nil {
			return nil, fmt.Errorf("%s: get network list: %s", err.Error(), u)
		}
		defer resp.Body.Close()

		log.Println("processing network list: " + u)

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			if expression.Match(scanner.Bytes()) {
				_, err := filtered.Write(append(scanner.Bytes(), 10))
				if err != nil {
					return nil, fmt.Errorf("%s: write to filtered buffer", err.Error())
				}
			}
		}
	}

	// parse as a "CSV" sort of...
	reader := csv.NewReader(&filtered)
	reader.Comma = '|'

	// slice out to allocation data
	var records []*GeoIpRecord
	for {
		rec, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("%s: read allocation data", err.Error())
		}
		ip := net.ParseIP(rec[3])

		network := net.IPNet{
			IP: ip,
		}

		mask, _ := strconv.ParseInt(rec[4], 10, 0)

		// IPv6 allocations have their mask as the 4th field
		// IPv6 do not have default mask, so use it for detection.
		if ip.DefaultMask() == nil {
			network.Mask = net.CIDRMask(int(mask), 128)
		} else {
			// masks in this file are represented as "addresses in the network"
			// 32 - Log2(<number of addresses>) = the netmask.
			network.Mask = net.CIDRMask(int(32-math.Log2(float64(mask))), 32)
		}
		record := &GeoIpRecord{
			CountryAlpha2: rec[1],
			IPNet:         network,
		}
		records = append(records, record)
	}

	return records, nil
}
