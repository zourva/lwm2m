package coap

import (
	"math"
	"strings"
)

type Option interface {
	Name() string
	IsElective() bool
	IsCritical() bool
	StringValue() string
	IntValue() int
	GetCode() OptionCode
	GetValue() interface{}
}

// Represents an Option for a CoAP Message
type CoapOption struct {
	Code  OptionCode
	Value interface{}
}

func (o *CoapOption) GetValue() interface{} {
	return o.Value
}

func (o *CoapOption) GetCode() OptionCode {
	return o.Code
}

func (o *CoapOption) Name() string {
	return "Name of option"
}

// Determines if an option is elective
func (o *CoapOption) IsElective() bool {
	if (int(o.Code) % 2) != 0 {
		return false
	}
	return true
}

// Determines if an option is critical
func (o *CoapOption) IsCritical() bool {
	if (int(o.Code) % 2) != 0 {
		return true
	}
	return false
}

// Returns the string value of an option
func (o *CoapOption) StringValue() string {
	return o.Value.(string)
}

func (o *CoapOption) IntValue() int {
	return int(o.Value.(uint32))
}

// Instantiates a New Option
func NewOption(optionNumber OptionCode, optionValue interface{}) *CoapOption {
	return &CoapOption{
		Code:  optionNumber,
		Value: optionValue,
	}
}

// Creates an array of options decomposed from a given path
func NewPathOptions(path string) []Option {
	opts := []Option{}
	ps := strings.Split(path, "/")

	for _, p := range ps {
		if p != "" {
			opt := NewOption(OptionURIPath, p)
			opts = append(opts, opt)
		}
	}
	return opts
}

// Checks if an option is repeatable
func IsRepeatableOption(opt Option) bool {

	switch opt.GetCode() {

	case OptionIfMatch, OptionEtag, OptionURIPort, OptionLocationPath, OptionURIPath, OptionURIQuery, OptionLocationQuery,
		OptionBlock2, OptionBlock1:
		return true

	default:
		return false
	}
}

// Checks if an option/option code is recognizable/valid
func IsValidOption(opt Option) bool {
	switch opt.GetCode() {

	case OptionIfNoneMatch, OptionURIHost,
		OptionEtag, OptionIfMatch, OptionObserve, OptionURIPort, OptionLocationPath,
		OptionURIPath, OptionContentFormat, OptionMaxAge, OptionURIQuery, OptionAccept,
		OptionLocationQuery, OptionBlock2, OptionBlock1, OptionProxyURI, OptionProxyScheme, OptionSize1:
		return true

	default:
		return false
	}
}

// Determines if an option is elective
func IsElectiveOption(opt Option) bool {
	i := int(opt.GetCode())

	if (i & 1) == 1 {
		return false
	}
	return true
}

// Determines if an option is critical
func IsCriticalOption(opt Option) bool {
	return !IsElectiveOption(opt)
}

func NewBlock1Option(bs BlockSizeType, more bool, seq uint32) *Block1Option {
	opt := &Block1Option{}
	opt.Code = OptionBlock1

	val := seq

	val = val << 4
	if more {
		val |= (1 << 3)
	}

	val |= (uint32(bs) << 0)

	opt.Value = val

	/*
		BLockSize
		val := o.Value.(uint)
		exp := val & 0x07

		return math.Exp2(float64(exp + 4))

		More
		val := o.Value.(uint)

		return ((val >> 3) & 0x01) == 1

	*/

	return opt
}

func Block1OptionFromOption(opt Option) *Block1Option {
	blockOpt := &Block1Option{}

	blockOpt.Value = opt.GetValue()
	blockOpt.Code = opt.GetCode()

	return blockOpt
}

type Block1Option struct {
	CoapOption
}

func (o *Block1Option) Sequence() uint32 {
	val := o.GetValue().(uint32)

	return val >> 4
}

func (o *Block1Option) Exponent() uint32 {
	val := o.GetValue().(uint32)

	return val & 0x07
}

func (o *Block1Option) BlockSizeLength() uint32 {
	sz := uint32(o.Size()) + 4

	return sz * sz
}

func (o *Block1Option) Size() BlockSizeType {
	val := o.GetValue().(uint32)
	exp := val & 0x07

	return BlockSizeType(byte(math.Exp2(float64(exp + 4))))
}

func (o *Block1Option) HasMore() bool {
	val := o.Value.(uint32)

	return ((val >> 3) & 0x01) == 1
}
