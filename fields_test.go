package graphqlselection_test

import (
	"testing"

	graphqlselection "github.com/rdfirtal/go-graphql-selection"
)

func TestToGraphQLFields(t *testing.T) {
	cases := map[string]struct {
		Type      any
		Expected  string
		ShouldErr bool
	}{
		"flat struct": {
			Type: struct {
				Name string `json:"name"`
			}{},
			Expected: "name",
		},
		"nested struct": {
			Type: struct {
				Person struct {
					Name string `json:"name,omitempty"`
					Age  uint   `json:"age,omitempty"`
				} `json:"person"`
			}{},
			Expected: "person { name age }",
		},
	}

	for n, c := range cases {
		t.Run(n, func(t *testing.T) {
			actual, err := graphqlselection.ToGraphQLFields(c.Type)

			if err != nil {
				if c.ShouldErr {
					return
				}

				t.Fatal(err)
			}

			if actual != c.Expected {
				t.Errorf("actual selection [%s] did not match expected [%s]", actual, c.Expected)
			}
		})
	}
}
