module github.com/taikoxyz/hive/simulators/taiko/sanity/l1-geth

go 1.18

//replace github.com/taikoxyz/hive/ => ../../../../
//
//replace github.com/taikoxyz/hive-taiko-clients => ../../../../../hive-taiko-clients

replace github.com/ethereum/go-ethereum v1.13.1 => github.com/ethereum/go-ethereum v1.11.5

replace github.com/ethereum/go-ethereum v1.13.5-0.20231031113925-bc42e88415d3 => github.com/ethereum/go-ethereum v1.11.5

replace github.com/cockroachdb/pebble => github.com/cockroachdb/pebble v0.0.0-20230404150825-93eff0a72e22

require (
	github.com/ethereum/go-ethereum v1.13.5-0.20231031113925-bc42e88415d3
	github.com/ethereum/hive v0.0.0-20230401205547-71595beab31d
	github.com/taikoxyz/hive/simulators/taiko/common v0.0.0-20240202100531-e788dd38dbe1

)

require (
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/shirou/gopsutil v3.21.11+incompatible // indirect
	github.com/syndtr/goleveldb v1.0.1-0.20220614013038-64ee5596c38a // indirect
	github.com/tklauser/go-sysconf v0.3.12 // indirect
	golang.org/x/sys v0.13.0 // indirect
)

require (
	github.com/DataDog/zstd v1.5.2 // indirect
	github.com/VictoriaMetrics/fastcache v1.12.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/btcsuite/btcd/btcec/v2 v2.3.2 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/cockroachdb/errors v1.9.1 // indirect
	github.com/cockroachdb/logtags v0.0.0-20230118201751-21c54148d20b // indirect
	github.com/cockroachdb/redact v1.1.3 // indirect
	github.com/deckarep/golang-set/v2 v2.3.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.1.0 // indirect
	github.com/getsentry/sentry-go v0.20.0 // indirect
	github.com/gofrs/flock v0.8.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/golang/snappy v0.0.5-0.20220116011046-fa5810519dcb // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/holiman/bloomfilter/v2 v2.0.3 // indirect
	github.com/holiman/uint256 v1.2.3 // indirect
	github.com/julienschmidt/httprouter v1.3.0 // indirect
	github.com/kilic/bls12-381 v0.1.0 // indirect
	github.com/klauspost/compress v1.16.3 // indirect
	github.com/klauspost/cpuid/v2 v2.2.4 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/minio/sha256-simd v1.0.0 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_golang v1.14.0 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.42.0 // indirect
	github.com/prometheus/procfs v0.9.0 // indirect
	github.com/protolambda/bls12-381-util v0.0.0-20220416220906-d8552aa452c7 // indirect
	github.com/protolambda/eth2api v0.0.0-20230316214135-5f8afbd6d05d // indirect
	github.com/protolambda/zrnt v0.30.0 // indirect
	github.com/protolambda/ztyp v0.2.2 // indirect
	github.com/rauljordan/engine-proxy v0.0.0-20230316220057-4c80c36c4c3a // indirect
	github.com/rivo/uniseg v0.4.4 // indirect
	github.com/rogpeppe/go-internal v1.10.0 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	github.com/tklauser/numcpus v0.6.1 // indirect
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/exp v0.0.0-20230905200255-921286631fa9 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	github.com/cockroachdb/pebble v0.0.0-20230928194634-aa077af62593 // indirect
	github.com/rs/cors v1.8.3 // indirect
	github.com/taikoxyz/hive-taiko-clients v0.0.0-20240131133338-ce619c823117 // indirect
	golang.org/x/net v0.17.0 // indirect
)
