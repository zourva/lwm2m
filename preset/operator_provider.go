package preset

import "github.com/zourva/lwm2m/core"

type OMAObjectOperatorProvider struct {
	kind      ProviderType
	operators map[core.ObjectID]core.InstanceOperator
}

func (o *OMAObjectOperatorProvider) Type() ProviderType {
	return o.kind
}

func (o *OMAObjectOperatorProvider) Get(n core.ObjectID) core.InstanceOperator {
	return o.operators[n]
}

func (o *OMAObjectOperatorProvider) GetAll() map[core.ObjectID]core.InstanceOperator {
	return o.operators
}

func (o *OMAObjectOperatorProvider) Set(n core.ObjectID, op core.InstanceOperator) {
	o.operators[n] = op
}

func (o *OMAObjectOperatorProvider) SetAll(all map[core.ObjectID]core.InstanceOperator) {
	o.operators = all
}

type NullOperator struct {
}

func (op *NullOperator) Create(class core.Object) core.ObjectInstance {
	return nil
}

// NewOMAObjectDummyOperatorProvider provides a way to add initial
// objects and instances a client wants to enable.
func NewOMAObjectDummyOperatorProvider() core.InstanceOperatorProvider {
	provider := &OMAObjectOperatorProvider{
		kind:      OMAObjects,
		operators: make(map[core.ObjectID]core.InstanceOperator),
	}

	all := GetAllPreset(OMAObjects)
	for id, _ := range all {
		provider.operators[id] = &NullOperator{}
	}

	return provider
}
