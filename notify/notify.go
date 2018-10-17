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

package notify

import (
	"time"
	"strconv"

	"github.com/ecoball/eballscan/data"
	"github.com/ecoball/eballscan/database"
	"github.com/ecoball/go-ecoball/common"
	"github.com/ecoball/go-ecoball/common/elog"
	"github.com/ecoball/go-ecoball/core/types"
	"github.com/ecoball/go-ecoball/spectator/info"
)

var (
	log = elog.NewLogger("notify", elog.DebugLog)
)

func Dispatch(one info.OneNotify) {
	switch one.InfoType {
	case info.InfoBlock:
		if err := handleBlock(one.Info); nil != err {
			log.Error("handleBlock error: ", err)
		}
	default:

	}
}

func handleBlock(info []byte) error {
	oneBlock := types.Block{}
	if err := oneBlock.Deserialize(info); nil != err {
		log.Fatal(err)
		return err
	}

	//add block
	if err := database.AddBlock(int(oneBlock.Height), int(oneBlock.CountTxs), int(oneBlock.TimeStamp), common.ToHex(oneBlock.Hash.Bytes()), common.ToHex(oneBlock.PrevHash.Bytes()),
		common.ToHex(oneBlock.MerkleHash.Bytes()), common.ToHex(oneBlock.StateHash.Bytes())); nil != err {
		//log.Fatal(err)
		//return err
	}

	data.AddBlock(int(oneBlock.Height), &data.BlockInfo{common.ToHex(oneBlock.Hash.Bytes()), common.ToHex(oneBlock.PrevHash.Bytes()),
		common.ToHex(oneBlock.MerkleHash.Bytes()), common.ToHex(oneBlock.StateHash.Bytes()), int(oneBlock.CountTxs), int(oneBlock.TimeStamp)})

	//add transactions
	for _, v := range oneBlock.Transactions {
		if err := database.AddTransaction(int(v.Type), int(v.TimeStamp), int(oneBlock.Height), common.ToHex(v.Hash.Bytes()),
			v.Permission, v.From.String(), v.Addr.String()); nil != err {
			//log.Fatal(err)
			//return err
		}
		data.AddTransaction(common.ToHex(v.Hash.Bytes()), &data.TransactionInfo{int(v.Type), time.Unix(v.TimeStamp/1000000000, 0).Format("2006-01-02 15:04:05"),
			v.Permission, v.From.String(), v.Addr.String(), int(oneBlock.Height)})
		
		if v.Type == 0x02 {//新增账号交易处理
			info := new(types.InvokeInfo)
			data, err := v.Payload.Serialize()
			if err != nil {
				continue
			}

			err = info.Deserialize(data)
			if err != nil {
				continue
			}

			if string(info.Method) == "new_account" {
				if err := database.AddAccount(info.Param[0], "ABA", int(v.TimeStamp), 0); nil != err {
					return err
				}
				
			}
		}

		if v.Type == 0x03 { //转账交易处理
			info := new(types.TransferInfo)
			data, err := v.Payload.Serialize()
			if err != nil {
				continue
			}

			err = info.Deserialize(data)
			if err != nil {
				continue
			}

			amount, err := strconv.Atoi(info.Value.String())
			if err != nil {
				continue
			}

			//from账户余额处理
			from := v.From.String()
			if from != "root" {
				from_balance, err := database.QueryAccountBalance(from)
				if err != nil {
					continue
				}
				balance := from_balance - amount
				err = database.UpdateAccountBalance(from, balance)
				if err != nil {
					continue
				}
			}
			
			//to账户余额处理
			to := v.Addr.String()
			if to != "root" {
				to := v.Addr.String()
				to_balance, err := database.QueryAccountBalance(to)
				if err != nil {
					continue
				}
				balance := to_balance + amount //to账户余额+
				err = database.UpdateAccountBalance(to, balance)
				if err != nil {
					continue
				}
			}
		}
	}

	return nil
}
