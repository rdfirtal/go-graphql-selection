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
		"struct pointer": {
			Type: &struct {
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
		"nested struct slice and no tags": {
			Type: struct {
				People []struct {
					Name string `json:"name,omitempty"`
					Age  uint   `json:"age,omitempty"`
					Pet struct{
						Name string
					} `graphql:"pet"`
				} `json:"people"`
			}{},
			Expected: "people { name age pet { Name } }",
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
