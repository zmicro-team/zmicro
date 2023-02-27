package form

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/spf13/cast"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var reg = regexp.MustCompile(`{[\\.\w]+}`)

// EncodeURL encode msg to url path.
// pathTemplate is a template of url path like http://helloworld.dev/{name}/sub/{sub.name},
func (c *Codec) EncodeURL(pathTemplate string, v any, needQuery bool) string {
	var repl func(in string) string

	if v == nil || (reflect.ValueOf(v).Kind() == reflect.Ptr && reflect.ValueOf(v).IsNil()) {
		return pathTemplate
	}

	pathParams := make(map[string]struct{})
	if mg, ok := v.(proto.Message); ok {
		repl = func(in string) string {
			// in: {xxx}
			if len(in) < 4 { //nolint:gomnd
				return in
			}
			key := in[1 : len(in)-1]
			vars := strings.Split(key, ".")
			if value, err := getValueFromProtoWithField(mg.ProtoReflect(), vars); err == nil {
				pathParams[key] = struct{}{}
				return value
			}
			return in
		}
	} else {
		repl = func(in string) string {
			// in: {xxx}
			if len(in) < 4 { //nolint:gomnd
				return in
			}
			key := in[1 : len(in)-1]
			fmt.Println(key)
			vars := strings.Split(key, ".")
			if value, err := getValueWithField(v, vars, c.TagName); err == nil {
				pathParams[key] = struct{}{}
				return value
			}
			return in
		}
	}
	path := reg.ReplaceAllStringFunc(pathTemplate, repl)

	if needQuery {
		values, err := c.Encode(v)
		if err == nil && len(values) > 0 {
			for key := range pathParams {
				delete(values, key)
			}
			query := values.Encode()
			if query != "" {
				path += "?" + query
			}
		}
	}
	return path
}

func getValueFromProtoWithField(v protoreflect.Message, fieldPath []string) (string, error) {
	var fd protoreflect.FieldDescriptor

	for i, fieldName := range fieldPath {
		fields := v.Descriptor().Fields()
		if fd = fields.ByJSONName(fieldName); fd == nil {
			fd = fields.ByName(protoreflect.Name(fieldName))
			if fd == nil {
				return "", fmt.Errorf("form: field path not found: %q", fieldName)
			}
		}
		if i == len(fieldPath)-1 {
			break
		}
		if fd.Message() == nil || fd.Cardinality() == protoreflect.Repeated {
			return "", fmt.Errorf("form: invalid path, %q is not a message", fieldName)
		}
		v = v.Get(fd).Message()
	}
	return EncodeField(fd, v.Get(fd))
}

func getValueWithField(s any, fieldPath []string, tagName string) (string, error) {
	v := reflect.ValueOf(s)
	// if pointer get the underlying element
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return "", errors.New("form: not struct")
	}
	for i, fieldName := range fieldPath {
		fields := findField(v, fieldName, tagName)
		if !fields.IsValid() {
			return "", fmt.Errorf("form: field path not found: %q", fieldName)
		}
		v = fields
		if i == len(fieldPath)-1 {
			break
		}
	}
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	return cast.ToString(v.Interface()), nil
}

func findField(v reflect.Value, searchName, tagName string) reflect.Value {
	if v.Kind() == reflect.Ptr && v.IsNil() {
		v = reflect.New(v.Type().Elem())
	}
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return reflect.Value{}
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		fv := v.Field(i)
		// we can't access the value of unexported fields
		if !fv.CanInterface() || field.PkgPath != "" {
			continue
		}
		// don't check if it's omitted
		tag := field.Tag.Get(tagName)
		if tag == "-" {
			continue
		}
		name := field.Name
		tagNamed, _ := parseTag(tag)
		if tagNamed != "" {
			name = tagNamed
		}
		if name == searchName {
			return v.FieldByName(field.Name)
		}
	}
	return reflect.Value{}
}
