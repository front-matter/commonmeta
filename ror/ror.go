// Package ror converts ROR (Research Organization Registry) metadata.
package ror

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"slices"

	"gopkg.in/yaml.v3"
	"github.com/hamba/avro/v2"

	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/front-matter/commonmeta/fileutils"
	"github.com/front-matter/commonmeta/utils"
)

// ROR represents the minimal ROR metadata record.
type ROR struct {
	ID    string `json:"id"`
	Locations []Location `json:"locations"`
	Names []Name `json:"names"`
	Admin struct {
		Created struct {
			Date          string `json:"date"`
			SchemaVersion string `json:"schema_version"`
		} `json:"created"`
		LastModified struct {
			Date          string `json:"date"`
			SchemaVersion string `json:"schema_version"`
		} `json:"last_modified"`
	}
}

// Content represents the full ROR metadata record.
type Content struct {
	*ROR
	Established   int            `json:"established"`
	ExternalIDs   []ExternalID   `json:"external_ids"`
	Links         []Link         `json:"links"`
	Relationships []Relationship `json:"relationships"`
	Types         []string       `json:"types"`
	Status        string         `json:"status"`
}

// InvenioRDM represents the ROR metadata record in InvenioRDM format.
type InvenioRDM struct {
	Acronym		  string       `avro:"acronym,omitempty" yaml:"acronym,omitempty"`
	ID          string       `avro:"id" yaml:"id"`
	Country 		string       `avro:"country,omitempty" yaml:"country,omitempty"`
	Identifiers []Identifier `avro:"identifiers" yaml:"identifiers"`
	Name        string       `avro:"name" yaml:"name"`
	Title       Title        `avro:"title" yaml:"title"`
}

type ExternalID struct {
	Type      string   `json:"type"`
	All       []string `json:"all"`
	Preferred string   `json:"preferred"`
}

type GeonamesDetails struct {
	ContinentCode          string  `json:"continent_code"`
	ContinentName          string  `json:"continent_name"`
	CountryCode            string  `json:"country_code"`
	CountryName            string  `json:"country_name"`
	CountrySubdivisionCode string  `json:"country_subdivision_code"`
	CountrySubdivisionName string  `json:"country_subdivision_name"`
	Lat                    float64 `json:"lat"`
	Lng                    float64 `json:"lng"`
	Name                   string  `json:"name"`
}

type Identifier struct {
	Identifier string `avro:"identifier" json:"identifier"`
	Scheme     string `avro:"scheme" json:"scheme"`
}

type Location struct {
	GeonamesID      int             `json:"geonames_id"`
	GeonamesDetails GeonamesDetails `json:"geonames_details"`
}

type Link struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Name struct {
	Value string   `json:"value"`
	Types []string `json:"types"`
	Lang  string   `json:"lang"`
}

type Relationship struct {
	Type  string `json:"type"`
	Label string `json:"label"`
	ID    string `json:"id"`
}

type Title struct {
	Aa string `avro:"aa,omitempty" yaml:"aa,omitempty"` // Afar
	Af string `avro:"af,omitempty" yaml:"af,omitempty"` // Afrikaans
	Am string `avro:"am,omitempty" yaml:"am,omitempty"` // Amharic
	Ar string `avro:"ar,omitempty" yaml:"ar,omitempty"` // Arabic
	As string `avro:"as,omitempty" yaml:"as,omitempty"` // Assamese
	Az string `avro:"az,omitempty" yaml:"az,omitempty"` // Azerbaijani
	Ba string `avro:"ba,omitempty" yaml:"ba,omitempty"` // Bashkir
	Be string `avro:"be,omitempty" yaml:"be,omitempty"` // Belgian
	Bg string `avro:"bu,omitempty" yaml:"bg,omitempty"` // Bulgarian
	Bi string `avro:"bi,omitempty" yaml:"bi,omitempty"` // Bislama
	Bn string `avro:"bn,omitempty" yaml:"bn,omitempty"` // Bengali
	Bs string `avro:"bo,omitempty" yaml:"bs,omitempty"` // Bosnian
	Ca string `avro:"ca,omitempty" yaml:"ca,omitempty"` // Catalan
	Ch string `avro:"ch,omitempty" yaml:"ch,omitempty"` // Chamorro
	Co string `avro:"co,omitempty" yaml:"co,omitempty"` // Corsican
	Cs string `avro:"cs,omitempty" yaml:"cs,omitempty"` // Czech
	Cu string `avro:"cu,omitempty" yaml:"cu,omitempty"` // Church Slavic
	Cy string `avro:"cy,omitempty" yaml:"cy,omitempty"` // Welsh
	Da string `avro:"da,omitempty" yaml:"da,omitempty"` // Danish
	De string `avro:"de,omitempty" yaml:"de,omitempty"` // German
	Dv string `avro:"dv,omitempty" yaml:"dv,omitempty"` // Divehi
	Dz string `avro:"dz,omitempty" yaml:"dz,omitempty"` // Dzongkha
	El string `avro:"el,omitempty" yaml:"el,omitempty"` // Greek
	En string `avro:"en,omitempty" yaml:"en,omitempty"` // English
	Es string `avro:"es,omitempty" yaml:"es,omitempty"` // Spanish
	Et string `avro:"et,omitempty" yaml:"et,omitempty"` // Estonian
	Eu string `avro:"eu,omitempty" yaml:"eu,omitempty"` // Basque
	Fa string `avro:"fa,omitempty" yaml:"fa,omitempty"` // Persian
	Fi string `avro:"fi,omitempty" yaml:"fi,omitempty"` // Finnish
	Fo string `avro:"fo,omitempty" yaml:"fo,omitempty"` // Faroese
	Fr string `avro:"fr,omitempty" yaml:"fr,omitempty"` // French
	Fy string `avro:"fy,omitempty" yaml:"fy,omitempty"` // Frisian
	Ga string `avro:"ga,omitempty" yaml:"ga,omitempty"` // Irish
	Gd string `avro:"gd,omitempty" yaml:"gd,omitempty"` // Scottish Gaelic
	Gl string `avro:"gl,omitempty" yaml:"gl,omitempty"` // Galician
	Gu string `avro:"gu,omitempty" yaml:"gu,omitempty"` // Gujarati
	Gv string `avro:"gv,omitempty" yaml:"gv,omitempty"` // Manx
	Ha string `avro:"ha,omitempty" yaml:"ha,omitempty"` // Hausa
	He string `avro:"he,omitempty" yaml:"he,omitempty"` // Hebrew
	Hi string `avro:"hi,omitempty" yaml:"hi,omitempty"` // Hindi
	Hr string `avro:"hr,omitempty" yaml:"hr,omitempty"` // Croatian
	Ht string `avro:"ht,omitempty" yaml:"ht,omitempty"` // Haitian
	Hu string `avro:"hu,omitempty" yaml:"hu,omitempty"` // Hungarian
	Hy string `avro:"hy,omitempty" yaml:"hy,omitempty"` // Armenian
	Id string `avro:"id,omitempty" yaml:"id,omitempty"` // Indonesian
	Is string `avro:"is,omitempty" yaml:"is,omitempty"` // Icelandic
	It string `avro:"it,omitempty" yaml:"it,omitempty"` // Italian
	Iu string `avro:"iu,omitempty" yaml:"iu,omitempty"` // Inuktitut
	Ja string `avro:"ja,omitempty" yaml:"ja,omitempty"` // Japanese
	Jv string `avro:"jv,omitempty" yaml:"jv,omitempty"` // Javanese
	Ka string `avro:"ka,omitempty" yaml:"ka,omitempty"` // Georgian
	Kg string `avro:"kg,omitempty" yaml:"kg,omitempty"` // Kongo
	Ki string `avro:"ki,omitempty" yaml:"ki,omitempty"` // Kikuyu
	Kk string `avro:"kk,omitempty" yaml:"kk,omitempty"` // Kazakh
	Kl string `avro:"kl,omitempty" yaml:"kl,omitempty"` // Greenlandic
	Km string `avro:"km,omitempty" yaml:"km,omitempty"` // Khmer
	Kn string `avro:"kn,omitempty" yaml:"kn,omitempty"` // Kannada
	Ko string `avro:"ko,omitempty" yaml:"ko,omitempty"` // Korean
	Kr string `avro:"kr,omitempty" yaml:"kr,omitempty"` // Kanuri
	Ku string `avro:"ku,omitempty" yaml:"ku,omitempty"` // Kurdish
	Ky string `avro:"ky,omitempty" yaml:"ky,omitempty"` // Kyrgyz
	La string `avro:"la,omitempty" yaml:"la,omitempty"` // Latin
	Lb string `avro:"lb,omitempty" yaml:"lb,omitempty"` // Luxembourgish
	Lo string `avro:"lo,omitempty" yaml:"lo,omitempty"` // Lao
	Lt string `avro:"lt,omitempty" yaml:"lt,omitempty"` // Lithuanian
	Lu string `avro:"lu,omitempty" yaml:"lu,omitempty"` // Luba-Katanga
	Lv string `avro:"lv,omitempty" yaml:"lv,omitempty"` // Latvian
	Mg string `avro:"mg,omitempty" yaml:"mg,omitempty"` // Malagasy
	Mi string `avro:"mi,omitempty" yaml:"mi,omitempty"` // Maori
	Mk string `avro:"mk,omitempty" yaml:"mk,omitempty"` // Macedonian
	Ml string `avro:"ml,omitempty" yaml:"ml,omitempty"` // Malayalam
	Mn string `avro:"mn,omitempty" yaml:"mn,omitempty"` // Mongolian
	Mr string `avro:"mr,omitempty" yaml:"mr,omitempty"` // Marathi
	Ms string `avro:"ms,omitempty" yaml:"ms,omitempty"` // Malay
	Mt string `avro:"mt,omitempty" yaml:"mt,omitempty"` // Maltese
	My string `avro:"my,omitempty" yaml:"my,omitempty"` // Burmese
	Na string `avro:"na,omitempty" yaml:"na,omitempty"` // Nauru
	Nb string `avro:"nb,omitempty" yaml:"nb,omitempty"` // Norwegian Bokmål
	Ne string `avro:"ne,omitempty" yaml:"ne,omitempty"` // Nepali
	Nl string `avro:"nl,omitempty" yaml:"nl,omitempty"` // Dutch
	Nn string `avro:"nn,omitempty" yaml:"nn,omitempty"` // Norwegian Nynorsk
	No string `avro:"no,omitempty" yaml:"no,omitempty"` // Norwegian
	Oc string `avro:"oc,omitempty" yaml:"oc,omitempty"` // Occitan
	Om string `avro:"om,omitempty" yaml:"om,omitempty"` // Oromo
	Or string `avro:"or,omitempty" yaml:"or,omitempty"` // Oriya
	Pa string `avro:"pa,omitempty" yaml:"pa,omitempty"` // Punjabi
	Pl string `avro:"pl,omitempty" yaml:"pl,omitempty"` // Polish
	Ps string `avro:"ps,omitempty" yaml:"ps,omitempty"` // Pashto
	Pt string `avro:"pt,omitempty" yaml:"pt,omitempty"` // Portuguese
	Rm string `avro:"rm,omitempty" yaml:"rm,omitempty"` // Romansh
	Ro string `avro:"ro,omitempty" yaml:"ro,omitempty"` // Romanian
	Ru string `avro:"ru,omitempty" yaml:"ru,omitempty"` // Russian
	Rw string `avro:"rw,omitempty" yaml:"rw,omitempty"` // Kinyarwanda
	Sa string `avro:"sa,omitempty" yaml:"sa,omitempty"` // Sanskrit
	Sd string `avro:"sd,omitempty" yaml:"sd,omitempty"` // Sindhi
	Se string `avro:"se,omitempty" yaml:"se,omitempty"` // Northern Sami
	Sh string `avro:"sh,omitempty" yaml:"sh,omitempty"` // Serbo-Croatian
	Si string `avro:"si,omitempty" yaml:"si,omitempty"` // Sinhalese
	Sk string `avro:"sk,omitempty" yaml:"sk,omitempty"` // Slovak
	Sl string `avro:"sl,omitempty" yaml:"sl,omitempty"` // Slovenian
	Sm string `avro:"sm,omitempty" yaml:"sm,omitempty"` // Samoan
	So string `avro:"so,omitempty" yaml:"so,omitempty"` // Somali
	Sq string `avro:"sq,omitempty" yaml:"sq,omitempty"` // Albanian
	Sr string `avro:"sr,omitempty" yaml:"sr,omitempty"` // Serbian
	St string `avro:"st,omitempty" yaml:"st,omitempty"` // Southern Sotho
	Sv string `avro:"sv,omitempty" yaml:"sv,omitempty"` // Swedish
	Sw string `avro:"sw,omitempty" yaml:"sw,omitempty"` // Swahili
	Ta string `avro:"ta,omitempty" yaml:"ta,omitempty"` // Tamil
	Te string `avro:"te,omitempty" yaml:"te,omitempty"` // Telugu
	Tg string `avro:"tg,omitempty" yaml:"tg,omitempty"` // Tajik
	Th string `avro:"th,omitempty" yaml:"th,omitempty"` // Thai
	Ti string `avro:"ti,omitempty" yaml:"ti,omitempty"` // Tigrinya
	Tk string `avro:"tk,omitempty" yaml:"tk,omitempty"` // Turkmen
	Tl string `avro:"tl,omitempty" yaml:"tl,omitempty"` // Tagalog
	Tr string `avro:"tr,omitempty" yaml:"tr,omitempty"` // Turkish
	Tt string `avro:"tt,omitempty" yaml:"tt,omitempty"` // Tatar
	Ug string `avro:"ug,omitempty" yaml:"ug,omitempty"` // Uighur
	Uk string `avro:"uk,omitempty" yaml:"uk,omitempty"` // Ukrainian
	Ur string `avro:"ur,omitempty" yaml:"ur,omitempty"` // Urdu
	Uz string `avro:"uz,omitempty" yaml:"uz,omitempty"` // Uzbek
	Vi string `avro:"vi,omitempty" yaml:"vi,omitempty"` // Vietnamese
	Xh string `avro:"xh,omitempty" yaml:"xh,omitempty"` // Xhosa
	Yo string `avro:"yo,omitempty" yaml:"yo,omitempty"` // Yoruba
	Zh string `avro:"zh,omitempty" yaml:"zh,omitempty"` // Chinese
	Zu string `avro:"zu,omitempty" yaml:"zu,omitempty"` // Zulu
}

// RORVersions contains the ROR versions and their release dates, published on Zenodo.
// The ROR version is the first part of the filename, e.g., v1.63-2025-04-03-ror-data_schema_v2.json
// Beginning with release v1.45 on 11 April 2024, data releases contain JSON and CSV files formatted
// according to both schema v1 and schema v2. Version 2 files have _schema_v2 appended to the end of
// the filename, e.g., v1.45-2024-04-11-ror-data_schema_v2.json.
var RORVersions = map[string]string{
	"v1.50": "2024-07-29",
	"v1.51": "2024-08-21",
	"v1.52": "2024-09-16",
	"v1.53": "2023-10-14",
	"v1.54": "2024-10-21",
	"v1.55": "2024-10-31",
	"v1.56": "2024-11-19",
	"v1.58": "2024-12-11",
	"v1.59": "2025-01-23",
	"v1.60": "2025-02-27",
	"v1.61": "2025-03-18",
	"v1.62": "2025-03-27",
	"v1.63": "2025-04-03",
}

var InvenioRDMSchema = `{
  "type": "array",
  "items": {
    "name": "InvenioRDM",
    "type": "record",
    "fields": [
      { "name": "acronym", "type": ["null", "string"], "default": null },
      { "name": "id", "type": "string" },
			{ "name": "country", "type": ["null", "string"], "default": null },
      {
        "name": "identifiers",
        "type": {
          "type": "array",
          "items": {
            "name": "identifier",
            "type": "record",
						"fields": [
							{ "name": "identifier", "type": "string" },
							{ "name": "scheme", "type": "string" }
						]
          }
        }
      },
      { "name": "name", "type": "string" },
      {
        "name": "title",
        "type": {
          "name": "title",
          "type": "record",
          "fields": [
            { "name": "aa", "type": ["null", "string"], "default": null },
						{ "name": "af", "type": ["null", "string"], "default": null },
						{ "name": "am", "type": ["null", "string"], "default": null },
						{ "name": "ar", "type": ["null", "string"], "default": null },
						{ "name": "as", "type": ["null", "string"], "default": null },
						{ "name": "az", "type": ["null", "string"], "default": null },
						{ "name": "ba", "type": ["null", "string"], "default": null },
						{ "name": "be", "type": ["null", "string"], "default": null },
						{ "name": "bg", "type": ["null", "string"], "default": null },
						{ "name": "bi", "type": ["null", "string"], "default": null },
						{ "name": "bn", "type": ["null", "string"], "default": null },
						{ "name": "bs", "type": ["null", "string"], "default": null },
						{ "name": "ca", "type": ["null", "string"], "default": null },
						{ "name": "ch", "type": ["null", "string"], "default": null },
						{ "name": "co", "type": ["null", "string"], "default": null },
						{ "name": "cs", "type": ["null", "string"], "default": null },
						{ "name": "cu", "type": ["null", "string"], "default": null },
						{ "name": "cy", "type": ["null", "string"], "default": null },
						{ "name": "da", "type": ["null", "string"], "default": null },
						{ "name": "de", "type": ["null", "string"], "default": null },
						{ "name": "dv", "type": ["null", "string"], "default": null },
						{ "name": "dz", "type": ["null", "string"], "default": null },
						{ "name": "el", "type": ["null", "string"], "default": null },
						{ "name": "en", "type": ["null", "string"], "default": null },
						{ "name": "es", "type": ["null", "string"], "default": null },
						{ "name": "et", "type": ["null", "string"], "default": null },
						{ "name": "eu", "type": ["null", "string"], "default": null },
						{ "name": "fa", "type": ["null", "string"], "default": null },
						{ "name": "fi", "type": ["null", "string"], "default": null },
						{ "name": "fo", "type": ["null", "string"], "default": null },
						{ "name": "fr", "type": ["null", "string"], "default": null },
						{ "name": "fy", "type": ["null", "string"], "default": null },
						{ "name": "ga", "type": ["null", "string"], "default": null },
						{ "name": "gd", "type": ["null", "string"], "default": null },
						{ "name": "gl", "type": ["null", "string"], "default": null },
						{ "name": "gu", "type": ["null", "string"], "default": null },
						{ "name": "ha", "type": ["null", "string"], "default": null },
						{ "name": "he", "type": ["null", "string"], "default": null },
						{ "name": "hi", "type": ["null", "string"], "default": null },
						{ "name": "hr", "type": ["null", "string"], "default": null },
						{ "name": "ht", "type": ["null", "string"], "default": null },
						{ "name": "hu", "type": ["null", "string"], "default": null },
						{ "name": "hy", "type": ["null", "string"], "default": null },
						{ "name": "id", "type": ["null", "string"], "default": null },
						{ "name": "is", "type": ["null", "string"], "default": null },
						{ "name": "it", "type": ["null", "string"], "default": null },
						{ "name": "iu", "type": ["null", "string"], "default": null },
						{ "name": "ja", "type": ["null", "string"], "default": null },
						{ "name": "jv", "type": ["null", "string"], "default": null },
						{ "name": "ka", "type": ["null", "string"], "default": null },
						{ "name": "kg", "type": ["null", "string"], "default": null },
						{ "name": "ki", "type": ["null", "string"], "default": null },
						{ "name": "kk", "type": ["null", "string"], "default": null },
						{ "name": "kl", "type": ["null", "string"], "default": null },
						{ "name": "km", "type": ["null", "string"], "default": null },
						{ "name": "kn", "type": ["null", "string"], "default": null },
						{ "name": "ko", "type": ["null", "string"], "default": null },
						{ "name": "kr", "type": ["null", "string"], "default": null },
						{ "name": "ku", "type": ["null", "string"], "default": null },
						{ "name": "ky", "type": ["null", "string"], "default": null },
						{ "name": "la", "type": ["null", "string"], "default": null },
						{ "name": "lb", "type": ["null", "string"], "default": null },
						{ "name": "lo", "type": ["null", "string"], "default": null },
						{ "name": "lt", "type": ["null", "string"], "default": null },
						{ "name": "lu", "type": ["null", "string"], "default": null },
						{ "name": "lv", "type": ["null", "string"], "default": null },
						{ "name": "mg", "type": ["null", "string"], "default": null },
						{ "name": "mi", "type": ["null", "string"], "default": null },
						{ "name": "mk", "type": ["null", "string"], "default": null },
						{ "name": "ml", "type": ["null", "string"], "default": null },
						{ "name": "mn", "type": ["null", "string"], "default": null },
						{ "name": "mr", "type": ["null", "string"], "default": null },
						{ "name": "ms", "type": ["null", "string"], "default": null },
						{ "name": "mt", "type": ["null", "string"], "default": null },
						{ "name": "my", "type": ["null", "string"], "default": null },
						{ "name": "na", "type": ["null", "string"], "default": null },
						{ "name": "nb", "type": ["null", "string"], "default": null },
						{ "name": "ne", "type": ["null", "string"], "default": null },
						{ "name": "nl", "type": ["null", "string"], "default": null },
						{ "name": "nn", "type": ["null", "string"], "default": null },
						{ "name": "no", "type": ["null", "string"], "default": null },
						{ "name": "oc", "type": ["null", "string"], "default": null },
						{ "name": "om", "type": ["null", "string"], "default": null },
						{ "name": "or", "type": ["null", "string"], "default": null },
						{ "name": "pa", "type": ["null", "string"], "default": null },
						{ "name": "pl", "type": ["null", "string"], "default": null },
						{ "name": "ps", "type": ["null", "string"], "default": null },
						{ "name": "pt", "type": ["null", "string"], "default": null },
						{ "name": "rm", "type": ["null", "string"], "default": null },
						{ "name": "ro", "type": ["null", "string"], "default": null },
						{ "name": "ru", "type": ["null", "string"], "default": null },
						{ "name": "rw", "type": ["null", "string"], "default": null },
						{ "name": "sa", "type": ["null", "string"], "default": null },
						{ "name": "sd", "type": ["null", "string"], "default": null },
						{ "name": "se", "type": ["null", "string"], "default": null },
						{ "name": "sh", "type": ["null", "string"], "default": null },
						{ "name": "si", "type": ["null", "string"], "default": null },
						{ "name": "sk", "type": ["null", "string"], "default": null },
						{ "name": "sl", "type": ["null", "string"], "default": null },
						{ "name": "sm", "type": ["null", "string"], "default": null },
						{ "name": "so", "type": ["null", "string"], "default": null },
						{ "name": "sq", "type": ["null", "string"], "default": null },
						{ "name": "sr", "type": ["null", "string"], "default": null },
						{ "name": "st", "type": ["null", "string"], "default": null },
						{ "name": "sv", "type": ["null", "string"], "default": null },
						{ "name": "sw", "type": ["null", "string"], "default": null },
						{ "name": "ta", "type": ["null", "string"], "default": null },
						{ "name": "te", "type": ["null", "string"], "default": null },
						{ "name": "tg", "type": ["null", "string"], "default": null },
						{ "name": "th", "type": ["null", "string"], "default": null },
						{ "name": "ti", "type": ["null", "string"], "default": null },
						{ "name": "tk", "type": ["null", "string"], "default": null },
						{ "name": "tl", "type": ["null", "string"], "default": null },
						{ "name": "tr", "type": ["null", "string"], "default": null },
						{ "name": "tt", "type": ["null", "string"], "default": null },
						{ "name": "ug", "type": ["null", "string"], "default": null },
						{ "name": "uk", "type": ["null", "string"], "default": null },
						{ "name": "ur", "type": ["null", "string"], "default": null },
						{ "name": "uz", "type": ["null", "string"], "default": null },
						{ "name": "vi", "type": ["null", "string"], "default": null },
						{ "name": "xh", "type": ["null", "string"], "default": null },
						{ "name": "yo", "type": ["null", "string"], "default": null },
						{ "name": "zh", "type": ["null", "string"], "default": null },
						{ "name": "zu", "type": ["null", "string"], "default": null }
          ]
        }
      }
    ]
  }
}`

var Extensions = []string{".json", ".yaml", ".avro"}
var RORTypes = []string{"archive", "company", "education", "facility", "funder", "government", "healthcare", "nonprofit", "other"}		

// LoadAll loads the metadata for a list of organizations from a ROR JSON file
func LoadAll(filename string, type_ string, country string) ([]ROR, error) {
	var data []ROR
	var content []Content
	var err error

	extension := path.Ext(filename)
	if extension == ".json" {
		file, err := os.Open(filename)
		if err != nil {
			return data, errors.New("error reading file")
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		err = decoder.Decode(&content)
		if err != nil {
			return data, err
		}
	} else if extension != ".json" {
		return data, errors.New("invalid file extension")
	} else {
		return data, errors.New("unsupported file format")
	}

	data, err = ReadAll(content, type_, country)
	if err != nil {
		return data, err
	}
	return data, nil
}

// LoadBuiltin loads the embedded ROR metadata from the ZIP file with all ROR records.
func LoadBuiltin() ([]byte, error) {
	output, err := fileutils.ReadZIPFile("affiliations_ror.yaml.zip")
	if err != nil {
		return nil, err
	}
	return output, err
}

// Read reads ROR full metadata and converts it into ROR minimal metadata.
func Read(content Content) (ROR, error) {
	var data ROR

	data.ID = content.ID
	data.Locations = content.Locations
	data.Names = content.Names
	data.Admin.LastModified.Date = content.Admin.LastModified.Date

	return data, nil
}

// ReadAll reads a list of ROR JSON organizations
func ReadAll(content []Content, type_ string, country string) ([]ROR, error) {
	var filtered []Content
	var data []ROR

  // optionally filter by type and/or country
	if type_ != "" || country != "" {
		for _, v := range content {
			if type_ != "" && !slices.Contains(v.Types, type_) {
				continue
			}
			if country != "" && !slices.ContainsFunc(v.Locations, func(l Location) bool {
				return l.GeonamesDetails.CountryCode == country
			}) {
        continue
			}
      filtered = append(filtered, v)
		}
	} else {
    filtered = append(filtered, content...)
	}

	for _, v := range filtered {
		d, err := Read(v)
		if err != nil {
			log.Println(err)
		}
		data = append(data, d)
	}
	return data, nil
}

// ExtractAll extracts ROR metadata from a JSON file in commonmeta format.
func ExtractAll(content []commonmeta.Data) ([]byte, error) {
	var data []InvenioRDM
	var extracted []InvenioRDM
	var ids []string
	var err error
  schema, err := avro.Parse(InvenioRDMSchema)
	if err != nil {
		return nil, err
	}

	// Load the ROR metadata from the embedded ZIP file with all ROR records
	out, err := fileutils.ReadZIPFile("affiliations_ror.yaml.zip")
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(out, &data)
	if err != nil {
		return nil, err
	}

	// Extract ROR IDs from the content
	for _, v := range content {
		if len(v.Contributors) > 0 {
			for _, c := range v.Contributors {
				if len(c.Affiliations) > 0 {
					for _, a := range c.Affiliations {
						if a.ID != "" && !slices.Contains(ids, a.ID) {
							id, _ := utils.ValidateROR(a.ID)
							idx := slices.IndexFunc(data, func(d InvenioRDM) bool { return d.ID == id })
							if idx != -1 {
								ids = append(ids, a.ID)
								extracted = append(extracted, data[idx])
							}
						}
					}
				}
			}
		}
	}

	output, err := avro.Marshal(schema, extracted)
	return output, err
}

// Convert converts ROR metadata into InvenioRDM format.
func Convert(data ROR, type_ string) (InvenioRDM, error) {
	var inveniordm InvenioRDM

	id, _ := utils.ValidateROR(data.ID)
	inveniordm.ID = id
	if type_ == "funder" {
		for _, location := range data.Locations {
      inveniordm.Country = location.GeonamesDetails.CountryCode
		}
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
		} else if type_ != "funder" && slices.Contains(name.Types, "acronym") && name.Value != "" {
			inveniordm.Acronym = name.Value
		}
	}
	inveniordm.Title = GetTitle(data.Names)
	return inveniordm, nil
}

// Write writes ROR metadata to InvenioRDM format.
func Write(data ROR, extension string, type_ string) ([]byte, error) {
	schema, err := avro.Parse(InvenioRDMSchema)
	if err != nil {
		return nil, err
	}
	inveniordm, err := Convert(data, type_)
	if err != nil {
		fmt.Println(err)
	}
	output, err := avro.Marshal(schema, inveniordm)
	return output, err
}

// WriteAll writes a list of ROR metadata in InvenioRDM YAML format.
func WriteAll(list []ROR, to string, extension string, type_ string) ([]byte, error) {
	var inveniordmList []InvenioRDM
	var err error
	var output []byte
		
	type InvenioRDM struct {
		Acronym		  string       `avro:"acronym,omitempty" yaml:"acronym,omitempty"`
		ID          string       `avro:"id" yaml:"id"`
		// Country 		string       `avro:"country,omitempty" yaml:"country,omitempty"`
		Name        string       `avro:"name" yaml:"name"`
		Title       Title        `avro:"title" yaml:"title"`
	}

	if to != "inveniordm" {
		return output, errors.New("unsupported output format")
	}

	for _, data := range list {
		inveniordm, err := Convert(data, type_)
		if err != nil {
			fmt.Println(err)
		}
		if inveniordm.ID != "" {
			inveniordmList = append(inveniordmList, inveniordm)
		}
	}
  if extension == ".yaml" {
		output, err = yaml.Marshal(inveniordmList)
	} else if extension == ".json" {
		output, err = json.Marshal(inveniordmList)
	} else if extension == ".avro" {
		schema, err := avro.Parse(InvenioRDMSchema)
		if err != nil {
			fmt.Println(err, "avro.Parse")
			return nil, err
		}
    output, err = avro.Marshal(schema, inveniordmList)
		if err != nil {
			fmt.Println(err, "avro.Marshal")
		}
	} else {
		return output, errors.New("unsupported file format")
	}
	if err != nil {
		return nil, err
	}
	return output, err
}

func GetTitle(names []Name) Title {
	var title Title
	for _, name := range names {
		if slices.Contains(name.Types, "label") {
			switch name.Lang {
			case "aa":
				title.Aa = name.Value
			case "ab":
				title.Aa = name.Value
			case "af":
				title.Af = name.Value
			case "am":
				title.Am = name.Value
			case "ar":
				title.Ar = name.Value
			case "as":
				title.As = name.Value
			case "az":
				title.Az = name.Value
			case "ba":
				title.Ba = name.Value
			case "be":
				title.Be = name.Value
			case "bg":
				title.Bg = name.Value
			case "bi":
				title.Bi = name.Value
			case "bn":
				title.Bn = name.Value
			case "bs":
				title.Bs = name.Value
			case "ca":
				title.Ca = name.Value
			case "ch":
				title.Ch = name.Value
			case "co":
				title.Co = name.Value
			case "cs":
				title.Cs = name.Value
			case "cu":
				title.Cu = name.Value
			case "cy":
				title.Cy = name.Value
			case "da":
				title.Da = name.Value
			case "de":
				title.De = name.Value
			case "dv":
				title.Dv = name.Value
			case "dz":
				title.Dz = name.Value
			case "el":
				title.El = name.Value
			case "en":
				title.En = name.Value
			case "es":
				title.Es = name.Value
			case "et":
				title.Et = name.Value
			case "eu":
				title.Eu = name.Value
			case "fa":
				title.Fa = name.Value
			case "fi":
				title.Fi = name.Value
			case "fo":
				title.Fo = name.Value
			case "fr":
				title.Fr = name.Value
			case "fy":
				title.Fy = name.Value
			case "ga":
				title.Ga = name.Value
			case "gd":
				title.Gd = name.Value
			case "gl":
				title.Gl = name.Value
			case "gu":
				title.Gu = name.Value
			case "gv":
				title.Gv = name.Value
			case "ha":
				title.Ha = name.Value
			case "he":
				title.He = name.Value
			case "hi":
				title.Hi = name.Value
			case "hr":
				title.Hr = name.Value
			case "ht":
				title.Ht = name.Value
			case "hu":
				title.Hu = name.Value
			case "hy":
				title.Hy = name.Value
			case "id":
				title.Id = name.Value
			case "is":
				title.Is = name.Value
			case "it":
				title.It = name.Value
			case "iu":
				title.Iu = name.Value
			case "ja":
				title.Ja = name.Value
			case "jv":
				title.Jv = name.Value
			case "ka":
				title.Ka = name.Value
			case "kg":
				title.Kg = name.Value
			case "ki":
				title.Ki = name.Value
			case "kk":
				title.Kk = name.Value
			case "kl":
				title.Kl = name.Value
			case "km":
				title.Km = name.Value
			case "kn":
				title.Kn = name.Value
			case "ko":
				title.Ko = name.Value
			case "kr":
				title.Kr = name.Value
			case "ku":
				title.Ku = name.Value
			case "ky":
				title.Ky = name.Value
			case "la":
				title.La = name.Value
			case "lb":
				title.Lb = name.Value
			case "lo":
				title.Lo = name.Value
			case "lt":
				title.Lt = name.Value
			case "lv":
				title.Lv = name.Value
			case "lu":
				title.Lu = name.Value
			case "mg":
				title.Mg = name.Value
			case "mi":
				title.Mi = name.Value
			case "mk":
				title.Mk = name.Value
			case "ml":
				title.Ml = name.Value
			case "mn":
				title.Mn = name.Value
			case "mr":
				title.Mr = name.Value
			case "ms":
				title.Ms = name.Value
			case "mt":
				title.Mt = name.Value
			case "my":
				title.My = name.Value
			case "na":
				title.Na = name.Value
			case "nb":
				title.Nb = name.Value
			case "ne":
				title.Ne = name.Value
			case "nl":
				title.Nl = name.Value
			case "nn":
				title.Nn = name.Value
			case "no":
				title.No = name.Value
			case "oc":
				title.Oc = name.Value
			case "om":
				title.Om = name.Value
			case "or":
				title.Or = name.Value
			case "pa":
				title.Pa = name.Value
			case "pl":
				title.Pl = name.Value
			case "ps":
				title.Ps = name.Value
			case "pt":
				title.Pt = name.Value
			case "rm":
				title.Rm = name.Value
			case "ro":
				title.Ro = name.Value
			case "ru":
				title.Ru = name.Value
			case "rw":
				title.Rw = name.Value
			case "sa":
				title.Sa = name.Value
			case "sd":
				title.Sd = name.Value
			case "se":
				title.Se = name.Value
			case "sh":
				title.Sh = name.Value
			case "si":
				title.Si = name.Value
			case "sk":
				title.Sk = name.Value
			case "sl":
				title.Sl = name.Value
			case "sm":
				title.Sm = name.Value
			case "so":
				title.So = name.Value
			case "sq":
				title.Sq = name.Value
			case "sr":
				title.Sr = name.Value
			case "st":
				title.St = name.Value
			case "sv":
				title.Sv = name.Value
			case "sw":
				title.Sw = name.Value
			case "ta":
				title.Ta = name.Value
			case "te":
				title.Te = name.Value
			case "tg":
				title.Tg = name.Value
			case "th":
				title.Th = name.Value
			case "ti":
				title.Ti = name.Value
			case "tk":
				title.Tk = name.Value
			case "tl":
				title.Tl = name.Value
			case "tr":
				title.Tr = name.Value
			case "tt":
				title.Tt = name.Value
			case "ug":
				title.Ug = name.Value
			case "uk":
				title.Uk = name.Value
			case "ur":
				title.Ur = name.Value
			case "uz":
				title.Uz = name.Value
			case "vi":
				title.Vi = name.Value
			case "xh":
				title.Xh = name.Value
			case "yo":
				title.Yo = name.Value
			case "zh":
				title.Zh = name.Value
			case "zu":
				title.Zu = name.Value
			default:
				title.En = name.Value
			}
		}
	}
	return title
}
