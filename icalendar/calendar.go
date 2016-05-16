package icalendar

import (
	"bufio"
	"fmt"
	"io"
	"reflect"
	"sort"
)

const (
	CLASS_PUBLIC  = "public"
	CLASS_PRIVATE = "private"
)

type Component interface {
	Write(tzid string, w io.Writer) error
}

type VCALENDAR struct {
	// The following are REQUIRED, but MUST NOT occur more than once.
	PRODID  string
	VERSION string
	// The following are OPTIONAL, byt MUST NOT occur more than once.
	CALSCALE string
	METHOD   string
	// The following are OPTIONAL, and MAY occur more than once.
	XPROP    string
	IANAPROP string

	components []Component
	tzid       string
}

func NewCalendar(tzid string) *VCALENDAR {
	return &VCALENDAR{
		CALSCALE: "GREGORIAN",
		VERSION:  "2.0",
		tzid:     tzid,
	}
}

func (cal *VCALENDAR) AddComponent(c Component) {
	cal.components = append(cal.components, c)
}

func (cal *VCALENDAR) Write(w io.Writer) error {
	t := reflect.TypeOf(cal).Elem()
	v := reflect.ValueOf(cal).Elem()
	bw := bufio.NewWriter(w)

	props := map[string]interface{}{}

	for _, k := range []string{
		"PRODID",
		"VERSION",
		"CALSCALE",
		"METHOD",
	} {
		if _, has := t.FieldByName(k); has {
			fv := v.FieldByName(k)
			val := fv.Interface()
			switch fv.Type() {
			case reflect.TypeOf(""):
				if val != "" {
					props[k] = val
				}
			default:
				props[k] = val
			}
		}
	}

	var keys []string
	for k, _ := range props {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Fprintf(bw, "BEGIN:%s\r\n", t.Name())
	for _, k := range keys {
		value := props[k]
		fmt.Fprintf(bw, "%s:%s\r\n", k, value)
	}
	for _, v := range cal.components {
		v.Write(cal.tzid, bw)
	}
	fmt.Fprintf(bw, "END:%s\r\n", t.Name())

	return bw.Flush()
}
