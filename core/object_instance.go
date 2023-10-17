package core

type InstanceID = uint16

const (
	NoneID uint16 = 0xFFFF
)

// Object defines an instance of
// an ObjectClass at runtime.
//
//	1 ObjectID -> 1 Object Instance Store
//	1 Object Instance Store -> 0/1/* Object Instances mapped by id
type Object interface {
	GetClass() ObjectClass
	InstanceID() InstanceID
	ResInstManager() *ResInstManager
	SetInstanceID(id InstanceID)
}

type ObjectImpl struct {
	class  ObjectClass
	instId InstanceID
	resMgr *ResInstManager
}

func NewObjectImpl(class ObjectClass, id InstanceID) *ObjectImpl {
	return &ObjectImpl{
		class:  class,
		instId: id,
		resMgr: NewResInstManager(),
	}
}

func (o *ObjectImpl) GetClass() ObjectClass {
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
