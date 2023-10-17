package server

import (
	"encoding/hex"
	"github.com/pborman/uuid"
)

// GuidProvider provides global uuid
// generation service.
type GuidProvider interface {
	GetGuid() string

	GetGuidWithHint(hint string) string
}

type UrnUuidProvider struct {
}

// GetGuid returns an id of format:
//
//	urn:uuid:########-####-####-####-############
func (p *UrnUuidProvider) GetGuid() string {
	return uuid.NewUUID().URN()
}

func (p *UrnUuidProvider) GetGuidWithHint(hint string) string {
	src := []byte(hint)
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
	return string(dst)
}

func NewUrnUuidProvider() GuidProvider {
	return &UrnUuidProvider{}
}
