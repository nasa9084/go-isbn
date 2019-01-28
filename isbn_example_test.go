package isbn_test

import (
	"fmt"
	"log"

	isbn "github.com/nasa9084/go-isbn"
)

func ExampleISBN() {
	n := "ISBN978-4-10-109205-8"

	code, err := isbn.Parse(n)
	if err != nil {
		log.Fatal(err)
	}
	if !code.IsValid() {
		log.Print("given ISBN code is invalid")
	}
	if code.IsLegacy {
		code, err = code.Update()
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Printf("%s", code.String())
	// Output:
	// ISBN978-4-10-109205-8
}
