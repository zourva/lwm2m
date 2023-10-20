package preset

import "github.com/zourva/lwm2m/core"

type OMAObjectInfoProvider struct {
	kind   ProviderType
	models map[core.ObjectID]core.Object
}

func NewOMAObjectInfoProvider() core.ObjectInfoProvider {
	provider := &OMAObjectInfoProvider{
		kind:   OMAObjects,
		models: GetAllPreset(OMAObjects),
	}

	return provider
}

func (o *OMAObjectInfoProvider) Type() ProviderType {
	return o.kind
}

func (o *OMAObjectInfoProvider) Get(n core.ObjectID) core.Object {
	return o.models[n]
}

func (o *OMAObjectInfoProvider) GetAll() map[core.ObjectID]core.Object {
	return o.models
}
