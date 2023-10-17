package endec

import (
	"bytes"
	"github.com/zourva/lwm2m/core"
)

// EncodeValue encodes the resource id and value and returns a byte array representation
func EncodeValue(resourceId core.ResourceID, allowMultipleValues bool, v core.Value) []byte {
	if v.Type() == core.ValueTypeMultiple {
		typeOfMultipleValue := v.ContainedType()
		if typeOfMultipleValue == core.ValueTypeInteger {

			// Resource Instances TLV
			resourceInstanceBytes := bytes.NewBuffer([]byte{})
			intValues := v.Get().([]core.Value)
			for i, intValue := range intValues {
				value := intValue.Get().(int)

				// Type Field Byte
				if allowMultipleValues {
					typeField := CreateTlvTypeField(TypeFieldTypeResourceInstance, value, uint16(i))
					resourceInstanceBytes.Write([]byte{typeField})
				} else {
					typeField := CreateTlvTypeField(TypeFieldTypeResource, value, uint16(i))
					resourceInstanceBytes.Write([]byte{typeField})
				}

				// Identifier Field
				identifierField := CreateTlvIdentifierField(uint16(i))
				resourceInstanceBytes.Write(identifierField)

				// Length Field
				lengthField := CreateTlvLengthField(value)
				resourceInstanceBytes.Write(lengthField)

				// Value Field
				valueField := CreateTlvValueField(value)
				resourceInstanceBytes.Write(valueField)
			}

			// Resource Root TLV
			resourceTlv := bytes.NewBuffer([]byte{})

			// Byte 7-6: identifier
			typeField := CreateTlvTypeField(128, resourceInstanceBytes.Bytes(), resourceId)
			resourceTlv.Write([]byte{typeField})

			// Identifier Field
			identifierField := CreateTlvIdentifierField(resourceId)
			resourceTlv.Write(identifierField)

			// Length Field
			lengthField := CreateTlvLengthField(resourceInstanceBytes.Bytes())
			resourceTlv.Write(lengthField)

			// Value Field, Append Resource Instances TLV to Resource TLV
			resourceTlv.Write(resourceInstanceBytes.Bytes())

			return resourceTlv.Bytes()
		}
	} else {
		return v.ToBytes()
	}
	return nil
}
