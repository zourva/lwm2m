package preset

import "github.com/zourva/lwm2m/core"

type OMAObjectInfoProvider struct {
	kind   core.ProviderType
	models map[core.ObjectID]core.ObjectClass
}

func NewOMAObjectInfoProvider() core.ClassInfoProvider {
	provider := &OMAObjectInfoProvider{
		kind:   core.OMAObjects,
		models: GetAllPresetClasses(core.OMAObjects),
	}

	return provider
}

func (o *OMAObjectInfoProvider) Type() core.ProviderType {
	return o.kind
}

func (o *OMAObjectInfoProvider) Get(n core.ObjectID) core.ObjectClass {
	return o.models[n]
}

func (o *OMAObjectInfoProvider) GetAll() map[core.ObjectID]core.ObjectClass {
	return o.models
}
