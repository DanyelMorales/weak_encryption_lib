package symbol

import (
	"fmt"
)

// Symbol represents a single element in an alphabet
type Symbol int64

//StrToSymbolArray create a new Symbol slice from a string
func StrToSymbolArray(str string) []Symbol{
	container := make([]Symbol,0)
	for _,char := range str {
		container = append(container, Symbol(char))
	}
	return container
}

// Equals  compares two symbols
func (c Symbol) Equals(comparable Symbol) bool {
	return c == comparable
}

// IsSpace evals if a symbol is a string
func (c Symbol) IsSpace() bool {
	return c == c.Space()
}

// AssocValue returns the associated value of a char
func (c Symbol) AssocValue(alphabet []Symbol) (Symbol, error) {
	for i := 1; i <= len(alphabet); i++ {
		if alphabet[i-1] == c {
			return Symbol(i), nil
		}
	}
	return ' ', fmt.Errorf("value not found %d", c)
}

// OriginalValue returns the original value of a symbol
func (c Symbol) OriginalValue(alphabet []Symbol) (Symbol, error) {
	index := int(c)
	if index > len(alphabet) || index < 0 {
		return ' ', fmt.Errorf("value not found %d", c)
	}
	return alphabet[index-1], nil
}

// Space returns the space symbol
func (c Symbol) Space() Symbol {
	return ' '
}
