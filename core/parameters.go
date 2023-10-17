package core

//see OMA-TS-LightweightM2M_Core-V1_2_1-20221209-A
//6.2.1.2. Behaviour with Current Transport Binding and Modes

type BindingMode = string

const (
	BindingModeUDP   BindingMode = "U"
	BindingModeMQTT  BindingMode = "M"
	BindingModeHTTP  BindingMode = "H"
	BindingModeTCP   BindingMode = "T"
	BindingModeSMS   BindingMode = "S"
	BindingModeNonIP BindingMode = "N"
)

// OpCode defines operations a Resource supports.
// The Execute Operation cannot be used together
// with the Read and Write operations.
//
// OpNone represents empty value, which means that
// a resource can only be accesses by
// the "Bootstrap" interface.
//
// see OMA-TS-LightweightM2M_Core-V1_2_1-20221209-A
// page 106(155) for details.
type OpCode int

const (
	OpNone      OpCode = 0 //empty value, which means that this field
	OpRead      OpCode = 1
	OpWrite     OpCode = 2
	OpReadWrite OpCode = 3
	OpExecute   OpCode = 4
)
