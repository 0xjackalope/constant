#!/usr/bin/env bash
if [ -z "$1" ]
then
    echo "Please enter node number to start"
    exit 0
fi

KEY1="112t8rnX2MipfoyC4v7DEZQuuddwEmmc6MJjMKpfXdVgeJMU4EzBy9qvH9KQyuc8uvWmBKecoBYAmnKWhp9yaV1acVFVHPeGXWYTKTnLuymj"
KEY2="112t8rnXTWWvWGBHRUb936UrdswQafhAWngq2Mnvh9B4tKPvoQ3y3ZYDbB6nv7QSh64CqNcS2RqeGKLZ3sJ6o9wYjY3kVneLGWavzRDuu5Yq"
KEY3="112t8rnXoyYJeCRStuF9uUA41nsSZUxLkeVCsEt5aiJVACZNcqfY7KmEgTE9tSocbvvJTbdoGDTSWwbwyJw4xwAri8SyYpYPibSPNwy7DbCR"
KEY4="112t8rnXtd7boamSNmoyy5Am884GP1DryaaFwkdPg51eY8riWtFyKXk4ixMHbLyf2DB5SNPrbyExaK2ZTfYLAoE7FgyYzw85hDpjrmkXkgjW"
KEY5="112t8rnYHkdSizoyW9t6Zu5b7mXg29PsxkM89Ss5MtQjBBTeDXbSZyEZpeS8BkkA2ckkUUPTNFeVCj9Yc3YhmAADxML13EL6kaWCUrtH7Dqd"
KEY6="112t8rnYk8h1vvEWqb7j49BRmbKyYankYSpzcnNwYmU5WafJDFeUAgmPuyLJqTrBKB5mFzBfSid7CNGpaTvDKiwnTjRw4mi3y1x3whA8QB7G"
KEY7="112t8rnYonRSSJnKRKvNj6zeU6bv9VJMyGtQpdDgfycj69m5AwvesJgx7wyQbxbqNsvrmbWQ2vqDYJKKmcodkPAbVgyii4QjdLCjYULsGXNd"
KEY8="112t8rnZEr2FYQ69nZSdvCtXyseUD13hPRZJparPaeFWyi4dApbQb1m1mVyrz9jJFMxPgWvxW1u2QaCgib2FVob8ryYxJdmoFHcFsgghPUs9"
KEY9="112t8rnZMQNYEqyJD4qzdcvVu21A5wn7fYmyC8wGz1qYxwbhUkYPUWuPHJt41sdxVEkGgNsSS7xd7ccCqDPZLXW5CcBiFVPyxVBkFKWpux2F"
KEY10="112t8rnZrt6pxiKpysc9KXLU3d8Dym7LyKCAMdhnNLrFXSQWXikhZqhfgofw7u99hNdupFaj4hkyQ6b92euxVjZkmj6615VrZxxJ3iG22jxL"
KEY11="112t8rnZyV1Z2nqatPpU2AwWT5e5GBb6prDPaAEJUfmfvrCanNMDwMK9QvkoBPt9QQ6BQT8URGSd1mS5Li1LUprwGFtbvhCFgbZD8DM1JurP"
KEY12="112t8rnaEVEfKgLa3uGfhk7PokNGrGjWdoipTuE7Ro2q7dhuApM2KuvuMcvbiiXHJZxC96QNByPsBevmkhd7kC3Rvi3NhbdetpmzqFrGcUsU"
KEY13="112t8rnaYsTKqa4dPQyJMCGiN1K3RytG4DE3rqoqmEHGDtQ1JAa4VDJLRyyhDDBe2i8QPJuQUyy4tteaRGu2ekTJtACCJVj1H3wkCYCvdSFr"
KEY14="112t8rnb1TFX2Xa2K3jLdaeTK1vRH1mtsjjF6F74MQpBqRzCBNpZ7Jrm14QoCfcv5mzB1VW87h5V1qMZDhME3JW6Lm1APzAUkpvHKAeDyhuk"
KEY15="112t8rnbA42ftvwiM27KNi3hQVFEphxb9gmZasrc3EogtiXdxR3CW481aem1xL3a4yYVK5VZXVq2Rs7HqPFKzz5VZwi9SzeCf613ANXXWrmi"
KEY16="112t8rnbd2Tw8uqXZ3cAJVbJhFkjRg6LaYFeGmNyKiuLXhPKCLxb3VM9WxENKjdj4pq94ryBtSxiL8p4C3c3LjQPSTaD1P81cQfKWAKjBeou"
KEY17="112t8rnbrN9JgNzQkki3S3PFeanFE9ARYq5iS4d59Cq8iu66ejkKMPDyZhu7jeuHR6xF77DnrEZC3bvEEDVUzaMJ8YvZsShgYpuA3DPK3zPG"
KEY18="112t8rnc9nkQy9UbdeFDiccnNyjhCSdC4jpbw9Mb6LMekh931cRSSWSvwFWNLziivCfWgpVRA1XiQDNRYQ8UkgPXuKyvXTguHmykyoc76dAp"
KEY19="112t8rncGVPyXz95MR8R6HgRMX3ti34dJWEXpdT8xnnvYY95kL1LdNGNjVekdpNrmrLygwwJ7AQ7EeWv1ybJoR47fuAHH91QjGKJKxjWuWfU"
KEY20="112t8rncXRmHSvMEVnPuBh9R8ekexYwa97TvbW5mpq4wubkNVfWqeJtjWx9riKDDcgGMdPmBGhehYLxKbX5rCDWds9AQ7yLdjXrnh912wPqQ"

rm -rf ./data/node-$1/mainnet/block
#rm -rf ./data/node-$1/mainnet/wallet
#rm -rf ./data/node-$1/mainnet/peer.json

mkdir -p ./data/node-$1
rm -rf ./cash-$1
go build
mv ./cash ./cash-$1
PORT=$((9430 + $1))
eval KEY=\${KEY$1}

export EXTERNAL_ADDRESS="127.0.0.1:$PORT"

if [ $1 != 1 ]
then
    ./cash-$1 --listen "127.0.0.1:$PORT" --discoverpeers --discoverpeersaddress "127.0.0.1:9330" --datadir "data/node-$1" --generate --sealerprivatekey $KEY --norpc --testnet
else
    ./cash-$1 --listen "127.0.0.1:$PORT" --discoverpeers --discoverpeersaddress "127.0.0.1:9330" --datadir "data/node-$1" --generate --sealerprivatekey $KEY --rpcuser "ad" --rpcpass "123" --enablewallet --walletpassphrase "12345678" --testnet
fi
