package ror

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"

	"github.com/front-matter/commonmeta/utils"
	"github.com/hamba/avro/v2"
	"gopkg.in/yaml.v3"
)

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
func Convert(data ROR) (InvenioRDM, error) {
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

	inveniordm, err := Convert(data)
	if err != nil {
		fmt.Println(err)
	}
	output, err = json.Marshal(inveniordm)
	return output, err
}

// WriteAllInvenioRDM writes a list of ROR metadata in InvenioRDM format.
func WriteAllInvenioRDM(list []ROR, extension string) ([]byte, error) {
	var inveniordmList []InvenioRDM
	var err error
	var output []byte

	for _, data := range list {
		inveniordm, err := Convert(data)
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
func FilterRecords(list []ROR, type_ string, country string, file string, number int, page int) ([]ROR, error) {
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

	// optionally filter by number and page
	if number > 0 && number != 10 {
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
