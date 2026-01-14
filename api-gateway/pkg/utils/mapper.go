package utils

import (
	"reflect"

	"github.com/gofiber/fiber/v2"
)

func GetFieldMaskPaths(c *fiber.Ctx, out interface{}) ([]string, error) {
	form, err := c.MultipartForm()
	if err != nil {
		return nil, err
	}

	v := reflect.ValueOf(out).Elem()
	t := v.Type()

	var paths []string
	visited := make(map[string]bool)

	for key := range form.Value {
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			formTag := field.Tag.Get("form")

			if formTag == key {
				maskTag := field.Tag.Get("mask")
				if maskTag != "" {
					paths = append(paths, maskTag)
				} else {
					paths = append(paths, formTag)
				}
				visited[maskTag] = true
				break
			}
		}
	}

	return paths, nil
}
