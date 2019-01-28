// Package isbn is a package for parsing, updating, validating International Standard Book Number code.
package isbn

import (
	"strconv"
	"strings"
)

// ISBN represents a parsed ISBN code.
type ISBN struct {
	Prefix            string
	RegistrationGroup string
	Registrant        string
	Publication       string
	Checksum          string
	IsLegacy          bool
}

// Error is package-specific error type
type Error string

func (err Error) Error() string {
	return string(err)
}

const (
	ErrInvalid       Error = "invalid ISBN"
	ErrEmpty         Error = "given string is empty"
	ErrInvalidLength Error = "invalid code length"
	ErrInvalidPrefix Error = "invalid prefix: prefix should be 978 or 979"
)

var null = ISBN{}

// Parse parses raw ISBN code string into ISBN struct.
func Parse(s string) (ISBN, error) {
	if s == "" {
		return null, ErrEmpty
	}
	s = strings.TrimPrefix(s, "ISBN")
	code := strings.Split(s, "-")
	codelen := len(strings.Join(code, ""))
	var prefix string
	var isLegacy bool
	var i int
	switch codelen {
	case 10:
		prefix = "978"
		isLegacy = true
		i = 1
	case 13:
		prefix = code[0]
		if prefix != "979" && prefix != "978" {
			return null, ErrInvalidPrefix
		}
	default:
		return null, ErrInvalidLength
	}
	return ISBN{
		Prefix:            prefix,
		RegistrationGroup: code[1-i],
		Registrant:        code[2-i],
		Publication:       code[3-i],
		Checksum:          code[4-i],
		IsLegacy:          isLegacy,
	}, nil
}

// Update updates a legacy ISBN code into current ISBN code.
// This function returns given ISBN itself when the given is not
// legacy one, and error when
func (isbn ISBN) Update() (ISBN, error) {
	if !isbn.IsLegacy {
		return isbn, nil
	}
	newChecksum := isbn.calcChecksum()
	if newChecksum < 0 {
		return null, ErrInvalid
	}
	isbn.Checksum = strconv.Itoa(newChecksum)
	isbn.IsLegacy = false
	return isbn, nil
}

// IsValid returns the checksum is valid or not.
// This function returns false when the ISBN code itself is invalid.
func (isbn ISBN) IsValid() bool {
	if isbn.IsLegacy {
		return isbn.isValidLegacy()
	}
	return isbn.isValid()
}

func (isbn ISBN) isValidLegacy() bool {
	checkDigit := isbn.calcChecksumLegacy()
	if checkDigit < 0 {
		return false
	}
	if checkDigit != 10 {
		return strconv.Itoa(checkDigit) == isbn.Checksum
	}
	return "X" == isbn.Checksum
}

func (isbn ISBN) calcChecksumLegacy() int {
	s := isbn.RegistrationGroup + isbn.Registrant + isbn.Publication
	i := 10
	checkDigit := 0
	for _, r := range s {
		c, err := strconv.Atoi(string([]rune{r}))
		if err != nil {
			return -1
		}
		checkDigit += c * i
		i--
	}
	checkDigit %= 11
	checkDigit = 11 - checkDigit
	return checkDigit
}

func (isbn ISBN) isValid() bool {
	checkDigit := isbn.calcChecksum()
	if checkDigit < 0 {
		return false
	}
	return strconv.Itoa(checkDigit) == isbn.Checksum
}

func (isbn ISBN) calcChecksum() int {
	s := isbn.Prefix + isbn.RegistrationGroup + isbn.Registrant + isbn.Publication
	b := []int{1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3}
	checkDigit := 0
	for i, r := range s {
		c, err := strconv.Atoi(string([]rune{r}))
		if err != nil {
			return -1
		}
		checkDigit += c * b[i]
	}
	checkDigit %= 10
	checkDigit = 10 - checkDigit
	return checkDigit
}
