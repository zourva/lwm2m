package coap

import (
	"github.com/plgd-dev/go-coap/v3/message"
	"regexp"
	"strings"
)

type MediaType = message.MediaType

// CoREAttribute is a key/value pair
// that describes the link or its target.
type CoREAttribute struct {
	Key   string
	Value any
}

type CoREAttributes []*CoREAttribute

// CoREResource contains a target URI, and optional target attributes.
//
//		e.g., </1/0>;ver=2.2;ct=40
//	    Target is </1/0>
//	    Attributes are ver=2.2 and ct=40 kv pairs.
type CoREResource struct {
	Target     string
	Attributes CoREAttributes
}

// AddAttribute adds an attribute (key/value) for a given core resource
func (c *CoREResource) AddAttribute(key string, value any) {
	c.Attributes = append(c.Attributes, &CoREAttribute{Key: key, Value: value})
}

// GetAttribute gets an attribute for a core resource
func (c *CoREResource) GetAttribute(key string) *CoREAttribute {
	for _, attr := range c.Attributes {
		if attr.Key == key {
			return attr
		}
	}
	return nil
}

// ParseCoRELinkString parses string data of a CoRE
// Resources Object into
func ParseCoRELinkString(strCoRELink string) []*CoREResource {
	var re = regexp.MustCompile(`(<[^>]+>\s*(;\s*\w+\s*(=\s*(\w+|"([^"\\]*(\\.[^"\\]*)*)")\s*)?)*)`)
	var elemRe = regexp.MustCompile(`<[^>]*>`)

	var resources []*CoREResource
	m := re.FindAllString(strCoRELink, -1)

	for _, match := range m {
		elemMatch := elemRe.FindString(match)
		target := elemMatch[1 : len(elemMatch)-1]

		resource := &CoREResource{
			Target: target,
		}

		if len(match) > len(elemMatch) {
			attrs := strings.Split(match[len(elemMatch)+1:], ";")
			for _, attr := range attrs {
				pair := strings.Split(attr, "=")
				resource.AddAttribute(pair[0], strings.Replace(pair[1], "\"", "", -1))
			}
		}
		resources = append(resources, resource)
	}
	return resources
}
