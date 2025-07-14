package commonmeta

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/front-matter/commonmeta/schemautils"
	"github.com/front-matter/commonmeta/utils"
	"gopkg.in/yaml.v3"
)

type Writer struct {
	w *bufio.Writer
}

// NewWriter returns a new Writer that writes to w.
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		w: bufio.NewWriter(w),
	}
}

// Write writes commonmeta metadata.
func Write(data Data) ([]byte, error) {
	output, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	err = schemautils.JSONSchemaErrors(output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// WriteAll writes commonmeta metadata in slice format into different serialization formats.
func WriteAll(list []Data, extension string) ([]byte, error) {
	var output []byte
	var err error
	switch extension {
	case ".yaml":
		output, err = yaml.Marshal(list)
		if err != nil {
			return nil, err
		}
	case ".json":
		output, err = json.Marshal(list)
		if err != nil {
			fmt.Println(err)
		}
		err = schemautils.JSONSchemaErrors(output)
		if err != nil {
			return nil, err
		}
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
			var countryCode, countryName, cityName string

			mainInsert := fmt.Sprintf("INSERT INTO organizations ("+
				"id, status, country_code, country_name, city, "+
				"wikipedia_url, website_url) "+
				"VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s');\n",
				utils.EscapeSQL(item.ID),
				utils.EscapeSQL(status),
				utils.EscapeSQL(countryCode),
				utils.EscapeSQL(countryName),
				utils.EscapeSQL(cityName),
				utils.EscapeSQL(wikipedia),
				utils.EscapeSQL(website))
			buffer.WriteString(mainInsert)
		}
		buffer.WriteString("\nCOMMIT;\n")

		output = buffer.Bytes()
	default:
		return output, errors.New("unsupported file format")
	}

	return output, nil
}
