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
	"time"
	//"encoding/json"
	//"fmt"
	"github.com/muesli/cache2go"
)

const (
	FINAL_BLOCK_SPAN time.Duration = 10 * time.Second
)

var (
	Final_blocks = cache2go.Cache("Final_blocks")

	//THashArray         []string
)

type Final_blockInfo struct {
	TimeStamp          int
	Hash               string
	PrevHash           string
	CMBlockHash        string
	TrxRootHash        string
	StateDeltaRootHash string
	MinorBlocksHash    string
	StateHashRoot      string
	MinorBlockCount    int
	TrxCount           int
	ProposalPubKey     string
	EpochNo            int
}
type Final_blockInfoH struct {
	Final_blockInfo
	Height int
}
