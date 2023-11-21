package core

import "github.com/zourva/lwm2m/objects"

//type OperationType = int
//
//// Operation Types
//const (
//	OperationRegister        OperationType = 0
//	OperationUpdate          OperationType = 1
//	OperationDeregister      OperationType = 2
//	OperationRead            OperationType = 3
//	OperationDiscover        OperationType = 4
//	OperationWrite           OperationType = 5
//	OperationWriteAttributes OperationType = 6
//	OperationExecute         OperationType = 7
//	OperationCreate          OperationType = 8
//	OperationDelete          OperationType = 9
//	OperationObserve         OperationType = 10
//	OperationNotify          OperationType = 11
//	OperationCancelObserve   OperationType = 12
//)

type ObjectMap = map[ObjectID]Object

// ObjectRegistry is a repository
// used to retrieve the object template.
//
// NOTE: write of object classes are implementation dependent.
type ObjectRegistry interface {
	// GetObject returns the class identified
	// by id or nil if not found.
	GetObject(id ObjectID) Object

	// GetObjects returns all classes defined.
	GetObjects() ObjectMap

	// GetMandatory returns all mandatory classes.
	GetMandatory() ObjectMap
}

// ObjectProvider provides object
// template info for ObjectRegistry.
type ObjectProvider interface {
	// Get returns the object classes
	// identified by the given id.
	Get(n ObjectID) Object

	// GetAll returns all object classes
	// loaded by this provider.
	GetAll() ObjectMap
}

// objectRegistry implements ObjectRegistry.
type objectRegistry struct {
	// objects loaded from providers
	objects   ObjectMap
	providers []ObjectProvider
}

func (m *objectRegistry) merge(objects ObjectMap) {
	for id, object := range objects {
		m.objects[id] = object
	}
}

func (m *objectRegistry) GetObject(n ObjectID) Object {
	for _, class := range m.objects {
		if class != nil && class.Id() == n {
			return class
		}
	}

	return nil
}

func (m *objectRegistry) GetObjects() ObjectMap {
	return m.objects
}

// GetMandatory returns all object classes which are declared mandatory.
func (m *objectRegistry) GetMandatory() ObjectMap {
	var mandatory ObjectMap

	for _, class := range m.objects {
		if class.Mandatory() {
			mandatory[class.Id()] = class
		}
	}

	return mandatory
}

type objectProvider struct {
	objects ObjectMap
}

func (o *objectProvider) build(descriptors []string) {
	for _, desc := range descriptors {
		obj := ParseObject(desc)
		o.objects[obj.Id()] = obj
	}
}

func (o *objectProvider) Get(n ObjectID) Object {
	return o.objects[n]
}

func (o *objectProvider) GetAll() ObjectMap {
	return o.objects
}

func NewObjectProvider(descriptors []string) ObjectProvider {
	provider := &objectProvider{
		objects: make(ObjectMap),
	}

	provider.build(descriptors)

	return provider
}

// NewObjectRegistry creates an object class registry
// and registers objects spawned according to the descriptors
// as the initial classes.
//
// Extended object classes by application are expected to be passed to
// the protocol by providing their own object descriptors.
//
// NOTE: OMA built-in object descriptors are initialized internally
// and thus no need to be passed in the descriptor group, and will be
// merged if they are still provided in case.
func NewObjectRegistry(descriptorsGroup ...[]string) ObjectRegistry {
	repo := &objectRegistry{
		objects: make(ObjectMap),
	}

	descriptorsGroup = append(descriptorsGroup, objects.GetOMAObjectDescriptors())

	var providers []ObjectProvider
	for _, descriptors := range descriptorsGroup {
		provider := NewObjectProvider(descriptors)
		providers = append(providers, provider)
	}

	for _, p := range providers {
		repo.merge(p.GetAll())
	}

	// save providers
	repo.providers = providers

	return repo
}
