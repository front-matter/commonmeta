package ror

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/front-matter/commonmeta/utils"
	"github.com/jszwec/csvutil"
	"gopkg.in/yaml.v3"
)

type RORCSV struct {
	ID                           string `csv:"id"`
	Name                         string `csv:"name"`
	Types                        string `csv:"types"`
	Status                       string `csv:"status"`
	Links                        string `csv:"links,omitempty"`
	Aliases                      string `csv:"aliases,omitempty"`
	Labels                       string `csv:"labels,omitempty"`
	Acronyms                     string `csv:"acronyms,omitempty"`
	WikipediaURL                 string `csv:"wikipedia_url,omitempty"`
	Established                  string `csv:"established,omitempty"`
	Latitude                     string `csv:"addresses[0].lat"`
	Longitude                    string `csv:"addresses[0].lng"`
	Place                        string `csv:"addresses[0].geonames_city.name"`
	GeonamesID                   string `csv:"addresses[0].geonames_city.id"`
	CountrySubdivisionName       string `csv:"addresses[0].geonames_city.geonames_admin1.name,omitempty"`
	CountrySubdivisionCode       string `csv:"addresses[0].geonames_city.geonames_admin1.code,omitempty"`
	CountryCode                  string `csv:"country.country_code"`
	CountryName                  string `csv:"country.country_name"`
	ExternalIDsGRIDPreferred     string `csv:"external_ids.GRID.preferred,omitempty"`
	ExternalIDsGRIDAll           string `csv:"external_ids.GRID.all,omitempty"`
	ExternalIDsISNIPreferred     string `csv:"external_ids.ISNI.preferred,omitempty"`
	ExternalIDsISNIAll           string `csv:"external_ids.ISNI.all,omitempty"`
	ExternalIDsFundrefPreferred  string `csv:"external_ids.FundRef.preferred,omitempty"`
	ExternalIDsFundrefAll        string `csv:"external_ids.FundRef.all,omitempty"`
	ExternalIDsWikidataPreferred string `csv:"external_ids.Wikidata.preferred,omitempty"`
	ExternalIDsWikidataAll       string `csv:"external_ids.Wikidata.all,omitempty"`
	Relationships                string `csv:"relationships,omitempty"`
}

// InvenioRDM represents the ROR metadata record in InvenioRDM format.
type InvenioRDM struct {
	Acronym     string       `json:"acronym,omitempty" yaml:"acronym,omitempty"`
	ID          string       `json:"id"`
	Country     string       `json:"country,omitempty" yaml:"country,omitempty"`
	Identifiers []Identifier `json:"identifiers"`
	Name        string       `json:"name"`
	Title       Title        `json:"title"`
}

type Identifier struct {
	Identifier string `json:"identifier"`
	Scheme     string `json:"scheme"`
}

type Title struct {
	Aa string `json:"aa,omitempty" yaml:"aa,omitempty"` // Afar
	Af string `json:"af,omitempty" yaml:"af,omitempty"` // Afrikaans
	Am string `json:"am,omitempty" yaml:"am,omitempty"` // Amharic
	Ar string `json:"ar,omitempty" yaml:"ar,omitempty"` // Arabic
	As string `json:"as,omitempty" yaml:"as,omitempty"` // Assamese
	Az string `json:"az,omitempty" yaml:"az,omitempty"` // Azerbaijani
	Ba string `json:"ba,omitempty" yaml:"ba,omitempty"` // Bashkir
	Be string `json:"be,omitempty" yaml:"be,omitempty"` // Belgian
	Bg string `json:"bg,omitempty" yaml:"bg,omitempty"` // Bulgarian
	Bi string `json:"bi,omitempty" yaml:"bi,omitempty"` // Bislama
	Bn string `json:"bn,omitempty" yaml:"bn,omitempty"` // Bengali
	Bs string `json:"bs,omitempty" yaml:"bs,omitempty"` // Bosnian
	Ca string `json:"ca,omitempty" yaml:"ca,omitempty"` // Catalan
	Ch string `json:"ch,omitempty" yaml:"ch,omitempty"` // Chamorro
	Co string `json:"co,omitempty" yaml:"co,omitempty"` // Corsican
	Cs string `json:"cs,omitempty" yaml:"cs,omitempty"` // Czech
	Cu string `json:"cu,omitempty" yaml:"cu,omitempty"` // Church Slavic
	Cy string `json:"cy,omitempty" yaml:"cy,omitempty"` // Welsh
	Da string `json:"da,omitempty" yaml:"da,omitempty"` // Danish
	De string `json:"de,omitempty" yaml:"de,omitempty"` // German
	Dv string `json:"dv,omitempty" yaml:"dv,omitempty"` // Divehi
	Dz string `json:"dz,omitempty" yaml:"dz,omitempty"` // Dzongkha
	El string `json:"el,omitempty" yaml:"el,omitempty"` // Greek
	En string `json:"en,omitempty" yaml:"en,omitempty"` // English
	Es string `json:"es,omitempty" yaml:"es,omitempty"` // Spanish
	Et string `json:"et,omitempty" yaml:"et,omitempty"` // Estonian
	Eu string `json:"eu,omitempty" yaml:"eu,omitempty"` // Basque
	Fa string `json:"fa,omitempty" yaml:"fa,omitempty"` // Persian
	Fi string `json:"fi,omitempty" yaml:"fi,omitempty"` // Finnish
	Fo string `json:"fo,omitempty" yaml:"fo,omitempty"` // Faroese
	Fr string `json:"fr,omitempty" yaml:"fr,omitempty"` // French
	Fy string `json:"fy,omitempty" yaml:"fy,omitempty"` // Frisian
	Ga string `json:"ga,omitempty" yaml:"ga,omitempty"` // Irish
	Gd string `json:"gd,omitempty" yaml:"gd,omitempty"` // Scottish Gaelic
	Gl string `json:"gl,omitempty" yaml:"gl,omitempty"` // Galician
	Gu string `json:"gu,omitempty" yaml:"gu,omitempty"` // Gujarati
	Gv string `json:"gv,omitempty" yaml:"gv,omitempty"` // Manx
	Ha string `json:"ha,omitempty" yaml:"ha,omitempty"` // Hausa
	He string `json:"he,omitempty" yaml:"he,omitempty"` // Hebrew
	Hi string `json:"hi,omitempty" yaml:"hi,omitempty"` // Hindi
	Hr string `json:"hr,omitempty" yaml:"hr,omitempty"` // Croatian
	Ht string `json:"ht,omitempty" yaml:"ht,omitempty"` // Haitian
	Hu string `json:"hu,omitempty" yaml:"hu,omitempty"` // Hungarian
	Hy string `json:"hy,omitempty" yaml:"hy,omitempty"` // Armenian
	Id string `json:"id,omitempty" yaml:"id,omitempty"` // Indonesian
	Is string `json:"is,omitempty" yaml:"is,omitempty"` // Icelandic
	It string `json:"it,omitempty" yaml:"it,omitempty"` // Italian
	Iu string `json:"iu,omitempty" yaml:"iu,omitempty"` // Inuktitut
	Ja string `json:"ja,omitempty" yaml:"ja,omitempty"` // Japanese
	Jv string `json:"jv,omitempty" yaml:"jv,omitempty"` // Javanese
	Ka string `json:"ka,omitempty" yaml:"ka,omitempty"` // Georgian
	Kg string `json:"kg,omitempty" yaml:"kg,omitempty"` // Kongo
	Ki string `json:"ki,omitempty" yaml:"ki,omitempty"` // Kikuyu
	Kk string `json:"kk,omitempty" yaml:"kk,omitempty"` // Kazakh
	Kl string `json:"kl,omitempty" yaml:"kl,omitempty"` // Greenlandic
	Km string `json:"km,omitempty" yaml:"km,omitempty"` // Khmer
	Kn string `json:"kn,omitempty" yaml:"kn,omitempty"` // Kannada
	Ko string `json:"ko,omitempty" yaml:"ko,omitempty"` // Korean
	Kr string `json:"kr,omitempty" yaml:"kr,omitempty"` // Kanuri
	Ku string `json:"ku,omitempty" yaml:"ku,omitempty"` // Kurdish
	Ky string `json:"ky,omitempty" yaml:"ky,omitempty"` // Kyrgyz
	La string `json:"la,omitempty" yaml:"la,omitempty"` // Latin
	Lb string `json:"lb,omitempty" yaml:"lb,omitempty"` // Luxembourgish
	Lo string `json:"lo,omitempty" yaml:"lo,omitempty"` // Lao
	Lt string `json:"lt,omitempty" yaml:"lt,omitempty"` // Lithuanian
	Lu string `json:"lu,omitempty" yaml:"lu,omitempty"` // Luba-Katanga
	Lv string `json:"lv,omitempty" yaml:"lv,omitempty"` // Latvian
	Mg string `json:"mg,omitempty" yaml:"mg,omitempty"` // Malagasy
	Mi string `json:"mi,omitempty" yaml:"mi,omitempty"` // Maori
	Mk string `json:"mk,omitempty" yaml:"mk,omitempty"` // Macedonian
	Ml string `json:"ml,omitempty" yaml:"ml,omitempty"` // Malayalam
	Mn string `json:"mn,omitempty" yaml:"mn,omitempty"` // Mongolian
	Mr string `json:"mr,omitempty" yaml:"mr,omitempty"` // Marathi
	Ms string `json:"ms,omitempty" yaml:"ms,omitempty"` // Malay
	Mt string `json:"mt,omitempty" yaml:"mt,omitempty"` // Maltese
	My string `json:"my,omitempty" yaml:"my,omitempty"` // Burmese
	Na string `json:"na,omitempty" yaml:"na,omitempty"` // Nauru
	Nb string `json:"nb,omitempty" yaml:"nb,omitempty"` // Norwegian Bokmål
	Ne string `json:"ne,omitempty" yaml:"ne,omitempty"` // Nepali
	Nl string `json:"nl,omitempty" yaml:"nl,omitempty"` // Dutch
	Nn string `json:"nn,omitempty" yaml:"nn,omitempty"` // Norwegian Nynorsk
	No string `json:"no,omitempty" yaml:"no,omitempty"` // Norwegian
	Oc string `json:"oc,omitempty" yaml:"oc,omitempty"` // Occitan
	Om string `json:"om,omitempty" yaml:"om,omitempty"` // Oromo
	Or string `json:"or,omitempty" yaml:"or,omitempty"` // Oriya
	Pa string `json:"pa,omitempty" yaml:"pa,omitempty"` // Punjabi
	Pl string `json:"pl,omitempty" yaml:"pl,omitempty"` // Polish
	Ps string `json:"ps,omitempty" yaml:"ps,omitempty"` // Pashto
	Pt string `json:"pt,omitempty" yaml:"pt,omitempty"` // Portuguese
	Rm string `json:"rm,omitempty" yaml:"rm,omitempty"` // Romansh
	Ro string `json:"ro,omitempty" yaml:"ro,omitempty"` // Romanian
	Ru string `json:"ru,omitempty" yaml:"ru,omitempty"` // Russian
	Rw string `json:"rw,omitempty" yaml:"rw,omitempty"` // Kinyarwanda
	Sa string `json:"sa,omitempty" yaml:"sa,omitempty"` // Sanskrit
	Sd string `json:"sd,omitempty" yaml:"sd,omitempty"` // Sindhi
	Se string `json:"se,omitempty" yaml:"se,omitempty"` // Northern Sami
	Sh string `json:"sh,omitempty" yaml:"sh,omitempty"` // Serbo-Croatian
	Si string `json:"si,omitempty" yaml:"si,omitempty"` // Sinhalese
	Sk string `json:"sk,omitempty" yaml:"sk,omitempty"` // Slovak
	Sl string `json:"sl,omitempty" yaml:"sl,omitempty"` // Slovenian
	Sm string `json:"sm,omitempty" yaml:"sm,omitempty"` // Samoan
	So string `json:"so,omitempty" yaml:"so,omitempty"` // Somali
	Sq string `json:"sq,omitempty" yaml:"sq,omitempty"` // Albanian
	Sr string `json:"sr,omitempty" yaml:"sr,omitempty"` // Serbian
	St string `json:"st,omitempty" yaml:"st,omitempty"` // Southern Sotho
	Sv string `json:"sv,omitempty" yaml:"sv,omitempty"` // Swedish
	Sw string `json:"sw,omitempty" yaml:"sw,omitempty"` // Swahili
	Ta string `json:"ta,omitempty" yaml:"ta,omitempty"` // Tamil
	Te string `json:"te,omitempty" yaml:"te,omitempty"` // Telugu
	Tg string `json:"tg,omitempty" yaml:"tg,omitempty"` // Tajik
	Th string `json:"th,omitempty" yaml:"th,omitempty"` // Thai
	Ti string `json:"ti,omitempty" yaml:"ti,omitempty"` // Tigrinya
	Tk string `json:"tk,omitempty" yaml:"tk,omitempty"` // Turkmen
	Tl string `json:"tl,omitempty" yaml:"tl,omitempty"` // Tagalog
	Tr string `json:"tr,omitempty" yaml:"tr,omitempty"` // Turkish
	Tt string `json:"tt,omitempty" yaml:"tt,omitempty"` // Tatar
	Ug string `json:"ug,omitempty" yaml:"ug,omitempty"` // Uighur
	Uk string `json:"uk,omitempty" yaml:"uk,omitempty"` // Ukrainian
	Ur string `json:"ur,omitempty" yaml:"ur,omitempty"` // Urdu
	Uz string `json:"uz,omitempty" yaml:"uz,omitempty"` // Uzbek
	Vi string `json:"vi,omitempty" yaml:"vi,omitempty"` // Vietnamese
	Xh string `json:"xh,omitempty" yaml:"xh,omitempty"` // Xhosa
	Yo string `json:"yo,omitempty" yaml:"yo,omitempty"` // Yoruba
	Zh string `json:"zh,omitempty" yaml:"zh,omitempty"` // Chinese
	Zu string `json:"zu,omitempty" yaml:"zu,omitempty"` // Zulu
}

// Convert converts ROR metadata into InvenioRDM format.
func ConvertInvenioRDM(data ROR) (InvenioRDM, error) {
	var inveniordm InvenioRDM

	id, _ := utils.ValidateROR(data.ID)
	inveniordm.ID = id
	if len(data.Locations) > 0 {
		inveniordm.Country = data.Locations[0].GeonamesDetails.CountryCode
	}
	inveniordm.Identifiers = []Identifier{
		{
			Identifier: id,
			Scheme:     "ror",
		},
	}
	for _, name := range data.Names {
		if slices.Contains(name.Types, "ror_display") {
			inveniordm.Name = name.Value
		} else if slices.Contains(name.Types, "acronym") && name.Value != "" {
			inveniordm.Acronym = name.Value
		}
	}
	inveniordm.Title = GetTitle(data.Names)
	return inveniordm, nil
}

// ConvertRORCSV converts ROR metadata into RORCSV format.
func ConvertRORCSV(data ROR) (RORCSV, error) {
	var rorcsv RORCSV
	var acronyms, aliases, labels, types, child, parent, related []string

	rorcsv.ID = data.ID
	for _, name := range data.Names {
		if slices.Contains(name.Types, "ror_display") {
			rorcsv.Name = name.Value
		} else if slices.Contains(name.Types, "acronym") && name.Value != "" {
			acronyms = append(acronyms, name.Value)
		} else if slices.Contains(name.Types, "alias") {
			aliases = append(aliases, name.Value)
		} else if slices.Contains(name.Types, "label") {
			if name.Lang != "" {
				labels = append(labels, fmt.Sprintf("%s: %s", name.Lang, name.Value))
			} else {
				labels = append(labels, name.Value)
			}
		}
	}
	for _, type_ := range data.Types {
		types = append(types, type_)
	}
	rorcsv.Types = strings.Join(slices.Compact(types), "; ")
	rorcsv.Status = data.Status
	for _, link := range data.Links {
		if link.Type == "website" {
			rorcsv.Links = link.Value
		} else if link.Type == "wikipedia" {
			rorcsv.WikipediaURL = link.Value
		}
	}
	rorcsv.Aliases = strings.Join(aliases, "; ")
	rorcsv.Labels = strings.Join(labels, "; ")
	rorcsv.Acronyms = strings.Join(acronyms, "; ")
	if data.Established != 0 {
		rorcsv.Established = strconv.Itoa(data.Established)
	}
	rorcsv.Latitude = fmt.Sprintf("%f", data.Locations[0].GeonamesDetails.Lat)
	rorcsv.Longitude = fmt.Sprintf("%f", data.Locations[0].GeonamesDetails.Lng)
	rorcsv.Place = data.Locations[0].GeonamesDetails.Name
	rorcsv.GeonamesID = strconv.Itoa(data.Locations[0].GeonamesID)
	rorcsv.CountrySubdivisionName = data.Locations[0].GeonamesDetails.CountrySubdivisionName
	rorcsv.CountrySubdivisionCode = data.Locations[0].GeonamesDetails.CountrySubdivisionCode
	rorcsv.CountryCode = data.Locations[0].GeonamesDetails.CountryCode
	rorcsv.CountryName = data.Locations[0].GeonamesDetails.CountryName
	for _, ext := range data.ExternalIDs {
		if ext.Type == "grid" {
			rorcsv.ExternalIDsGRIDPreferred = ext.Preferred
			rorcsv.ExternalIDsGRIDAll = strings.Join(ext.All, ";")
		} else if ext.Type == "isni" {
			rorcsv.ExternalIDsISNIPreferred = ext.Preferred
			rorcsv.ExternalIDsISNIAll = strings.Join(ext.All, ";")
		} else if ext.Type == "fundref" {
			rorcsv.ExternalIDsFundrefPreferred = ext.Preferred
			rorcsv.ExternalIDsFundrefAll = strings.Join(ext.All, ";")
		} else if ext.Type == "wikidata" {
			rorcsv.ExternalIDsWikidataPreferred = ext.Preferred
			rorcsv.ExternalIDsWikidataAll = strings.Join(ext.All, ";")
		}
	}

	for _, relation := range data.Relationships {
		if relation.Type == "child" {
			child = append(child, relation.ID)
		} else if relation.Type == "parent" {
			parent = append(parent, relation.ID)
		} else if relation.Type == "related" {
			related = append(related, relation.ID)
		}
	}
	if len(child) > 0 {
		rorcsv.Relationships += "Child: " + strings.Join(child, ", ")
	}
	if len(parent) > 0 {
		rorcsv.Relationships += "Parent: " + strings.Join(parent, ", ")
	}
	if len(related) > 0 {
		rorcsv.Relationships += "Related: " + strings.Join(related, ", ")
	}
	return rorcsv, nil
}

// Write writes ROR metadata.
func Write(data ROR) ([]byte, error) {
	var err error
	var output []byte

	output, err = json.Marshal(data)
	return output, err
}

// WriteAll writes a list of ROR metadata, optionally filtered by type and/or country.
func WriteAll(list []ROR, extension string) ([]byte, error) {
	var err error
	var output []byte

	switch extension {
	case ".yaml":
		output, err = yaml.Marshal(list)
	case ".json":
		output, err = json.Marshal(list)
	case ".jsonl":
		buffer := &bytes.Buffer{}
		encoder := json.NewEncoder(buffer)
		for _, item := range list {
			err = encoder.Encode(item)
			if err != nil {
				fmt.Println(err)
			}
		}
		output = buffer.Bytes()
	case ".csv":
		var rorcsvList []RORCSV
		// convert ROR to RORCSV, a custom lossy mapping to CSV
		for _, item := range list {
			rorcsv, err := ConvertRORCSV(item)
			if err != nil {
				fmt.Println(err)
			}
			rorcsvList = append(rorcsvList, rorcsv)
		}
		output, err = csvutil.Marshal(rorcsvList)
		if err != nil {
			fmt.Println(err, "csvutil.Marshal")
		}
	case ".sql":
		buffer := &bytes.Buffer{}

		// Create a TABLE definition for ROR organizations optimized for SQLite
		tableDef := `-- ROR Organizations SQL Schema
-- This schema is optimized for SQLite and includes indices for faster queries
DROP TABLE IF EXISTS organizations;
CREATE TABLE organizations (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    status TEXT NOT NULL,
    established INTEGER,
    types JSON,
    names JSON,
    country_code TEXT,
    country_name TEXT,
    latitude REAL,
    longitude REAL,
    city TEXT,
    wikipedia_url TEXT,
    website_url TEXT,
    external_ids JSON,
    relationships JSON,
    created_at TEXT,
    updated_at TEXT
);

-- Indices for faster queries (SQLite syntax)
CREATE INDEX idx_organizations_name ON organizations(name);
CREATE INDEX idx_organizations_country ON organizations(country_code);
CREATE INDEX idx_organizations_types ON organizations(json_extract(types, '$'));
CREATE INDEX idx_organizations_external_ids ON organizations(json_extract(external_ids, '$.grid.preferred'));
`
		buffer.WriteString(tableDef)
		buffer.WriteString("BEGIN TRANSACTION;\n\n")

		for _, item := range list {
			var status, website, wikipedia string
			var lat, lng float64
			var countryCode, countryName, cityName string
			var established int

			status = item.Status
			established = item.Established

			for _, link := range item.Links {
				if link.Type == "website" {
					website = link.Value
				} else if link.Type == "wikipedia" {
					wikipedia = link.Value
				}
			}

			if len(item.Locations) > 0 {
				lat = item.Locations[0].GeonamesDetails.Lat
				lng = item.Locations[0].GeonamesDetails.Lng
				countryCode = item.Locations[0].GeonamesDetails.CountryCode
				countryName = item.Locations[0].GeonamesDetails.CountryName
				cityName = item.Locations[0].GeonamesDetails.Name
			}

			typesJSON, _ := json.Marshal(item.Types)
			namesJSON, _ := json.Marshal(item.Names)

			externalIDs := make(map[string]map[string]interface{})
			for _, extID := range item.ExternalIDs {
				externalIDs[extID.Type] = map[string]interface{}{
					"preferred": extID.Preferred,
					"all":       extID.All,
				}
			}
			externalIDsJSON, _ := json.Marshal(externalIDs)
			relationshipsJSON, _ := json.Marshal(item.Relationships)

			createdAt := item.Admin.Created.Date
			updatedAt := item.Admin.LastModified.Date

			mainInsert := fmt.Sprintf("INSERT INTO organizations ("+
				"id, name, status, established, types, names, "+
				"country_code, country_name, latitude, longitude, city, "+
				"wikipedia_url, website_url, "+
				"external_ids, relationships, created_at, updated_at) "+
				"VALUES ('%s', '%s', '%s', %d, '%s', '%s', '%s', '%s', %f, %f, '%s', '%s', '%s', '%s', '%s', '%s', '%s');\n",
				utils.EscapeSQL(item.ID),
				utils.EscapeSQL(GetDisplayName(item)),
				utils.EscapeSQL(status),
				established,
				utils.EscapeSQL(string(typesJSON)),
				utils.EscapeSQL(string(namesJSON)),
				utils.EscapeSQL(countryCode),
				utils.EscapeSQL(countryName),
				lat,
				lng,
				utils.EscapeSQL(cityName),
				utils.EscapeSQL(wikipedia),
				utils.EscapeSQL(website),
				utils.EscapeSQL(string(externalIDsJSON)),
				utils.EscapeSQL(string(relationshipsJSON)),
				utils.EscapeSQL(createdAt),
				utils.EscapeSQL(updatedAt))
			buffer.WriteString(mainInsert)
		}
		buffer.WriteString("\nCOMMIT;\n")

		output = buffer.Bytes()
	default:
		return output, errors.New("unsupported file format")
	}
	return output, err
}

// WriteInvenioRDM writes ROR metadata in InvenioRDM format.
func WriteInvenioRDM(data ROR) ([]byte, error) {
	var err error
	var output []byte

	inveniordm, err := ConvertInvenioRDM(data)
	if err != nil {
		fmt.Println(err)
	}
	output, err = yaml.Marshal(inveniordm)
	return output, err
}

// WriteAllInvenioRDM writes a ROR catalog in InvenioRDM format.
func WriteAllInvenioRDM(list []ROR, extension string) ([]byte, error) {
	var inveniordmList []InvenioRDM
	var err error
	var output []byte

	for _, item := range list {
		inveniordm, err := ConvertInvenioRDM(item)
		if err != nil {
			fmt.Println(err)
		}
		inveniordmList = append(inveniordmList, inveniordm)
	}
	switch extension {
	case ".yaml":
		output, err = yaml.Marshal(inveniordmList)
	case ".json":
		output, err = json.Marshal(inveniordmList)
	default:
		return output, errors.New("unsupported file format")
	}
	return output, err
}

// FilterList filters a ROR list by type and/or country.
func FilterList(list []ROR, type_ string, country string, dateUpdated string, file string, number int, page int) ([]ROR, error) {
	var filtered []ROR

	if file == "funders.yaml" {
		type_ = "funder"
	}
	// fmt.Printf("type: %s, country: %s, dateUpdated: %s, file: %s, number: %d, page: %d", type_, country, dateUpdated, file, number, page)

	// optionally filter by type, country, and/or date updated
	if type_ != "" || country != "" || file != "" || dateUpdated != "" {
		for _, v := range list {
			if type_ != "" && !slices.Contains(v.Types, type_) {
				continue
			}
			if country != "" && !slices.ContainsFunc(v.Locations, func(l Location) bool {
				return l.GeonamesDetails.CountryCode == strings.ToUpper(country)
			}) {
				continue
			}
			if dateUpdated != "" && v.Admin.LastModified.Date < dateUpdated {
				continue
			}
			if file == "funders.yaml" {
				// remove acronyms
				v.Names = slices.DeleteFunc(v.Names, func(n Name) bool {
					return slices.ContainsFunc(n.Types, func(t string) bool {
						return t == "acronym"
					})
				})
			} else if file == "affiliations_ror.yaml" {
				// remove country
				v.Locations = nil
			}
			filtered = append(filtered, v)
		}
	} else {
		filtered = list
	}

	// optionally filter by number and page
	if number > 0 {
		page = max(page, 1)
		start := (page - 1) * number
		end := min(start+number, len(filtered))

		// check if start is greater than the length of keys
		if start >= len(filtered) {
			return filtered, nil
		}
		filtered = filtered[start:end]
	}
	return filtered, nil
}

// GetTitle extracts the title from a list of names.
func GetTitle(names []Name) Title {
	var title Title

	titleValue := reflect.ValueOf(&title).Elem()

	for _, name := range names {
		if slices.Contains(name.Types, "label") {
			lang := name.Lang

			field := titleValue.FieldByNameFunc(func(fieldName string) bool {
				return strings.EqualFold(fieldName, lang)
			})

			if field.IsValid() && field.CanSet() {
				field.SetString(name.Value)
				// } else if lang != "en" {
				// 	enField := titleValue.FieldByName("En")
				// 	if enField.IsValid() && enField.CanSet() && enField.String() == "" {
				// 		enField.SetString(name.Value)
				// 	}
			}
		}
	}

	return title
}
