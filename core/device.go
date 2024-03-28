package core

import "github.com/zourva/lwm2m/coap"

// DeviceControlServer defines operations on registered client
// on server side using proxy/delegation pattern.
// It the Device Management and Service Enablement Interface,
// which is used by the LwM2M Server to access object instances and
// resources available from a registered client
type DeviceControlServer interface {
	Create(oid ObjectID, newValue Value) error
	Read(oid ObjectID, oiId InstanceID, rid ResourceID, riId InstanceID) ([]byte, error)
	Write(oid ObjectID, oiId InstanceID, rid ResourceID, riId InstanceID, newValue Value) ([]byte, error)
	Delete(oid ObjectID, oiId InstanceID, rid ResourceID, riId InstanceID) error
	Execute(oid ObjectID, oiId InstanceID, rid ResourceID, args string) error
	Discover(oid ObjectID, oiId InstanceID, rid ResourceID, depth int) ([]*coap.CoREResource, error)
	//ReadComposite()
	//WriteComposite()
	//WriteAttributes()

	//	//Create(client RegisteredClient, oid ObjectID, newValue Value) error
	//	//Read(client RegisteredClient, oid ObjectID, instId InstanceID, resId ResourceID, resInstId InstanceID) error
	//	//Write(client RegisteredClient, oid ObjectID, instId InstanceID, resId ResourceID, resInstId InstanceID, newValue Value) error
	//	//Delete(client RegisteredClient, oid ObjectID, instId InstanceID, resId ResourceID, resInstId InstanceID) error
	//	//Execute(client RegisteredClient, oid ObjectID, instId InstanceID, resId ResourceID, args string) error
	//	//Discover(client RegisteredClient, oid ObjectID, instId InstanceID, resId ResourceID, depth int) error
}

// DeviceControlClient defines client side operations
// of the Device Management and Service Enablement Interface.
type DeviceControlClient interface {
	// OnPack implements Bootstrap-Pack operation
	// method: POST
	// path  : non
	// see: 6.1.7.7. Bootstrap-Pack-Request Operation
	OnPack(newValue []byte) error

	// OnCreate implements Create operation
	//  method: POST
	//  format: LwM2M CBOR, SenML CBOR, SenML JSON, or TLV
	//  path: /{Object ID}
	//  code may be responded:
	//    2.01 Created "Create" operation is completed successfully
	//    4.00 Bad Request Target (i.e., Object) already exists or
	//              Mandatory Resources are not specified or
	//              Content Format is not specified
	//    4.01 Unauthorized Access Right Permission Denied
	//    4.04 Not Found URI of "Create" operation is not found
	//    4.05 Method Not Allowed Target is not allowed for "Create" operation
	//    4.06 Not Acceptable The specified Content-Format is not supported
	//    4.15 Unsupported content format The specified format is not supported
	OnCreate(oid ObjectID, newValue []byte) error

	// OnRead implements Read operation
	//  method: GET
	//  path: /{Object ID}
	//        /{Object ID}/{Object Instance ID}
	//        /{Object ID}/{Object Instance ID}/{Resource ID}
	//        /{Object ID}/{Object Instance ID}/{Resource ID}/{Resource Instance ID}
	OnRead(oid ObjectID, oiId InstanceID, rid ResourceID, riId InstanceID) ([]byte, error)

	// OnWrite implements Write operation
	//  method: POST/PUT
	//  format: LwM2M CBOR, SenML CBOR, SenML JSON, or TLV
	//  path: /{Object ID}/{Object Instance ID}
	//        /{Object ID}/{Object Instance ID}/{Resource ID}
	//        /{Object ID}/{Object Instance ID}/{Resource ID}/{Resource Instance ID}
	OnWrite(oid ObjectID, oiId InstanceID, rid ResourceID, riId InstanceID, newValue []byte) ([]byte, error)

	// OnDelete implements Delete operation
	//  method: DELETE
	//  path: /{Object ID}/{Object Instance ID}
	//        /{Object ID}/{Object Instance ID}/{Resource ID}/{Resource Instance ID}
	OnDelete(oid ObjectID, oiId InstanceID, rid ResourceID, riId InstanceID) error

	// OnExecute implements Write operation
	//  method: POST/PUT
	//  format: none, or text/plain
	//  path: /{Object ID}/{Object Instance ID}/{Resource ID}
	OnExecute(oid ObjectID, oiId InstanceID, rid ResourceID, args string) error

	// OnDiscover implements Discover operation
	//  method: GET
	//  path: /{Object ID}<Depth>
	//        /{Object ID}/{Object Instance ID}<Depth>
	//        /{Object ID}/{Object Instance ID}/{Resource ID}<Depth>
	//        Depth: ?depth={Value}
	OnDiscover(oid ObjectID, oiId InstanceID, rid ResourceID, depth int) error
	//OnReadComposite()
	//OnWriteComposite()
	//OnWriteAttributes()

}
