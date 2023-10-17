package core

// see OMA-TS-LightweightM2M_Core-V1_2_1-20221209-A Chapter 6 for details.

type BootstrapClient interface {
	BootstrapRequest()
	BootstrapPackRequest()

	OnBootstrapRead()
	OnBootstrapWrite()
	OnBootstrapDelete()
	OnBootstrapDiscover()
	OnBootstrapFinish()
	OnBootstrapPack()
}

type BootstrapServer interface {
	OnBootstrapRequest()
	OnBootstrapPackRequest()

	BootstrapRead()
	BootstrapWrite()
	BootstrapDelete()
	BootstrapDiscover()
	BootstrapFinish()
	BootstrapPack()
}
