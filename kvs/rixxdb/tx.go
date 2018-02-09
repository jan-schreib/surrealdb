// Copyright © 2016 Abcum Ltd
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rixxdb

import (
	"github.com/abcum/rixxdb"
	"github.com/abcum/surreal/kvs"
)

type TX struct {
	pntr *rixxdb.TX
}

func one(res *rixxdb.KV, err error) (kvs.KV, error) {

	switch err {
	case nil:
		break
	default:
		return nil, &kvs.DBError{}
	case rixxdb.ErrTxNotExpectedValue:
		return nil, &kvs.KVError{}
	}

	return res, err

}

func many(res []*rixxdb.KV, err error) ([]kvs.KV, error) {

	switch err {
	case nil:
		break
	default:
		return nil, &kvs.DBError{}
	case rixxdb.ErrTxNotExpectedValue:
		return nil, &kvs.KVError{}
	}

	var out = make([]kvs.KV, len(res))

	for i, v := range res {
		out[i] = v
	}

	return out, err

}

func (tx *TX) Closed() bool {
	return tx.pntr.Closed()
}

func (tx *TX) Cancel() error {
	return tx.pntr.Cancel()
}

func (tx *TX) Commit() error {
	return tx.pntr.Commit()
}

func (tx *TX) Clr(key []byte) (kvs.KV, error) {
	res, err := tx.pntr.Clr(key)
	return one(res, err)
}

func (tx *TX) ClrP(key []byte, max uint64) ([]kvs.KV, error) {
	res, err := tx.pntr.ClrP(key, max)
	return many(res, err)
}

func (tx *TX) ClrR(beg []byte, end []byte, max uint64) ([]kvs.KV, error) {
	res, err := tx.pntr.ClrR(beg, end, max)
	return many(res, err)
}

func (tx *TX) Get(ver int64, key []byte) (kvs.KV, error) {
	res, err := tx.pntr.Get(uint64(ver), key)
	return one(res, err)
}

func (tx *TX) GetP(ver int64, key []byte, max uint64) ([]kvs.KV, error) {
	res, err := tx.pntr.GetP(uint64(ver), key, max)
	return many(res, err)
}

func (tx *TX) GetR(ver int64, beg []byte, end []byte, max uint64) ([]kvs.KV, error) {
	res, err := tx.pntr.GetR(uint64(ver), beg, end, max)
	return many(res, err)
}

func (tx *TX) Del(ver int64, key []byte) (kvs.KV, error) {
	res, err := tx.pntr.Del(uint64(ver), key)
	return one(res, err)
}

func (tx *TX) DelC(ver int64, key []byte, exp []byte) (kvs.KV, error) {
	res, err := tx.pntr.DelC(uint64(ver), key, exp)
	return one(res, err)
}

func (tx *TX) DelP(ver int64, key []byte, max uint64) ([]kvs.KV, error) {
	res, err := tx.pntr.DelP(uint64(ver), key, max)
	return many(res, err)
}

func (tx *TX) DelR(ver int64, beg []byte, end []byte, max uint64) ([]kvs.KV, error) {
	res, err := tx.pntr.DelR(uint64(ver), beg, end, max)
	return many(res, err)
}

func (tx *TX) Put(ver int64, key []byte, val []byte) (kvs.KV, error) {
	res, err := tx.pntr.Put(uint64(ver), key, val)
	return one(res, err)
}

func (tx *TX) PutC(ver int64, key []byte, val []byte, exp []byte) (kvs.KV, error) {
	res, err := tx.pntr.PutC(uint64(ver), key, val, exp)
	return one(res, err)
}
