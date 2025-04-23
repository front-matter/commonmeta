package geonames

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	rtcache "github.com/ArthurHlt/go-roundtripper-cache"
	"github.com/front-matter/commonmeta/fileutils"
)

// Countr represents Geonames Country Info
type Country struct {
	GeonameID       int64   `json:"geoname_id"`
	Name            string  `json:"name"`
	ISO             string  `json:"iso"`
	ISO3            string  `json:"iso3"`
	ISONumeric      string  `json:"isonumeric"`
	FIPS            string  `json:"fips"`
	ContinentCode   string  `json:"continent_code"`
	Capital         string  `json:"capital"`
	AreaKm2         float64 `json:"areakm2"`
	Population      uint64  `json:"population"`
	TLD             string  `json:"tld"`
	CurrencyCode    string  `json:"currency_code"`
	CurrencyName    string  `json:"currency_name"`
	Phone           string  `json:"phone"`
	PostalCodeRegex string  `json:"postalcode_regex"`
	Languages       string  `json:"languages"`
	Neighbours      string  `json:"neighbours"`
}

// Feature represents a single Geonames administrative object - for example, a city
type Feature struct {
	GeonameID        int64     `json:"geoname_id"`        // geoname_id        : integer id of record in geonames database
	Name             string    `json:"name"`              // name              : name of geographical point (utf8) varchar(200)
	ASCIIName        string    `json:"ascii_name"`        // asciiname         : name of geographical point in plain ascii characters, varchar(200)
	AlternateNames   []string  `json:"alternate_names"`   // alternate names   : alternate names, comma separated, ascii names automatically transliterated, convenience attribute from alternatename table, varchar(10000)
	Latitude         float64   `json:"latitude"`          // latitude          : latitude in decimal degrees (wgs84)
	Longitude        float64   `json:"longitude"`         // longitude         : longitude in decimal degrees (wgs84)
	Class            string    `json:"class"`             // feature class     : see http://www.geonames.org/export/codes.html, char(1)
	Code             string    `json:"code"`              // feature code      : see http://www.geonames.org/export/codes.html, varchar(10)l, varchar(10)
	CountryCode      string    `json:"country_code"`      // country code      : ISO-3166 2-letter country code, 2 characters
	Cc2              string    `json:"cc2"`               // cc2               : alternate country codes, comma separated, ISO-3166 2-letter country code, 200 characters
	Admin1Code       string    `json:"admin1_code"`       // admin1 code       : fipscode (subject to change to iso code), see exceptions below, see file admin1Codes.txt for display names of this code; varchar(20)ames of this code; varchar(20)
	Admin2Code       string    `json:"admin2_code"`       // admin2 code       : code for the second administrative division, a county in the US, see file admin2Codes.txt; varchar(80)
	Admin3Code       string    `json:"admin3_code"`       // admin3 code       : code for third level administrative division, varchar(20)
	Admin4Code       string    `json:"admin4_code"`       // admin4 code       : code for fourth level administrative division, varchar(20)el administrative division, varchar(20)
	Population       *int      `json:"population"`        // population        : bigint (8 byte int))
	Elevation        *int      `json:"elevation"`         // elevation         : in meters, integer
	Dem              int       `json:"dem"`               // dem               : digital elevation model, srtm3 or gtopo30, average elevation of 3''x3'' (ca 90mx90m) or 30''x30'' (ca 900mx900m) area in meters, integer. srtm processed by cgiar/ciat.elevation of 3''x3'' (ca 90mx90m) or 30''x30'' (ca 900mx900m) area in meters, integer. srtm processed by cgiar/ciat.
	TimeZone         string    `json:"time_zone"`         // timezone          : the timezone id (see file timeZone.txt) varchar(40)r(40)
	ModificationDate time.Time `json:"modification_date"` // modification date : date of last modification in yyyy-MM-dd formatModificationDate time.Time
}

const (
	geonamesURL    = "http://download.geonames.org/export/dump/"
	countryInfoURL = "countryInfo.txt"
	commentSymbol  = byte('#')
	newLineSymbol  = byte('\n')
	delimSymbol    = byte('\t')
	boolTrue       = "1"
)

// LoadGeonamesCountries loads countries from geonamesnames
func LoadGeonamesCountries() (map[string]Country, error) {
	url := geonamesURL + countryInfoURL
	bytes, err := fileutils.DownloadFile(url)
	if err != nil {
		return nil, fmt.Errorf("error downloading country info: %w", err)
	}
	countries := make(map[string]Country)

	parse(bytes, 0, func(fields [][]byte) bool {
		fmt.Println(len(fields))
		if len(fields) < 18 {
			return true
		}

		geonameID, err := strconv.ParseInt(string(fields[0]), 10, 64)
		if err != nil {
			return true
		}

		population, _ := strconv.ParseUint(string(fields[7]), 10, 64)
		area, _ := strconv.ParseFloat(string(fields[6]), 64)

		country := Country{
			GeonameID:       geonameID,
			ISO:             string(fields[0]),
			ISO3:            string(fields[1]),
			ISONumeric:      string(fields[2]),
			FIPS:            string(fields[3]),
			Name:            string(fields[4]),
			Capital:         string(fields[5]),
			AreaKm2:         area,
			Population:      population,
			ContinentCode:   string(fields[8]),
			TLD:             string(fields[9]),
			CurrencyCode:    string(fields[10]),
			CurrencyName:    string(fields[11]),
			Phone:           string(fields[12]),
			PostalCodeRegex: string(fields[13]),
			Languages:       string(fields[15]),
			Neighbours:      string(fields[17]),
		}

		// add country to map, using ISO code as key
		countries[string(fields[0])] = country
		return true
	})

	return countries, nil
}

// LoadGeonamesCities loads cities from geonameses
func LoadGeonamesCities() (map[int]Feature, error) {
	httpClient := &http.Client{
		Timeout:   60 * time.Second,
		Transport: rtcache.NewRoundTripperCache(24 * time.Hour),
	}

	url := "http://download.geonames.org/export/dump/cities15000.zip"
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	bodyReader := bytes.NewReader(body)
	zipReader, err := zip.NewReader(bodyReader, resp.ContentLength)
	if err != nil {
		return nil, fmt.Errorf("error reading zip: %w", err)
	}

	var cityFile *zip.File
	for _, file := range zipReader.File {
		if file.Name == "cities15000.txt" {
			cityFile = file
			break
		}
	}

	if cityFile == nil {
		return nil, fmt.Errorf("cities15000.txt not found in zip")
	}

	fileReader, err := cityFile.Open()
	if err != nil {
		return nil, fmt.Errorf("error opening city file: %w", err)
	}
	defer fileReader.Close()

	data, err := io.ReadAll(fileReader)
	if err != nil {
		return nil, fmt.Errorf("error reading city data: %w", err)
	}

	features := make(map[int]Feature)

	parse(data, 0, func(fields [][]byte) bool {
		if len(fields) < 19 {
			return true
		}

		geonameID, err := strconv.Atoi(string(fields[0]))
		if err != nil {
			return true
		}

		var population *int
		if popVal, err := strconv.Atoi(string(fields[14])); err == nil {
			population = &popVal
		}

		var elevation *int
		if elevVal, err := strconv.Atoi(string(fields[15])); err == nil {
			elevation = &elevVal
		}

		dem, _ := strconv.Atoi(string(fields[16]))
		modDate, _ := time.Parse("2006-01-02", string(fields[18]))

		feature := Feature{
			GeonameID:        int64(geonameID),
			Name:             string(fields[1]),
			ASCIIName:        string(fields[2]),
			AlternateNames:   strings.Split(string(fields[3]), ","),
			Latitude:         parseFloat(fields[4]),
			Longitude:        parseFloat(fields[5]),
			Class:            string(fields[6]),
			Code:             string(fields[7]),
			CountryCode:      string(fields[8]),
			Cc2:              string(fields[9]),
			Admin1Code:       string(fields[10]),
			Admin2Code:       string(fields[11]),
			Admin3Code:       string(fields[12]),
			Admin4Code:       string(fields[13]),
			Population:       population,
			Elevation:        elevation,
			Dem:              dem,
			TimeZone:         string(fields[17]),
			ModificationDate: modDate,
		}

		features[geonameID] = feature
		return true
	})

	return features, nil
}

// parse parses the data and calls the callback for each line
func parse(data []byte, skipLines int, callback func([][]byte) bool) {
	var (
		lineStart  = 0
		fieldStart = 0
		line       = 0
		fields     [][]byte
	)

	if len(data) == 0 {
		return
	}

	for i, b := range data {
		// skip comment lines
		if i == lineStart && b == commentSymbol {
			for i < len(data) && data[i] != newLineSymbol {
				i++
			}
			lineStart = i + 1
			fieldStart = lineStart
			continue
		}

		if b == newLineSymbol || i == len(data)-1 {
			if i == len(data)-1 && b != newLineSymbol {
				fields = append(fields, data[fieldStart:i+1])
			}

			line++
			if line > skipLines && len(fields) > 0 {
				if !callback(fields) {
					return
				}
			}

			lineStart = i + 1
			fieldStart = lineStart
			fields = fields[:0]
			continue
		}

		// delimiter
		if b == delimSymbol {
			fields = append(fields, data[fieldStart:i])
			fieldStart = i + 1
		}
	}
}

// parseFloat parses a byte slice to a float64
func parseFloat(data []byte) float64 {
	f, _ := strconv.ParseFloat(string(data), 64)
	return f
}
