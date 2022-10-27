package graphqlselection

import (
	"fmt"
	"reflect"
	"strings"
)

func ToGraphQLFields(v any) (string, error) {

	t, ok := v.(reflect.Type)

	if ! ok {
		t = reflect.TypeOf(v)

		if t.Kind() == reflect.Pointer {
			t = t.Elem()
		}

		if t.Kind() != reflect.Struct {
			return "", fmt.Errorf("expected struct, got %s", t.Kind())
		}
	}

	var selection []string

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		name := getGraphQLFieldName(field)

		if name == "" || name == "-" {
			continue
		}

		ft := unwrapFieldType(field.Type) 
		switch ft.Kind() {
		case reflect.Struct:
			subSelection, err := ToGraphQLFields(ft)

			if err != nil {
				return "", err
			}

			selection = append(selection, fmt.Sprintf(`%s { %s }`, name, subSelection))
		case reflect.Bool:
			fallthrough
		case reflect.String:
			fallthrough
		case reflect.Int:
			fallthrough
		case reflect.Int16:
			fallthrough
		case reflect.Int32:
			fallthrough
		case reflect.Int64:
			fallthrough
		case reflect.Uint:
			fallthrough
		case reflect.Uint16:
			fallthrough
		case reflect.Uint32:
			fallthrough
		case reflect.Uint64:
			fallthrough
		case reflect.Float32:
			fallthrough
		case reflect.Float64:
			selection = append(selection, name)
		}
	}
	
	if len(selection) == 0 {
		return "", fmt.Errorf("no GraphQL fields in struct")
	}

	return strings.Join(selection, " "), nil
}

func unwrapFieldType(t reflect.Type) reflect.Type {
	switch t.Kind(){
	case reflect.Pointer:
		fallthrough
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		return unwrapFieldType(t.Elem())

	default:
		return t
	}
}

func getGraphQLFieldName(f reflect.StructField) string {
	if n, ok := f.Tag.Lookup("graphql"); ok {
		return n
	}

	if n, ok := f.Tag.Lookup("json"); ok {
		return strings.SplitN(n, ",", 2)[0]
	}

	return f.Name
}