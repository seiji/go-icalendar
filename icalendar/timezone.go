package icalendar

import (
	"bufio"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"
	"time"
)

type VTIMEZONE struct {
	// tzid are REQUIRED, but MUST NOT occur more than once.
	TZID string
	// 'last-mod' and 'tzurl' are OPTIONAL, and MAY occur more than once.
	LASTMODIFIED string
	TZURL        string
	// One of 'standardc' or 'daylightc' MUST occur and each MAY occur more than once.
	STANDARD TZPROP
	DAYLIGHT TZPROP
	// The following are OPTIONAL, and MAY occur more than once.
	XPROP    string
	IANAPROP string
}

type TZPROP struct {
	// The following are REQUIRED, but MUST NOT occur more than once.
	DTSTART      time.Time
	TZOFFSETTO   string
	TZOFFSETFROM string
	// The following is OPTIONAL, but SHOULD NOT occur more than once.
	RRULE string
	// The following are OPTIONAL, and MAY occur more than once.
	COMMENT  string
	RDATE    time.Time
	TZNAME   string
	XPROP    string
	IANAPROP string
}

func NewTimezone(tzid string) *VTIMEZONE {
	return &VTIMEZONE{
		TZID: tzid,
	}
}

func (c *VTIMEZONE) Write(w io.Writer) error {
	t := reflect.TypeOf(c).Elem()
	v := reflect.ValueOf(c).Elem()
	bw := bufio.NewWriter(w)

	props := map[string]interface{}{}

	for _, k := range []string{
		"TZID",
	} {
		if _, has := t.FieldByName(k); has {
			fv := v.FieldByName(k)
			val := fv.Interface()
			switch fv.Type() {
			case reflect.TypeOf(time.Time{}):
				t := val.(time.Time)
				if k == "DTEND" && t.Unix() <= 0 {
					start, _ := time.Parse("20060102T150405Z", props["DTSTART"].(string))
					delete(props, "DTSTART")
					props["DTSTART;VALUE=DATE"] = start.UTC().Format("20060102")
				} else {
					props[k] = t.UTC().Format("20060102T150405Z")
				}
			case reflect.TypeOf(""):
				if val != "" {
					props[k] = val
				}
			case reflect.TypeOf([]string{}):
				if s := val.([]string); len(s) > 0 {
					props[k] = strings.Join(s, ",")
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
		fmt.Fprintf(bw, "%s:%s\r\n", k, props[k])
	}
	fmt.Fprintf(bw, "END:%s\r\n", t.Name())
	return nil
}
