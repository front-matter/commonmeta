![GitHub](https://img.shields.io/github/license/front-matter/commonmeta-go?logo=MIT)

# commonmeta-go

commonmeta-go is a Go library to implement Commonmeta, the common Metadata Model for Scholarly Metadata. Use commonmeta-go to convert scholarly metadata, in a variety of formats, listed below. Commonmeta-go is work in progress, the first release was on April 19, 2024. Implementations in other languages are also available ([Ruby](https://github.com/front-matter/commonmeta-ruby), [Python](https://github.com/front-matter/commonmeta-py)).

commonmeta-go uses semantic versioning. Currently, its major version number is still at 0, meaning the API is not yet stable, and breaking changes are expected in the internal API and commonmeta JSON format.


## Supported Metadata Formats

Commometa-go reads and/or writes these metadata formats:

| Format                                                                                           | Name          | Content Type                           | Read    | Write   |
| ------------------------------------------------------------------------------------------------ | ------------- | -------------------------------------- | ------- | ------- |
| Commonmeta  | commonmeta    | application/vnd.commonmeta+json        | later     | yes     |
| [CrossRef XML](https://www.crossref.org/schema/documentation/unixref1.1/unixref1.1.html) | crossref_xml      | application/vnd.crossref.unixref+xml   | later | later |
| [Crossref](https://api.crossref.org)                                                             | crossref | application/vnd.crossref+json          | yes     | n/a     |
| [DataCite](https://api.datacite.org/)                                                            | datacite | application/vnd.datacite.datacite+json | yes     | later |
| [Schema.org (in JSON-LD)](http://schema.org/)                                                    | schema_org    | application/vnd.schemaorg.ld+json      | later     | later   |
| [RDF XML](http://www.w3.org/TR/rdf-syntax-grammar/)                                              | rdf_xml       | application/rdf+xml                    | no      | later   |
| [RDF Turtle](http://www.w3.org/TeamSubmission/turtle/)                                           | turtle        | text/turtle                            | no      | later   |
| [CSL-JSON](https://citationstyles.org/)                                                     | csl      | application/vnd.citationstyles.csl+json | later | later   |
| [Formatted text citation](https://citationstyles.org/)                                           | citation      | text/x-bibliography                    | n/a     | yes     |
| [Codemeta](https://codemeta.github.io/)                                                          | codemeta      | application/vnd.codemeta.ld+json       | later | later |
| [Citation File Format (CFF)](https://citation-file-format.github.io/)                            | cff           | application/vnd.cff+yaml               | later | later |
| [JATS](https://jats.nlm.nih.gov/)                                                                | jats          | application/vnd.jats+xml               | later   | later   |
| [CSV](ttps://en.wikipedia.org/wiki/Comma-separated_values)                                       | csv           | text/csv                               | no      | later   |
| [BibTex](http://en.wikipedia.org/wiki/BibTeX)                                                    | bibtex        | application/x-bibtex                   | later | later   |
| [RIS](http://en.wikipedia.org/wiki/RIS_(file_format))                                            | ris           | application/x-research-info-systems    | later | later   |
| [InvenioRDM](https://inveniordm.docs.cern.ch/reference/metadata/)                                | inveniordm    | application/vnd.inveniordm.v1+json     | later | later   |
| [JSON Feed](https://www.jsonfeed.org/)                                                           | json_feed_item     | application/feed+json    | later | later     |

_commonmeta_: the Commonmeta format is the native format for the library and used internally.
_Planned_: we plan to implement this format for the v1.0 public release.  
_Later_: we plan to implement this format in a later release.

## Meta

Please note that this project is released with a [Contributor Code of Conduct](https://github.com/front-matter/commonmeta-go/blob/main/CODE_OF_CONDUCT.md). By participating in this project you agree to abide by its terms.  

License: [MIT](https://github.com/front-matter/commonmeta-go/blob/main/LICENSE)
