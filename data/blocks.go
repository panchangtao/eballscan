// Copyright 2018 The eballscan Authors
// This file is part of the eballscan.
//
// The eballscan is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The eballscan is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the eballscan. If not, see <http://www.gnu.org/licenses/>.

package data

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ecoball/go-ecoball/common/elog"
	"github.com/muesli/cache2go"
)

const (
	BLOCK_SPAN time.Duration = 10 * time.Second
)

var (
	Blocks          = cache2go.Cache("Blocks")
	log             = elog.NewLogger("data", elog.DebugLog)
	
	Length          int
)

type BlockInfo struct {
	Hash       string
	PrevHash   string
	MerkleHash string
	StateHash  string
	CountTxs   int
	Timestamp  int 
	NumTransaction int
}
type BlockInfoh struct {
	BlockInfo
	Height int
}


func AddBlock(hight int, info *BlockInfo) {
	Blocks.Add(hight, BLOCK_SPAN, info)

}
func PrintBlock() string {
	Blocks.RLock()
	defer Blocks.RUnlock()

	var BlockInfoHArray []BlockInfoh

	for i := 1; i <= Length; i++ {
		res, err := Blocks.Value(i)

		if err == nil {
			One := BlockInfoh{}
			One.Height = i
			One.Hash = res.Data().(*BlockInfo).Hash
			One.PrevHash = res.Data().(*BlockInfo).PrevHash
			One.MerkleHash = res.Data().(*BlockInfo).MerkleHash
			One.StateHash = res.Data().(*BlockInfo).StateHash
			One.CountTxs = res.Data().(*BlockInfo).CountTxs
			One.Timestamp = res.Data().(*BlockInfo).Timestamp
			One.NumTransaction = res.Data().(*BlockInfo).NumTransaction
			BlockInfoHArray = append(BlockInfoHArray, One)
		} else {
			fmt.Println("Error retrieving value from cache:", err)
		}
	}
	buf, _ := json.Marshal(BlockInfoHArray)
	result := string(buf)
	return result

}


