package blockchain

import (
	"time"
)

// DNSSeed identifies a DNS seed.
type DNSSeed struct {
	// Host defines the hostname of the seed.
	Host string

	// HasFiltering defines whether the seed supports filtering
	// by service flags (wire.ServiceFlag).
	HasFiltering bool
}

/*
Params defines a network by its params. These params may be used by Applications
to differentiate network as well as addresses and keys for one network
from those intended for use on another network
*/
type Params struct {
	// Name defines a human-readable identifier for the network.
	Name string

	// Net defines the magic bytes used to identify the network.
	Net uint32

	// DefaultPort defines the default peer-to-peer port for the network.
	DefaultPort string

	// DNSSeeds defines a list of DNS seeds for the network that are used
	// as one method to discover peers.
	DNSSeeds []string

	// GenesisBlock defines the first block of the chain.
	GenesisBlock *Block

	// SubsidyReductionInterval is the interval of blocks before the subsidy
	// is reduced.
	SubsidyReductionInterval int32

	// TargetTimespan is the desired amount of time that should elapse
	// before the block difficulty requirement is examined to determine how
	// it should be changed in order to maintain the desired block
	// generation rate.
	TargetTimespan time.Duration

	// TargetTimePerBlock is the desired amount of time to generate each
	// block.
	TargetTimePerBlock time.Duration

	// RetargetAdjustmentFactor is the adjustment factor used to limit
	// the minimum and maximum amount of adjustment that can occur between
	// difficulty retargets.
	RetargetAdjustmentFactor int64

	// ReduceMinDifficulty defines whether the network should reduce the
	// minimum required difficulty after a long enough period of time has
	// passed without finding a block.  This is really only useful for test
	// networks and should not be set on a main network.
	ReduceMinDifficulty bool

	// MinDiffReductionTime is the amount of time after which the minimum
	// required difficulty should be reduced when a block hasn't been found.
	//
	// NOTE: This only applies if ReduceMinDifficulty is true.
	MinDiffReductionTime time.Duration

	// GenerateSupported specifies whether or not CPU mining is allowed.
	GenerateSupported bool
}

// TODO: generate for mainnet
var preSelectValidatorsMainnet = []string{}

// MainNetParams defines the network parameters for the main coin network.
var MainNetParams = Params{
	Name:        MainetName,
	Net:         Mainnet,
	DefaultPort: MainnetDefaultPort,
	DNSSeeds: []string{
		//"/ip4/127.0.0.1/tcp/9333/ipfs/QmRuvXN7BpTqxqpPLSftDFbKEYiRZRUb7iqZJcz2CxFvVS",
	},

	// blockChain parameters
	GenesisBlock: GenesisBlockGenerator{}.CreateGenesisBlockPoSParallel(1, MainnetGenesisblockPaymentAddress, preSelectValidatorsMainnet, MainnetInitFundSalary, 0, 0),
}

// TestNetParams defines the network parameters for the test coin network.
var preSelectValidatorsTestnet = []string{
	"1YYse7WkP3yjHrX6zfdCwtpyzgNVCAgiXnu1jbtaLqYmPU9tErFd2vm4JAUYXjo8LJFJs3ngZJERYYzVe3XJAC7GB6HZy8",
	"1CQfmbeysZj3V6GPCkVKkdcLDMt9tgHAk8RuF3vnV8shaDibZvREo5CccwVnxypMVRPkEJ5joK79wwCvDvDaZJVQ2uBQ2M",
	"1UJwNHRTnTmGLAxaLirX9DkbH6jk2yhvLUTep5AnQD5pHzHXGh48vN2uJAnwLAJ7YFbsRRYujtizwyVDu87qUkdgdJq5Xt",
	"1aFvm9XfRjzkZxNuXkAtuTcAsTVNqd3vbikXqy49ZfeA4ZVxMhReXWt5uL7w259kM7fkUwzoYs5krDDmC8NpQBK3YoRbA8",
	"15WzsoCCovjyQKtucFRsVdpzs9ncwp6zZ1Cewyk2b2mXrMXpM2W8nGHTJp9H3EiHzJsicCw8zohNmsiVRrVx38DzZkDGUJ",
	"1BXiijqgYs8uRUMBVdYzdgzQEjSUy1X5VAXzJRvA41NZQTXUqaPdKjfx2AunZxdJ8HDzjFtMmdXcEkmgMpg1hVrfqDpAfM",
	"1WpfJdAkLF6VeGZBG5hhUmvdH7FWDFYGswd24oguUu8mXvVdiMq8ayZG1dxKrmLvmNmUn3CHaVBBGe2Z4W81DJDnag96QS",
	"1NsWjmo7nxRHvzKs5pNSRvnF2632vnzy7musPxJt7qRXzuSsvG5wL3FbvHgMEaHUGSRCrmn3Fr4sACwyn2w7LYAQE5wE3R",
	"1UGEwAatRb3MSs8WysByLsATp4KSASQwD9wsGP3g9RyDBC25mPSPgpHYJ7Hf2B61K7bhDnwvnUnWqZNkX3k1tpz8u91uqy",
	"17ErvSN2GSXt6EBsyHyctYvB5uKDfYvB3Zk3K8MNUmbnRhcakcq96Px5QzFqo4sYe29d99GYXXfo9n5hut6ZqxQNjpEsdd",
	"1VbYQf4MHPevQLKsE8hniDZM2kcxHKM6UnwdkkuoZQzTghSAkCGWeg4gdzk4RBVfiGwKmFJpgnpoAh5qWoQmNY46PijndF",
	"1QRE1SJcnVuE9WooTLrcdCWhFh5xdwH7GFiLygW5JqeN2g3PZUc17nR1TLAFbDFkWeZrZRTia1gadRUMPkLSw4dRygXzct",
	"1SESRhnQEhpYuaJt6tzNT8cjos4wvgkLe7YBAHLmhtj9dfPVJuYHsgzbP5tyj93WcPvzXtfmBKw6pHCLpcZMkRS7p5n7EE",
	"1DPxJgGSZqeJnBD7snQd8EwAQABKy2wCEpffAUcUKZuVS1Q6fAHqTJJCbhzZBHUcC9Jjygj64JCkqQjZNcsFsnWSv3Wkqf",
	"1Xks2EnEMaQfvMxr1cfqJ1qXdPD6UaUERQPjVL1gUjuCiFJrDWqCqdf6RqKDkei6TKTC2kdL4D6Nqn244s9wKa6KzpCm8h",
	"1YKgSVcnGsdyW2ywwbKVebpsnGdN6QvNY8186P9SYj6mRAkroEcq5WtL2z8PphVUX7GnVSNQazisqRFNbGUXV32gswFqdh",
	"1HbJh7KfCHx69qNfdvmVtCZeiupjYqnoBuFSZCJbfcboHQpWbTvKXUPfF68ARzENQokn1aW5so9KFAdsgSRDGLwvcAAqob",
	"1W3ZgZ7FFY3QDQ2X1i34FzDyfPZJVgywrdYmMcEeqGQu4NBkULFetf1FDFASr4ZHx7bXkbJc2uUqxnnbimnGZHDpXGPqQk",
	"1H1sFe3LUbtDu3ZkbPgDAMHiqz9ryUG1VLhLsPGA2nXdMdkqDkh8sER9AfoDyBhTaHLx9jT4bxUADZe2dvbGX6NsBKirCF",
	"1TuVSHsrsmfgHYtTb42UKCghdvyR6GyccjGF3mk53j9f5Qs6DcHfRV2ePpVTLzpNstTSciS3AAYhE91eNogrk7h5gSb7R7",
}
var TestNetParams = Params{
	Name:        TestnetName,
	Net:         Testnet,
	DefaultPort: TestnetDefaultPort,
	DNSSeeds: []string{
		//"/ip4/127.0.0.1/tcp/9333/ipfs/QmRuvXN7BpTqxqpPLSftDFbKEYiRZRUb7iqZJcz2CxFvVS",
	},

	// blockChain parameters
	GenesisBlock: GenesisBlockGenerator{}.CreateGenesisBlockPoSParallel(1, TestnetGenesisBlockPaymentAddress, preSelectValidatorsTestnet, TestnetInitFundSalary, 1000, 1000),
}
