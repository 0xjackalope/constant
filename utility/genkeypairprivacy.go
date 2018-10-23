package main

import (
	"github.com/ninjadotorg/cash/wallet"
	"os"
)

func main() {
	argsWithoutProg := os.Args[1:]

	mnemonicGen := wallet.MnemonicGenerator{}
	Entropy, _ := mnemonicGen.NewEntropy(128)
	Mnemonic, _ := mnemonicGen.NewMnemonic(Entropy)
	Seed := mnemonicGen.NewSeed(Mnemonic, argsWithoutProg[0])
}
