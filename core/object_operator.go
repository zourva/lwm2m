package core

// Operator defines operations that
// can be applied to resources defined
// in an Object.
type Operator interface {
	// Construct creates and saves an instance.
	Construct(class Object) ObjectInstance
}

type OperatorMap = map[ObjectID]Operator

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
