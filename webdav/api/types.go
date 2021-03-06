// Package api has type definitions for webdav
package api

import (
	"encoding/xml"
	"regexp"
	"strconv"
	"time"
)

const (
	// Wed, 27 Sep 2017 14:28:34 GMT
	timeFormat = time.RFC1123
)

// Multistatus contains responses returned from an HTTP 207 return code
type Multistatus struct {
	Responses []Response `xml:"response"`
}

// Response contains an Href the response it about and its properties
type Response struct {
	Href  string `xml:"href"`
	Props Prop   `xml:"propstat"`
}

// Prop is the properties of a response
type Prop struct {
	Status   string   `xml:"DAV: status"`
	Name     string   `xml:"DAV: prop>displayname,omitempty"`
	Type     xml.Name `xml:"DAV: prop>resourcetype>collection,omitempty"`
	Size     int64    `xml:"DAV: prop>getcontentlength,omitempty"`
	Modified Time     `xml:"DAV: prop>getlastmodified,omitempty"`
}

// Parse a status of the form "HTTP/1.1 200 OK",
var parseStatus = regexp.MustCompile(`^HTTP/[0-9.]+\s+(\d+)\s+(.*)$`)

// StatusOK examines the Status and returns an OK flag
func (p *Prop) StatusOK() bool {
	match := parseStatus.FindStringSubmatch(p.Status)
	if len(match) < 3 {
		return false
	}
	code, err := strconv.Atoi(match[1])
	if err != nil {
		return false
	}
	if code >= 200 && code < 300 {
		return true
	}
	return false
}

// PropValue is a tagged name and value
type PropValue struct {
	XMLName xml.Name `xml:""`
	Value   string   `xml:",chardata"`
}

// Error is used to desribe webdav errors
//
// <d:error xmlns:d="DAV:" xmlns:s="http://sabredav.org/ns">
//   <s:exception>Sabre\DAV\Exception\NotFound</s:exception>
//   <s:message>File with name Photo could not be located</s:message>
// </d:error>
type Error struct {
	Exception  string `xml:"exception,omitempty"`
	Message    string `xml:"message,omitempty"`
	Status     string
	StatusCode int
}

// Error returns a string for the error and statistifes the error interface
func (e *Error) Error() string {
	if e.Message != "" {
		return e.Message
	}
	if e.Exception != "" {
		return e.Exception
	}
	if e.Status != "" {
		return e.Status
	}
	return "Webdav Error"
}

// Time represents represents date and time information for the
// webdav API marshalling to and from timeFormat
type Time time.Time

// MarshalXML turns a Time into XML
func (t *Time) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	timeString := (*time.Time)(t).Format(timeFormat)
	return e.EncodeElement(timeString, start)
}

// UnmarshalXML turns XML into a Time
func (t *Time) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	err := d.DecodeElement(&v, &start)
	if err != nil {
		return err
	}
	newT, err := time.Parse(timeFormat, v)
	if err != nil {
		return err
	}
	*t = Time(newT)
	return nil
}
