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

	fieldMap := make(map[string]string)
	v := reflect.ValueOf(out)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()

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

		fieldMap[inputKey] = maskTarget
	}

	uniquePaths := make(map[string]bool)
	var paths []string

	for _, key := range inputKeys {
		if target, found := fieldMap[key]; found {
			if !uniquePaths[target] {
				paths = append(paths, target)
				uniquePaths[target] = true
			}
		}
	}

	return paths, nil
}
