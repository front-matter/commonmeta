package ror

import (
	"regexp"
	"sort"
	"strings"
	"unicode"

	"github.com/google/go-cmp/cmp"
	"github.com/texttheater/golang-levenshtein/levenshtein"
)

// Constants
const (
	MIN_CHOSEN_SCORE   = 0.9
	MIN_MATCHING_SCORE = 0.5

	MATCHING_TYPE_PHRASE     = "PHRASE"
	MATCHING_TYPE_COMMON     = "COMMON TERMS"
	MATCHING_TYPE_FUZZY      = "FUZZY"
	MATCHING_TYPE_HEURISTICS = "HEURISTICS"
	MATCHING_TYPE_ACRONYM    = "ACRONYM"
	MATCHING_TYPE_EXACT      = "EXACT"

	SPECIAL_CHARS_REGEX = `[\+\-\=\|\>\<\!\(\)\\\{\}\[\]\^"\~\*\?\:\/\.\,\;]`
	DO_NOT_MATCH        = "university hospital"
)

var NODE_MATCHING_TYPES = []string{
	MATCHING_TYPE_PHRASE,
	MATCHING_TYPE_COMMON,
	MATCHING_TYPE_FUZZY,
	MATCHING_TYPE_HEURISTICS,
}

// Global variables (equivalent to Python module-level variables)
var (
	GEONAMES_COUNTRIES map[string]map[string]interface{}
	GEONAMES_CITIES    map[string]map[string]interface{}
	COUNTRIES          [][2]string
)

// RORCountries is a map of country codes to their names used in ror.
// https://github.com/ror-community/ror-api/blob/master/rorapi/common/countries.txt
var RORCountries = map[string][]string{
	"ad": {"andorra"},
	"ae": {"al imarat al arabiyah al muttahidah", "arab emirates", "uae", "united arab emirates"},
	"af": {"afghanestan", "afghanistan"},
	"ag": {"antigua", "antigua and barbuda", "barbuda"},
	"ai": {"anguilla"},
	"al": {"albania", "shqiperia"},
	"am": {"armenia", "hayastan"},
	"an": {"netherlands antilles"},
	"ao": {"angola"},
	"aq": {"antarctica"},
	"ar": {"ar", "argentina", "argentine"},
	"as": {"american samoa"},
	"at": {"austria", "aut", "oesterreich", "osterreich"},
	"au": {"au", "aus", "australia"},
	"aw": {"aruba"},
	"ax": {"aland islands"},
	"az": {"azarbaycan respublikasi", "aze", "azerbaijan"},
	"ba": {"bosna i hercegovina", "bosnia", "bosnia and herzegovina", "herzegovina"},
	"bb": {"barbados"},
	"bd": {"bangladesh"},
	"be": {"bel", "belgie", "belgique", "belgium"},
	"bf": {"bfa", "burkina faso"},
	"bg": {"bulariya", "bulgaria", "republika bulgariya"},
	"bh": {"al bahrayn", "bahrain"},
	"bi": {"burundi"},
	"bj": {"ben", "benin"},
	"bl": {"saint barthelemy"},
	"bm": {"bermuda"},
	"bn": {"brunei", "brunei darussalam"},
	"bo": {"bolivia"},
	"bq": {"bonaire", "bonaire, sint eustatius and saba"},
	"br": {"br", "bra", "brasil", "brazil"},
	"bs": {"bahamas", "the bahamas"},
	"bt": {"bhutan", "drukyul"},
	"bv": {"bouvet island"},
	"bw": {"botswana"},
	"by": {"belarus", "byelarus"},
	"bz": {"belice", "belize"},
	"ca": {"ca", "can", "canada"},
	"cc": {"cocos islands", "keeling islands"},
	"cf": {"central african republic", "republique centrafricaine"},
	"cg": {"congo", "democratic republic of the congo", "republic of the congo", "republique democratique du congo", "republique du congo"},
	"ch": {"che", "schweiz", "suisse", "svizzera", "switzerland"},
	"ci": {"cote d'ivoire"},
	"ck": {"cook islands"},
	"cl": {"chile"},
	"cm": {"cameroon", "cameroun"},
	"cn": {"china", "china (people's republic of)", "chn", "cn", "people's republic of china", "p r china", "p. r. china", "pr china", "republic of china", "r. o. c.", "zhong guo"},
	"co": {"col", "colombia"},
	"cr": {"costa rica", "cri"},
	"cu": {"cuba"},
	"cv": {"cabo verde", "cape verde"},
	"cw": {"curacao"},
	"cx": {"christmas island"},
	"cy": {"cyprus", "kibris", "kypros"},
	"cz": {"ceska republika", "cze", "czech republic"},
	"de": {"de", "deu", "deutschland", "federal republic of germany", "germany"},
	"dj": {"djibouti"},
	"dk": {"danmark", "denmark", "dnk"},
	"dm": {"dominica"},
	"do": {"dominican republic", "republica dominicana"},
	"dz": {"algeria", "al jaza'ir"},
	"ec": {"ecuador"},
	"ee": {"eesti", "est", "estonia"},
	"eg": {"egypt", "misr"},
	"eh": {"western sahara"},
	"er": {"eritrea", "ertra"},
	"es": {"esp", "espana", "spain"},
	"et": {"ethiopia", "yeityop'iya"},
	"eu": {"eec"},
	"fi": {"fin", "finland", "suomi"},
	"fj": {"fiji"},
	"fk": {"falkland islands"},
	"fm": {"federated states of micronesia", "micronesia"},
	"fo": {"faroe islands"},
	"fr": {"fr", "fra", "france", "republique francaise"},
	"ga": {"gabon"},
	"gd": {"grenada"},
	"ge": {"georgia", "sak'art'velo"},
	"gf": {"french guiana"},
	"gg": {"guernsey"},
	"gh": {"ghana"},
	"gi": {"gibraltar"},
	"gl": {"greenland"},
	"gm": {"gambia"},
	"gn": {"guinea", "guinea ecuatorial", "guinee"},
	"gp": {"guadeloupe"},
	"gq": {"equatorial guinea"},
	"gr": {"ellas", "grc", "greece"},
	"gs": {"south georgia"},
	"gt": {"guatemala"},
	"gu": {"guam"},
	"gw": {"guinea-bissau", "guine-bissau"},
	"gy": {"guyana"},
	"hk": {"hong kong"},
	"hm": {"heard island", "heard island and mcdonald islands", "mcdonald islands"},
	"hn": {"honduras"},
	"hr": {"croatia", "hrvatska"},
	"ht": {"haiti"},
	"hu": {"hun", "hungary", "magyarorszag"},
	"id": {"idn", "indonesia"},
	"ie": {"eire", "ireland", "irl"},
	"il": {"isr", "israel", "yisra'el"},
	"im": {"isle of man"},
	"in": {"bharat", "ind", "india"},
	"io": {"british indian ocean territory"},
	"iq": {"al iraq", "iraq"},
	"ir": {"iran", "persia"},
	"is": {"iceland", "isl", "island"},
	"it": {"ita", "italia", "italy"},
	"je": {"jersey"},
	"jm": {"jamaica"},
	"jo": {"al urdun", "jordan"},
	"jp": {"japan", "jpn", "nippon", "tokyo"},
	"ke": {"kenya"},
	"kg": {"kyrgyz respublikasy", "kyrgyzstan"},
	"kh": {"cambodia", "kampuchea"},
	"ki": {"kiribati"},
	"km": {"comores", "comoros"},
	"kn": {"nevis", "saint kitts", "saint kitts and nevis"},
	"ko": {"kosova", "kosovo"},
	"kp": {"choson", "choson-minjujuui-inmin-konghwaguk", "north korea"},
	"kr": {"kor", "korea", "republic of korea", "korea (republic of)", "south korea", "taehan-min'guk"},
	"kw": {"al kuwayt", "kuwait"},
	"ky": {"cayman islands"},
	"kz": {"kazakhstan", "qazaqstan"},
	"la": {"lao (people's democratic republic)", "lao people's democratic republic", "laos", "sathalanalat paxathipatai paxaxon lao"},
	"lb": {"lebanon", "lubnan"},
	"lc": {"saint lucia"},
	"li": {"liechtenstein"},
	"lk": {"sri lanka"},
	"lr": {"liberia"},
	"ls": {"lesotho"},
	"lt": {"lietuva", "lithuania", "ltu"},
	"lu": {"lux", "luxembourg"},
	"lv": {"latvia", "latvija"},
	"ly": {"libya", "libyan arab jamahiriya"},
	"ma": {"al maghrib", "morocco"},
	"mc": {"monaco"},
	"md": {"moldova", "moldova (republic of)", "republic of moldova"},
	"me": {"montenegro"},
	"mf": {"saint martin", "saint martin (french part)"},
	"mg": {"madagascar"},
	"mh": {"marshall islands"},
	"mk": {"federal republic of yugoslavia", "former yugoslav republic of macedonia", "macedonia", "makedonija", "yugoslavia"},
	"ml": {"mali"},
	"mm": {"burma", "myanma naingngandaw", "myanmar"},
	"mn": {"mongolia", "mongol uls"},
	"mo": {"macao"},
	"mp": {"northern mariana islands"},
	"mq": {"martinique"},
	"mr": {"mauritania", "muritaniyah"},
	"ms": {"montserrat"},
	"mt": {"malta"},
	"mu": {"mauritius"},
	"mv": {"dhivehi raajje", "maldives"},
	"mw": {"malawi"},
	"mx": {"mex", "mexico"},
	"my": {"malaysia"},
	"mz": {"mocambique", "mozambique"},
	"na": {"namibia"},
	"nc": {"new caledonia"},
	"ne": {"niger"},
	"nf": {"norfolk island"},
	"ng": {"nigeria"},
	"ni": {"nicaragua"},
	"nl": {"nederland", "netherlands", "nl", "nld", "the netherlands"},
	"no": {"nor", "norge", "norway"},
	"np": {"nepal"},
	"nr": {"nauru"},
	"nu": {"niue"},
	"nz": {"new zealand", "nzl"},
	"om": {"oman", "uman"},
	"pa": {"pan", "panama"},
	"pe": {"peru"},
	"pf": {"french polynesia", "pyf"},
	"pg": {"papua new guinea"},
	"ph": {"philippines", "pilipinas"},
	"pk": {"pakistan"},
	"pl": {"pol", "poland", "polska"},
	"pm": {"miquelon", "saint pierre", "saint pierre and miquelon"},
	"pn": {"pitcairn"},
	"pr": {"puerto rico"},
	"ps": {"palestine"},
	"pt": {"portugal"},
	"pw": {"palau"},
	"py": {"paraguay"},
	"qa": {"qatar"},
	"re": {"reunion"},
	"ro": {"romania", "rou"},
	"rs": {"serbia", "srbija-crna gora"},
	"ru": {"rossiya", "rus", "russia", "russian federation"},
	"rw": {"rwanda"},
	"sa": {"al arabiyah as suudiyah", "saudi arabia"},
	"sb": {"solomon islands"},
	"sc": {"seychelles"},
	"sd": {"as-sudan", "sudan"},
	"se": {"sverige", "swe", "sweden"},
	"sg": {"sgp", "singapore"},
	"sh": {"saint helena"},
	"si": {"slovenia", "slovenija", "svn"},
	"sj": {"jan mayen", "svalbard"},
	"sk": {"slovakia", "slovensko", "svk"},
	"sl": {"sierra leone"},
	"sm": {"san marino"},
	"sn": {"sen", "senegal"},
	"so": {"somalia"},
	"sr": {"suriname"},
	"ss": {"south sudan"},
	"st": {"principe", "sao tome", "sao tome and principe", "sao tome e principe"},
	"sv": {"el salvador", "slv"},
	"sx": {"sint maarten"},
	"sy": {"suriyah", "syria", "syrian arab republic"},
	"sz": {"swaziland"},
	"tc": {"caicos islands", "turks", "turks and caicos islands"},
	"td": {"chad", "tchad"},
	"tf": {"french southern territories"},
	"tg": {"togo"},
	"th": {"muang thai", "thailand"},
	"tj": {"jumhurii tojikistan", "tajikistan"},
	"tk": {"tokelau"},
	"tl": {"timor-leste"},
	"tm": {"turkmenistan"},
	"tn": {"tunisia"},
	"to": {"tonga"},
	"tr": {"tur", "turkey", "turkiye"},
	"tt": {"tobago", "trinidad", "trinidad and tobago"},
	"tv": {"tuvalu"},
	"tw": {"taiwan", "t'ai-wan", "taiwan r.o.c", "twn"},
	"tz": {"tanzania", "tanzania (united republic of)"},
	"ua": {"ukr", "ukraine", "ukrayina"},
	"ug": {"uganda"},
	"uk": {"gb", "gbr", "great britain", "northern ireland", "london", "scotland", "uk", "u. k.", "u.k.", "u. k", "u.k", "united kingdom"},
	"um": {"united states minor outlying islands"},
	"us": {"ak", "al", "alabama", "alaska", "ar", "arizona", "arkansas", "az", "ca", "california", "co", "colorado", "connecticut", "ct", "de", "delaware", "fl", "florida", "ga", "georgia", "hawaii", "hi", "ia", "id", "idaho", "il", "ill", "illinois", "in", "indiana", "iowa", "kansas", "kentucky", "ks", "ky", "la", "louisiana", "ma", "maine", "maryland", "massachusetts", "me", "mi", "michigan", "minnesota", "mississippi", "missouri", "mn", "mo", "montana", "ms", "mt", "nc", "nd", "ne", "nebraska", "nevada", "new hampshire", "new jersey", "new mexico", "new york", "nh", "nj", "nm", "north carolina", "north dakota", "nv", "ny", "oh", "ohio", "ok", "okla", "oklahoma", "or", "oregon", "pa", "pennsylvania", "rhode island", "ri", "sc", "sd", "south carolina", "south dakota", "tennessee", "texas", "tn", "tx", "united states", "united states of america", "us", "u. s.", "u.s.", "u. s", "u.s", "usa", "u. s. a.", "u.s.a.", "ut", "utah", "va", "vermont", "virginia", "vt", "wa", "washington", "west virginia", "wi", "wisconsin", "wv", "wy", "wyoming"},
	"uy": {"uruguay"},
	"uz": {"uzbekistan", "uzbekiston respublikasi"},
	"va": {"citta del vaticano", "holy see", "santa sede", "vatican", "vatican city"},
	"vc": {"saint vincent", "saint vincent and the grenadines", "the grenadines"},
	"ve": {"venezuela"},
	"vg": {"virgin islands"},
	"vn": {"vietnam", "viet nam"},
	"vu": {"vanuatu"},
	"wf": {"futuna", "wallis", "wallis and futuna"},
	"ws": {"samoa"},
	"ye": {"al yaman", "yemen"},
	"yt": {"mayotte"},
	"za": {"south africa", "zaf"},
	"zm": {"zambia"},
	"zw": {"zimbabwe"},
}

// MatchedOrganization represents a matched organization from a match query.
type MatchedOrganization struct {
	Chosen       bool
	Substring    string
	MatchingType string
	Score        float64
	Organization ROR
}

// Options represents the options for Levenshhein matching.
type Options struct {
	InsCost int
	DelCost int
	SubCost int
	Matches MatchFunction
}

// MatchFunction is a function type used for Levenshtein matching.
type MatchFunction func(rune, rune) bool

// #####################################################################
// # Country extraction                                                #
// #####################################################################

// ToRegion maps country code to "region" string
func ToRegion(c string) string {
	regionMap := map[string]string{
		"GB": "GB-UK",
		"UK": "GB-UK",
		"CN": "CN-HK-TW",
		"HK": "CN-HK-TW",
		"TW": "CN-HK-TW",
		"PR": "US-PR",
		"US": "US-PR",
	}

	if region, ok := regionMap[c]; ok {
		return region
	}
	return c
}

// GetCountryCodes extracts country codes from the string
func GetCountryCodes(s string) []string {
	// Normalize and clean input string
	s = normalizeString(s)

	// Create different string variants for matching
	lower := createLowerCase(s)
	lowerAlpha := createLowerAlphaOnly(s)
	alpha := createAlphaOnly(s)

	// Store found codes in a map to ensure uniqueness
	codesMap := make(map[string]bool)

	// Threshold for considering a match
	const threshold float64 = 90.0

	// Search for countries based on fuzzy matching
	for code, names := range RORCountries {
		if matchCountryNames(names, lower, lowerAlpha, alpha, threshold) {
			codesMap[strings.ToUpper(code)] = true
		}
	}

	// Convert map keys to slice
	return mapKeysToSlice(codesMap)
}

// normalizeString normalizes the input string by trimming spaces
func normalizeString(s string) string {
	return strings.TrimSpace(s)
}

// createLowerCase converts to lowercase and normalizes whitespace
func createLowerCase(s string) string {
	return regexp.MustCompile(`\s+`).ReplaceAllString(strings.ToLower(s), " ")
}

// createLowerAlphaOnly creates lowercase version with only a-z and spaces
func createLowerAlphaOnly(s string) string {
	return regexp.MustCompile(`\s+`).ReplaceAllString(
		regexp.MustCompile(`[^a-z]`).ReplaceAllString(strings.ToLower(s), " "), " ")
}

// createAlphaOnly creates version with only a-zA-Z and spaces
func createAlphaOnly(s string) string {
	return regexp.MustCompile(`\s+`).ReplaceAllString(
		regexp.MustCompile(`[^a-zA-Z]`).ReplaceAllString(s, " "), " ")
}

// matchCountryNames tries to match country names against string variants
func matchCountryNames(names []string, lower, lowerAlpha, alpha string, threshold float64) bool {
	var op levenshtein.Options

	for _, name := range names {
		var score float64

		// Check if name contains non-a-z characters
		if regexp.MustCompile(`[^a-z]`).MatchString(name) {
			// For names with non-alphabetic characters, use partial ratio
			score = levenshtein.RatioForStrings([]rune(name), []rune(lower), op)
		} else if len(name) == 2 {
			// For 2-letter names, compare with each token in alpha
			score = calculateMaxTokenScore(strings.ToUpper(name), alpha)
		} else {
			// For other names, compare with each token in lower_alpha
			score = calculateMaxTokenScore(name, lowerAlpha)
		}

		// If score is high enough, consider it a match
		if score >= threshold {
			return true
		}
	}

	return false
}

// calculateMaxTokenScore finds the maximum matching score across all tokens
func calculateMaxTokenScore(name, tokenString string) float64 {
	maxScore := 0.0
	var op levenshtein.Options

	for _, token := range strings.Split(tokenString, " ") {
		if len(token) > 0 {
			r := levenshtein.RatioForStrings([]rune(name), []rune(token), op)
			if r > maxScore {
				maxScore = r
			}
		}
	}

	return maxScore
}

// mapKeysToSlice converts map keys to a slice
func mapKeysToSlice(m map[string]bool) []string {
	var slice []string
	for key := range m {
		slice = append(slice, key)
	}
	return slice
}

// GetCountries extracts country codes and maps to regions
func GetCountries(s string) []string {
	codes := GetCountryCodes(s)
	regions := make([]string, len(codes))

	for i, code := range codes {
		regions[i] = ToRegion(code)
	}

	return regions
}

// #####################################################################
// # Similarity                                                        #
// #####################################################################

// CheckLatinChars checks if all characters are Latin
func CheckLatinChars(s string) bool {
	for _, ch := range s {
		if unicode.IsLetter(ch) {
			// This is a simplification - in Python, it checks unicodedata.name(ch) for "LATIN"
			if ch > 127 {
				return false
			}
		}
	}
	return true
}

// Normalize normalizes string for matching
func Normalize(s string) string {
	// Mock implementation - would use more regex replacements in real implementation
	s = strings.ToLower(strings.TrimSpace(s))

	// Various replacements would happen here

	return s
}

// GetSimilarity calculates similarity between affiliation substring and candidate name
func GetSimilarity(affSub, candName string) float64 {
	// Mock implementation - would use fuzzy matching in real implementation
	affSub = Normalize(affSub)
	candName = Normalize(candName)

	// Compare the strings and return a similarity score

	return 0.0
}

// GetScore calculates similarity between affiliation substring and candidate
func GetScore(candidate map[string]interface{}, affSub string, countries []string, version string) float64 {
	// Mock implementation - in real implementation would extract names and calculate scores

	return 0.0
}

// #####################################################################
// # Matching                                                          #
// #####################################################################

// MatchByQuery matches affiliation text using specific ES query
func MatchByQuery(text, matchingType string, query interface{}, countries []string, version string) (MatchedOrganization, []MatchedOrganization) {
	// Mock implementation - would execute query and calculate scores in real implementation

	return MatchedOrganization{}, []MatchedOrganization{}
}

// MatchByType matches affiliation text using specific matching mode/type
func MatchByType(text, matchingType string, countries []string, version string) (MatchedOrganization, []MatchedOrganization) {
	// Implementation based on the Python code
	// var fieldsV1 = []string{"name.norm", "aliases.norm", "labels.label.norm"}
	// var fieldsV2 = []string{"names.value.norm"}
	var substrings []string

	if matchingType == MATCHING_TYPE_HEURISTICS {
		h1Regex := regexp.MustCompile(`University of ([^\s]+)`)
		h1Match := h1Regex.FindStringSubmatch(text)
		if h1Match != nil {
			substrings = append(substrings, h1Match[0])
			substrings = append(substrings, h1Match[1]+" University")
		}

		h2Regex := regexp.MustCompile(`([^\s]+) University`)
		h2Match := h2Regex.FindStringSubmatch(text)
		if h2Match != nil {
			substrings = append(substrings, h2Match[0])
			substrings = append(substrings, "University of "+h2Match[1])
		}
	} else if matchingType == MATCHING_TYPE_ACRONYM {
		var iso3Substrings []string
		acronymRegex := regexp.MustCompile(`[A-Z]{3,}`)
		allSubstrings := acronymRegex.FindAllString(text, -1)

		for _, substring := range allSubstrings {
			for _, country := range GEONAMES_COUNTRIES {
				iso3, ok := country["iso3"].(string)
				if ok && strings.EqualFold(substring, iso3) {
					iso3Substrings = append(iso3Substrings, substring)
				}
			}
		}

		for _, x := range allSubstrings {
			isISO3 := false
			for _, iso3 := range iso3Substrings {
				if x == iso3 {
					isISO3 = true
					break
				}
			}
			if !isISO3 {
				substrings = append(substrings, x)
			}
		}
	} else {
		substrings = append(substrings, text)
	}

	// Mock query builders
	type ESQueryBuilder struct {
		version string
	}

	queries := make([]ESQueryBuilder, len(substrings))
	for i := range queries {
		queries[i] = ESQueryBuilder{version: version}
	}

	// var fields []string
	// if version == "v2" {
	// 	fields = fieldsV2
	// } else {
	// 	fields = fieldsV1
	// }

	// Mock execution of queries
	var matched []struct {
		chosen     MatchedOrganization
		allMatched []MatchedOrganization
	}

	// If no matches, return empty result
	if len(matched) == 0 {
		emptyOrg := MatchedOrganization{
			Substring:    text,
			MatchingType: matchingType,
		}
		return emptyOrg, []MatchedOrganization{}
	}

	var allMatched []MatchedOrganization
	for _, m := range matched {
		allMatched = append(allMatched, m.allMatched...)
	}

	maxScore := 0.0
	for _, m := range matched {
		if m.chosen.Score > maxScore {
			maxScore = m.chosen.Score
		}
	}

	var chosen MatchedOrganization
	for _, m := range matched {
		if m.chosen.Score == maxScore {
			chosen = m.chosen
			break
		}
	}

	return chosen, allMatched
}

// MatchingNode represents a substring of the original affiliation
type MatchingNode struct {
	Text       string
	Version    string
	Matched    *MatchedOrganization
	AllMatched []MatchedOrganization
}

// NewMatchingNode creates a new MatchingNode
func NewMatchingNode(text, version string) *MatchingNode {
	return &MatchingNode{
		Text:       text,
		Version:    version,
		AllMatched: []MatchedOrganization{},
	}
}

// Match tries to match the node text to an organization using different matching types
func (node *MatchingNode) Match(countries []string, minScore float64) {
	for _, matchingType := range NODE_MATCHING_TYPES {
		chosen, allMatched := MatchByType(node.Text, matchingType, countries, node.Version)
		node.AllMatched = append(node.AllMatched, allMatched...)

		if node.Matched == nil {
			node.Matched = &chosen
		}

		if node.Matched != nil &&
			chosen.Score > node.Matched.Score &&
			node.Matched.Score < minScore {
			node.Matched = &chosen
		}
	}
}

// CleanSearchString cleans the search string by removing special characters
func CleanSearchString(searchString string) string {
	re := regexp.MustCompile(SPECIAL_CHARS_REGEX)
	cleaned := re.ReplaceAllString(searchString, " ")

	// Replace multiple spaces with a single space
	spaceRe := regexp.MustCompile(`\s+`)
	cleaned = spaceRe.ReplaceAllString(cleaned, " ")

	// Remove postal codes
	postalRe := regexp.MustCompile(`\d{5}`)
	cleaned = postalRe.ReplaceAllString(cleaned, "")

	return strings.TrimSpace(cleaned)
}

// CheckDoNotMatch checks if the search string should not be matched
func CheckDoNotMatch(searchString string) bool {
	if strings.EqualFold(searchString, DO_NOT_MATCH) {
		return true
	}

	for _, country := range GEONAMES_COUNTRIES {
		name, ok1 := country["name"].(string)
		iso, ok2 := country["iso"].(string)
		iso3, ok3 := country["iso3"].(string)

		if (ok1 && strings.EqualFold(searchString, name)) ||
			(ok2 && strings.EqualFold(searchString, iso)) ||
			(ok3 && strings.EqualFold(searchString, iso3)) {
			return true
		}
	}

	for _, city := range GEONAMES_CITIES {
		name, ok := city["name"].(string)
		if ok && strings.EqualFold(searchString, name) {
			return true
		}
	}

	return false
}

// MatchingGraph represents the entire input affiliation
type MatchingGraph struct {
	Nodes       []*MatchingNode
	Version     string
	Affiliation string
}

// NewMatchingGraph creates a new MatchingGraph
func NewMatchingGraph(affiliation, version string) *MatchingGraph {
	graph := &MatchingGraph{
		Nodes:       []*MatchingNode{},
		Version:     version,
		Affiliation: affiliation,
	}

	// Replace &amp; with &
	affiliation = strings.ReplaceAll(affiliation, "&amp;", "&")
	affiliationCleaned := CleanSearchString(affiliation)

	n := NewMatchingNode(affiliationCleaned, version)
	graph.Nodes = append(graph.Nodes, n)

	// Split by commas, semicolons, or colons
	re := regexp.MustCompile(`[,;:]`)
	parts := re.Split(affiliation, -1)

	for _, part := range parts {
		partCleaned := CleanSearchString(strings.TrimSpace(part))
		doNotMatch := CheckDoNotMatch(partCleaned)

		// Do not perform search if substring exactly matches a country name or ISO code
		if !doNotMatch {
			n = NewMatchingNode(partCleaned, version)
			graph.Nodes = append(graph.Nodes, n)
		}
	}

	return graph
}

// RemoveLowScores removes nodes with scores below the minimum
func (graph *MatchingGraph) RemoveLowScores(minScore float64) {
	for _, node := range graph.Nodes {
		if node.Matched != nil && node.Matched.Score < minScore {
			node.Matched = nil
		}
	}
}

// Match tries to match all nodes in the graph
func (graph *MatchingGraph) Match(countries []string, minScore float64) ([]MatchedOrganization, []MatchedOrganization) {
	for _, node := range graph.Nodes {
		node.Match(countries, minScore)
	}

	graph.RemoveLowScores(minScore)

	var chosen []MatchedOrganization
	var allMatched []MatchedOrganization

	for _, node := range graph.Nodes {
		allMatched = append(allMatched, node.AllMatched...)

		if node.Matched != nil {
			// Check if organization ID is already in chosen
			alreadyChosen := false
			// for _, m := range chosen {
			// 	if m.Organization["id"] == node.Matched.Organization["id"] {
			// 		alreadyChosen = true
			// 		break
			// 	}
			// }

			if !alreadyChosen {
				chosen = append(chosen, *node.Matched)
			}
		}
	}

	// acrChosen, acrAllMatched := MatchByType(graph.Affiliation, MATCHING_TYPE_ACRONYM, countries, graph.Version)
	// allMatched = append(allMatched, acrAllMatched...)

	return chosen, allMatched
}

// GetOutput processes matched organizations and returns final output
func GetOutput(chosen interface{}, allMatched []MatchedOrganization, activeOnly bool) []MatchedOrganization {
	// Don't allow multiple results with chosen=True
	var chosenList []MatchedOrganization

	switch v := chosen.(type) {
	case []MatchedOrganization:
		if len(v) > 1 {
			chosenList = []MatchedOrganization{}
		} else {
			chosenList = v
		}
	case MatchedOrganization:
		chosenList = []MatchedOrganization{v}
	}

	typeMap := map[string]int{
		MATCHING_TYPE_EXACT:      5,
		MATCHING_TYPE_PHRASE:     4,
		MATCHING_TYPE_COMMON:     3,
		MATCHING_TYPE_FUZZY:      2,
		MATCHING_TYPE_HEURISTICS: 1,
		MATCHING_TYPE_ACRONYM:    0,
	}

	// Filter by score
	// var filtered []MatchedOrganization
	// for _, m := range allMatched {
	// 	if m.Score > MIN_MATCHING_SCORE {
	// 		if !activeOnly || m.Organization["status"] == "active" {
	// 			filtered = append(filtered, m)
	// 		}
	// 	}
	// }

	// Group by organization ID
	orgGroups := make(map[string][]MatchedOrganization)
	// for _, m := range filtered {
	// 	orgID := m.Organization["id"].(string)
	// 	orgGroups[orgID] = append(orgGroups[orgID], m)
	// }

	// Sort organization IDs
	var orgIDs []string
	for id := range orgGroups {
		orgIDs = append(orgIDs, id)
	}
	sort.Strings(orgIDs)

	var output []MatchedOrganization
	for _, orgID := range orgIDs {
		g := orgGroups[orgID]
		best := g[0]

		for _, c := range g {
			// Check if c is in chosenList
			isChosen := false
			for _, chosen := range chosenList {
				if cmp.Equal(c, chosen) {
					isChosen = true
					break
				}
			}

			if isChosen {
				best = MatchedOrganization{
					Chosen:       true,
					Substring:    c.Substring,
					Score:        c.Score,
					MatchingType: c.MatchingType,
					Organization: c.Organization,
				}
				break
			}

			if c.Score == 1.0 &&
				typeMap[best.MatchingType] == typeMap[MATCHING_TYPE_EXACT] &&
				typeMap[c.MatchingType] == typeMap[MATCHING_TYPE_EXACT] {
				best = MatchedOrganization{
					Chosen:       true,
					Substring:    c.Substring,
					Score:        c.Score,
					MatchingType: c.MatchingType,
					Organization: c.Organization,
				}
				break
			}

			if best.Score < c.Score {
				best = c
			}

			if best.Score == c.Score &&
				typeMap[best.MatchingType] < typeMap[c.MatchingType] {
				best = c
			}

			if best.Score == c.Score &&
				typeMap[best.MatchingType] == typeMap[c.MatchingType] &&
				len(best.Substring) >= len(c.Substring) {
				best = c
			}
		}

		output = append(output, best)
	}

	// Sort output by score in descending order
	sort.Slice(output, func(i, j int) bool {
		return output[i].Score > output[j].Score
	})

	// Return only the top 100 results
	if len(output) > 100 {
		return output[:100]
	}

	return output
}

// CheckExactMatch checks for exact match of affiliation
func CheckExactMatch(affiliation string, countries []string, version string) (MatchedOrganization, []MatchedOrganization) {
	// Mock implementation

	return MatchedOrganization{}, []MatchedOrganization{}
}

// MatchAffiliation matches an affiliation string
func MatchAffiliation(affiliation string, activeOnly bool, version string) []MatchedOrganization {
	countries := GetCountries(affiliation)
	exactChosen, exactAllMatched := CheckExactMatch(affiliation, countries, version)

	if exactChosen.Score == 1.0 {
		return GetOutput(exactChosen, exactAllMatched, activeOnly)
	} else {
		graph := NewMatchingGraph(affiliation, version)
		chosen, allMatched := graph.Match(countries, MIN_CHOSEN_SCORE)
		return GetOutput(chosen, allMatched, activeOnly)
	}
}

// MatchOrganizations matches organizations based on parameters
func MatchOrganizations(params map[string]string, version string) (interface{}, interface{}) {
	// Mock implementation

	return nil, nil
}
