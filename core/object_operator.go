package core

// Operator defines operations that
// can be applied to resources defined
// in an Object.
type Operator interface {
	// Class returns the class
	// bounded with this operator.
	Class() Object
	SetClass(o Object)

	// constructor and destructor

	Construct(inst ObjectInstance) error //Constructs an instance
	Destruct(inst ObjectInstance) error  //destruct an instance

	// field operators

	Add(inst ObjectInstance, rid ResourceID, riId InstanceID, field Field) error
	Update(inst ObjectInstance, rid ResourceID, riId InstanceID, field Field) error
	Get(inst ObjectInstance, rid ResourceID, riId InstanceID) (Field, error)
	GetAll(inst ObjectInstance, rid ResourceID) (*Fields, error)
	Delete(inst ObjectInstance, rid ResourceID, riId InstanceID) error
	Execute(inst ObjectInstance, rid ResourceID, riId InstanceID) error
}

type OperatorMap = map[ObjectID]Operator

// BaseOperator as a base
// impl of Operator interface.
type BaseOperator struct {
	//object class bounded
	class Object
}

func (b *BaseOperator) Class() Object {
	return b.class
}

func (b *BaseOperator) SetClass(object Object) {
	b.class = object
}

func (b *BaseOperator) Construct(inst ObjectInstance) error {
	return ErrorNone
}

func (b *BaseOperator) Destruct(inst ObjectInstance) error {
	return ErrorNone
}

func (b *BaseOperator) Add(inst ObjectInstance, rid ResourceID, riId InstanceID, field Field) error {
	return ErrorNone
}

func (b *BaseOperator) Update(inst ObjectInstance, rid ResourceID, riId InstanceID, field Field) error {
	return ErrorNone
}

func (b *BaseOperator) Get(inst ObjectInstance, rid ResourceID, riId InstanceID) (Field, error) {
	return nil, ErrorNone
}

func (b *BaseOperator) GetAll(inst ObjectInstance, rid ResourceID) (*Fields, error) {
	return nil, ErrorNone
}

func (b *BaseOperator) Delete(inst ObjectInstance, id InstanceID, riId InstanceID) error {
	return ErrorNone
}

func (b *BaseOperator) Execute(inst ObjectInstance, id InstanceID, riId InstanceID) error {
	return ErrorNone
}

var _ Operator = &BaseOperator{}

func NewBaseOperator() *BaseOperator {
	return &BaseOperator{}
}

//// OperatorProvider provides an easy
//// way to set operators.
//type OperatorProvider interface {
//	// Get returns the object classes
//	// operators identified by the given id.
//	Get(n ObjectID) Operator
//
//	Set(n ObjectID, op Operator)
//
//	// GetAll returns all operators
//	// covered by this provider.
//	GetAll() OperatorMap
//
//	SetAll(all OperatorMap)
//
//	Merge(p OperatorProvider)
//}
//
//type operatorProvider struct {
//	operators OperatorMap
//}
//
//func (o *operatorProvider) Get(n ObjectID) Operator {
//	return o.operators[n]
//}
//
//func (o *operatorProvider) GetAll() OperatorMap {
//	return o.operators
//}
//
//func (o *operatorProvider) Set(n ObjectID, op Operator) {
//	o.operators[n] = op
//}
//
//func (o *operatorProvider) SetAll(all OperatorMap) {
//	o.operators = all
//}
//
//func (o *operatorProvider) Merge(p OperatorProvider) {
//	for id, operator := range p.GetAll() {
//		o.operators[id] = operator
//	}
//}
//
//func NewOperatorProvider(ops OperatorMap) OperatorProvider {
//	return &operatorProvider{
//		operators: ops,
//	}
//}
