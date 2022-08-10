module github.com/drgomesp/peerforge

go 1.18

replace (
	github.com/drgomesp/git-remote-go => ../git-remote-go
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	github.com/ipfs-shipyard/git-remote-ipld => ../../ipfs-shipyard/git-remote-ipld
	github.com/keybase/go-keychain => github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4
	google.golang.org/grpc => google.golang.org/grpc v1.33.2
)

require (
	github.com/drgomesp/git-remote-go v0.0.0-00010101000000-000000000000
	github.com/go-git/go-git/v5 v5.4.2
	github.com/ipfs-shipyard/git-remote-ipld v0.0.0-20220621091744-304abe54d48e
	github.com/ipfs/go-cid v0.2.0
	github.com/ipfs/go-ipfs-api v0.3.0
	github.com/joho/godotenv v1.4.0
	github.com/rs/zerolog v1.27.0
	github.com/stretchr/testify v1.8.0
	github.com/urfave/cli/v2 v2.11.1
	gopkg.in/src-d/go-git.v4 v4.13.1
)

require (
	github.com/AndreasBriese/bbloom v0.0.0-20190825152654-46b345b51c96 // indirect
	github.com/Microsoft/go-winio v0.5.2 // indirect
	github.com/ProtonMail/go-crypto v0.0.0-20210428141323-04723f9f07d7 // indirect
	github.com/acomagu/bufpipe v1.0.3 // indirect
	github.com/btcsuite/btcd/btcec/v2 v2.1.3 // indirect
	github.com/cespare/xxhash v1.1.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/crackcomm/go-gitignore v0.0.0-20170627025303-887ab5e44cc3 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1 // indirect
	github.com/dgraph-io/badger v1.6.2 // indirect
	github.com/dgraph-io/ristretto v0.0.3 // indirect
	github.com/dgryski/go-farm v0.0.0-20200201041132-a6ae2369ad13 // indirect
	github.com/dustin/go-humanize v1.0.1-0.20200219035652-afde56e7acac // indirect
	github.com/emirpasic/gods v1.12.0 // indirect
	github.com/go-git/gcfg v1.5.0 // indirect
	github.com/go-git/go-billy/v5 v5.3.1 // indirect
	github.com/gogo/protobuf v1.3.3 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/imdario/mergo v0.3.13 // indirect
	github.com/ipfs/go-block-format v0.0.2 // indirect
	github.com/ipfs/go-ipfs-files v0.0.9 // indirect
	github.com/ipfs/go-ipfs-util v0.0.2 // indirect
	github.com/ipfs/go-ipld-format v0.0.1 // indirect
	github.com/ipfs/go-ipld-git v0.0.2 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/kevinburke/ssh_config v0.0.0-20201106050909-4977a11b4351 // indirect
	github.com/klauspost/cpuid/v2 v2.0.4 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/libp2p/go-buffer-pool v0.1.0 // indirect
	github.com/libp2p/go-flow-metrics v0.0.3 // indirect
	github.com/libp2p/go-libp2p-core v0.19.1 // indirect
	github.com/libp2p/go-openssl v0.0.7 // indirect
	github.com/matryer/is v1.4.0 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/minio/blake2b-simd v0.0.0-20160723061019-3f5f724cb5b1 // indirect
	github.com/minio/sha256-simd v1.0.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mr-tron/base58 v1.2.0 // indirect
	github.com/multiformats/go-base32 v0.0.3 // indirect
	github.com/multiformats/go-base36 v0.1.0 // indirect
	github.com/multiformats/go-multiaddr v0.4.1 // indirect
	github.com/multiformats/go-multibase v0.0.3 // indirect
	github.com/multiformats/go-multicodec v0.4.1 // indirect
	github.com/multiformats/go-multihash v0.0.15 // indirect
	github.com/multiformats/go-varint v0.0.6 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/remeh/sizedwaitgroup v0.0.0-20180822144253-5e7302b12cce // indirect
	github.com/rogpeppe/go-internal v1.8.1 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/sergi/go-diff v1.1.0 // indirect
	github.com/spacemonkeygo/spacelog v0.0.0-20180420211403-2296661a0572 // indirect
	github.com/src-d/gcfg v1.4.0 // indirect
	github.com/whyrusleeping/tar-utils v0.0.0-20180509141711-8c6c8ba81d5c // indirect
	github.com/xanzy/ssh-agent v0.3.0 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	golang.org/x/crypto v0.0.0-20220525230936-793ad666bf5e // indirect
	golang.org/x/net v0.0.0-20220617184016-355a448f1bc9 // indirect
	golang.org/x/sys v0.0.0-20220702020025-31831981b65f // indirect
	golang.org/x/term v0.0.0-20220526004731-065cf7ba2467 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/src-d/go-billy.v4 v4.3.2 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
