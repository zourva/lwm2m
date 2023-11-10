package core

type ReportingServer interface {
	// Observe implements Observe operation
	//  method: GET with Observe option = 0
	//  path: /{Object ID}<Attributes>
	//        /{Object ID}/{Object Instance ID}<Attributes>
	//        /{Object ID}/{Object Instance ID}/{Resource ID}<Attributes>
	//        /{Object ID}/{Object Instance ID}/{Resource ID}/{Resource Instance ID}<Attributes>
	//      Attributes: ?pmin={minimum period}&pmax={maximum period}&gt={greater than}&lt={less than}
	//       &st={step}&epmin={minimum evaluation period}&epmax={maximum evaluation period}&edge={0 or 1}
	//      &con={0 or 1}&hqmax={maximum historical queue}
	//  code may be responded:
	//    2.05 Content operation is completed successfully
	//    4.00 Undetermined error occurred
	//    4.01 Unauthorized Access Right Permission Denied
	//    4.04 Not Found URI of Operation is not found
	//    4.05 Method Not Allowed Target is not allowed for "Create" operation
	//    4.06 Not Acceptable The specified Content-Format is not supported
	Observe(oid ObjectID, oiId InstanceID, rid ResourceID, riId InstanceID, attrs map[string]any) error

	// CancelObservation implements Cancel Observation operation
	//  method: GET with Observe option= 1
	//  path: /{Object ID}
	//        /{Object ID}/{Object Instance ID}
	//        /{Object ID}/{Object Instance ID}/{Resource ID}
	//        /{Object ID}/{Object Instance ID}/{Resource ID}/{Resource Instance ID}
	//  code may be responded:
	//    2.05 Content operation is completed successfully
	//    4.00 Undetermined error occurred
	//    4.01 Unauthorized Access Right Permission Denied
	//    4.04 Not Found URI of Operation is not found
	//    4.05 Method Not Allowed Target is not allowed for "Create" operation
	//    4.06 Not Acceptable The specified Content-Format is not supported
	CancelObservation(oid ObjectID, oiId InstanceID, rid ResourceID, riId InstanceID) error

	// ObserveComposite implements ObserveComposite operation
	//  method: FETCH with Observe option = 0
	//  format: SenML-ETCH JSON, SenML-ETCH CBOR, SenML CBOR or SenML JSON
	//  path: /?pmin={minimum period}&pmax={maximum period}&epmin= {minimum evaluation period}&
	//          epmax={maximum evaluation period}&con={0 or 1}
	//      URI paths for resources to be observed are provided in request payload
	//  code may be responded:
	//    2.05 Content operation is completed successfully
	//    4.00 Undetermined error occurred
	//    4.01 Unauthorized Access Right Permission Denied
	//    4.04 Not Found URI of Operation is not found
	//    4.05 Method Not Allowed Target is not allowed for "Create" operation
	//    4.06 Not Acceptable The specified Content-Format is not supported
	//    4.15 Unsupported content format The specified format is not supported
	ObserveComposite() error

	// CancelObservationComposite implements Cancel ObservationComposite operation
	//  method: FETCH with Observe option= 1
	//  path: provided in request payload, different for each other.
	//  code may be responded:
	//    2.05 Content operation is completed successfully
	//    4.00 Undetermined error occurred
	//    4.01 Unauthorized Access Right Permission Denied
	//    4.04 Not Found URI of Operation is not found
	//    4.05 Method Not Allowed Target is not allowed for "Create" operation
	//    4.06 Not Acceptable The specified Content-Format is not supported
	CancelObservationComposite() error

	OnNotify(c RegisteredClient, value []byte) error

	// OnSend may respond code:
	//    2.04 Changed "Send" operation completed successfully
	//    4.00 Undetermined error occurred
	//    4.04 Not Found Reported Object was not registered to the LwM2M Server
	OnSend(c RegisteredClient, value []byte) ([]byte, error)
}

type ReportingClient interface {
	// OnObserve implements server side logic of Observe operation defined in coap.
	// observationId must have the format of /oid/oiid/rid/riid
	OnObserve(observationId string, attrs map[string]any) error
	OnCancelObservation(observationId string) error
	OnObserveComposite() error
	OnCancelObservationComposite() error

	// Notify implements Notify operation
	//  method: N/A since it's defined as an Asynchronous Response
	//  format: LwM2M CBOR, SenML CBOR, SenML JSON
	//  code may be responded:
	//    2.05 Content "Notify" operation completed successfully
	// observationId should be one provided by observing client.
	Notify(observationId string, value []byte) error

	// Send implements Send operation
	//  method: POST
	//  format: LwM2M CBOR, SenML CBOR, SenML JSON
	//  path: /dp
	//  code may be responded:
	//    2.04 Changed "Send" operation completed successfully
	//    4.00 Undetermined error occurred
	//    4.04 Not Found URI of "Create" operation is not found
	Send(value []byte) ([]byte, error)
}
