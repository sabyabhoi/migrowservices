package data

import "testing"

func TestValidation(t *testing.T) {
	p := &Product{
		Name: "Latte",
		Price: 100,
		SKU: "aaa-bbb-cccc",
	}

	err := p.Validate()
	if err != nil {
		t.Fatal(err)
	}
}
