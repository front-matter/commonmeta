package data

// import (
// 	"commonmeta/cff"
// 	"commonmeta/codemeta"
// 	"commonmeta/commonmeta"
// 	"commonmeta/crossref"
// 	"commonmeta/crossrefxml"
// 	"commonmeta/csl"
// 	"commonmeta/datacite"
// 	"commonmeta/inveniordm"
// 	"commonmeta/jsonfeed"
// 	"commonmeta/schemaorg"
// 	"commonmeta/types"
// 	"encoding/json"
// )

// type Content interface {
// 	crossref.Content |
// 		crossrefxml.Content |
// 		codemeta.Content |
// 		schemaorg.Content |
// 		csl.Content |
// 		datacite.Content |
// 		jsonfeed.Content |
// 		inveniordm.Content |
// 		cff.Content
// }

// type Data struct {
// 	*types.Data
// }

// // func (c *Content) NewMetadata(str string, via string) (types.Data, error) {
// // 	var data types.Data
// // 	id := utils.NormalizeID(str)
// // 	res := GetMetadata(id, "")
// // 	// p := utils.Params{
// // 	// 	id: id,
// // 	// }
// // 	// via := utils.FindFromFormat(p)
// // 	return data, nil
// // }

// func GetMetadata(id string, str string) (Content, error) {
// 	var content Content
// 	if id != "" {
// 		if c.Via == "schema_org" {
// 			return schemaorg.GetSchemaOrg(id)
// 		} else if c.Via == "datacite" {
// 			return datacite.GetDatacite(id)
// 		} else if c.Via == "crossref" || c.Via == "op" {
// 			return crossref.GetCrossref(id)
// 		} else if c.Via == "crossref_xml" {
// 			return crossrefxml.GetCrossrefXML(id)
// 		} else if c.Via == "codemeta" {
// 			return codemeta.GetCodemeta(id)
// 		} else if c.Via == "cff" {
// 			return cff.GetCFF(id)
// 		} else if c.Via == "json_feed_item" {
// 			return jsonfeed.GetJsonFeedItem(id)
// 		} else if c.Via == "inveniordm" {
// 			return inveniordm.GetInvenioRDM(id)
// 		}
// 	} else if str != "" {
// 		if c.Via == "datacite_xml" {
// 			panic("Datacite XML not supported")
// 			// return ParseXML(str)
// 		} else if c.Via == "crossref_xml" {
// 			panic("Crossref XML not supported")
// 			// return ParseXML(str, "crossref")
// 		} else if c.Via == "cff" {
// 			panic("CFF not supported")
// 			// return ParseYAML(str)
// 		} else if c.Via == "bibtex" {
// 			panic("Bibtex not supported")
// 		} else if c.Via == "ris" {
// 			panic("RIS not supported")
// 			// return ParseRIS(str)
// 		} else if c.Via == "commonmeta" || c.Via == "crossref" || c.Via == "datacite" || c.Via == "schema_org" || c.Via == "csl" || c.Via == "json_feed_item" || c.Via == "codemeta" || c.Via == "kbase" || c.Via == "inveniordm" {
// 			err := json.Unmarshal([]byte(str), &content)
// 			if err != nil {
// 				return content, err
// 			}
// 			return content, nil
// 		}
// 	}
// 	return content, nil
// }

// // Parse metadata into Commonmeta format
// func ReadMetadata(content Content) (types.Data, error) {
// 	var data types.Data
// 	if content.Via == "commonmeta" {
// 		return commonmeta.ReadCommonmeta(content)
// 	}
// 	if content.Via == "schema_org" {
// 		return schemaorg.ReadSchemaorg(content)
// 	}
// 	if content.Via == "datacite" {
// 		return datacite.ReadDatacite(content)
// 	}
// 	if content.Via == "crossref" || content.Via == "op" {
// 		return crossref.ReadCrossref(content)
// 	}
// 	if content.Via == "crossref_xml" {
// 		return crossrefxml.ReadCrossrefXML(content)
// 	}
// 	if content.Via == "csl" {
// 		return csl.ReadCsl(content)
// 	}
// 	if content.Via == "codemeta" {
// 		return codemeta.ReadCodemeta(content)
// 	}
// 	if content.Via == "cff" {
// 		return cff.ReadCFF(content)
// 	}
// 	if content.Via == "json_feed_item" {
// 		return jsonfeed.ReadJsonFeedItem(content)
// 	}
// 	if content.Via == "inveniordm" {
// 		return inveniordm.ReadInvenioRDM(content)
// 	}
// 	return data, nil
// }
