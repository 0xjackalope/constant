package blockchain

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ninjadotorg/constant/blockchain/params"
	"github.com/ninjadotorg/constant/cashec"
	"github.com/ninjadotorg/constant/common"
	"github.com/ninjadotorg/constant/common/base58"
	"github.com/ninjadotorg/constant/metadata"
	"github.com/ninjadotorg/constant/privacy"
)

/*
	Load beststate + block of current block from cache to create new block
	Because final beststate height should behind highest block 1
	For example: current height block: 91, final beststate should be 90, new block height is 92

	Create Block (Body and Header)
	* Header:
		1. Create Producer: public key of child 0 or from config
		2. Create Version: load current version
		3. Create Height: prev block height + 1
		4. Create Epoch: Epoch ++ if height % epoch 0
		5. Create timestamp: now
		6. Attach previous block hash

	* Header & Body
		7. Create Shard State:
			- Shard State Vulue from beaconblockpool
			- Shard State Hash
			- Get new staker from shard(beacon or pool) -> help to create Instruction
			- Swap validator from shard -> help to create Instruction
		8. Create Instruction:
			- Instruction value -> body
			- Instruction Hash -> Header
		9. Process Instruction with best state:
			- Create Validator Root -> Header
			- Create BeaconCandidate Root -> Header
	Sign:
		Sign block and update validator index, agg sig
*/
func (blkTmplGenerator *BlkTmplGenerator) NewBlockBeacon(payToAddress *privacy.PaymentAddress, privateKey *privacy.SpendingKey, round int, shardsToBeacon map[byte]uint64) (*BeaconBlock, error) {
	beaconBlock := &BeaconBlock{}
	beaconBestState := BestStateBeacon{}
	// lock blockchain
	blkTmplGenerator.chain.chainLock.Lock()
	// fmt.Printf("Beacon Produce: BeaconBestState Original %+v \n", blkTmplGenerator.chain.BestState.Beacon)
	// produce new block with current beststate
	tempMarshal, err := json.Marshal(*blkTmplGenerator.chain.BestState.Beacon)
	if err != nil {
		return nil, NewBlockChainError(MashallJsonError, err)
	}
	err = json.Unmarshal(tempMarshal, &beaconBestState)
	if err != nil {
		return nil, NewBlockChainError(UnmashallJsonBlockError, err)
	}
	beaconBestState.CandidateShardWaitingForCurrentRandom = blkTmplGenerator.chain.BestState.Beacon.CandidateShardWaitingForCurrentRandom
	beaconBestState.CandidateShardWaitingForNextRandom = blkTmplGenerator.chain.BestState.Beacon.CandidateShardWaitingForNextRandom
	beaconBestState.CandidateBeaconWaitingForCurrentRandom = blkTmplGenerator.chain.BestState.Beacon.CandidateBeaconWaitingForCurrentRandom
	beaconBestState.CandidateBeaconWaitingForNextRandom = blkTmplGenerator.chain.BestState.Beacon.CandidateBeaconWaitingForNextRandom
	// fmt.Printf("Beacon Produce: BeaconBestState Copy %+v \n", beaconBestState)
	// fmt.Printf("Beacon Produce: BeaconBestState Compare %+v \n", reflect.DeepEqual(beaconBestState, *blkTmplGenerator.chain.BestState.Beacon))
	if reflect.DeepEqual(beaconBestState, BestStateBeacon{}) {
		panic(NewBlockChainError(BeaconError, errors.New("problem with beststate in producing new block")))
	}

	// unlock blockchain
	blkTmplGenerator.chain.chainLock.Unlock()

	//==========Create header
	beaconBlock.Header.Producer = base58.Base58Check{}.Encode(payToAddress.Pk, byte(0x00))
	beaconBlock.Header.Version = VERSION
	beaconBlock.Header.Height = beaconBestState.BeaconHeight + 1
	beaconBlock.Header.Epoch = beaconBestState.Epoch
	beaconBlock.Header.Round = round
	// Eg: Epoch is 200 blocks then increase epoch at block 201, 401, 601
	if beaconBlock.Header.Height%common.EPOCH == 1 {
		beaconBlock.Header.Epoch++
	}
	beaconBlock.Header.Timestamp = time.Now().Unix()
	beaconBlock.Header.PrevBlockHash = beaconBestState.BestBlockHash
	tempShardState, staker, swap, stabilityInstructions := blkTmplGenerator.GetShardState(&beaconBestState, shardsToBeacon)
	tempInstruction := beaconBestState.GenerateInstruction(beaconBlock, staker, swap, beaconBestState.CandidateShardWaitingForCurrentRandom, stabilityInstructions)

	//Add Voting instruction
	votingInstruction, err := blkTmplGenerator.chain.generateVotingInstruction(privateKey, 0)
	if err != nil {
		return nil, NewBlockChainError(BeaconError, err)
	}
	tempInstruction = append(tempInstruction, votingInstruction...)

	//==========Create Body
	beaconBlock.Body.Instructions = tempInstruction
	beaconBlock.Body.ShardState = tempShardState
	//==========End Create Body
	//============Process new block with beststate
	fmt.Println("Beacon Candidate", beaconBestState.CandidateBeaconWaitingForCurrentRandom)
	fmt.Println("Beacon Produce: Beacon Instruction", beaconBlock.Body.Instructions)
	fmt.Printf("Beacon Produce/AfterUpdate: Beacon Pending Validator %+v \n, Beacon Committee %+v \n, Beacon Validator Root %+v \n", beaconBestState.BeaconPendingValidator, beaconBestState.BeaconCommittee, beaconBlock.Header.ValidatorsRoot)
	fmt.Println("=======================================")
	fmt.Printf("Beacon Produce/Before Update: Shard Pending Validator %+v \n, ShardCommitee %+v \n, Shard Validator Root %+v \n", beaconBestState.ShardPendingValidator, beaconBestState.ShardCommittee, beaconBlock.Header.ShardValidatorsRoot)
	beaconBestState.Update(beaconBlock)
	//============End Process new block with beststate
	//==========Create Hash in Header
	// BeaconValidator root: beacon committee + beacon pending committee
	validatorArr := append(beaconBestState.BeaconCommittee, beaconBestState.BeaconPendingValidator...)
	beaconBlock.Header.ValidatorsRoot, err = GenerateHashFromStringArray(validatorArr)
	fmt.Printf("Beacon Produce/AfterUpdate: Beacon Pending Validator %+v , Beacon Committee %+v, Beacon Validator Root %+v \n", beaconBestState.BeaconPendingValidator, beaconBestState.BeaconCommittee, beaconBlock.Header.ValidatorsRoot)
	if err != nil {
		panic(err)
	}
	// BeaconCandidate root: beacon current candidate + beacon next candidate
	beaconCandidateArr := append(beaconBestState.CandidateBeaconWaitingForCurrentRandom, beaconBestState.CandidateBeaconWaitingForNextRandom...)
	beaconBlock.Header.BeaconCandidateRoot, err = GenerateHashFromStringArray(beaconCandidateArr)
	if err != nil {
		panic(err)
	}
	// Shard candidate root: shard current candidate + shard next candidate
	shardCandidateArr := append(beaconBestState.CandidateShardWaitingForCurrentRandom, beaconBestState.CandidateShardWaitingForNextRandom...)
	beaconBlock.Header.ShardCandidateRoot, err = GenerateHashFromStringArray(shardCandidateArr)
	if err != nil {
		panic(err)
	}
	// Shard Validator root
	beaconBlock.Header.ShardValidatorsRoot, err = GenerateHashFromMapByteString(beaconBestState.ShardPendingValidator, beaconBestState.ShardCommittee)
	fmt.Printf("Beacon Produce/AfterUpdate: Shard Pending Validator %+v , ShardCommitee %+v, Shard Validator Root %+v \n", beaconBestState.ShardPendingValidator, beaconBestState.ShardCommittee, beaconBlock.Header.ShardValidatorsRoot)
	if err != nil {
		panic(err)
	}
	// Shard state hash
	tempShardStateHash, err := GenerateHashFromShardState(tempShardState)
	if err != nil {
		Logger.log.Error(err)
		return nil, err
	}
	beaconBlock.Header.ShardStateHash = tempShardStateHash
	// Instruction Hash
	tempInstructionArr := []string{}
	for _, strs := range tempInstruction {
		tempInstructionArr = append(tempInstructionArr, strs...)
	}
	tempInstructionHash, err := GenerateHashFromStringArray(tempInstructionArr)
	if err != nil {
		Logger.log.Error(err)
		return nil, err
	}
	beaconBlock.Header.InstructionHash = tempInstructionHash
	//===============End Create Header
	//===============Generate Signature
	// Signature of producer, sign on hash of header
	blockHash := beaconBlock.Header.Hash()
	keySet := &cashec.KeySet{}
	keySet.ImportFromPrivateKey(privateKey)
	producerSig, err := keySet.SignDataB58(blockHash.GetBytes())
	if err != nil {
		Logger.log.Error(err)
		return nil, err
	}
	beaconBlock.ProducerSig = producerSig
	//================End Generate Signature
	return beaconBlock, nil
}

// return param:
// #1: shard state
// #2: valid stakers
// #3: swap validator => map[byte][][]string
func (blkTmplGenerator *BlkTmplGenerator) GetShardState(beaconBestState *BestStateBeacon, shardsToBeacon map[byte]uint64) (map[byte][]ShardState, [][]string, map[byte][][]string, [][]string) {
	shardStates := make(map[byte][]ShardState)
	validStakers := [][]string{}
	validSwap := make(map[byte][][]string)
	//Get shard to beacon block from pool
	shardsBlocks := blkTmplGenerator.shardToBeaconPool.GetValidPendingBlock(shardsToBeacon)
	//Shard block is a map ShardId -> array of shard block
	stabilityInstructions := [][]string{}
	accumulativeValues := &accumulativeValues{
		saleDataMap: map[string]*params.SaleData{},
	}
	for shardID, shardBlocks := range shardsBlocks {
		// Only accept block in one epoch
		totalBlock := 0
		//UNCOMMENT FOR TESTING
		fmt.Println("Beacon Producer/ Got These Block from pool")
		for _, shardBlocks := range shardBlocks {
			fmt.Printf(" %+v ", shardBlocks.Header.Height)
		}
		fmt.Println()
		//=======
		for index, shardBlock := range shardBlocks {
			currentCommittee := beaconBestState.ShardCommittee[shardID]
			hash := shardBlock.Header.Hash()
			err1 := ValidateAggSignature(shardBlock.ValidatorsIdx, currentCommittee, shardBlock.AggregatedSig, shardBlock.R, &hash)
			fmt.Println("Beacon Producer/ Validate Agg Signature for shard", shardID, err1 == nil)
			if err1 != nil {
				break
			}
			stabilityInstructionsPerBlock, err := blkTmplGenerator.buildStabilityInstructions(
				shardID,
				shardBlock.Instructions,
				beaconBestState,
				accumulativeValues,
			)
			if err != nil {
				panic(fmt.Sprintf("Build stability instructions failed: %s", err.Error()))
			}
			stabilityInstructions = append(stabilityInstructions, stabilityInstructionsPerBlock...)
			for _, l := range shardBlock.Instructions {
				if l[0] == "swap" {
					if l[3] != "shard" || l[4] != strconv.Itoa(int(shardID)) {
						panic("Swap instruction is invalid")
					}
				}
			}
			if index != 0 && err1 != nil {
				break
			}
			totalBlock = index
		}
		fmt.Printf("Beacon Producer/ AFTER FILTER, ONLY GET %+v block \n", totalBlock)
		fmt.Println("Beacon Producer/ FILTER and ONLY GET These Block from pool")
		// for _, shardBlocks := range shardBlocks[:totalBlock+1] {
		// 	fmt.Printf(" %+v ", shardBlocks.Header.Height)
		// }
		fmt.Println()
		for _, shardBlock := range shardBlocks[:totalBlock+1] {
			stakers := [][]string{}
			swaps := [][]string{}
			fmt.Printf("Becon Produce: Got Shard Block %+v Shard %+v \n", shardBlock.Header.Height, shardID)
			// for each shard block, create a corresponding shard state
			instructions := shardBlock.Instructions
			shardState := ShardState{}
			shardState.CrossShard = make([]byte, len(shardBlock.Header.CrossShards))
			copy(shardState.CrossShard, shardBlock.Header.CrossShards)
			fmt.Println("Beacon Producer/ CrossShard Byte Arr", shardState.CrossShard)
			shardState.Hash = shardBlock.Header.Hash()
			shardState.Height = shardBlock.Header.Height
			shardStates[shardID] = append(shardStates[shardID], shardState)

			fmt.Printf("\n \n Instruction in shardBlock %+v, %+v \n \n", shardBlock.Header.Height, instructions)
			for _, l := range instructions {
				if l[0] == "stake" {
					fmt.Println("Beacon Producer/ Stake Instructions", l)
					stakers = append(stakers, l)
				} else if l[0] == "swap" {
					swaps = append(swaps, l)
				}
			}
			// ["stake" "pubkey1,pubkey2,..." "shard"]
			// ["stake" "pubkey1,pubkey2,..." "beacon"]
			stakeBeacon := []string{}
			stakeShard := []string{}
			fmt.Println("Beacon Producer/ Process Stakers List", stakers)
			for _, staker := range stakers {
				var tempStaker []string
				newBeaconCandidate, newShardCandidate := getStakeValidatorArrayString(staker)
				assignShard := true
				if !reflect.DeepEqual(newBeaconCandidate, []string{}) {
					tempStaker = make([]string, len(newBeaconCandidate))
					copy(tempStaker, newBeaconCandidate[:])
					assignShard = false
				} else {
					tempStaker = make([]string, len(newShardCandidate))
					copy(tempStaker, newShardCandidate[:])
				}
				tempStaker = blkTmplGenerator.chain.BestState.Beacon.GetValidStakers(tempStaker)
				tempStaker = metadata.GetValidStaker(stakeShard, tempStaker)
				tempStaker = metadata.GetValidStaker(stakeBeacon, tempStaker)
				if len(tempStaker) > 0 {
					if assignShard {
						stakeShard = append(stakeShard, tempStaker...)
					} else {
						stakeBeacon = append(stakeBeacon, tempStaker...)
					}
				}
			}
			if len(stakeShard) > 0 {
				validStakers = append(validStakers, []string{"stake", strings.Join(stakeShard, ","), "shard"})
			}
			if len(stakeBeacon) > 0 {
				validStakers = append(validStakers, []string{"stake", strings.Join(stakeBeacon, ","), "beacon"})
			}
			// format
			// ["swap" "inPubkey1,inPubkey2,..." "outPupkey1, outPubkey2,..." "shard" "shardID"]
			// ["swap" "inPubkey1,inPubkey2,..." "outPupkey1, outPubkey2,..." "beacon"]
			// Validate swap instruction => extract only valid swap instruction
			//TODO: define error handler scheme
			for _, swap := range swaps {
				if swap[3] == "beacon" {
					continue
				} else if swap[3] == "shard" {
					temp, err := strconv.Atoi(swap[4])
					if err != nil {
						continue
					}
					swapShardID := byte(temp)
					if swapShardID != shardID {
						continue
					}
					validSwap[shardID] = append(validSwap[shardID], swap)
				} else {
					continue
				}
			}
		}
	}
	return shardStates, validStakers, validSwap, stabilityInstructions
}

//todo @0xjackalope process instruction without create new tx (eg: update db)
//should be merge with buildStabilityInstruction
func (blkTmplGenerator *BlkTmplGenerator) processInstruction(beaconBestState *BestStateBeacon, instruction []string) error {
	//bestBlock := beaconBestState.BestBlock
	metaType, err := strconv.Atoi(instruction[0])
	if err != nil {
		return err
	}
	contentBytes, err := base64.StdEncoding.DecodeString(instruction[1])
	_ = contentBytes
	if err != nil {
		return err
	}
	switch metaType {
	// process some instruction without create tx (update params,...)
	default:
		return nil
	}

	return nil
}

/*
	- set instruction
	- del instruction
	- swap instruction -> ok
	+ format
	+ ["swap" "inPubkey1,inPubkey2,..." "outPupkey1, outPubkey2,..." "shard" "shardID"]
	+ ["swap" "inPubkey1,inPubkey2,..." "outPupkey1, outPubkey2,..." "beacon"]
	- random instruction -> ok
	- stake instruction -> ok
*/
func (bestStateBeacon *BestStateBeacon) GenerateInstruction(
	block *BeaconBlock,
	stakers [][]string,
	swap map[byte][][]string,
	shardCandidates []string,
	stabilityInstructions [][]string,
) [][]string {
	instructions := [][]string{}
	instructions = append(instructions, stabilityInstructions...)
	//=======Swap
	// Shard Swap: both abnormal or normal swap
	for _, swapInstruction := range swap {
		instructions = append(instructions, swapInstruction...)
	}
	// TODO: beacon unexpeted swap -> pbft
	// Beacon normal swap
	if block.Header.Height%common.EPOCH == 0 {
		swapBeaconInstructions := []string{}
		_, _, swappedValidator, beaconNextCommittee, _ := SwapValidator(bestStateBeacon.BeaconPendingValidator, bestStateBeacon.BeaconCommittee, bestStateBeacon.BeaconCommitteeSize, common.OFFSET)
		if len(swappedValidator) > 0 || len(beaconNextCommittee) > 0 {
			swapBeaconInstructions = append(swapBeaconInstructions, "swap")
			swapBeaconInstructions = append(swapBeaconInstructions, strings.Join(beaconNextCommittee, ","))
			swapBeaconInstructions = append(swapBeaconInstructions, strings.Join(swappedValidator, ","))
			swapBeaconInstructions = append(swapBeaconInstructions, "beacon")
			instructions = append(instructions, swapBeaconInstructions)
		}
	}
	//=======Stake
	// ["stake", "pubkey.....", "shard" or "beacon"]
	// beaconStaker := []string{}
	// shardStaker := []string{}
	instructions = append(instructions, stakers...)
	//=======Random and Assign if random number is detected
	// Time to get random number and no block in this epoch get it
	fmt.Printf("RandomTimestamp %+v \n", bestStateBeacon.CurrentRandomTimeStamp)
	fmt.Printf("=========Epoch %+v \n", block.Header.Epoch)
	fmt.Printf("============height epoch: %+v, RANDOM TIME: %+v \n", block.Header.Height%common.EPOCH+1, common.RANDOM_TIME)
	fmt.Printf("============IsGetRandomNumber %+v \n", bestStateBeacon.IsGetRandomNumber)
	fmt.Printf("===================ShardCandidate %+v \n", shardCandidates)
	if block.Header.Height%common.EPOCH > common.RANDOM_TIME && !bestStateBeacon.IsGetRandomNumber {
		var err error
		// COMMENT FOR TESTING
		// chainTimeStamp, err := btcapi.GetCurrentChainTimeStamp()
		// UNCOMMENT FOR TESTING
		chainTimeStamp := bestStateBeacon.CurrentRandomTimeStamp + 1
		fmt.Printf("============chainTimeStamp %+v \n", chainTimeStamp)
		if err != nil {
			panic(err)
		}
		assignedCandidates := make(map[byte][]string)
		if chainTimeStamp > bestStateBeacon.CurrentRandomTimeStamp {
			var wg sync.WaitGroup
			wg.Add(1)
			randomInstruction, rand := generateRandomInstruction(bestStateBeacon.CurrentRandomTimeStamp, &wg)
			wg.Wait()
			instructions = append(instructions, randomInstruction)
			Logger.log.Infof("RandomNumber %+v", randomInstruction)
			for _, candidate := range shardCandidates {
				shardID := calculateCandidateShardID(candidate, rand, bestStateBeacon.ActiveShards)
				assignedCandidates[shardID] = append(assignedCandidates[shardID], candidate)
			}
			for shardId, candidates := range assignedCandidates {
				shardAssingInstruction := []string{"assign"}
				shardAssingInstruction = append(shardAssingInstruction, strings.Join(candidates, ","))
				shardAssingInstruction = append(shardAssingInstruction, "shard")
				shardAssingInstruction = append(shardAssingInstruction, strconv.Itoa(int(shardId)))
				instructions = append(instructions, shardAssingInstruction)
			}
		}
	}
	return instructions
}

func (bestStateBeacon *BestStateBeacon) GetValidStakers(tempStaker []string) []string {
	for _, committees := range bestStateBeacon.ShardCommittee {
		tempStaker = metadata.GetValidStaker(committees, tempStaker)
	}
	for _, validators := range bestStateBeacon.ShardPendingValidator {
		tempStaker = metadata.GetValidStaker(validators, tempStaker)
	}
	tempStaker = metadata.GetValidStaker(bestStateBeacon.BeaconCommittee, tempStaker)
	tempStaker = metadata.GetValidStaker(bestStateBeacon.BeaconPendingValidator, tempStaker)
	tempStaker = metadata.GetValidStaker(bestStateBeacon.CandidateBeaconWaitingForCurrentRandom, tempStaker)
	tempStaker = metadata.GetValidStaker(bestStateBeacon.CandidateBeaconWaitingForNextRandom, tempStaker)
	tempStaker = metadata.GetValidStaker(bestStateBeacon.CandidateShardWaitingForCurrentRandom, tempStaker)
	tempStaker = metadata.GetValidStaker(bestStateBeacon.CandidateShardWaitingForNextRandom, tempStaker)
	tempStaker = metadata.GetValidStaker(bestStateBeacon.CandidateShardWaitingForNextRandom, tempStaker)
	return tempStaker
}

//===================================Util for Beacon=============================

// ["random" "{blockheight}" "{bitcointimestamp}" "{nonce}" "{timestamp}"]
func generateRandomInstruction(timestamp int64, wg *sync.WaitGroup) ([]string, int64) {
	//COMMENT FOR TESTING
	// msg := make(chan string)
	// go btcapi.GenerateRandomNumber(timestamp, msg)
	// res := <-msg
	// reses := strings.Split(res, (","))
	strs := []string{}
	//UNCOMMENT FOR TESTTING
	reses := []string{"1000", strconv.Itoa(int(timestamp) + 1), "1000"}
	strs = append(strs, "random")
	strs = append(strs, reses...)
	strs = append(strs, strconv.Itoa(int(timestamp)))
	nonce, _ := strconv.Atoi(reses[2])
	wg.Done()
	return strs, int64(nonce)
}

func getStakeValidatorArrayString(v []string) ([]string, []string) {
	beacon := []string{}
	shard := []string{}
	if v[0] == "stake" && v[2] == "beacon" {
		beacon = strings.Split(v[1], ",")
	}
	if v[0] == "stake" && v[2] == "shard" {
		shard = strings.Split(v[1], ",")
	}
	return beacon, shard
}
