module github.com/axone-protocol/axone-mcp

go 1.24

require (
	github.com/axone-protocol/axone-sdk v1.2.0
	github.com/justinas/alice v1.2.0
	github.com/mark3labs/mcp-go v0.21.1
	github.com/mattn/go-isatty v0.0.20
	github.com/rs/zerolog v1.34.0
	github.com/smartystreets/goconvey v1.8.1
	github.com/spf13/cobra v1.9.1
	github.com/spf13/viper v1.20.1
	go.uber.org/mock v0.5.1
	google.golang.org/grpc v1.71.1
	gopkg.in/yaml.v2 v2.4.0
	resenje.org/casbab v0.1.3
)

// Use cutom fork of prolog interpreter, only needed to build
replace github.com/ichiban/prolog => github.com/axone-protocol/prolog v1.0.0

require (
	cosmossdk.io/api v0.7.6 // indirect
	cosmossdk.io/collections v0.4.0 // indirect
	cosmossdk.io/core v0.11.1 // indirect
	cosmossdk.io/depinject v1.1.0 // indirect
	cosmossdk.io/errors v1.0.1 // indirect
	cosmossdk.io/log v1.4.1 // indirect
	cosmossdk.io/math v1.4.0 // indirect
	cosmossdk.io/store v1.1.1 // indirect
	cosmossdk.io/x/tx v0.13.7 // indirect
	cosmossdk.io/x/upgrade v0.1.4 // indirect
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4 // indirect
	github.com/99designs/keyring v1.2.1 // indirect
	github.com/CosmWasm/wasmd v0.53.2 // indirect
	github.com/CosmWasm/wasmvm/v2 v2.1.5 // indirect
	github.com/DataDog/datadog-go v3.2.0+incompatible // indirect
	github.com/DataDog/zstd v1.5.5 // indirect
	github.com/VictoriaMetrics/fastcache v1.5.7 // indirect
	github.com/axone-protocol/axone-contract-schema/go/cognitarium-schema/v6 v6.0.0-20250411103805-21486d26bb1e // indirect
	github.com/axone-protocol/axone-contract-schema/go/dataverse-schema/v6 v6.0.0-20250411103805-21486d26bb1e // indirect
	github.com/axone-protocol/axone-contract-schema/go/law-stone-schema/v6 v6.0.0-20250411103805-21486d26bb1e // indirect
	github.com/axone-protocol/axoned/v10 v10.0.0 // indirect
	github.com/bahlo/generic-list-go v0.2.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bgentry/speakeasy v0.1.1-0.20220910012023-760eaf8b6816 // indirect
	github.com/btcsuite/btcd v0.22.0-beta // indirect
	github.com/btcsuite/btcutil v1.0.3-0.20201208143702-a53e38424cce // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/cenkalti/backoff/v4 v4.1.3 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cockroachdb/apd/v3 v3.2.1 // indirect
	github.com/cockroachdb/errors v1.11.3 // indirect
	github.com/cockroachdb/fifo v0.0.0-20240606204812-0bbfbd93a7ce // indirect
	github.com/cockroachdb/logtags v0.0.0-20230118201751-21c54148d20b // indirect
	github.com/cockroachdb/pebble v1.1.2 // indirect
	github.com/cockroachdb/redact v1.1.5 // indirect
	github.com/cockroachdb/tokenbucket v0.0.0-20230807174530-cc333fc44b06 // indirect
	github.com/cometbft/cometbft v0.38.17 // indirect
	github.com/cometbft/cometbft-db v0.14.1 // indirect
	github.com/cosmos/btcutil v1.0.5 // indirect
	github.com/cosmos/cosmos-db v1.1.1 // indirect
	github.com/cosmos/cosmos-proto v1.0.0-beta.5 // indirect
	github.com/cosmos/cosmos-sdk v0.50.13 // indirect
	github.com/cosmos/go-bip39 v1.0.0 // indirect
	github.com/cosmos/gogogateway v1.2.0 // indirect
	github.com/cosmos/gogoproto v1.7.0 // indirect
	github.com/cosmos/iavl v1.2.2 // indirect
	github.com/cosmos/ibc-go/modules/capability v1.0.1 // indirect
	github.com/cosmos/ibc-go/v8 v8.7.0 // indirect
	github.com/cosmos/ics23/go v0.11.0 // indirect
	github.com/cosmos/ledger-cosmos-go v0.14.0 // indirect
	github.com/danieljoos/wincred v1.1.2 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.3.0 // indirect
	github.com/desertbit/timer v0.0.0-20180107155436-c41aec40b27f // indirect
	github.com/dgraph-io/badger/v4 v4.2.0 // indirect
	github.com/dgraph-io/ristretto v0.1.1 // indirect
	github.com/distribution/reference v0.5.0 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/dustinxie/ecc v0.0.0-20210511000915-959544187564 // indirect
	github.com/dvsekhvalnov/jose2go v1.6.0 // indirect
	github.com/emicklei/dot v1.6.2 // indirect
	github.com/fatih/color v1.15.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/getsentry/sentry-go v0.27.0 // indirect
	github.com/go-jose/go-jose/v3 v3.0.1 // indirect
	github.com/go-kit/kit v0.13.0 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/godbus/dbus v0.0.0-20190726142602-4481cbc300e2 // indirect
	github.com/gogo/googleapis v1.4.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/glog v1.2.4 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/btree v1.1.3 // indirect
	github.com/google/flatbuffers v1.12.1 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/tink/go v1.7.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/gorilla/handlers v1.5.2 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/gsterjov/go-libsecret v0.0.0-20161001094733-a6f4afe4910c // indirect
	github.com/hashicorp/go-hclog v1.5.0 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-metrics v0.5.3 // indirect
	github.com/hashicorp/go-plugin v1.5.2 // indirect
	github.com/hashicorp/golang-lru v1.0.2 // indirect
	github.com/hashicorp/yamux v0.1.1 // indirect
	github.com/hdevalence/ed25519consensus v0.1.0 // indirect
	github.com/huandu/skiplist v1.2.0 // indirect
	github.com/hyperledger/aries-framework-go v0.3.2 // indirect
	github.com/hyperledger/aries-framework-go/component/kmscrypto v0.0.0-20230427134832-0c9969493bd3 // indirect
	github.com/hyperledger/aries-framework-go/component/log v0.0.0-20230427134832-0c9969493bd3 // indirect
	github.com/hyperledger/aries-framework-go/component/models v0.0.0-20230501135648-a9a7ad029347 // indirect
	github.com/hyperledger/aries-framework-go/spi v0.0.0-20230427134832-0c9969493bd3 // indirect
	github.com/hyperledger/ursa-wrapper-go v0.3.1 // indirect
	github.com/iancoleman/strcase v0.3.0 // indirect
	github.com/ichiban/prolog v1.2.0 // indirect
	github.com/improbable-eng/grpc-web v0.15.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jmhodges/levigo v1.0.0 // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/kilic/bls12-381 v0.1.1-0.20210503002446-7b7597926c69 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/linxGnu/grocksdb v1.8.14 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mr-tron/base58 v1.2.0 // indirect
	github.com/mtibben/percent v0.2.1 // indirect
	github.com/multiformats/go-base32 v0.1.0 // indirect
	github.com/multiformats/go-base36 v0.1.0 // indirect
	github.com/multiformats/go-multibase v0.2.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/oasisprotocol/curve25519-voi v0.0.0-20230904125328-1f23a7beb09a // indirect
	github.com/oklog/run v1.1.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/petermattis/goid v0.0.0-20240813172612-4fcff4a6cae7 // indirect
	github.com/piprate/json-gold v0.5.1-0.20230111113000-6ddbe6e6f19f // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/pquerna/cachecontrol v0.1.0 // indirect
	github.com/prometheus/client_golang v1.20.5 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.62.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	github.com/rs/cors v1.11.1 // indirect
	github.com/rs/xid v1.6.0 // indirect
	github.com/sagikazarmark/locafero v0.7.0 // indirect
	github.com/samber/lo v1.39.0 // indirect
	github.com/sasha-s/go-deadlock v0.3.5 // indirect
	github.com/shamaton/msgpack/v2 v2.2.0 // indirect
	github.com/smarty/assertions v1.15.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.12.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/syndtr/goleveldb v1.0.1-0.20220721030215-126854af5e6d // indirect
	github.com/tendermint/go-amino v0.16.0 // indirect
	github.com/teserakt-io/golang-ed25519 v0.0.0-20210104091850-3888c087a4c8 // indirect
	github.com/tidwall/btree v1.7.0 // indirect
	github.com/wk8/go-ordered-map/v2 v2.1.8 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.2.0 // indirect
	github.com/yosida95/uritemplate/v3 v3.0.2 // indirect
	github.com/zondax/hid v0.9.2 // indirect
	github.com/zondax/ledger-go v0.14.3 // indirect
	go.etcd.io/bbolt v1.4.0-alpha.0.0.20240404170359-43604f3112c5 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.35.0 // indirect
	golang.org/x/exp v0.0.0-20240719175910-8a7402abbf56 // indirect
	golang.org/x/net v0.36.0 // indirect
	golang.org/x/sync v0.11.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/term v0.29.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	google.golang.org/genproto v0.0.0-20241118233622-e639e219e697 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250106144421-5f5ef82da422 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250115164207-1a7da9e5054f // indirect
	google.golang.org/protobuf v1.36.4 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gotest.tools/v3 v3.5.1 // indirect
	nhooyr.io/websocket v1.8.6 // indirect
	pgregory.net/rapid v1.1.0 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)
