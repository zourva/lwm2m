package core

// see OMA-TS-LightweightM2M_Core-V1_2_1-20221209-A Chapter 6 for details.

// BootstrapClient defines methods for
// Client Initiated Bootstrap mode.
type BootstrapClient interface {
	// Request implements BootstrapRequest operation
	//  method: POST
	//  path: /bs?ep={Endpoint Client Name}&pct={Preferred Content Format}
	//  code may be responded:
	//    2.04 Changed. Operation is completed successfully
	//    4.00 Bad Request. Unknown Endpoint Client Name
	//    4.15 Unsupported content format. The specified format is not supported
	Request() error

	// PackRequest implements BootstrapPackRequest operation
	//  method: GET
	//  format: SenML CBOR, SenML JSON, or LwM2M CBOR
	//  path: /bspack?ep={Endpoint Client Name}
	//  code may be responded:
	//    2.05 Content. The response includes the Bootstrap-Pack
	//    4.00 Bad Request. Undetermined error occurred
	//    4.01 Unauthorized
	//    4.04 Not Found
	//    4.05 Method Not Allowed
	//    4.06 Not Acceptable
	//    5.01 Not Implemented
	PackRequest() error

	OnRead() (*ResourceField, error)
	OnWrite() error
	OnDelete() error
	OnDiscover() (*ResourceField, error)
	OnFinish() error

	// BootstrapServerBootstrapInfo returns bootstrap
	// information for the bootstrap server.
	BootstrapServerBootstrapInfo() *BootstrapServerBootstrapInfo
	//SecurityCredentials() *SecurityCredentials
}

type BootstrapServer interface {
	OnRequest(ep, addr string) error
	OnPackRequest(ep string) ([]byte, error)

	//// Read implements Read operation
	////  method: GET
	////  format: TLV, LwM2M CBOR, SenML CBOR or SenML JSON
	////  path: /{Object ID} in LwM2M 1.1 and thereafter, Object ID MUST be '1'
	////     (Server Object) or '2' (Access Control Object)
	//Read(oid ObjectID)
	//
	//// Write implements Write operation
	////  method: PUT
	////  path: /{Object ID}
	////        /{Object ID}/{optional Object Instance ID}
	////        /{Object ID}/{optional Object Instance ID}/{optional Resource ID}
	//Write(oid ObjectID, oiId InstanceID, rid ResourceID, value Value)
	//
	//// Delete implements Delete operation
	////  method: DELETE
	////  path: /{Object ID}/{Object Instance ID}
	//Delete(oid ObjectID, oiId InstanceID)
	//
	//// Discover implements Discover operation
	////  method: GET
	////  path: /{Object ID}
	//Discover(oid ObjectID)
	//
	//// Finish implements Finish operation
	////  method: POST
	////  path: /bs
	//Finish()
}

// BootstrapServerBootstrapInfo is used by the LwM2M Client to contact the
// LwM2M BootstrapServer to get the LwM2M Server Bootstrap Information.
//
//	The LwM2M Client SHOULD have the LwM2M Bootstrap-Server Bootstrap Information
//	The LwM2M Client MUST have the LwM2M Server Bootstrap Information after the bootstrap sequence
//	The LwM2M Client MUST have at most one LwM2M Bootstrap-Server Account
type BootstrapServerBootstrapInfo struct {
	//pre-provisioned Bootstrap-Server Account
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
