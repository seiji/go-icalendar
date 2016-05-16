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

type VEVENT struct {
	// The following are REQUIRED, but MUST NOT occur more than once.
	DTSTAMP time.Time
	UID     string
	// The following is REQUIRED if the component appears in an iCalendar object that doesn't
	//  specify the "METHOD" property; otherwise, it is OPTIONAL; in any case, it MUST NOT occur more than once.
	DTSTART time.Time
	// The following are OPTIONAL, and MAY occur more than once.
	CLASS       string // "PUBLIC" / "PRIVATE" / "CONFIDENTIAL"
	CREATED     string
	DESCRIPTION string
	GEO         string
	LASTMOD     string
	LOCATION    string
	ORGANIZER   string
	PRIORITY    string
	SEQ         string
	STATUS      string
	SUMMARY     string
	TRASP       string
	URL         string
	RECURID     string
	// The following are OPTIONAL, but MUST NOT occur more than once.
	RRULE string
	// Either 'dtend' of 'duration' MAY appear in a 'eventprop',
	// but 'dtend' and 'duration; MUST NOT occur in save 'eventprop',
	DTEND time.Time
	// *  DURATION time.Duration
	// The following are OPTIONAL, and MAY occur more than once.
	ATTACH     string
	ATTENDEE   string
	CATEGORIES []string
	COMMENT    string
	CONTACT    string
	EXDATE     string
	RSTATUS    string
	RELATED    string
	RESOURCES  string
	RDATE      string
	XPROP      string
	IANAPROP   string
}

func NewEvent(title string, start, end time.Time) *VEVENT {
	return &VEVENT{
		CLASS:   CLASS_PUBLIC,
		DTSTAMP: time.Now(),
		DTSTART: start,
		DTEND:   end,
		SUMMARY: title,
	}
}

func (c *VEVENT) Write(tzid string, w io.Writer) error {
	t := reflect.TypeOf(c).Elem()
	v := reflect.ValueOf(c).Elem()
	bw := bufio.NewWriter(w)

	props := map[string]interface{}{}

	for _, k := range []string{
		"UID",
		"DTSTAMP",
		"DTSTART",
		"DTEND",
		"SUMMARY",

		"CLASS",
		"CATEGORIES",
		"DESCRIPTION",
		"LOCATION",
		"GEO",
	} {
		if _, has := t.FieldByName(k); has {
			fv := v.FieldByName(k)
			val := fv.Interface()
			switch fv.Type() {
			case reflect.TypeOf(time.Time{}):
				t := val.(time.Time)
				if k == "DTEND" && t.Unix() <= 0 {
					start, _ := time.Parse("20060102T150405", props["DTSTART"].(string))
					delete(props, "DTSTART")
					props["DTSTART;VALUE=DATE"] = start.Format("20060102")
				} else {
					props[k] = t.Format("20060102T150405")
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
		value := props[k]
		if k == "DTSTART" || k == "DTEND" {
			k = fmt.Sprintf("%s;TZID=%s", k, tzid)
		}
		fmt.Fprintf(bw, "%s:%s\r\n", k, value)
	}
	fmt.Fprintf(bw, "END:%s\r\n", t.Name())
	return nil
}

