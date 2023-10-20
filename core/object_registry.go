package core

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

// ObjectRegistry is a repository
// used to retrieve the object template.
//
// NOTE: write of object classes are implementation dependent.
type ObjectRegistry interface {
	// GetClass returns the class identified by id.
	GetClass(id ObjectID) Object

	// GetClasses returns all classes defined.
	GetClasses() map[ObjectID]Object

	// GetMandatory returns all mandatory classes.
	GetMandatory() map[ObjectID]Object
}

// ObjectInfoProvider provides object
// template info for ObjectRegistry.
type ObjectInfoProvider interface {
	//// Type returns type of this provider.
	//Type() ProviderType

	// Get returns the object classes
	// identified by the given id.
	Get(n ObjectID) Object

	// GetAll returns all object classes
	// loaded by this provider.
	GetAll() map[ObjectID]Object
}

// NewObjectRegistry creates an object class registry
// and registers the given object class providers as the initial classes.
//
// Extended object classes from applications are expected to be passed to
// the protocol by providing extended providers here.
func NewObjectRegistry(providers ...ObjectInfoProvider) ObjectRegistry {
	repo := &DefaultObjectRegistry{
		classes: make(map[ObjectID]Object),
	}

	for _, p := range providers {
		repo.copy(p.GetAll())
	}

	return repo
}

// DefaultObjectRegistry implements a registry
// using map as an object class tree.
type DefaultObjectRegistry struct {
	// classes loaded from providers
	classes map[ObjectID]Object
}

func (m *DefaultObjectRegistry) copy(classes map[ObjectID]Object) {
	m.classes = classes
}

func (m *DefaultObjectRegistry) GetClass(n ObjectID) Object {
	for _, class := range m.classes {
		if class != nil && class.Id() == n {
			return class
		}
	}

	return nil
}

func (m *DefaultObjectRegistry) GetClasses() map[ObjectID]Object {
	return m.classes
}

// GetMandatory returns all object classes which are declared mandatory.
func (m *DefaultObjectRegistry) GetMandatory() map[ObjectID]Object {
	var mandatory map[ObjectID]Object

	for _, class := range m.classes {
		if class.Mandatory() {
			mandatory[class.Id()] = class
		}
	}

	return mandatory
}
