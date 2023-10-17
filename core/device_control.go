package core

// DeviceControlProxy defines operations on registered client
// on server side using proxy/delegation pattern.
// It the Device Management and Service Enablement Interface,
// which is used by the LwM2M Server to access object instances and
// resources available from a registered client
type DeviceControlProxy interface {
	Create(oid ObjectID, newValue Value) error
	Read(oid ObjectID, oiId InstanceID, rid ResourceID, riId InstanceID) error
	Write(oid ObjectID, oiId InstanceID, rid ResourceID, riId InstanceID, newValue Value) error
	Delete(oid ObjectID, oiId InstanceID, rid ResourceID, riId InstanceID) error
	Execute(oid ObjectID, oiId InstanceID, rid ResourceID, args string) error
	Discover(oid ObjectID, oiId InstanceID, rid ResourceID, depth int) error
	//ReadComposite()
	//WriteComposite()
	//WriteAttributes()
}

// DeviceControlClient defines client side operations
// of the Device Management and Service Enablement Interface.
type DeviceControlClient interface {
	// OnCreate implements Create operation
	//  method: POST
	//  format: LwM2M CBOR, SenML CBOR, SenML JSON, or TLV
	//  path: /{Object ID}
	OnCreate(oid ObjectID, newValue Value) ErrorType

	// OnRead implements Read operation
	//  method: GET
	//  path: /{Object ID}
	//        /{Object ID}/{Object Instance ID}
	//        /{Object ID}/{Object Instance ID}/{Resource ID}
	//        /{Object ID}/{Object Instance ID}/{Resource ID}/{Resource Instance ID}
	OnRead(oid ObjectID, oiId InstanceID, rid ResourceID, riId InstanceID) (*ResourceValue, ErrorType)

	// OnWrite implements Write operation
	//  method: POST/PUT
	//  format: LwM2M CBOR, SenML CBOR, SenML JSON, or TLV
	//  path: /{Object ID}/{Object Instance ID}
	//        /{Object ID}/{Object Instance ID}/{Resource ID}
	//        /{Object ID}/{Object Instance ID}/{Resource ID}/{Resource Instance ID}
	OnWrite(oid ObjectID, oiId InstanceID, rid ResourceID, riId InstanceID, newValue Value) ErrorType

	// OnDelete implements Delete operation
	//  method: DELETE
	//  path: /{Object ID}/{Object Instance ID}
	//        /{Object ID}/{Object Instance ID}/{Resource ID}/{Resource Instance ID}
	OnDelete(oid ObjectID, oiId InstanceID, rid ResourceID, riId InstanceID) ErrorType

	// OnExecute implements Write operation
	//  method: POST/PUT
	//  format: none, or text/plain
	//  path: /{Object ID}/{Object Instance ID}/{Resource ID}
	OnExecute(oid ObjectID, oiId InstanceID, rid ResourceID, args string) ErrorType

	// OnDiscover implements Discover operation
	//  method: GET
	//  path: /{Object ID}<Depth>
	//        /{Object ID}/{Object Instance ID}<Depth>
	//        /{Object ID}/{Object Instance ID}/{Resource ID}<Depth>
	//        Depth: ?depth={Value}
	OnDiscover(oid ObjectID, oiId InstanceID, rid ResourceID, depth int) ErrorType
	//OnReadComposite()
	//OnWriteComposite()
	//OnWriteAttributes()
}
