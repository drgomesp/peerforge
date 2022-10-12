module github.com/peerforge/peerforge

go 1.18

replace github.com/peerforge/git-remote-ipldprime => ../git-remote-ipldprime

replace github.com/peerforge/go-ipld-gitprime => ../go-ipld-gitprime

require (
	github.com/dgraph-io/badger/v3 v3.2103.2
	github.com/go-git/go-git/v5 v5.4.2
	github.com/google/uuid v1.3.0
	github.com/ipfs/go-cid v0.3.2
	github.com/ipfs/go-ipfs-api v0.3.0
	github.com/ipld/go-ipld-prime v0.18.0
	github.com/joho/godotenv v1.4.0
	github.com/libp2p/go-libp2p v0.23.2
	github.com/multiformats/go-multihash v0.2.1
	github.com/nuts-foundation/go-did v0.3.0
	github.com/peerforge/git-remote-ipldprime v0.0.0-00010101000000-000000000000
	github.com/peerforge/go-ipld-gitprime v0.0.0-00010101000000-000000000000
	github.com/rs/zerolog v1.28.0
	github.com/stretchr/testify v1.8.0
	github.com/tendermint/tendermint v0.35.9
)

require (
	github.com/AndreasBriese/bbloom v0.0.0-20190825152654-46b345b51c96 // indirect
	github.com/Microsoft/go-winio v0.6.0 // indirect
	github.com/ProtonMail/go-crypto v0.0.0-20220930113650-c6815a8c17ad // indirect
	github.com/acomagu/bufpipe v1.0.3 // indirect
	github.com/benbjohnson/clock v1.3.0 // indirect
	github.com/btcsuite/btcd v0.22.1 // indirect
	github.com/cespare/xxhash v1.1.0 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/cloudflare/circl v1.2.0 // indirect
	github.com/crackcomm/go-gitignore v0.0.0-20170627025303-887ab5e44cc3 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.1.0 // indirect
	github.com/dgraph-io/badger v1.6.2 // indirect
	github.com/dgraph-io/badger/v2 v2.2007.4 // indirect
	github.com/dgraph-io/ristretto v0.1.0 // indirect
	github.com/dgryski/go-farm v0.0.0-20200201041132-a6ae2369ad13 // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/go-git/gcfg v1.5.0 // indirect
	github.com/go-git/go-billy/v5 v5.3.1 // indirect
	github.com/goccy/go-json v0.9.11 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/glog v1.0.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/btree v1.1.2 // indirect
	github.com/google/flatbuffers v1.12.1 // indirect
	github.com/imdario/mergo v0.3.13 // indirect
	github.com/ipfs/go-block-format v0.0.3 // indirect
	github.com/ipfs/go-ipfs-files v0.1.1 // indirect
	github.com/ipfs/go-ipfs-util v0.0.2 // indirect
	github.com/ipfs/go-ipld-git v0.1.1 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/kevinburke/ssh_config v1.2.0 // indirect
	github.com/klauspost/compress v1.15.11 // indirect
	github.com/klauspost/cpuid/v2 v2.1.2 // indirect
	github.com/lestrrat-go/backoff/v2 v2.0.8 // indirect
	github.com/lestrrat-go/blackmagic v1.0.1 // indirect
	github.com/lestrrat-go/httpcc v1.0.1 // indirect
	github.com/lestrrat-go/iter v1.0.2 // indirect
	github.com/lestrrat-go/jwx v1.2.25 // indirect
	github.com/lestrrat-go/option v1.0.0 // indirect
	github.com/lib/pq v1.10.7 // indirect
	github.com/libp2p/go-buffer-pool v0.1.0 // indirect
	github.com/libp2p/go-flow-metrics v0.1.0 // indirect
	github.com/libp2p/go-libp2p-core v0.20.1 // indirect
	github.com/libp2p/go-openssl v0.1.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/mattn/go-pointer v0.0.1 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2 // indirect
	github.com/minio/sha256-simd v1.0.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mr-tron/base58 v1.2.0 // indirect
	github.com/multiformats/go-base32 v0.1.0 // indirect
	github.com/multiformats/go-base36 v0.1.0 // indirect
	github.com/multiformats/go-multiaddr v0.7.0 // indirect
	github.com/multiformats/go-multibase v0.1.1 // indirect
	github.com/multiformats/go-multicodec v0.6.0 // indirect
	github.com/multiformats/go-varint v0.0.6 // indirect
	github.com/oasisprotocol/curve25519-voi v0.0.0-20221003100820-41fad3beba17 // indirect
	github.com/ockam-network/did v0.1.4-0.20210103172416-02ae01ce06d8 // indirect
	github.com/petermattis/goid v0.0.0-20220824145935-af5520614cb6 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/remeh/sizedwaitgroup v1.0.0 // indirect
	github.com/sasha-s/go-deadlock v0.3.1 // indirect
	github.com/sergi/go-diff v1.2.0 // indirect
	github.com/shengdoushi/base58 v1.0.0 // indirect
	github.com/spacemonkeygo/spacelog v0.0.0-20180420211403-2296661a0572 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/stretchr/objx v0.4.0 // indirect
	github.com/tendermint/tm-db v0.6.7 // indirect
	github.com/whyrusleeping/tar-utils v0.0.0-20201201191210-20a61371de5b // indirect
	github.com/xanzy/ssh-agent v0.3.2 // indirect
	go.opencensus.io v0.23.0 // indirect
	golang.org/x/crypto v0.0.0-20221010152910-d6f0a8c073c2 // indirect
	golang.org/x/mod v0.6.0-dev.0.20220419223038-86c51ed26bb4 // indirect
	golang.org/x/net v0.0.0-20221004154528-8021a29435af // indirect
	golang.org/x/sys v0.0.0-20221010170243-090e33056c14 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/tools v0.1.12 // indirect
	google.golang.org/genproto v0.0.0-20221010155953-15ba04fc1c0e // indirect
	google.golang.org/grpc v1.50.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	lukechampine.com/blake3 v1.1.7 // indirect
)
