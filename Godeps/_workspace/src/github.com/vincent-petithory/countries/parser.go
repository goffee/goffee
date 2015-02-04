// +build ignore

package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"go/format"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/goffee/goffee/Godeps/_workspace/src/github.com/vincent-petithory/countries"
)

func main() {
	goFile := os.Getenv("GOFILE")
	goPkg := os.Getenv("GOPACKAGE")
	if goPkg == "" {
		log.Fatal("GOPACKAGE env is empty")
	}
	if goFile == "" {
		log.Fatal("GOFILE env is empty")
	}
	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	// Parse template for go source
	t, err := template.New("_").Funcs(template.FuncMap{
		"countrySrc": countrySrc,
	}).Parse(tmpl)
	if err != nil {
		log.Fatal(err)
	}

	// Parse csv data
	f, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	csvr := csv.NewReader(f)
	csvr.Comma = ';'
	csvr.Comment = '#'
	csvr.FieldsPerRecord = 27

	var allCountries []countries.Country
	for {
		record, err := csvr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		tr := make([]string, len(record))
		for i, rec := range record {
			tr[i] = strings.TrimSpace(rec)
		}

		iso3166OneNumeric, err := strconv.ParseInt(tr[2], 10, 64)
		if err != nil {
			log.Fatalf("iso3166OneNumeric: %v", err)
		}
		country := countries.Country{
			ISO3166OneAlphaTwo:   tr[0],
			ISO3166OneAlphaThree: tr[1],
			ISO3166OneNumeric:    int(iso3166OneNumeric),

			ISO3166OneEnglishShortNameGazetteerOrder:          tr[3],
			ISO3166OneEnglishShortNameReadingOrder:            tr[4],
			ISO3166OneEnglishRomanizedShortNameGazetteerOrder: tr[5],
			ISO3166OneEnglishRomanizedShortNameReadingOrder:   tr[6],
			ISO3166OneFrenchShortNameGazetteerOrder:           tr[7],
			ISO3166OneFrenchShortNameReadingOrder:             tr[8],
			ISO3166OneSpanishShortNameGazetteerOrder:          tr[9],

			UNGEGNEnglishFormalName: tr[10],
			UNGEGNFrenchFormalName:  tr[11],
			UNGEGNSpanishFormalName: tr[12],
			UNGEGNRussianShortName:  tr[13],
			UNGEGNRussianFormalName: tr[14],
			UNGEGNLocalShortName:    tr[15],
			UNGEGNLocalFormalName:   tr[16],

			BGNEnglishShortNameGazetteerOrder: tr[17],
			BGNEnglishShortNameReadingOrder:   tr[18],
			BGNEnglishLongName:                tr[19],
			BGNLocalShortName:                 tr[20],
			BGNLocalLongName:                  tr[21],

			PCGNEnglishShortNameGazetteerOrder: tr[22],
			PCGNEnglishShortNameReadingOrder:   tr[23],
			PCGNEnglishLongName:                tr[24],

			FAOItalianLongName: tr[25],
			FFOGermanShortName: tr[26],
		}

		allCountries = append(allCountries, country)
	}

	// Prepare template context
	ctx := struct {
		PkgName   string
		Countries []countries.Country
	}{
		PkgName:   goPkg,
		Countries: allCountries,
	}

	var buf bytes.Buffer
	// Exec template
	if err := t.Execute(&buf, ctx); err != nil {
		log.Fatal(err)
	}

	src, err := format.Source(buf.Bytes())
	if err != nil {
		os.Stderr.Write(buf.Bytes())
		log.Fatalf("fmt: %v", err)
	}

	outf, err := os.Create("countries.gen.go")
	if err != nil {
		log.Fatal(err)
	}
	defer outf.Close()

	if _, err = io.Copy(outf, bytes.NewReader(src)); err != nil {
		log.Fatal(err)
	}
}

func countrySrc(country countries.Country) string {
	src := strings.Replace(fmt.Sprintf("%#v", country), "countries.", "", 1)
	return src
}

const tmpl = `package {{ .PkgName }}

// All countries
var (
{{ range .Countries }}  {{ .ISO3166OneAlphaTwo }} = {{ countrySrc . }}
{{ end }}
)

// Countries defines all countries with they ISO 3166-1 Alpha-2 code as key.
var Countries = map[string]Country{
{{ range .Countries }}"{{ .ISO3166OneAlphaTwo }}": {{ .ISO3166OneAlphaTwo }},
{{ end }}
}
`
