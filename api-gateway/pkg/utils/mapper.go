package utils

import (
	"reflect"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func GetFieldMaskPaths(c *fiber.Ctx, out interface{}) ([]string, error) {
	var inputKeys []string

	contentType := c.Get("Content-Type")
	if strings.HasPrefix(contentType, "multipart/form-data") {
		form, err := c.MultipartForm()
		if err != nil {
			return nil, err
		}

		for k := range form.Value {
			inputKeys = append(inputKeys, k)
		}
		for k := range form.File {
			inputKeys = append(inputKeys, k)
		}
	}

	if len(inputKeys) == 0 {
		return []string{}, nil
	}

	v := reflect.ValueOf(out)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	fieldMap := make(map[string]string)
	buildFieldMap(v.Type(), "", fieldMap)

	uniquePaths := make(map[string]bool)
	var paths []string

	for _, key := range inputKeys {
		if target, ok := fieldMap[key]; ok {
			if !uniquePaths[target] {
				paths = append(paths, target)
				uniquePaths[target] = true
			}
		}
	}

	return paths, nil
}

func buildFieldMap(t reflect.Type, prefix string, fieldMap map[string]string) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		inputKey := field.Tag.Get("json")
		if inputKey == "" {
			inputKey = field.Tag.Get("form")
		}

		inputKey = strings.Split(inputKey, ",")[0]
		if inputKey == "" || inputKey == "-" {
			continue
		}

		maskTarget := field.Tag.Get("mask")
		if maskTarget == "" {
			maskTarget = inputKey
		}

		fullInputKey := inputKey
		fullMask := maskTarget

		if prefix != "" {
			fullInputKey = prefix + "." + inputKey
			fullMask = prefix + "." + maskTarget
		}

		if field.Type.Kind() == reflect.Struct &&
			field.Type.String() != "time.Time" {

			buildFieldMap(field.Type, fullInputKey, fieldMap)
			continue
		}

		fieldMap[fullInputKey] = fullMask
	}
}
