package core

type OperationType = int

// Operation Types
const (
	OperationRegister        OperationType = 0
	OperationUpdate          OperationType = 1
	OperationDeregister      OperationType = 2
	OperationRead            OperationType = 3
	OperationDiscover        OperationType = 4
	OperationWrite           OperationType = 5
	OperationWriteAttributes OperationType = 6
	OperationExecute         OperationType = 7
	OperationCreate          OperationType = 8
	OperationDelete          OperationType = 9
	OperationObserve         OperationType = 10
	OperationNotify          OperationType = 11
	OperationCancelObserve   OperationType = 12
)

type ProviderType = int

const (
	OMAObjects ProviderType = iota
	IPSOObjects
	QJWLObjects //customized
)

// ObjectClassStore is a repository
// used to access the object template.
type ObjectClassStore interface {
	// GetClass returns the class identified by id.
	GetClass(id ObjectID) ObjectClass

	// GetClasses returns all classes defined.
	GetClasses() map[ObjectID]ObjectClass

	// GetMandatory returns all mandatory classes.
	GetMandatory() map[ObjectID]ObjectClass
}

// ClassInfoProvider provides object
// information for ObjectClassStore.
type ClassInfoProvider interface {
	// Type returns type of this provider.
	Type() ProviderType

	// Get returns the object classes
	// identified by the given id.
	Get(n ObjectID) ObjectClass

	// GetAll returns all object classes
	// loaded by this provider.
	GetAll() map[ObjectID]ObjectClass
}

// NewClassStore creates an object class repository
// and registers the given object providers as the initial classes.
func NewClassStore(providers ...ClassInfoProvider) ObjectClassStore {
	repo := &InMemoryObjectClassStore{
		//providers: make(map[ProviderType]ClassInfoProvider),
		classes: make(map[ObjectID]ObjectClass),
	}

	for _, p := range providers {
		//repo.providers[p.Type()] = p
		repo.copy(p.GetAll())
	}

	return repo
}

// InMemoryObjectClassStore caches object class info in memory.
type InMemoryObjectClassStore struct {
	// classes provided from providers
	classes map[ObjectID]ObjectClass

	//providers map[ProviderType]ClassInfoProvider
}

func (m *InMemoryObjectClassStore) copy(classes map[ObjectID]ObjectClass) {
	m.classes = classes
}

func (m *InMemoryObjectClassStore) GetClass(n ObjectID) ObjectClass {
	for _, class := range m.classes {
		if class != nil && class.Id() == n {
			return class
		}
	}

	return nil
}

func (m *InMemoryObjectClassStore) GetClasses() map[ObjectID]ObjectClass {
	return m.classes
}

// GetMandatory returns all object classes which are mandatory.
func (m *InMemoryObjectClassStore) GetMandatory() map[ObjectID]ObjectClass {
	var mandatory map[ObjectID]ObjectClass

	for _, class := range m.classes {
		if class.Mandatory() {
			mandatory[class.Id()] = class
		}
	}

	return mandatory
}
