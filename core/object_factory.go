package core

type ObjectFactory interface {
	// Create creates a new instance
	// of the given object class.
	Create(id ObjectID) Object

	// ClassStore returns the class store
	// referenced when creating instances.
	ClassStore() ObjectClassStore
}

type ObjectPersistence interface {
	Load() (map[ObjectID]*InstanceStore, error)
	Flush(objects map[ObjectID]*InstanceStore) error
}

type DefaultObjectFactory struct {
	classStore ObjectClassStore
}

func NewObjectFactory(repo ObjectClassStore) ObjectFactory {
	return &DefaultObjectFactory{classStore: repo}
}

func (f *DefaultObjectFactory) Create(id ObjectID) Object {
	class := f.classStore.GetClass(id)
	return NewObjectImpl(class, id)
}

func (f *DefaultObjectFactory) ClassStore() ObjectClassStore {
	return f.classStore
}
