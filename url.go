package wbem

import (
	"bytes"
	"errors"
	"strconv"

	"unicode"
)

const (
	state_init = iota
	state_property_name_begin
	state_property_name
	state_property_value_begin
	state_property_value_typed_begin
	state_property_value_typed_end
	state_property_value_qouted_end
	state_property_value_qouted
	state_property_value_qouted_escaped
	state_property_value_unqouted
)

func is_white(c rune) bool {
	return unicode.IsSpace(c)
}

func is_name_char(c rune) bool {
	if c == '_' || c == '-' ||
		c >= 'a' && c <= 'z' ||
		c >= 'A' && c <= 'Z' ||
		c >= '0' && c <= '9' {
		return true
	}
	return false
}

func Parse(s string) (namespace string, class_name string, keyBindings CimKeyBindings, e error) {

	var buf bytes.Buffer
	namespace_last := 0
	var propertyName string
	var propertyType string

	state := state_init
	for idx, c := range s {
		switch state {
		case state_init:
			if is_name_char(c) {
				buf.WriteRune(c)
				continue
			}
			if '/' == c {
				namespace_last = buf.Len()
				buf.WriteRune(c)
				continue
			}
			if '.' == c {
				if 0 != namespace_last {
					byteArray := buf.Bytes()
					namespace = string(byteArray[:namespace_last])
					class_name = string(byteArray[namespace_last+1:])
				} else {
					class_name = buf.String()
				}
				buf.Reset()
				state = state_property_name_begin
				continue
			}
			e = errors.New("invalid classpath - `" + s + "` at " + strconv.FormatInt(int64(idx), 10))
			return
		case state_property_name_begin:
			if ',' == c {
				e = errors.New("invalid property - `" + s + "` at " + strconv.FormatInt(int64(idx), 10))
				return
			}
			propertyName = ""
			propertyType = ""
			state = state_property_name
			fallthrough
		case state_property_name:
			if is_name_char(c) {
				buf.WriteRune(c)
				continue
			}

			if '=' == c {
				propertyName = buf.String()
				buf.Reset()
				state = state_property_value_begin
				continue
			}
			e = errors.New("invalid property name - `" + s + "` at " + strconv.FormatInt(int64(idx), 10))
			return
		case state_property_value_begin:
			if '"' == c {
				state = state_property_value_qouted
				continue
			}
			if '(' == c {
				state = state_property_value_typed_begin
				continue
			}
			buf.WriteRune(c)
			state = state_property_value_unqouted
		case state_property_value_typed_begin:
			if ')' == c {
				propertyType = buf.String()
				buf.Reset()
				state = state_property_value_typed_end
				continue
			}
			buf.WriteRune(c)
		case state_property_value_typed_end:
			if '"' == c {
				state = state_property_value_qouted
				continue
			}
			if is_name_char(c) {
				buf.WriteRune(c)
				state = state_property_value_unqouted
				continue
			}
			e = errors.New("invalid property value - `" + s + "` at " + strconv.FormatInt(int64(idx), 10))
			return
		case state_property_value_qouted_end:
			if ',' == c {
				state = state_property_name_begin
				continue
			}
			e = errors.New("invalid property value - `" + s + "` at " + strconv.FormatInt(int64(idx), 10))
			return
		case state_property_value_qouted:
			if '"' == c {
				keyBindings = append(keyBindings, CimKeyBinding{
					Name:     propertyName,
					KeyValue: &CimKeyValue{Type: propertyType, Value: buf.String()},
				})
				buf.Reset()
				state = state_property_value_qouted_end
				continue
			}

			if '\'' == c {
				state = state_property_value_qouted_escaped
				continue
			}
			buf.WriteRune(c)
		case state_property_value_qouted_escaped:
			if '"' == c || '\'' == c {
				buf.WriteRune(c)
				state = state_property_value_qouted
				continue
			}
			e = errors.New("invalid property value, invalid escaped - `" + s + "` at " + strconv.FormatInt(int64(idx), 10))
			return
		case state_property_value_unqouted:
			if ',' == c {
				keyBindings = append(keyBindings, CimKeyBinding{
					Name:     propertyName,
					KeyValue: &CimKeyValue{Type: propertyType, Value: buf.String()},
				})
				buf.Reset()
				state = state_property_name_begin
				continue
			}

			if is_name_char(c) {
				buf.WriteRune(c)
				continue
			}
			e = errors.New("invalid property value - `" + s + "` at " + strconv.FormatInt(int64(idx), 10))
			return
		default:
			e = errors.New("unknow state - " + strconv.FormatInt(int64(state), 10))
			return
		}
	}

	switch state {
	case state_init:
		if 0 != namespace_last {
			byteArray := buf.Bytes()
			namespace = string(byteArray[:namespace_last])
			class_name = string(byteArray[namespace_last+1:])
		} else {
			class_name = buf.String()
		}
	//case  state_property_name_begin:
	case state_property_name,
		state_property_value_begin,
		state_property_value_typed_begin,
		state_property_value_typed_end:
		e = errors.New("property value is missing - `" + s + "`")

	//case   state_property_value_qouted_end:
	case state_property_value_qouted, state_property_value_qouted_escaped:
		e = errors.New("qouted is missing - `" + s + "`")
	case state_property_value_unqouted:

		keyBindings = append(keyBindings, CimKeyBinding{
			Name:     propertyName,
			KeyValue: &CimKeyValue{Type: propertyType, Value: buf.String()},
		})
		//buf.Reset()
		//state = state_property_name_begin
	}
	return
}
