package coap

// Instantiates a new core-attribute with a given key/value
func NewCoreAttribute(key string, value interface{}) *CoreAttribute {
	return &CoreAttribute{
		Key:   key,
		Value: value,
	}
}

type CoreAttribute struct {
	Key   string
	Value interface{}
}

// Instantiates a new Core Resource Object
func NewCoreResource() *CoreResource {
	c := &CoreResource{}

	return c
}

type CoreAttributes []*CoreAttribute
type CoreResource struct {
	Target     string
	Attributes CoreAttributes
}

// Adds an attribute (key/value) for a given core resource
func (c *CoreResource) AddAttribute(key string, value interface{}) {
	c.Attributes = append(c.Attributes, NewCoreAttribute(key, value))
}

// Gets an attribute for a core resource
func (c *CoreResource) GetAttribute(key string) *CoreAttribute {
	for _, attr := range c.Attributes {
		if attr.Key == key {
			return attr
		}
	}
	return nil
}
