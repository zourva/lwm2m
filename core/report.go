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
	Observe(oid ObjectID, oiId InstanceID, rid ResourceID, riId InstanceID, attrs map[string]any) error

	// CancelObservation implements Cancel Observation operation
	//  method: GET with Observe option= 1
	//  path: /{Object ID}
	//        /{Object ID}/{Object Instance ID}
	//        /{Object ID}/{Object Instance ID}/{Resource ID}
	//        /{Object ID}/{Object Instance ID}/{Resource ID}/{Resource Instance ID}
	CancelObservation(oid ObjectID, oiId InstanceID, rid ResourceID, riId InstanceID) error

	// ObserveComposite implements ObserveComposite operation
	//  method: FETCH with Observe option = 0
	//  format: SenML-ETCH JSON, SenML-ETCH CBOR, SenML CBOR or SenML JSON
	//  path: /?pmin={minimum period}&pmax={maximum period}&epmin= {minimum evaluation period}&
	//          epmax={maximum evaluation period}&con={0 or 1}
	//      URI paths for resources to be observed are provided in request payload
	ObserveComposite() error

	// CancelObservationComposite implements Cancel ObservationComposite operation
	//  method: FETCH with Observe option= 1
	//  path: provided in request payload, different for each other.
	CancelObservationComposite() error

	OnNotify(value []byte) error

	OnSend(value []byte) ([]byte, error)
}

type ReportingClient interface {
	OnObserve(oid ObjectID, oiId InstanceID, rid ResourceID, riId InstanceID, attrs map[string]any) ErrorType
	OnCancelObservation(oid ObjectID, oiId InstanceID, rid ResourceID, riId InstanceID) ErrorType
	OnObserveComposite() ErrorType
	OnCancelObservationComposite() ErrorType

	// Notify implements Notify operation
	//  method: N/A since it's defined as an Asynchronous Response
	//  format: LwM2M CBOR, SenML CBOR, SenML JSON
	Notify(updated *Value) error

	// Send implements Send operation
	//  method: POST
	//  format: LwM2M CBOR, SenML CBOR, SenML JSON
	//  path: /dp
	Send(value []byte) ([]byte, error)
}
