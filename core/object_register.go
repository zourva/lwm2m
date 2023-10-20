package core

// RegisteredObject defines the delegation
// of an object on the server side.
// It is just like a link to the real
// object instance in the client.
type RegisteredObject interface {
	ObjectID() ObjectID
	InstanceID() InstanceID
	ObjectClass() Object
}

type RegisteredObjectImpl struct {
	class  Object
	instId InstanceID
}

func (r *RegisteredObjectImpl) ObjectID() ObjectID {
	return r.class.Id()
}

func (r *RegisteredObjectImpl) InstanceID() InstanceID {
	return r.instId
}

func (r *RegisteredObjectImpl) ObjectClass() Object {
	return r.class
}

func NewRegisteredObject(class Object, oiId InstanceID) RegisteredObject {
	ro := &RegisteredObjectImpl{
		class:  class,
		instId: oiId,
	}

	return ro
}

//
//type ObjectFactory interface {
//	// Create creates a new instance
//	// of the given object class.
//	Create(id ObjectID) ObjectInstance
//
//	// ClassStore returns the class store
//	// referenced when creating instances.
//	ClassStore() ObjectRegistry
//}
//
//type DefaultObjectFactory struct {
//	classStore ObjectRegistry
//}
//
//// NewObjectFactory creates a factory
//func NewObjectFactory(repo ObjectRegistry) ObjectFactory {
//	return &DefaultObjectFactory{classStore: repo}
//}
//
//func (f *DefaultObjectFactory) Create(id ObjectID) ObjectInstance {
//	class := f.classStore.GetClass(id)
//	return NewObjectImpl(class, id)
//}
//
//func (f *DefaultObjectFactory) ClassStore() ObjectRegistry {
//	return f.classStore
//}
