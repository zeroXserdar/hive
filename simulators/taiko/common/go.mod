module github.com/taikoxyz/hive/simulators/taiko/common

go 1.18

//replace github.com/taikoxyz/hive-taiko-clients => ../../../../hive-taiko-clients

//replace github.com/taikoxyz/hive/simulators/taiko/common => ./

//replace github.com/taikoxyz/hive => ../../../

replace github.com/ethereum/go-ethereum v1.13.1 => github.com/ethereum/go-ethereum v1.11.5

replace github.com/ethereum/go-ethereum v1.13.5-0.20231031113925-bc42e88415d3 => github.com/ethereum/go-ethereum v1.11.5

require (
	github.com/ethereum/go-ethereum v1.13.5-0.20231031113925-bc42e88415d3
	github.com/google/uuid v1.3.0
	github.com/herumi/bls-eth-go-binary v1.29.1
	github.com/holiman/uint256 v1.2.3
	github.com/pkg/errors v0.9.1
	github.com/protolambda/bls12-381-util v0.0.0-20220416220906-d8552aa452c7
	//github.com/protolambda/eth2api v0.0.0-20230316214135-5f8afbd6d05d
	github.com/protolambda/go-keystorev4 v0.0.0-20211007151826-f20444f6d564
	github.com/protolambda/zrnt v0.30.0
	github.com/protolambda/ztyp v0.2.2
	github.com/rauljordan/engine-proxy v0.0.0-20230316220057-4c80c36c4c3a
	github.com/taikoxyz/hive-taiko-clients v0.0.0-20240216123906-e77fdbfbe38e
	github.com/tyler-smith/go-bip39 v1.1.0
	github.com/wealdtech/go-eth2-util v1.8.1
	golang.org/x/exp v0.0.0-20230905200255-921286631fa9
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/kr/pretty v0.3.1
	github.com/protolambda/eth2api v0.0.0-20230316214135-5f8afbd6d05d
)

require (
	github.com/Azure/go-ansiterm v0.0.0-20210617225240-d185dfc1b5a1 // indirect
	github.com/DataDog/zstd v1.5.2 // indirect
	github.com/Microsoft/go-winio v0.6.1 // indirect
	github.com/Microsoft/hcsshim v0.9.6 // indirect
	github.com/VictoriaMetrics/fastcache v1.12.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/btcsuite/btcd/btcec/v2 v2.3.2 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/cockroachdb/errors v1.9.1 // indirect
	github.com/cockroachdb/logtags v0.0.0-20230118201751-21c54148d20b // indirect
	github.com/cockroachdb/pebble v0.0.0-20230928194634-aa077af62593 // indirect
	github.com/cockroachdb/redact v1.1.3 // indirect
	github.com/containerd/cgroups v1.0.4 // indirect
	github.com/containerd/containerd v1.6.18 // indirect
	github.com/deckarep/golang-set/v2 v2.3.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.1.0 // indirect
	github.com/docker/docker v24.0.7+incompatible // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/edsrzf/mmap-go v1.1.0 // indirect
	github.com/ferranbt/fastssz v0.1.3 // indirect
	github.com/fsouza/go-dockerclient v1.9.8 // indirect
	github.com/getsentry/sentry-go v0.20.0 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/gofrs/flock v0.8.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/golang/snappy v0.0.5-0.20220116011046-fa5810519dcb // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/hashicorp/yamux v0.1.1 // indirect
	github.com/holiman/bloomfilter/v2 v2.0.3 // indirect
	github.com/julienschmidt/httprouter v1.3.0 // indirect
	github.com/kilic/bls12-381 v0.1.0 // indirect
	github.com/klauspost/compress v1.16.3 // indirect
	github.com/klauspost/cpuid/v2 v2.2.4 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/minio/sha256-simd v1.0.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/moby/patternmatcher v0.6.0 // indirect
	github.com/moby/sys/mount v0.3.3 // indirect
	github.com/moby/sys/mountinfo v0.6.2 // indirect
	github.com/moby/sys/sequential v0.5.0 // indirect
	github.com/moby/term v0.0.0-20210619224110-3f7ff695adc6 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.3-0.20211202183452-c5a74bcca799 // indirect
	github.com/opencontainers/runc v1.1.5 // indirect
	github.com/prometheus/client_golang v1.14.0 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.42.0 // indirect
	github.com/prometheus/procfs v0.9.0 // indirect
	github.com/rivo/uniseg v0.4.4 // indirect
	github.com/rogpeppe/go-internal v1.10.0 // indirect
	github.com/shirou/gopsutil v3.21.11+incompatible // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	github.com/syndtr/goleveldb v1.0.1-0.20220614013038-64ee5596c38a // indirect
	github.com/taikoxyz/hive v0.0.0-20240222070312-5dffa4a5a6c4 // indirect
	github.com/taikoxyz/hive/hiveproxy v0.0.0-20240216041901-c05b16635cb1 // indirect
	github.com/tklauser/go-sysconf v0.3.12 // indirect
	github.com/tklauser/numcpus v0.6.1 // indirect
	github.com/wealdtech/go-bytesutil v1.2.1 // indirect
	github.com/wealdtech/go-eth2-types/v2 v2.8.1 // indirect
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	go.opencensus.io v0.23.0 // indirect
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/mod v0.12.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	golang.org/x/tools v0.13.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/inconshreveable/log15.v2 v2.0.0-20200109203555-b30bc20e4fd1 // indirect
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
