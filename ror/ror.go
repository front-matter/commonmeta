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

	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/front-matter/commonmeta/fileutils"
	"github.com/front-matter/commonmeta/utils"
)

// ROR represents the minimal ROR metadata record.
type ROR struct {
	ID    string `json:"id"`
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
	Locations     []Location     `json:"locations"`
	Established   int            `json:"established"`
	ExternalIDs   []ExternalID   `json:"external_ids"`
	Links         []Link         `json:"links"`
	Relationships []Relationship `json:"relationships"`
	Status        string         `json:"status"`
	Types         []string       `json:"types"`
}

// InvenioRDM represents the ROR metadata record in InvenioRDM format.
type InvenioRDM struct {
	ID          string       `json:"id"`
	Identifiers []Identifier `json:"identifiers"`
	Name        string       `json:"name"`
	Title       Title        `json:"title"`
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
	Identifier string `json:"identifier"`
	Scheme     string `json:"scheme"`
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
	Aa string `yaml:"aa,omitempty"` // Afar
	Af string `yaml:"af,omitempty"` // Afrikaans
	Am string `yaml:"am,omitempty"` // Amharic
	Ar string `yaml:"ar,omitempty"` // Arabic
	As string `yaml:"as,omitempty"` // Assamese
	Az string `yaml:"az,omitempty"` // Azerbaijani
	Ba string `yaml:"ba,omitempty"` // Bashkir
	Be string `yaml:"be,omitempty"` // Belgian
	Bg string `yaml:"bg,omitempty"` // Bulgarian
	Bi string `yaml:"bi,omitempty"` // Bislama
	Bn string `yaml:"bn,omitempty"` // Bengali
	Bs string `yaml:"bs,omitempty"` // Bosnian
	Ca string `yaml:"ca,omitempty"` // Catalan
	Ch string `yaml:"ch,omitempty"` // Chamorro
	Co string `yaml:"co,omitempty"` // Corsican
	Cs string `yaml:"cs,omitempty"` // Czech
	Cu string `yaml:"cu,omitempty"` // Church Slavic
	Cy string `yaml:"cy,omitempty"` // Welsh
	Da string `yaml:"da,omitempty"` // Danish
	De string `yaml:"de,omitempty"` // German
	Dv string `yaml:"dv,omitempty"` // Divehi
	Dz string `yaml:"dz,omitempty"` // Dzongkha
	El string `yaml:"el,omitempty"` // Greek
	En string `yaml:"en,omitempty"` // English
	Es string `yaml:"es,omitempty"` // Spanish
	Et string `yaml:"et,omitempty"` // Estonian
	Eu string `yaml:"eu,omitempty"` // Basque
	Fa string `yaml:"fa,omitempty"` // Persian
	Fo string `yaml:"fo,omitempty"` // Faroese
	Fi string `yaml:"fi,omitempty"` // Finnish
	Fr string `yaml:"fr,omitempty"` // French
	Fy string `yaml:"fy,omitempty"` // Frisian
	Ga string `yaml:"ga,omitempty"` // Irish
	Gd string `yaml:"gd,omitempty"` // Scottish Gaelic
	Gl string `yaml:"gl,omitempty"` // Galician
	Gu string `yaml:"gu,omitempty"` // Gujarati
	Gv string `yaml:"gv,omitempty"` // Manx
	Ha string `yaml:"ha,omitempty"` // Hausa
	He string `yaml:"he,omitempty"` // Hebrew
	Hi string `yaml:"hi,omitempty"` // Hindi
	Hr string `yaml:"hr,omitempty"` // Croatian
	Ht string `yaml:"ht,omitempty"` // Haitian
	Hu string `yaml:"hu,omitempty"` // Hungarian
	Hy string `yaml:"hy,omitempty"` // Armenian
	Id string `yaml:"id,omitempty"` // Indonesian
	Is string `yaml:"is,omitempty"` // Icelandic
	It string `yaml:"it,omitempty"` // Italian
	Iu string `yaml:"iu,omitempty"` // Inuktitut
	Ja string `yaml:"ja,omitempty"` // Japanese
	Jv string `yaml:"jv,omitempty"` // Javanese
	Ka string `yaml:"ka,omitempty"` // Georgian
	Kg string `yaml:"kg,omitempty"` // Kongo
	Ki string `yaml:"ki,omitempty"` // Kikuyu
	Kk string `yaml:"kk,omitempty"` // Kazakh
	Kl string `yaml:"kl,omitempty"` // Greenlandic
	Km string `yaml:"km,omitempty"` // Khmer
	Kn string `yaml:"kn,omitempty"` // Kannada
	Ko string `yaml:"ko,omitempty"` // Korean
	Kr string `yaml:"kr,omitempty"` // Kanuri
	Ku string `yaml:"ku,omitempty"` // Kurdish
	Ky string `yaml:"ky,omitempty"` // Kyrgyz
	La string `yaml:"la,omitempty"` // Latin
	Lb string `yaml:"lb,omitempty"` // Luxembourgish
	Lo string `yaml:"lo,omitempty"` // Lao
	Lt string `yaml:"lt,omitempty"` // Lithuanian
	Lu string `yaml:"lu,omitempty"` // Luba-Katanga
	Lv string `yaml:"lv,omitempty"` // Latvian
	Mg string `yaml:"mg,omitempty"` // Malagasy
	Mi string `yaml:"mi,omitempty"` // Maori
	Mk string `yaml:"mk,omitempty"` // Macedonian
	Ml string `yaml:"ml,omitempty"` // Malayalam
	Mn string `yaml:"mn,omitempty"` // Mongolian
	Mr string `yaml:"mr,omitempty"` // Marathi
	Ms string `yaml:"ms,omitempty"` // Malay
	Mt string `yaml:"mt,omitempty"` // Maltese
	My string `yaml:"my,omitempty"` // Burmese
	Na string `yaml:"na,omitempty"` // Nauru
	Nb string `yaml:"nb,omitempty"` // Norwegian Bokmål
	Ne string `yaml:"ne,omitempty"` // Nepali
	Nl string `yaml:"nl,omitempty"` // Dutch
	Nn string `yaml:"nn,omitempty"` // Norwegian Nynorsk
	No string `yaml:"no,omitempty"` // Norwegian
	Oc string `yaml:"oc,omitempty"` // Occitan
	Om string `yaml:"om,omitempty"` // Oromo
	Or string `yaml:"or,omitempty"` // Oriya
	Pa string `yaml:"pa,omitempty"` // Punjabi
	Pl string `yaml:"pl,omitempty"` // Polish
	Ps string `yaml:"ps,omitempty"` // Pashto
	Pt string `yaml:"pt,omitempty"` // Portuguese
	Rm string `yaml:"rm,omitempty"` // Romansh
	Ro string `yaml:"ro,omitempty"` // Romanian
	Ru string `yaml:"ru,omitempty"` // Russian
	Rw string `yaml:"rw,omitempty"` // Kinyarwanda
	Sa string `yaml:"sa,omitempty"` // Sanskrit
	Sd string `yaml:"sd,omitempty"` // Sindhi
	Se string `yaml:"se,omitempty"` // Northern Sami
	Sh string `yaml:"sh,omitempty"` // Serbo-Croatian
	Si string `yaml:"si,omitempty"` // Sinhalese
	Sk string `yaml:"sk,omitempty"` // Slovak
	Sl string `yaml:"sl,omitempty"` // Slovenian
	Sm string `yaml:"sm,omitempty"` // Samoan
	So string `yaml:"so,omitempty"` // Somali
	Sq string `yaml:"sq,omitempty"` // Albanian
	Sr string `yaml:"sr,omitempty"` // Serbian
	St string `yaml:"st,omitempty"` // Southern Sotho
	Sv string `yaml:"sv,omitempty"` // Swedish
	Sw string `yaml:"sw,omitempty"` // Swahili
	Ta string `yaml:"ta,omitempty"` // Tamil
	Te string `yaml:"te,omitempty"` // Telugu
	Tg string `yaml:"tg,omitempty"` // Tajik
	Th string `yaml:"th,omitempty"` // Thai
	Ti string `yaml:"ti,omitempty"` // Tigrinya
	Tk string `yaml:"tk,omitempty"` // Turkmen
	Tl string `yaml:"tl,omitempty"` // Tagalog
	Tr string `yaml:"tr,omitempty"` // Turkish
	Tt string `yaml:"tt,omitempty"` // Tatar
	Ug string `yaml:"ug,omitempty"` // Uighur
	Uk string `yaml:"uk,omitempty"` // Ukrainian
	Ur string `yaml:"ur,omitempty"` // Urdu
	Uz string `yaml:"uz,omitempty"` // Uzbek
	Vi string `yaml:"vi,omitempty"` // Vietnamese
	Xh string `yaml:"xh,omitempty"` // Xhosa
	Yo string `yaml:"yo,omitempty"` // Yoruba
	Zh string `yaml:"zh,omitempty"` // Chinese
	Zu string `yaml:"zu,omitempty"` // Zulu
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

// LoadAll loads the metadata for a list of organizations from a ROR JSON file
func LoadAll(filename string) ([]ROR, error) {
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

	data, err = ReadAll(content)
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
	data.Names = content.Names
	data.Admin.LastModified.Date = content.Admin.LastModified.Date

	return data, nil
}

// ReadAll reads a list of ROR JSON organizations
func ReadAll(content []Content) ([]ROR, error) {
	var data []ROR
	for _, v := range content {
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

	output, err := yaml.Marshal(extracted)
	return output, err
}

// Convert converts ROR metadata into InvenioRDM format.
func Convert(data ROR) (InvenioRDM, error) {
	var inveniordm InvenioRDM

	id, _ := utils.ValidateROR(data.ID)
	inveniordm.ID = id
	inveniordm.Identifiers = []Identifier{
		{
			Identifier: id,
			Scheme:     "ror",
		},
	}
	for _, name := range data.Names {
		if slices.Contains(name.Types, "ror_display") {
			inveniordm.Name = name.Value
		}
	}
	inveniordm.Title = GetTitle(data.Names)
	return inveniordm, nil
}

// Write writes ROR metadata to InvenioRDM YAML format.
func Write(data ROR) ([]byte, error) {
	inveniordm, err := Convert(data)
	if err != nil {
		fmt.Println(err)
	}
	output, err := yaml.Marshal(inveniordm)
	return output, err
}

// WriteAll writes a list of ROR metadata in InvenioRDM YAML format.
func WriteAll(list []ROR, to string) ([]byte, error) {
	var inveniordmList []InvenioRDM
	var err error
	var output []byte

	if to != "inveniordm" {
		return output, errors.New("unsupported output format")
	}

	for _, data := range list {
		inveniordm, err := Convert(data)
		if err != nil {
			fmt.Println(err)
		}
		if inveniordm.ID != "" {
			inveniordmList = append(inveniordmList, inveniordm)
		}
	}

	output, err = yaml.Marshal(inveniordmList)
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
