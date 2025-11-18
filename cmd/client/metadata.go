package main

import (
	"errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// unmarshalFlags maps command flags into the given struct based on its json tags.
func unmarshalFlags(cmd *cobra.Command, out interface{}) error {
	v := reflect.ValueOf(out)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return errors.New("out must be a non-nil pointer to a struct")
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return errors.New("out must point to a struct")
	}

	// Build: json-tag -> field index PATH (supports embedded structs).
	tagToPath := map[string][]int{}
	var collect func(rt reflect.Type, path []int)
	collect = func(rt reflect.Type, path []int) {
		for i := range rt.NumField() {
			sf := rt.Field(i)

			// Recurse into anonymous embedded structs (including pointers to them).
			if sf.Anonymous {
				ft := sf.Type
				if ft.Kind() == reflect.Ptr {
					ft = ft.Elem()
				}

				if ft.Kind() == reflect.Struct {
					collect(ft, append(path, i))
					continue
				}
			}

			tag := sf.Tag.Get("json")
			if tag == "" || tag == "-" {
				continue
			}
			tag = strings.Split(tag, ",")[0] // strip ",omitempty"
			tagToPath[tag] = append(append([]int(nil), path...), i)
		}
	}
	collect(v.Type(), nil)

	apply := func(fs *pflag.FlagSet) {
		fs.VisitAll(func(flag *pflag.Flag) {
			// Only apply flags the user actually changed.
			if !flag.Changed {
				return
			}

			// Skip global flags you don't want to map.
			switch flag.Name {
			case "debug", "format", "writetofile":
				return
			}

			// Normalize: drop a single "Enabled" suffix if present.
			key := strings.Replace(flag.Name, "Enabled", "", 1)

			idxPath, ok := tagToPath[key]
			if !ok {
				return
			}

			field := v.FieldByIndex(idxPath)
			if !field.CanSet() {
				return
			}

			val := flag.Value.String()
			switch field.Kind() { //nolint:exhaustive
			case reflect.String:
				field.SetString(val)
			case reflect.Bool:
				if b, err := strconv.ParseBool(val); err == nil {
					field.SetBool(b)
				}
			case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
				if i64, err := strconv.ParseInt(val, 10, field.Type().Bits()); err == nil {
					field.SetInt(i64)
				}
			case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
				if u64, err := strconv.ParseUint(val, 10, field.Type().Bits()); err == nil {
					field.SetUint(u64)
				}
			case reflect.Float32, reflect.Float64:
				if f64, err := strconv.ParseFloat(val, field.Type().Bits()); err == nil {
					field.SetFloat(f64)
				}
			default:
				log.Debug().Str("flag", flag.Name).Str("kind", field.Kind().String()).Msg("Unsupported flag type, skipping")
			}
		})
	}

	// Apply persistent (inherited) first, then local.
	apply(cmd.InheritedFlags())
	apply(cmd.LocalFlags())

	return nil
}
