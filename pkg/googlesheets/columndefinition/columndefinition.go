package columndefinition

import (
	"strings"

	"google.golang.org/api/sheets/v4"
)

// ColumnDefinition represents a spreadsheet column definition.
type ColumnDefinition struct {
	Header      string
	ColumnIndex int
	types       map[string]bool
	units       map[string]bool
}

// New creates a new ColumnDefinition.
func New(header string, index int) *ColumnDefinition {
	return &ColumnDefinition{
		Header:      header,
		ColumnIndex: index,
		types:       map[string]bool{},
		units:       map[string]bool{},
	}
}

// CheckCell checks a ColumnDefinition's cell.
func (cd *ColumnDefinition) CheckCell(cell *sheets.CellData) {
	cd.checkType(cell)
	cd.checkUnit(cell)
}

// GetType gets the type of a ColumnDefinition.
func (cd *ColumnDefinition) GetType() string {
	if len(cd.types) == 1 {
		for columnType := range cd.types {
			return columnType
		}
	}

	//The column has mixed or no data types - fallback to string
	return "STRING"
}

// GetUnit gets the unit of a ColumnDefinition.
func (cd *ColumnDefinition) GetUnit() string {
	if len(cd.units) == 1 {
		for unit := range cd.units {
			return unit
		}
	}

	return ""
}

// HasMixedTypes returns whether a ColumnDefinition has mixed types.
func (cd *ColumnDefinition) HasMixedTypes() bool {
	return len(cd.types) > 1
}

// HasMixedUnits returns whether a ColumnDefinition has mixed units.
func (cd *ColumnDefinition) HasMixedUnits() bool {
	return len(cd.units) > 1
}

func (cd *ColumnDefinition) checkType(cell *sheets.CellData) {
	if cell != nil {
		if cell.UserEnteredFormat != nil && cell.UserEnteredFormat.NumberFormat != nil {
			switch cell.UserEnteredFormat.NumberFormat.Type {
			case "DATE", "DATE_TIME":
				cd.types["TIME"] = true
			case "NUMBER", "PERCENT", "CURRENCY":
				cd.types["NUMBER"] = true
			}
		} else if cell.FormattedValue != "" {
			cd.types["STRING"] = true
		}
	}
}

var unitMappings = map[string]string{
	"$":   "currencyUSD",
	"£":   "currencyGBP",
	"€":   "currencyEUR",
	"¥":   "currencyJPY",
	"₽":   "currencyRUB",
	"₴":   "currencyUAH",
	"R$":  "currencyBRL",
	"kr.": "currencyDKK",
	"kr":  "currencySEK",
	"czk": "currencyCZK",
	"CHF": "currencyCHF",
	"PLN": "currencyPLN",
	"฿":   "currencyBTC",
	"R":   "currencyZAR",
	"₹":   "currencyINR",
	"₩":   "currencyKRW",
}

// A lot more that can be done/improved here. For example it should be possible to extract
// the number of decimals from the pattern. Read more here: https://developers.google.com/sheets/api/guides/formats
func (cd *ColumnDefinition) checkUnit(cellData *sheets.CellData) {
	if cellData == nil || cellData.UserEnteredFormat == nil || cellData.UserEnteredFormat.NumberFormat == nil {
		return
	}

	switch cellData.UserEnteredFormat.NumberFormat.Type {
	case "NUMBER":
		for unit, unitID := range unitMappings {
			if strings.Contains(cellData.UserEnteredFormat.NumberFormat.Pattern, unit) {
				cd.units[unitID] = true
			}
		}
	case "PERCENT":
		cd.units["percent"] = true
	case "CURRENCY":
		for unit, unitID := range unitMappings {
			if strings.Contains(cellData.FormattedValue, unit) {
				cd.units[unitID] = true
			}
		}
	}
}
