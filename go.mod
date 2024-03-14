module github.com/zourva/lwm2m

go 1.20

require (
	github.com/asdine/storm/v3 v3.2.1
	github.com/pborman/uuid v1.2.1
	github.com/pion/dtls/v2 v2.2.8-0.20240201071732-2597464081c8
	github.com/plgd-dev/go-coap/v3 v3.3.3
	github.com/redis/go-redis/v9 v9.5.1
	github.com/sirupsen/logrus v1.9.3
	github.com/stretchr/testify v1.9.0
	github.com/valyala/fastjson v1.6.4
	github.com/vmihailenco/msgpack/v5 v5.4.1
	github.com/zourva/pareto v0.0.0-20240307081036-01ea22f36ff3
	go.etcd.io/bbolt v1.3.9
)

// note: github.com/pion/dtls/v2/pkg/net package is not yet available in release branches,
// so we force to the use of the pinned master branch
replace github.com/pion/dtls/v2 => github.com/pion/dtls/v2 v2.2.8-0.20240201071732-2597464081c8

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/dsnet/golib/memfile v1.0.0 // indirect
	github.com/fxamacker/cbor/v2 v2.6.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/kr/pretty v0.2.1 // indirect
	github.com/pion/logging v0.2.2 // indirect
	github.com/pion/transport/v3 v3.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	golang.org/x/crypto v0.21.0 // indirect
	golang.org/x/exp v0.0.0-20240222234643-814bf88cf225 // indirect
	golang.org/x/net v0.22.0 // indirect
	golang.org/x/sync v0.6.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
