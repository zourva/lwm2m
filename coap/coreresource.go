package coap

// CoreAttribute is a key/value pair
// that describes the link or its target.
type CoreAttribute struct {
	Key   string
	Value interface{}
}

// NewCoreAttribute creates a new core-attribute with a given key/value
func NewCoreAttribute(key string, value interface{}) *CoreAttribute {
	return &CoreAttribute{
		Key:   key,
		Value: value,
	}
}

type CoreAttributes []*CoreAttribute

// NewCoreResource creates a new Core Resource Object
func NewCoreResource() *CoreResource {
	c := &CoreResource{}
	return c
}

// CoreResource contains a target URI, and optional target attributes.
//
//		e.g., </1/0>;ver=2.2;ct=40
//	    Target is </1/0>
//	    Attributes are ver=2.2 and ct=40 kv pairs.
type CoreResource struct {
	Target     string
	Attributes CoreAttributes
}

// AddAttribute adds an attribute (key/value) for a given core resource
func (c *CoreResource) AddAttribute(key string, value interface{}) {
	c.Attributes = append(c.Attributes, NewCoreAttribute(key, value))
}

// GetAttribute gets an attribute for a core resource
func (c *CoreResource) GetAttribute(key string) *CoreAttribute {
	for _, attr := range c.Attributes {
		if attr.Key == key {
			return attr
		}
	}
	return nil
}
