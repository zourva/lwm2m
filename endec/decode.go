package endec

import (
	"encoding/binary"
	"errors"
	"github.com/zourva/lwm2m/core"
	"log"
)

// DecodeTypeField extracts/decodes the TLV type field from a byte array
func DecodeTypeField(typeField byte) (idType byte, idLen byte, typeLen byte, valLen byte) {
	idType = typeField & TlvFieldIdentifierType
	idLen = typeField & TlvFieldIdentifierLength
	typeLen = typeField & TlvFieldTypeOfLength
	valLen = typeField & TlvFieldLengthOfValue

	return
}

// DecodeIdentifierField decodes the identifier field and returns the type and length
func DecodeIdentifierField(b []byte, pos int) (identifier core.ResourceID, typeLength int) {
	_, typeFieldLengthOfIdentifier, _, _ := DecodeTypeField(b[0])

	if typeFieldLengthOfIdentifier == 0 {
		_identifier, _ := binary.Uvarint(b[pos : pos+1])
		identifier = core.ResourceID(_identifier)
		typeLength = 1
	} else {
		_identifier, _ := binary.Uvarint(b[pos : pos+2])
		identifier = core.ResourceID(_identifier)
		typeLength = 2
	}
	return
}

// DecodeLengthField decodes the length field and returns the type and value length
func DecodeLengthField(b []byte, pos int) (valueLength uint64, typeLength int) {
	_, _, typeFieldTypeOfLength, typeFieldLengthOfValue := DecodeTypeField(b[0])

	typeLength = 0
	if typeFieldTypeOfLength == 0 {
		valueLength = uint64(typeFieldLengthOfValue)
	} else if typeFieldTypeOfLength == 8 {
		valueLength, _ = binary.Uvarint(b[pos : pos+1])
		typeLength = 1
	} else if typeFieldTypeOfLength == 16 {
		valueLength, _ = binary.Uvarint(b[pos : pos+2])
		typeLength = 2
	} else if typeFieldTypeOfLength == 24 {
		valueLength, _ = binary.Uvarint(b[pos : pos+3])
		typeLength = 3
	} else {
		// Invalid type of Length
	}
	return
}

// DecodeResourceValue decodes the resource value
func DecodeResourceValue(resourceId core.ResourceID, b []byte, resourceDef core.Resource) (core.Value, error) {
	if resourceDef.Multiple() {

		err := ValidResourceTypeField(b)
		if err != nil {
			log.Println(err)
		}

		typeFieldTypeOfIdentifier, _, _, _ := DecodeTypeField(b[0])

		if typeFieldTypeOfIdentifier == TypeFieldTypeMultipleResource {
			valueOffset := 1
			identifier, identifierTypeLength := DecodeIdentifierField(b, valueOffset)
			valueOffset += identifierTypeLength

			_, valueTypeLength := DecodeLengthField(b, valueOffset)
			valueOffset += valueTypeLength
			bytesValue := b[valueOffset:]

			bytesLeft := bytesValue
			var resourceBytes [][]byte
			for len(bytesLeft) > 0 {
				err := ValidResourceTypeField(bytesLeft)
				if err != nil {
					log.Println(err)
				}

				valueOffset := 1
				_, identifierTypeLength := DecodeIdentifierField(bytesLeft, valueOffset)
				valueOffset += identifierTypeLength

				valueFieldLength, valueTypeLength := DecodeLengthField(bytesLeft, valueOffset)
				valueOffset += valueTypeLength

				actualValueLength := uint64(valueOffset) + valueFieldLength

				bytesValue := bytesLeft[:actualValueLength]
				resourceBytes = append(resourceBytes, bytesValue)
				bytesLeft = bytesLeft[actualValueLength:]
			}

			var decodedValues []*core.ResourceValue
			for _, r := range resourceBytes {
				v, _ := DecodeResourceValue(identifier, r, resourceDef)
				decodedValues = append(decodedValues, v.(*core.ResourceValue))
			}

			return core.NewMultipleResourceValue(identifier, decodedValues), nil
		} else {
			valueOffset := 1
			_, identifierTypeLength := DecodeIdentifierField(b, valueOffset)
			valueOffset += identifierTypeLength

			_, valueTypeLength := DecodeLengthField(b, valueOffset)
			valueOffset += valueTypeLength

			bytesValue := b[valueOffset:]
			return core.NewResourceValue(resourceId, ValueFromBytes(bytesValue, resourceDef.Type())), nil
		}
	} else {
		return core.NewResourceValue(resourceId, ValueFromBytes(b, resourceDef.Type())), nil
	}
}

// ValueFromBytes extracts value from a lwm2m byte fragment
func ValueFromBytes(b []byte, v core.ValueType) core.Value {
	if len(b) == 0 {
		return core.Empty()
	}

	switch v {
	case core.ValueTypeString:
		return core.String(string(b))

	case core.ValueTypeInteger:
		return core.BytesToIntegerValue(b)

	case core.ValueTypeTime:
		return core.String("")
	}

	return core.Empty()
}

// ValidResourceTypeField Checks if a type field is of a valid type (resource instance, multiple resource etc)
func ValidResourceTypeField(b []byte) error {
	typeField := b[0]

	typeFieldTypeOfIdentifier, _, _, _ := DecodeTypeField(typeField)

	if typeFieldTypeOfIdentifier != TypeFieldTypeResourceInstance && typeFieldTypeOfIdentifier != TypeFieldTypeMultipleResource && typeFieldTypeOfIdentifier != TypeFieldTypeResource {
		return errors.New("invalid resource identifier")
	}
	return nil
}
