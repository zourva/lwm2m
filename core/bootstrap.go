package core

// see OMA-TS-LightweightM2M_Core-V1_2_1-20221209-A Chapter 6 for details.

// BootstrapClient defines methods for
// Client Initiated Bootstrap mode.
type BootstrapClient interface {
	// BootstrapRequest implements BootstrapRequest operation
	//  method: POST
	//  path: /bs?ep={Endpoint Client Name}&pct={Preferred Content Format}
	BootstrapRequest() error

	// BootstrapPackRequest implements BootstrapPackRequest operation
	//  method: GET
	//  format: SenML CBOR, SenML JSON, or LwM2M CBOR
	//  path: /bspack?ep={Endpoint Client Name}
	BootstrapPackRequest() error

	OnBootstrapRead() (*ResourceField, ErrorType)
	OnBootstrapWrite() ErrorType
	OnBootstrapDelete() ErrorType
	OnBootstrapDiscover() (*ResourceField, ErrorType)
	OnBootstrapFinish() ErrorType

	// BootstrapServerBootstrapInfo returns bootstrap information
	// for the bootstrap server.
	BootstrapServerBootstrapInfo() *BootstrapServerBootstrapInfo
	//SecurityCredentials() *SecurityCredentials
}

type BootstrapServer interface {
	OnBootstrapRequest()
	OnBootstrapPackRequest()

	// BootstrapRead implements BootstrapRead operation
	//  method: GET
	//  format: TLV, LwM2M CBOR, SenML CBOR or SenML JSON
	//  path: /{Object ID} in LwM2M 1.1 and thereafter, Object ID MUST be '1'
	//     (Server Object) or '2' (Access Control Object)
	BootstrapRead()

	// BootstrapWrite implements BootstrapWrite operation
	//  method: PUT
	//  path: /{Object ID}
	//        /{Object ID}/{Object Instance ID}
	//        /{Object ID}/{Object Instance ID}/{Resource ID}
	BootstrapWrite()

	// BootstrapDelete implements BootstrapDelete operation
	//  method: DELETE
	//  path: /{Object ID}/{Object Instance ID}
	BootstrapDelete()

	// BootstrapDiscover implements BootstrapDiscover operation
	//  method: GET
	//  path: /{Object ID}
	BootstrapDiscover()

	// BootstrapFinish implements BootstrapFinish operation
	//  method: POST
	//  path: /bs
	BootstrapFinish()
}

// BootstrapServerBootstrapInfo is used by the LwM2M Client to contact the
// LwM2M BootstrapServer to get the LwM2M Server Bootstrap Information.
//
//	The LwM2M Client SHOULD have the LwM2M Bootstrap-Server Bootstrap Information
//	The LwM2M Client MUST have the LwM2M Server Bootstrap Information after the bootstrap sequence
//	The LwM2M Client MUST have at most one LwM2M Bootstrap-Server Account
type BootstrapServerBootstrapInfo struct {
	BootstrapServerAccount *BootstrapServerAccount
}

// ServerBootstrapInfo
//
//	The LwM2M Server Bootstrap Information MUST contain at least one LwM2M Server Account
//	The LwM2M Client MAY be configured to use one or more LwM2M Server Account(s)
type ServerBootstrapInfo struct {
	ActiveServerAccounts []*ServerAccount //accounts configured to use
}

// BootstrapServerAccount defines LwM2M
// Security Object Instance with
// Bootstrap-Server Resource true.
type BootstrapServerAccount struct {
	SecurityObjectInstance ObjectInstance
}

// ServerAccount defines LwM2M Server
// Object Instance and associated LwM2M
// Security Object Instance with Bootstrap-Server
// Resource false.
type ServerAccount struct {
	SecurityObjectInstance ObjectInstance
	ServerObjectInstance   ObjectInstance
	//potentially OSCORE ObjectInstance
	//COSE ObjectInstance
	//MQTT Server Object Instance
}

//type SecurityCredentials struct {
//}
