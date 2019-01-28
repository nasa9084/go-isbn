package isbn_test

import (
	"log"
	"reflect"
	"testing"

	isbn "github.com/nasa9084/go-isbn"
)

func TestParse(t *testing.T) {
	candidates := []struct {
		input    string
		expected isbn.ISBN
		err      error
	}{
		{
			"", isbn.ISBN{}, isbn.ErrEmpty,
		},
		{
			"4-00-310101-4",
			isbn.ISBN{
				Prefix:            "978",
				RegistrationGroup: "4",
				Registrant:        "00",
				Publication:       "310101",
				Checksum:          "4",
				IsLegacy:          true,
			},
			nil,
		},
		{
			"978-4-00-310101-8",
			isbn.ISBN{
				Prefix:            "978",
				RegistrationGroup: "4",
				Registrant:        "00",
				Publication:       "310101",
				Checksum:          "8",
				IsLegacy:          false,
			},
			nil,
		},
		{
			"978-4-00-3101011-8",
			isbn.ISBN{},
			isbn.ErrInvalidLength,
		},
		{
			"971-4-00-310101-8",
			isbn.ISBN{},
			isbn.ErrInvalidPrefix,
		},
	}
	for _, c := range candidates {
		out, err := isbn.Parse(c.input)
		if err != c.err {
			t.Errorf("%v != %v", err, c.err)
			return
		}
		if !reflect.DeepEqual(out, c.expected) {
			t.Errorf("%+v != %+v", out, c.expected)
			return
		}
	}
}

func TestIsValid(t *testing.T) {
	candidates := []struct {
		label    string
		input    isbn.ISBN
		expected bool
	}{
		{
			"legacy-valid",
			isbn.ISBN{
				Prefix:            "978",
				RegistrationGroup: "4",
				Registrant:        "00",
				Publication:       "310101",
				Checksum:          "4",
				IsLegacy:          true,
			},
			true,
		},
		{
			"legacy-valid-X",
			isbn.ISBN{
				Prefix:            "978",
				RegistrationGroup: "4",
				Registrant:        "10",
				Publication:       "109215",
				Checksum:          "X",
				IsLegacy:          true,
			},
			true,
		},
		{
			"legacy-invalid",
			isbn.ISBN{
				Prefix:            "978",
				RegistrationGroup: "4",
				Registrant:        "00",
				Publication:       "310101",
				Checksum:          "3",
				IsLegacy:          true,
			},
			false,
		},
		{
			"valid",
			isbn.ISBN{
				Prefix:            "978",
				RegistrationGroup: "4",
				Registrant:        "00",
				Publication:       "310101",
				Checksum:          "8",
				IsLegacy:          false,
			},
			true,
		},
		{
			"invalid",
			isbn.ISBN{
				Prefix:            "978",
				RegistrationGroup: "4",
				Registrant:        "00",
				Publication:       "310101",
				Checksum:          "7",
				IsLegacy:          false,
			},
			false,
		},
	}
	for _, c := range candidates {
		if b := c.input.IsValid(); b != c.expected {
			t.Errorf("%s: %t != %t", c.label, b, c.expected)
			return
		}
	}
}

func TestUpdate(t *testing.T) {
	candidates := []struct {
		input       isbn.ISBN
		newChecksum string
	}{
		{
			isbn.ISBN{
				Prefix:            "978",
				RegistrationGroup: "4",
				Registrant:        "00",
				Publication:       "310101",
				Checksum:          "4",
				IsLegacy:          true,
			},
			"8",
		},
		{
			isbn.ISBN{
				Prefix:            "978",
				RegistrationGroup: "4",
				Registrant:        "10",
				Publication:       "109205",
				Checksum:          "X",
				IsLegacy:          true,
			},
			"8",
		},
		{
			isbn.ISBN{
				Prefix:            "978",
				RegistrationGroup: "4",
				Registrant:        "00",
				Publication:       "310101",
				Checksum:          "8",
				IsLegacy:          false,
			},
			"8",
		},
	}
	for _, c := range candidates {
		out, err := c.input.Update()
		if err != nil {
			t.Error(err)
			return
		}
		if out.IsLegacy {
			t.Error("out.IsLegacy should be false")
			return
		}
		if out.Checksum != c.newChecksum {
			t.Errorf("%s != %s", out.Checksum, c.newChecksum)
			return
		}
	}
}

func TestString(t *testing.T) {
	candidates := []struct {
		label    string
		input    isbn.ISBN
		expected string
	}{
		{"legacy", isbn.ISBN{Prefix: "978", RegistrationGroup: "4", Registrant: "10", Publication: "109205", Checksum: "2", IsLegacy: true}, "ISBN4-10-109205-2"},
		{"current", isbn.ISBN{Prefix: "978", RegistrationGroup: "4", Registrant: "00", Publication: "310101", Checksum: "8", IsLegacy: false}, "ISBN978-4-00-310101-8"},
	}
	for _, c := range candidates {
		out := c.input.String()
		if out != c.expected {
			log.Printf("%s: %s != %s", c.label, out, c.expected)
			return
		}
	}
}
