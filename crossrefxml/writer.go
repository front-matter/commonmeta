package crossrefxml

import (
	"encoding/xml"
	"fmt"

	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/front-matter/commonmeta/dateutils"
	"github.com/xeipuuv/gojsonschema"
)

type StringMap map[string]string

// CMToCRMappings maps Commonmeta types to Crossref types
// source: http://api.crossref.org/types
var CMToCRMappings = map[string]string{
	"Article":            "PostedContent",
	"BookChapter":        "BookChapter",
	"BookSeries":         "BookSeries",
	"Book":               "Book",
	"Component":          "Component",
	"Dataset":            "Dataset",
	"Dissertation":       "Dissertation",
	"Grant":              "Grant",
	"JournalArticle":     "JournalArticle",
	"JournalIssue":       "JournalIssue",
	"JournalVolume":      "JournalVolume",
	"Journal":            "Journal",
	"ProceedingsArticle": "ProceedingsArticle",
	"ProceedingsSeries":  "ProceedingsSeries",
	"Proceedings":        "Proceedings",
	"ReportComponent":    "ReportComponent",
	"ReportSeries":       "ReportSeries",
	"Report":             "Report",
	"Review":             "PeerReview",
	"Other":              "Other",
}

// Convert converts Commonmeta metadata to Crossrefxml metadata
func Convert(data commonmeta.Data) (*Crossref, error) {
	c := &Crossref{
		Xmlns:          "http://www.crossref.org/schema/5.3.1",
		SchemaLocation: "http://www.crossref.org/schema/5.3.1 ",
		Version:        "5.3.1",
	}
	abstract := []Abstract{}
	if len(data.Descriptions) > 0 {
		for _, description := range data.Descriptions {
			if description.Type == "Abstract" {
				abstract = append(abstract, Abstract{
					Text: description.Description,
				})
			}
		}
	}
	contributors := &Contributors{}
	doiData := DOIData{
		DOI:      data.ID,
		Resource: data.URL,
	}
	titles := &Titles{}
	if len(data.Titles) > 0 {
		for _, title := range data.Titles {
			if title.Type == "Title" {
				titles.Title = title.Title
			} else if title.Type == "Subtitle" {
				titles.Subtitle = title.Title
			} else if title.Type == "TranslatedTitle" {
				titles.OriginalLanguageTitle.Text = title.Title
				titles.OriginalLanguageTitle.Language = title.Language
			}
		}
	}

	switch data.Type {
	case "Article":
		var groupTitle string
		if len(data.Subjects) > 0 {
			groupTitle = data.Subjects[0].Subject
		}
		var postedDate PostedDate
		if len(data.Date.Published) > 0 {
			datePublished := dateutils.GetDateStruct(data.Date.Published)
			postedDate = PostedDate{
				MediaType: "online",
				Year:      datePublished.Year,
				Month:     datePublished.Month,
				Day:       datePublished.Day,
			}
		}
		c.PostedContent = &PostedContent{
			Type:         "other",
			Language:     data.Language,
			GroupTitle:   groupTitle,
			Contributors: contributors,
			Titles:       titles,
			PostedDate:   postedDate,
			Abstract:     &abstract,
			DOIData:      doiData,
		}
	case "JournalArticle":
		c.Journal = &Journal{}
	}
	// required properties

	// optional properties

	return c, nil
}

// Write writes Crossrefxml metadata.
func Write(data commonmeta.Data) ([]byte, []gojsonschema.ResultError) {
	crossref, err := Convert(data)
	if err != nil {
		fmt.Println(err)
	}
	output, err := xml.MarshalIndent(crossref, "", "  ")
	if err == nil {
		fmt.Println(err)
	}
	output = []byte(xml.Header + string(output))
	return output, nil
}
