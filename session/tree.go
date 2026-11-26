package session

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

func Tree(node any) string {
	var b strings.Builder
	walk(&b, reflect.ValueOf(node), 0)
	return b.String()
}

func walk(b *strings.Builder, v reflect.Value, depth int) {
	// Follow pointers and interfaces down to the value they point at
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		t := v.Type()
		var attrs []string
		var children []reflect.Value
		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			if !field.IsExported() || field.Name == "Token" {
				continue
			}
			fv := v.Field(i)
			switch fv.Kind() {
			case reflect.String, reflect.Bool, reflect.Int, reflect.Int64, reflect.Float64:
				attrs = append(attrs, fmt.Sprintf("%s=%v", field.Name, fv.Interface()))
			default:
				children = append(children, fv)
			}
		}
		line := t.Name()
		if len(attrs) > 0 {
			line += " " + strings.Join(attrs, " ")
		}
		fmt.Fprintf(b, "%s%s\n", strings.Repeat("  ", depth), line)
		for _, child := range children {
			walk(b, child, depth+1)
		}

	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			walk(b, v.Index(i), depth)
		}

	case reflect.Map:
		// Maps have no order, so sort the keys to get the same tree every time
		keys := v.MapKeys()
		sort.Slice(keys, func(i, j int) bool {
			return fmt.Sprint(keys[i].Interface()) < fmt.Sprint(keys[j].Interface())
		})
		for _, key := range keys {
			walk(b, key, depth)
			walk(b, v.MapIndex(key), depth)
		}
	}
}