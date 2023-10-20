package core

type InstanceID = uint16

const (
	NoneID uint16 = 0xFFFF
)

// ObjectInstance defines an instance of
// an Object at runtime.
//
//	1 ObjectID -> 1 Object Instance Store
//	1 Object Instance Store -> 0/1/* Object Instances mapped by id
type ObjectInstance interface {
	GetClass() Object
	ResInstManager() *ResInstManager
	InstanceID() InstanceID
	SetInstanceID(id InstanceID)
}

type InstanceOperator interface {
	// Create creates and saves an instance.
	Create(class Object) ObjectInstance
}

// InstanceOperatorProvider provides object instance
// access information for objects defined in ObjectRegistry.
type InstanceOperatorProvider interface {
	//// Type returns type of this provider.
	//Type() ProviderType

	// Get returns the object classes
	// operators identified by the given id.
	Get(n ObjectID) InstanceOperator

	Set(n ObjectID, op InstanceOperator)

	// GetAll returns all operators
	// covered by this provider.
	GetAll() map[ObjectID]InstanceOperator

	SetAll(all map[ObjectID]InstanceOperator)
}

type InstanceStorageManager interface {
	Load() (map[ObjectID]*InstanceManager, error)
	Flush(objects map[ObjectID]*InstanceManager) error
}

type ObjectImpl struct {
	class  Object
	instId InstanceID
	resMgr *ResInstManager
}

func NewObjectImpl(class Object, id InstanceID) *ObjectImpl {
	return &ObjectImpl{
		class:  class,
		instId: id,
		resMgr: NewResInstManager(),
	}
}

func (o *ObjectImpl) GetClass() Object {
	return o.class
}

func (o *ObjectImpl) InstanceID() InstanceID {
	return o.instId
}

func (o *ObjectImpl) SetInstanceID(id InstanceID) {
	o.instId = id
}

func (o *ObjectImpl) ResInstManager() *ResInstManager {
	return o.resMgr
}
