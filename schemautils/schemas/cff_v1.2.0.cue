package compose

import (
	"strings"
	"list"
	"time"
	"net"
)

#schema: {
	// Citation File Format
	//
	// A file with citation metadata for software or datasets.
	@jsonschema(schema="http://json-schema.org/draft-07/schema#")
	close({
		@jsonschema(id="https://citation-file-format.github.io/1.2.0/schema.json")

		// A description of the software or dataset.
		abstract?: strings.MinRunes(1)

		// The author(s) of the software or dataset.
		authors!: list.UniqueItems() & [...matchN(>=1, [#person, #entity])] & [_, ...]

		// The version of CFF used for providing the citation metadata.
		"cff-version"!: =~"^1\\.2\\.0$"
		commit?:        #commit

		// The contact person, group, company, etc. for the software or
		// dataset.
		contact?: list.UniqueItems() & [...matchN(>=1, [#person, #entity])] & [_, ...]
		"date-released"?: #date
		doi?:             #doi

		// The identifiers of the software or dataset.
		identifiers?: list.UniqueItems() & [...#identifier] & [_, ...]

		// Keywords that describe the work.
		keywords?: list.UniqueItems() & [...strings.MinRunes(1)] & [_, ...]
		license?:       #license
		"license-url"?: #url

		// A message to the human reader of the file to let them know what
		// to do with the citation metadata.
		message!:              strings.MinRunes(1) | *"If you use this software, please cite it using the metadata from this file."
		"preferred-citation"?: #reference

		// Reference(s) to other creative works.
		references?: list.UniqueItems() & [...#reference] & [_, ...]
		repository?:            #url
		"repository-artifact"?: #url
		"repository-code"?:     #url

		// The name of the software or dataset.
		title!: strings.MinRunes(1)

		// The type of the work.
		type?:    "dataset" | "software" | *"software"
		url?:     #url
		version?: #version
	})

	#address: strings.MinRunes(1)

	#alias: strings.MinRunes(1)

	#city: strings.MinRunes(1)

	#commit: strings.MinRunes(1)

	#country: "AD" | "AE" | "AF" | "AG" | "AI" | "AL" | "AM" | "AO" | "AQ" | "AR" | "AS" | "AT" | "AU" | "AW" | "AX" | "AZ" | "BA" | "BB" | "BD" | "BE" | "BF" | "BG" | "BH" | "BI" | "BJ" | "BL" | "BM" | "BN" | "BO" | "BQ" | "BR" | "BS" | "BT" | "BV" | "BW" | "BY" | "BZ" | "CA" | "CC" | "CD" | "CF" | "CG" | "CH" | "CI" | "CK" | "CL" | "CM" | "CN" | "CO" | "CR" | "CU" | "CV" | "CW" | "CX" | "CY" | "CZ" | "DE" | "DJ" | "DK" | "DM" | "DO" | "DZ" | "EC" | "EE" | "EG" | "EH" | "ER" | "ES" | "ET" | "FI" | "FJ" | "FK" | "FM" | "FO" | "FR" | "GA" | "GB" | "GD" | "GE" | "GF" | "GG" | "GH" | "GI" | "GL" | "GM" | "GN" | "GP" | "GQ" | "GR" | "GS" | "GT" | "GU" | "GW" | "GY" | "HK" | "HM" | "HN" | "HR" | "HT" | "HU" | "ID" | "IE" | "IL" | "IM" | "IN" | "IO" | "IQ" | "IR" | "IS" | "IT" | "JE" | "JM" | "JO" | "JP" | "KE" | "KG" | "KH" | "KI" | "KM" | "KN" | "KP" | "KR" | "KW" | "KY" | "KZ" | "LA" | "LB" | "LC" | "LI" | "LK" | "LR" | "LS" | "LT" | "LU" | "LV" | "LY" | "MA" | "MC" | "MD" | "ME" | "MF" | "MG" | "MH" | "MK" | "ML" | "MM" | "MN" | "MO" | "MP" | "MQ" | "MR" | "MS" | "MT" | "MU" | "MV" | "MW" | "MX" | "MY" | "MZ" | "NA" | "NC" | "NE" | "NF" | "NG" | "NI" | "NL" | "NO" | "NP" | "NR" | "NU" | "NZ" | "OM" | "PA" | "PE" | "PF" | "PG" | "PH" | "PK" | "PL" | "PM" | "PN" | "PR" | "PS" | "PT" | "PW" | "PY" | "QA" | "RE" | "RO" | "RS" | "RU" | "RW" | "SA" | "SB" | "SC" | "SD" | "SE" | "SG" | "SH" | "SI" | "SJ" | "SK" | "SL" | "SM" | "SN" | "SO" | "SR" | "SS" | "ST" | "SV" | "SX" | "SY" | "SZ" | "TC" | "TD" | "TF" | "TG" | "TH" | "TJ" | "TK" | "TL" | "TM" | "TN" | "TO" | "TR" | "TT" | "TV" | "TW" | "TZ" | "UA" | "UG" | "UM" | "US" | "UY" | "UZ" | "VA" | "VC" | "VE" | "VG" | "VI" | "VN" | "VU" | "WF" | "WS" | "YE" | "YT" | "ZA" | "ZM" | "ZW"

	#date: time.Format("2006-01-02") & =~"^[0-9]{4}-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])$"

	#doi: =~"^10\\.\\d{4,9}(\\.\\d+)?/[A-Za-z0-9:/_;\\-\\.\\(\\)\\[\\]\\\\]+$"

	#email: =~"^[\\S]+@[\\S]+\\.[\\S]{2,}$"

	#entity: close({
		address?:      #address
		alias?:        #alias
		city?:         #city
		country?:      #country
		"date-end"?:   #date
		"date-start"?: #date
		email?:        #email
		fax?:          #fax

		// The entity's location, e.g., when the entity is a conference.
		location?: strings.MinRunes(1)

		// The entity's name.
		name!:        strings.MinRunes(1)
		orcid?:       #orcid
		"post-code"?: #["post-code"]
		region?:      #region
		tel?:         #tel
		website?:     #url
	})

	#fax: strings.MinRunes(1)

	#identifier: matchN(>=1, [close({
		description?: #["identifier-description"]
		type!:        "doi"
		value!:       #doi
	}), close({
		description?: #["identifier-description"]
		type!:        "url"
		value!:       #url
	}), close({
		description?: #["identifier-description"]
		type!:        "swh"
		value!:       #["swh-identifier"]
	}), close({
		description?: #["identifier-description"]
		type!:        "other"
		value!:       strings.MinRunes(1)
	})])

	#: "identifier-description": strings.MinRunes(1)

	#license: matchN(1, [#["license-enum"], list.UniqueItems() & [...#["license-enum"]] & [_, ...]])

	#: "license-enum": "0BSD" | "AAL" | "Abstyles" | "Adobe-2006" | "Adobe-Glyph" | "ADSL" | "AFL-1.1" | "AFL-1.2" | "AFL-2.0" | "AFL-2.1" | "AFL-3.0" | "Afmparse" | "AGPL-1.0" | "AGPL-1.0-only" | "AGPL-1.0-or-later" | "AGPL-3.0" | "AGPL-3.0-only" | "AGPL-3.0-or-later" | "Aladdin" | "AMDPLPA" | "AML" | "AMPAS" | "ANTLR-PD" | "ANTLR-PD-fallback" | "Apache-1.0" | "Apache-1.1" | "Apache-2.0" | "APAFML" | "APL-1.0" | "APSL-1.0" | "APSL-1.1" | "APSL-1.2" | "APSL-2.0" | "Artistic-1.0" | "Artistic-1.0-cl8" | "Artistic-1.0-Perl" | "Artistic-2.0" | "Bahyph" | "Barr" | "Beerware" | "BitTorrent-1.0" | "BitTorrent-1.1" | "blessing" | "BlueOak-1.0.0" | "Borceux" | "BSD-1-Clause" | "BSD-2-Clause" | "BSD-2-Clause-FreeBSD" | "BSD-2-Clause-NetBSD" | "BSD-2-Clause-Patent" | "BSD-2-Clause-Views" | "BSD-3-Clause" | "BSD-3-Clause-Attribution" | "BSD-3-Clause-Clear" | "BSD-3-Clause-LBNL" | "BSD-3-Clause-Modification" | "BSD-3-Clause-No-Nuclear-License" | "BSD-3-Clause-No-Nuclear-License-2014" | "BSD-3-Clause-No-Nuclear-Warranty" | "BSD-3-Clause-Open-MPI" | "BSD-4-Clause" | "BSD-4-Clause-Shortened" | "BSD-4-Clause-UC" | "BSD-Protection" | "BSD-Source-Code" | "BSL-1.0" | "BUSL-1.1" | "bzip2-1.0.5" | "bzip2-1.0.6" | "C-UDA-1.0" | "CAL-1.0" | "CAL-1.0-Combined-Work-Exception" | "Caldera" | "CATOSL-1.1" | "CC-BY-1.0" | "CC-BY-2.0" | "CC-BY-2.5" | "CC-BY-3.0" | "CC-BY-3.0-AT" | "CC-BY-3.0-US" | "CC-BY-4.0" | "CC-BY-NC-1.0" | "CC-BY-NC-2.0" | "CC-BY-NC-2.5" | "CC-BY-NC-3.0" | "CC-BY-NC-4.0" | "CC-BY-NC-ND-1.0" | "CC-BY-NC-ND-2.0" | "CC-BY-NC-ND-2.5" | "CC-BY-NC-ND-3.0" | "CC-BY-NC-ND-3.0-IGO" | "CC-BY-NC-ND-4.0" | "CC-BY-NC-SA-1.0" | "CC-BY-NC-SA-2.0" | "CC-BY-NC-SA-2.5" | "CC-BY-NC-SA-3.0" | "CC-BY-NC-SA-4.0" | "CC-BY-ND-1.0" | "CC-BY-ND-2.0" | "CC-BY-ND-2.5" | "CC-BY-ND-3.0" | "CC-BY-ND-4.0" | "CC-BY-SA-1.0" | "CC-BY-SA-2.0" | "CC-BY-SA-2.0-UK" | "CC-BY-SA-2.1-JP" | "CC-BY-SA-2.5" | "CC-BY-SA-3.0" | "CC-BY-SA-3.0-AT" | "CC-BY-SA-4.0" | "CC-PDDC" | "CC0-1.0" | "CDDL-1.0" | "CDDL-1.1" | "CDL-1.0" | "CDLA-Permissive-1.0" | "CDLA-Sharing-1.0" | "CECILL-1.0" | "CECILL-1.1" | "CECILL-2.0" | "CECILL-2.1" | "CECILL-B" | "CECILL-C" | "CERN-OHL-1.1" | "CERN-OHL-1.2" | "CERN-OHL-P-2.0" | "CERN-OHL-S-2.0" | "CERN-OHL-W-2.0" | "ClArtistic" | "CNRI-Jython" | "CNRI-Python" | "CNRI-Python-GPL-Compatible" | "Condor-1.1" | "copyleft-next-0.3.0" | "copyleft-next-0.3.1" | "CPAL-1.0" | "CPL-1.0" | "CPOL-1.02" | "Crossword" | "CrystalStacker" | "CUA-OPL-1.0" | "Cube" | "curl" | "D-FSL-1.0" | "diffmark" | "DOC" | "Dotseqn" | "DRL-1.0" | "DSDP" | "dvipdfm" | "ECL-1.0" | "ECL-2.0" | "eCos-2.0" | "EFL-1.0" | "EFL-2.0" | "eGenix" | "Entessa" | "EPICS" | "EPL-1.0" | "EPL-2.0" | "ErlPL-1.1" | "etalab-2.0" | "EUDatagrid" | "EUPL-1.0" | "EUPL-1.1" | "EUPL-1.2" | "Eurosym" | "Fair" | "Frameworx-1.0" | "FreeBSD-DOC" | "FreeImage" | "FSFAP" | "FSFUL" | "FSFULLR" | "FTL" | "GD" | "GFDL-1.1" | "GFDL-1.1-invariants-only" | "GFDL-1.1-invariants-or-later" | "GFDL-1.1-no-invariants-only" | "GFDL-1.1-no-invariants-or-later" | "GFDL-1.1-only" | "GFDL-1.1-or-later" | "GFDL-1.2" | "GFDL-1.2-invariants-only" | "GFDL-1.2-invariants-or-later" | "GFDL-1.2-no-invariants-only" | "GFDL-1.2-no-invariants-or-later" | "GFDL-1.2-only" | "GFDL-1.2-or-later" | "GFDL-1.3" | "GFDL-1.3-invariants-only" | "GFDL-1.3-invariants-or-later" | "GFDL-1.3-no-invariants-only" | "GFDL-1.3-no-invariants-or-later" | "GFDL-1.3-only" | "GFDL-1.3-or-later" | "Giftware" | "GL2PS" | "Glide" | "Glulxe" | "GLWTPL" | "gnuplot" | "GPL-1.0" | "GPL-1.0-only" | "GPL-1.0-or-later" | "GPL-1.0+" | "GPL-2.0" | "GPL-2.0-only" | "GPL-2.0-or-later" | "GPL-2.0-with-autoconf-exception" | "GPL-2.0-with-bison-exception" | "GPL-2.0-with-classpath-exception" | "GPL-2.0-with-font-exception" | "GPL-2.0-with-GCC-exception" | "GPL-2.0+" | "GPL-3.0" | "GPL-3.0-only" | "GPL-3.0-or-later" | "GPL-3.0-with-autoconf-exception" | "GPL-3.0-with-GCC-exception" | "GPL-3.0+" | "gSOAP-1.3b" | "HaskellReport" | "Hippocratic-2.1" | "HPND" | "HPND-sell-variant" | "HTMLTIDY" | "IBM-pibs" | "ICU" | "IJG" | "ImageMagick" | "iMatix" | "Imlib2" | "Info-ZIP" | "Intel" | "Intel-ACPI" | "Interbase-1.0" | "IPA" | "IPL-1.0" | "ISC" | "JasPer-2.0" | "JPNIC" | "JSON" | "LAL-1.2" | "LAL-1.3" | "Latex2e" | "Leptonica" | "LGPL-2.0" | "LGPL-2.0-only" | "LGPL-2.0-or-later" | "LGPL-2.0+" | "LGPL-2.1" | "LGPL-2.1-only" | "LGPL-2.1-or-later" | "LGPL-2.1+" | "LGPL-3.0" | "LGPL-3.0-only" | "LGPL-3.0-or-later" | "LGPL-3.0+" | "LGPLLR" | "Libpng" | "libpng-2.0" | "libselinux-1.0" | "libtiff" | "LiLiQ-P-1.1" | "LiLiQ-R-1.1" | "LiLiQ-Rplus-1.1" | "Linux-OpenIB" | "LPL-1.0" | "LPL-1.02" | "LPPL-1.0" | "LPPL-1.1" | "LPPL-1.2" | "LPPL-1.3a" | "LPPL-1.3c" | "MakeIndex" | "MirOS" | "MIT" | "MIT-0" | "MIT-advertising" | "MIT-CMU" | "MIT-enna" | "MIT-feh" | "MIT-Modern-Variant" | "MIT-open-group" | "MITNFA" | "Motosoto" | "mpich2" | "MPL-1.0" | "MPL-1.1" | "MPL-2.0" | "MPL-2.0-no-copyleft-exception" | "MS-PL" | "MS-RL" | "MTLL" | "MulanPSL-1.0" | "MulanPSL-2.0" | "Multics" | "Mup" | "NAIST-2003" | "NASA-1.3" | "Naumen" | "NBPL-1.0" | "NCGL-UK-2.0" | "NCSA" | "Net-SNMP" | "NetCDF" | "Newsletr" | "NGPL" | "NIST-PD" | "NIST-PD-fallback" | "NLOD-1.0" | "NLPL" | "Nokia" | "NOSL" | "Noweb" | "NPL-1.0" | "NPL-1.1" | "NPOSL-3.0" | "NRL" | "NTP" | "NTP-0" | "Nunit" | "O-UDA-1.0" | "OCCT-PL" | "OCLC-2.0" | "ODbL-1.0" | "ODC-By-1.0" | "OFL-1.0" | "OFL-1.0-no-RFN" | "OFL-1.0-RFN" | "OFL-1.1" | "OFL-1.1-no-RFN" | "OFL-1.1-RFN" | "OGC-1.0" | "OGDL-Taiwan-1.0" | "OGL-Canada-2.0" | "OGL-UK-1.0" | "OGL-UK-2.0" | "OGL-UK-3.0" | "OGTSL" | "OLDAP-1.1" | "OLDAP-1.2" | "OLDAP-1.3" | "OLDAP-1.4" | "OLDAP-2.0" | "OLDAP-2.0.1" | "OLDAP-2.1" | "OLDAP-2.2" | "OLDAP-2.2.1" | "OLDAP-2.2.2" | "OLDAP-2.3" | "OLDAP-2.4" | "OLDAP-2.5" | "OLDAP-2.6" | "OLDAP-2.7" | "OLDAP-2.8" | "OML" | "OpenSSL" | "OPL-1.0" | "OSET-PL-2.1" | "OSL-1.0" | "OSL-1.1" | "OSL-2.0" | "OSL-2.1" | "OSL-3.0" | "Parity-6.0.0" | "Parity-7.0.0" | "PDDL-1.0" | "PHP-3.0" | "PHP-3.01" | "Plexus" | "PolyForm-Noncommercial-1.0.0" | "PolyForm-Small-Business-1.0.0" | "PostgreSQL" | "PSF-2.0" | "psfrag" | "psutils" | "Python-2.0" | "Qhull" | "QPL-1.0" | "Rdisc" | "RHeCos-1.1" | "RPL-1.1" | "RPL-1.5" | "RPSL-1.0" | "RSA-MD" | "RSCPL" | "Ruby" | "SAX-PD" | "Saxpath" | "SCEA" | "Sendmail" | "Sendmail-8.23" | "SGI-B-1.0" | "SGI-B-1.1" | "SGI-B-2.0" | "SHL-0.5" | "SHL-0.51" | "SimPL-2.0" | "SISSL" | "SISSL-1.2" | "Sleepycat" | "SMLNJ" | "SMPPL" | "SNIA" | "Spencer-86" | "Spencer-94" | "Spencer-99" | "SPL-1.0" | "SSH-OpenSSH" | "SSH-short" | "SSPL-1.0" | "StandardML-NJ" | "SugarCRM-1.1.3" | "SWL" | "TAPR-OHL-1.0" | "TCL" | "TCP-wrappers" | "TMate" | "TORQUE-1.1" | "TOSL" | "TU-Berlin-1.0" | "TU-Berlin-2.0" | "UCL-1.0" | "Unicode-DFS-2015" | "Unicode-DFS-2016" | "Unicode-TOU" | "Unlicense" | "UPL-1.0" | "Vim" | "VOSTROM" | "VSL-1.0" | "W3C" | "W3C-19980720" | "W3C-20150513" | "Watcom-1.0" | "Wsuipa" | "WTFPL" | "wxWindows" | "X11" | "Xerox" | "XFree86-1.1" | "xinetd" | "Xnet" | "xpp" | "XSkat" | "YPL-1.0" | "YPL-1.1" | "Zed" | "Zend-2.0" | "Zimbra-1.3" | "Zimbra-1.4" | "Zlib" | "zlib-acknowledgement" | "ZPL-1.1" | "ZPL-2.0" | "ZPL-2.1"

	#orcid: net.AbsURL & =~"https://orcid\\.org/[0-9]{4}-[0-9]{4}-[0-9]{4}-[0-9]{3}[0-9X]{1}"

	#person: close({
		address?: #address

		// The person's affilitation.
		affiliation?: strings.MinRunes(1)
		alias?:       #alias
		city?:        #city
		country?:     #country
		email?:       #email

		// The person's family names.
		"family-names"?: strings.MinRunes(1)
		fax?:            #fax

		// The person's given names.
		"given-names"?: strings.MinRunes(1)

		// The person's name particle, e.g., a nobiliary particle or a
		// preposition meaning 'of' or 'from' (for example 'von' in
		// 'Alexander von Humboldt').
		"name-particle"?: strings.MinRunes(1)

		// The person's name-suffix, e.g. 'Jr.' for Sammy Davis Jr. or
		// 'III' for Frank Edwin Wright III.
		"name-suffix"?: strings.MinRunes(1)
		orcid?:         #orcid
		"post-code"?:   #["post-code"]
		region?:        #region
		tel?:           #tel
		website?:       #url
	})

	#: "post-code": matchN(>=1, [strings.MinRunes(1), number])

	#reference: close({
		// The abbreviation of a work.
		abbreviation?: strings.MinRunes(1)

		// The abstract of a work.
		abstract?: strings.MinRunes(1)

		// The author(s) of a work.
		authors!: list.UniqueItems() & [...matchN(>=1, [#person, #entity])] & [_, ...]
		"collection-doi"?: #doi

		// The title of a collection or proceedings.
		"collection-title"?: strings.MinRunes(1)

		// The type of a collection.
		"collection-type"?: strings.MinRunes(1)
		commit?:            #commit
		conference?:        #entity

		// The contact person, group, company, etc. for a work.
		contact?: list.UniqueItems() & [...matchN(>=1, [#person, #entity])] & [_, ...]

		// The copyright information pertaining to the work.
		copyright?: strings.MinRunes(1)

		// The data type of a data set.
		"data-type"?: strings.MinRunes(1)

		// The name of the database where a work was accessed/is stored.
		database?:            strings.MinRunes(1)
		"database-provider"?: #entity
		"date-accessed"?:     #date
		"date-downloaded"?:   #date
		"date-published"?:    #date
		"date-released"?:     #date

		// The department where a work has been produced.
		department?: strings.MinRunes(1)
		doi?:        #doi

		// The edition of the work.
		edition?: strings.MinRunes(1)

		// The editor(s) of a work.
		editors?: list.UniqueItems() & [...matchN(>=1, [#person, #entity])] & [_, ...]

		// The editor(s) of a series in which a work has been published.
		"editors-series"?: list.UniqueItems() & [...matchN(>=1, [#person, #entity])] & [_, ...]

		// The end page of the work.
		end?: matchN(>=1, [int, strings.MinRunes(1)])

		// An entry in the collection that constitutes the work.
		entry?: strings.MinRunes(1)

		// The name of the electronic file containing the work.
		filename?: strings.MinRunes(1)

		// The format in which a work is represented.
		format?: strings.MinRunes(1)

		// The identifier(s) of the work.
		identifiers?: list.UniqueItems() & [...#identifier] & [_, ...]
		institution?: #entity

		// The ISBN of the work.
		isbn?: =~"^[0-9\\- ]{10,17}X?$"

		// The ISSN of the work.
		issn?: =~"^\\d{4}-\\d{3}[\\dxX]$"

		// The issue of a periodical in which a work appeared.
		issue?: matchN(>=1, [strings.MinRunes(1), number])

		// The publication date of the issue of a periodical in which a
		// work appeared.
		"issue-date"?: strings.MinRunes(1)

		// The name of the issue of a periodical in which the work
		// appeared.
		"issue-title"?: strings.MinRunes(1)

		// The name of the journal/magazine/newspaper/periodical where the
		// work was published.
		journal?: strings.MinRunes(1)

		// Keywords pertaining to the work.
		keywords?: list.UniqueItems() & [...strings.MinRunes(1)] & [_, ...]

		// The language identifier(s) of the work according to ISO 639
		// language strings.
		languages?: list.UniqueItems() & [...strings.MaxRunes(3) & strings.MinRunes(2) & =~"^[a-z]{2,3}$"] & [_, ...]
		license?:       #license
		"license-url"?: #url

		// The line of code in the file where the work ends.
		"loc-end"?: matchN(>=1, [int, strings.MinRunes(1)])

		// The line of code in the file where the work starts.
		"loc-start"?: matchN(>=1, [int, strings.MinRunes(1)])
		location?: #entity

		// The medium of the work.
		medium?: strings.MinRunes(1)

		// The month in which a work has been published.
		month?: matchN(>=1, [int & <=12 & >=1, "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9" | "10" | "11" | "12"])

		// The NIHMSID of a work.
		nihmsid?: strings.MinRunes(1)

		// Notes pertaining to the work.
		notes?: strings.MinRunes(1)

		// The accession number for a work.
		"number"?: matchN(>=1, [strings.MinRunes(1), number])

		// The number of volumes making up the collection in which the
		// work has been published.
		"number-volumes"?: matchN(>=1, [int, strings.MinRunes(1)])

		// The number of pages of the work.
		pages?: matchN(>=1, [int, strings.MinRunes(1)])

		// The states for which a patent is granted.
		"patent-states"?: list.UniqueItems() & [...strings.MinRunes(1)] & [_, ...]

		// The PMCID of a work.
		pmcid?:     =~"^PMC[0-9]{7}$"
		publisher?: #entity

		// The recipient(s) of a personal communication.
		recipients?: list.UniqueItems() & [...matchN(>=1, [#entity, #person])] & [_, ...]
		repository?:            #url
		"repository-artifact"?: #url
		"repository-code"?:     #url

		// The scope of the reference, e.g., the section of the work it
		// adheres to.
		scope?: strings.MinRunes(1)

		// The section of a work that is referenced.
		section?: matchN(>=1, [strings.MinRunes(1), number])

		// The sender(s) of a personal communication.
		senders?: list.UniqueItems() & [...matchN(>=1, [#entity, #person])] & [_, ...]

		// The start page of the work.
		start?: matchN(>=1, [int, strings.MinRunes(1)])

		// The publication status of the work.
		status?: "abstract" | "advance-online" | "in-preparation" | "in-press" | "preprint" | "submitted"

		// The term being referenced if the work is a dictionary or
		// encyclopedia.
		term?: strings.MinRunes(1)

		// The type of the thesis that is the work.
		"thesis-type"?: strings.MinRunes(1)

		// The title of the work.
		title!: strings.MinRunes(1)

		// The translator(s) of a work.
		translators?: list.UniqueItems() & [...matchN(>=1, [#entity, #person])] & [_, ...]

		// The type of the work.
		type!:    "art" | "article" | "audiovisual" | "bill" | "blog" | "book" | "catalogue" | "conference-paper" | "conference" | "data" | "database" | "dictionary" | "edited-work" | "encyclopedia" | "film-broadcast" | "generic" | "government-document" | "grant" | "hearing" | "historical-work" | "legal-case" | "legal-rule" | "magazine-article" | "manual" | "map" | "multimedia" | "music" | "newspaper-article" | "pamphlet" | "patent" | "personal-communication" | "proceedings" | "report" | "serial" | "slides" | "software-code" | "software-container" | "software-executable" | "software-virtual-machine" | "software" | "sound-recording" | "standard" | "statute" | "thesis" | "unpublished" | "video" | "website"
		url?:     #url
		version?: #version

		// The volume of the periodical in which a work appeared.
		volume?: matchN(>=1, [int, strings.MinRunes(1)])

		// The title of the volume in which the work appeared.
		"volume-title"?: strings.MinRunes(1)

		// The year in which a work has been published.
		year?: matchN(>=1, [int, strings.MinRunes(1)])

		// The year of the original publication.
		"year-original"?: matchN(>=1, [int, strings.MinRunes(1)])
	})

	#region: strings.MinRunes(1)

	#: "swh-identifier": =~"^swh:1:(snp|rel|rev|dir|cnt):[0-9a-fA-F]{40}$"

	#tel: strings.MinRunes(1)

	#url: net.AbsURL & =~"^(https|http|ftp|sftp)://.+"

	#version: matchN(>=1, [strings.MinRunes(1), number])
}
