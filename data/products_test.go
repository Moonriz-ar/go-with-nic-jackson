package data

import "testing"

func TestChecksValidation(t *testing.T) {
	p := &Product{
		Name:  "bubble tea",
		Price: 10.00,
		SKU:   "abf-asd-fdd",
	}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
