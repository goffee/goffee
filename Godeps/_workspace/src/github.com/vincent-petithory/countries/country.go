// Package countries defines country names and codes according to ISO 3166.
// The source is http://opengeocode.org/download/countrynames.txt.
//
// From countrynames.txt header:
//
//      Country Codes to Country Names
//      Created by OpenGeoCode.Org, Submitted into the Public Domain Jan 26, 2014 (version 7)
//
//      Abbreviations:
//      ISO : International Standards Organization
//      BGN : U.S. Board on Geographic Names
//      UNGEGN : United Nations Group of Experts on Geographic Names
//      PCGN : UK Permanent Committee on Geographic Names
//      FAO  : United Nations Food & Agriculture Organization
//      FFO  : German Federal Foreign Office
//
//      Metadata (one entry per line)
//      ISO 3166-1 alpha-2, ISO 3166-1 alpha-3; ISO 3166-1 numeric;
//      ISO 3166-1 English short name (Gazetteer order); ISO 3166-1 English short name (proper reading order); ISO 3166-1 English romanized short name (Gazetteer order); ISO 3166-1 English romanized short name (proper reading oorder);
//      ISO 3166-1 French short name (Gazetteer order); ISO 3166-1 French short name (proper reading order); ISO 3166-1 Spanish short name (Gazetteer order);
//      UNGEGN English formal name; UNGEGN French formal name; UNGEGN Spanish formal name; UNGEGN Russian short and formal name; UNGEGN local short name; UNGEGN local formal name;
//      BGN English short name (Gazetteer order); BGN English short name (proper reading order); BGN English long name; BGN local short name; BGN local long name
//      PCGN English short name (Gazetteer order); PCGN English short name (proper reading order); PCGN English long name; FAO Italian long name; FFO German short name
//
//      NOTES:
//      UNGEGN and BGN local names: when there is more than one language local name, each local name is followed by the 639-1 alpha-2 language code within paranthesis (xx) and separated by a slash (/).
//      Ex. Canada(en)/le Canada(fr)
package countries

//go:generate go run parser.go countrynames.txt

// Country holds fields for a country as defined by ISO 3166.
type Country struct {
	ISO3166OneAlphaTwo   string
	ISO3166OneAlphaThree string
	ISO3166OneNumeric    int

	ISO3166OneEnglishShortNameGazetteerOrder          string
	ISO3166OneEnglishShortNameReadingOrder            string
	ISO3166OneEnglishRomanizedShortNameGazetteerOrder string
	ISO3166OneEnglishRomanizedShortNameReadingOrder   string
	ISO3166OneFrenchShortNameGazetteerOrder           string
	ISO3166OneFrenchShortNameReadingOrder             string
	ISO3166OneSpanishShortNameGazetteerOrder          string

	UNGEGNEnglishFormalName string
	UNGEGNFrenchFormalName  string
	UNGEGNSpanishFormalName string
	UNGEGNRussianShortName  string
	UNGEGNRussianFormalName string
	UNGEGNLocalShortName    string
	UNGEGNLocalFormalName   string

	BGNEnglishShortNameGazetteerOrder string
	BGNEnglishShortNameReadingOrder   string
	BGNEnglishLongName                string
	BGNLocalShortName                 string
	BGNLocalLongName                  string

	PCGNEnglishShortNameGazetteerOrder string
	PCGNEnglishShortNameReadingOrder   string
	PCGNEnglishLongName                string

	FAOItalianLongName string
	FFOGermanShortName string
}
