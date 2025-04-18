package ror

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/front-matter/commonmeta/utils"
	"github.com/hamba/avro/v2"
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
	Acronym     string       `avro:"acronym,omitempty" json:"acronym,omitempty" yaml:"acronym,omitempty"`
	ID          string       `avro:"id" json:"id"`
	Country     string       `avro:"country,omitempty" json:"country,omitempty" yaml:"country,omitempty"`
	Identifiers []Identifier `avro:"identifiers" json:"identifiers"`
	Name        string       `avro:"name" json:"name"`
	Title       Title        `avro:"title" json:"title"`
}

type Identifier struct {
	Identifier string `avro:"identifier" json:"identifier"`
	Scheme     string `avro:"scheme" json:"scheme"`
}

type Title struct {
	Aa string `avro:"aa,omitempty" json:"aa,omitempty" yaml:"aa,omitempty"` // Afar
	Af string `avro:"af,omitempty" json:"af,omitempty" yaml:"af,omitempty"` // Afrikaans
	Am string `avro:"am,omitempty" json:"am,omitempty" yaml:"am,omitempty"` // Amharic
	Ar string `avro:"ar,omitempty" json:"ar,omitempty" yaml:"ar,omitempty"` // Arabic
	As string `avro:"as,omitempty" json:"as,omitempty" yaml:"as,omitempty"` // Assamese
	Az string `avro:"az,omitempty" json:"az,omitempty" yaml:"az,omitempty"` // Azerbaijani
	Ba string `avro:"ba,omitempty" json:"aa,omitempty" yaml:"ba,omitempty"` // Bashkir
	Be string `avro:"be,omitempty" json:"aa,omitempty" yaml:"be,omitempty"` // Belgian
	Bg string `avro:"bu,omitempty" json:"aa,omitempty" yaml:"bg,omitempty"` // Bulgarian
	Bi string `avro:"bi,omitempty" json:"aa,omitempty" yaml:"bi,omitempty"` // Bislama
	Bn string `avro:"bn,omitempty" json:"aa,omitempty" yaml:"bn,omitempty"` // Bengali
	Bs string `avro:"bo,omitempty" json:"aa,omitempty" yaml:"bs,omitempty"` // Bosnian
	Ca string `avro:"ca,omitempty" json:"aa,omitempty" yaml:"ca,omitempty"` // Catalan
	Ch string `avro:"ch,omitempty" json:"aa,omitempty" yaml:"ch,omitempty"` // Chamorro
	Co string `avro:"co,omitempty" json:"aa,omitempty" yaml:"co,omitempty"` // Corsican
	Cs string `avro:"cs,omitempty" json:"aa,omitempty" yaml:"cs,omitempty"` // Czech
	Cu string `avro:"cu,omitempty" json:"aa,omitempty" yaml:"cu,omitempty"` // Church Slavic
	Cy string `avro:"cy,omitempty" json:"aa,omitempty" yaml:"cy,omitempty"` // Welsh
	Da string `avro:"da,omitempty" json:"da,omitempty" yaml:"da,omitempty"` // Danish
	De string `avro:"de,omitempty" json:"de,omitempty" yaml:"de,omitempty"` // German
	Dv string `avro:"dv,omitempty" json:"aa,omitempty" yaml:"dv,omitempty"` // Divehi
	Dz string `avro:"dz,omitempty" json:"aa,omitempty" yaml:"dz,omitempty"` // Dzongkha
	El string `avro:"el,omitempty" json:"aa,omitempty" yaml:"el,omitempty"` // Greek
	En string `avro:"en,omitempty" json:"en,omitempty" yaml:"en,omitempty"` // English
	Es string `avro:"es,omitempty" json:"es,omitempty" yaml:"es,omitempty"` // Spanish
	Et string `avro:"et,omitempty" json:"aa,omitempty" yaml:"et,omitempty"` // Estonian
	Eu string `avro:"eu,omitempty" json:"aa,omitempty" yaml:"eu,omitempty"` // Basque
	Fa string `avro:"fa,omitempty" json:"aa,omitempty" yaml:"fa,omitempty"` // Persian
	Fi string `avro:"fi,omitempty" json:"aa,omitempty" yaml:"fi,omitempty"` // Finnish
	Fo string `avro:"fo,omitempty" json:"aa,omitempty" yaml:"fo,omitempty"` // Faroese
	Fr string `avro:"fr,omitempty" json:"fr,omitempty" yaml:"fr,omitempty"` // French
	Fy string `avro:"fy,omitempty" json:"aa,omitempty" yaml:"fy,omitempty"` // Frisian
	Ga string `avro:"ga,omitempty" json:"aa,omitempty" yaml:"ga,omitempty"` // Irish
	Gd string `avro:"gd,omitempty" json:"aa,omitempty" yaml:"gd,omitempty"` // Scottish Gaelic
	Gl string `avro:"gl,omitempty" json:"aa,omitempty" yaml:"gl,omitempty"` // Galician
	Gu string `avro:"gu,omitempty" json:"aa,omitempty" yaml:"gu,omitempty"` // Gujarati
	Gv string `avro:"gv,omitempty" json:"aa,omitempty" yaml:"gv,omitempty"` // Manx
	Ha string `avro:"ha,omitempty" json:"aa,omitempty" yaml:"ha,omitempty"` // Hausa
	He string `avro:"he,omitempty" json:"aa,omitempty" yaml:"he,omitempty"` // Hebrew
	Hi string `avro:"hi,omitempty" json:"aa,omitempty" yaml:"hi,omitempty"` // Hindi
	Hr string `avro:"hr,omitempty" json:"aa,omitempty" yaml:"hr,omitempty"` // Croatian
	Ht string `avro:"ht,omitempty" json:"aa,omitempty" yaml:"ht,omitempty"` // Haitian
	Hu string `avro:"hu,omitempty" json:"aa,omitempty" yaml:"hu,omitempty"` // Hungarian
	Hy string `avro:"hy,omitempty" json:"aa,omitempty" yaml:"hy,omitempty"` // Armenian
	Id string `avro:"id,omitempty" json:"aa,omitempty" yaml:"id,omitempty"` // Indonesian
	Is string `avro:"is,omitempty" json:"aa,omitempty" yaml:"is,omitempty"` // Icelandic
	It string `avro:"it,omitempty" json:"aa,omitempty" yaml:"it,omitempty"` // Italian
	Iu string `avro:"iu,omitempty" json:"aa,omitempty" yaml:"iu,omitempty"` // Inuktitut
	Ja string `avro:"ja,omitempty" json:"aa,omitempty" yaml:"ja,omitempty"` // Japanese
	Jv string `avro:"jv,omitempty" json:"aa,omitempty" yaml:"jv,omitempty"` // Javanese
	Ka string `avro:"ka,omitempty" json:"aa,omitempty" yaml:"ka,omitempty"` // Georgian
	Kg string `avro:"kg,omitempty" json:"aa,omitempty" yaml:"kg,omitempty"` // Kongo
	Ki string `avro:"ki,omitempty" json:"aa,omitempty" yaml:"ki,omitempty"` // Kikuyu
	Kk string `avro:"kk,omitempty" json:"aa,omitempty" yaml:"kk,omitempty"` // Kazakh
	Kl string `avro:"kl,omitempty" json:"aa,omitempty" yaml:"kl,omitempty"` // Greenlandic
	Km string `avro:"km,omitempty" json:"aa,omitempty" yaml:"km,omitempty"` // Khmer
	Kn string `avro:"kn,omitempty" json:"aa,omitempty" yaml:"kn,omitempty"` // Kannada
	Ko string `avro:"ko,omitempty" json:"aa,omitempty" yaml:"ko,omitempty"` // Korean
	Kr string `avro:"kr,omitempty" json:"aa,omitempty" yaml:"kr,omitempty"` // Kanuri
	Ku string `avro:"ku,omitempty" json:"aa,omitempty" yaml:"ku,omitempty"` // Kurdish
	Ky string `avro:"ky,omitempty" json:"aa,omitempty" yaml:"ky,omitempty"` // Kyrgyz
	La string `avro:"la,omitempty" json:"aa,omitempty" yaml:"la,omitempty"` // Latin
	Lb string `avro:"lb,omitempty" json:"aa,omitempty" yaml:"lb,omitempty"` // Luxembourgish
	Lo string `avro:"lo,omitempty" json:"aa,omitempty" yaml:"lo,omitempty"` // Lao
	Lt string `avro:"lt,omitempty" json:"aa,omitempty" yaml:"lt,omitempty"` // Lithuanian
	Lu string `avro:"lu,omitempty" json:"aa,omitempty" yaml:"lu,omitempty"` // Luba-Katanga
	Lv string `avro:"lv,omitempty" json:"aa,omitempty" yaml:"lv,omitempty"` // Latvian
	Mg string `avro:"mg,omitempty" json:"aa,omitempty" yaml:"mg,omitempty"` // Malagasy
	Mi string `avro:"mi,omitempty" json:"aa,omitempty" yaml:"mi,omitempty"` // Maori
	Mk string `avro:"mk,omitempty" json:"aa,omitempty" yaml:"mk,omitempty"` // Macedonian
	Ml string `avro:"ml,omitempty" json:"aa,omitempty" yaml:"ml,omitempty"` // Malayalam
	Mn string `avro:"mn,omitempty" json:"aa,omitempty" yaml:"mn,omitempty"` // Mongolian
	Mr string `avro:"mr,omitempty" json:"aa,omitempty" yaml:"mr,omitempty"` // Marathi
	Ms string `avro:"ms,omitempty" json:"aa,omitempty" yaml:"ms,omitempty"` // Malay
	Mt string `avro:"mt,omitempty" json:"aa,omitempty" yaml:"mt,omitempty"` // Maltese
	My string `avro:"my,omitempty" json:"aa,omitempty" yaml:"my,omitempty"` // Burmese
	Na string `avro:"na,omitempty" json:"aa,omitempty" yaml:"na,omitempty"` // Nauru
	Nb string `avro:"nb,omitempty" json:"aa,omitempty" yaml:"nb,omitempty"` // Norwegian Bokmål
	Ne string `avro:"ne,omitempty" json:"aa,omitempty" yaml:"ne,omitempty"` // Nepali
	Nl string `avro:"nl,omitempty" json:"aa,omitempty" yaml:"nl,omitempty"` // Dutch
	Nn string `avro:"nn,omitempty" json:"aa,omitempty" yaml:"nn,omitempty"` // Norwegian Nynorsk
	No string `avro:"no,omitempty" json:"aa,omitempty" yaml:"no,omitempty"` // Norwegian
	Oc string `avro:"oc,omitempty" json:"aa,omitempty" yaml:"oc,omitempty"` // Occitan
	Om string `avro:"om,omitempty" json:"aa,omitempty" yaml:"om,omitempty"` // Oromo
	Or string `avro:"or,omitempty" json:"aa,omitempty" yaml:"or,omitempty"` // Oriya
	Pa string `avro:"pa,omitempty" json:"aa,omitempty" yaml:"pa,omitempty"` // Punjabi
	Pl string `avro:"pl,omitempty" json:"aa,omitempty" yaml:"pl,omitempty"` // Polish
	Ps string `avro:"ps,omitempty" json:"aa,omitempty" yaml:"ps,omitempty"` // Pashto
	Pt string `avro:"pt,omitempty" json:"aa,omitempty" yaml:"pt,omitempty"` // Portuguese
	Rm string `avro:"rm,omitempty" json:"aa,omitempty" yaml:"rm,omitempty"` // Romansh
	Ro string `avro:"ro,omitempty" json:"aa,omitempty" yaml:"ro,omitempty"` // Romanian
	Ru string `avro:"ru,omitempty" json:"aa,omitempty" yaml:"ru,omitempty"` // Russian
	Rw string `avro:"rw,omitempty" json:"aa,omitempty" yaml:"rw,omitempty"` // Kinyarwanda
	Sa string `avro:"sa,omitempty" json:"aa,omitempty" yaml:"sa,omitempty"` // Sanskrit
	Sd string `avro:"sd,omitempty" json:"aa,omitempty" yaml:"sd,omitempty"` // Sindhi
	Se string `avro:"se,omitempty" json:"aa,omitempty" yaml:"se,omitempty"` // Northern Sami
	Sh string `avro:"sh,omitempty" json:"aa,omitempty" yaml:"sh,omitempty"` // Serbo-Croatian
	Si string `avro:"si,omitempty" json:"aa,omitempty" yaml:"si,omitempty"` // Sinhalese
	Sk string `avro:"sk,omitempty" json:"aa,omitempty" yaml:"sk,omitempty"` // Slovak
	Sl string `avro:"sl,omitempty" json:"aa,omitempty" yaml:"sl,omitempty"` // Slovenian
	Sm string `avro:"sm,omitempty" json:"aa,omitempty" yaml:"sm,omitempty"` // Samoan
	So string `avro:"so,omitempty" json:"aa,omitempty" yaml:"so,omitempty"` // Somali
	Sq string `avro:"sq,omitempty" json:"aa,omitempty" yaml:"sq,omitempty"` // Albanian
	Sr string `avro:"sr,omitempty" json:"aa,omitempty" yaml:"sr,omitempty"` // Serbian
	St string `avro:"st,omitempty" json:"aa,omitempty" yaml:"st,omitempty"` // Southern Sotho
	Sv string `avro:"sv,omitempty" json:"aa,omitempty" yaml:"sv,omitempty"` // Swedish
	Sw string `avro:"sw,omitempty" json:"aa,omitempty" yaml:"sw,omitempty"` // Swahili
	Ta string `avro:"ta,omitempty" json:"aa,omitempty" yaml:"ta,omitempty"` // Tamil
	Te string `avro:"te,omitempty" json:"aa,omitempty" yaml:"te,omitempty"` // Telugu
	Tg string `avro:"tg,omitempty" json:"aa,omitempty" yaml:"tg,omitempty"` // Tajik
	Th string `avro:"th,omitempty" json:"aa,omitempty" yaml:"th,omitempty"` // Thai
	Ti string `avro:"ti,omitempty" json:"aa,omitempty" yaml:"ti,omitempty"` // Tigrinya
	Tk string `avro:"tk,omitempty" json:"aa,omitempty" yaml:"tk,omitempty"` // Turkmen
	Tl string `avro:"tl,omitempty" json:"aa,omitempty" yaml:"tl,omitempty"` // Tagalog
	Tr string `avro:"tr,omitempty" json:"aa,omitempty" yaml:"tr,omitempty"` // Turkish
	Tt string `avro:"tt,omitempty" json:"aa,omitempty" yaml:"tt,omitempty"` // Tatar
	Ug string `avro:"ug,omitempty" json:"aa,omitempty" yaml:"ug,omitempty"` // Uighur
	Uk string `avro:"uk,omitempty" json:"aa,omitempty" yaml:"uk,omitempty"` // Ukrainian
	Ur string `avro:"ur,omitempty" json:"aa,omitempty" yaml:"ur,omitempty"` // Urdu
	Uz string `avro:"uz,omitempty" json:"aa,omitempty" yaml:"uz,omitempty"` // Uzbek
	Vi string `avro:"vi,omitempty" json:"aa,omitempty" yaml:"vi,omitempty"` // Vietnamese
	Xh string `avro:"xh,omitempty" json:"aa,omitempty" yaml:"xh,omitempty"` // Xhosa
	Yo string `avro:"yo,omitempty" json:"aa,omitempty" yaml:"yo,omitempty"` // Yoruba
	Zh string `avro:"zh,omitempty" json:"zh,omitempty" yaml:"zh,omitempty"` // Chinese
	Zu string `avro:"zu,omitempty" json:"aa,omitempty" yaml:"zu,omitempty"` // Zulu
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

	if extension == ".yaml" {
		output, err = yaml.Marshal(list)
	} else if extension == ".json" {
		output, err = json.Marshal(list)
	} else if extension == ".csv" {
		var rorcsvList []RORCSV
		// convert ROR to RORCSV, a custom lossy mapping to CSV
		for _, data := range list {
			rorcsv, err := ConvertRORCSV(data)
			if err != nil {
				fmt.Println(err)
			}
			rorcsvList = append(rorcsvList, rorcsv)
		}
		output, err = csvutil.Marshal(rorcsvList)
		if err != nil {
			fmt.Println(err, "csvutil.Marshal")
		}
	} else if extension == ".avro" {
		schema, err := avro.Parse(RORSchema)
		if err != nil {
			fmt.Println(err, "avro.Parse")
			return nil, err
		}
		output, err = avro.Marshal(schema, list)
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

// WriteAllInvenioRDM writes a list of ROR metadata in InvenioRDM format.
func WriteAllInvenioRDM(list []ROR, extension string) ([]byte, error) {
	var inveniordmList []InvenioRDM
	var err error
	var output []byte

	for _, data := range list {
		inveniordm, err := ConvertInvenioRDM(data)
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

// FilterRecords filters a list of ROR records by type and/or country.
func FilterRecords(list []ROR, type_ string, country string, dateUpdated string, file string, number int, page int) ([]ROR, error) {
	var filtered []ROR

	if file == "funders.yaml" {
		type_ = "funder"
	}

	// optionally filter by type and/or country
	if type_ != "" || country != "" || file != "" {
		for _, v := range list {
			if type_ != "" && !slices.Contains(v.Types, type_) {
				continue
			}
			if country != "" && !slices.ContainsFunc(v.Locations, func(l Location) bool {
				return l.GeonamesDetails.CountryCode == country
			}) {
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
		filtered = append(filtered, list...)
	}

	// optionally filter by date updated
	if dateUpdated != "" {
		// validate date format
		_, err := time.Parse("2006-01-02", dateUpdated)
		if err != nil {
			return nil, fmt.Errorf("invalid date format: %v", err)
		}
		filtered = slices.DeleteFunc(filtered, func(r ROR) bool {
			return r.Admin.LastModified.Date < dateUpdated
		})
	}

	// optionally filter by number and page
	if number > 0 {
		page = max(page, 1)
		start := (page - 1) * number
		end := min(start+number, len(filtered))
		if start > len(filtered) {
			start = len(filtered)
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
