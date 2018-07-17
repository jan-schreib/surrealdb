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

package db

import (
	"fmt"

	"context"

	"github.com/abcum/surreal/sql"
	"github.com/abcum/surreal/util/data"
	"github.com/abcum/surreal/util/guid"
	"github.com/abcum/surreal/util/keys"
)

func (e *executor) executeRelate(ctx context.Context, stm *sql.RelateStatement) ([]interface{}, error) {

	var from, with sql.Exprs

	for _, val := range stm.From {
		f, err := e.fetch(ctx, val, nil)
		if err != nil {
			return nil, err
		}
		from = append(from, f)
	}

	for _, val := range stm.With {
		w, err := e.fetch(ctx, val, nil)
		if err != nil {
			return nil, err
		}
		with = append(with, w)
	}

	i := newIterator(e, ctx, stm, false)

	switch what := stm.Type.(type) {

	default:
		return nil, fmt.Errorf("Can not execute RELATE query using value '%v'", what)

	case *sql.Table:
		key := &keys.Thing{KV: stm.KV, NS: stm.NS, DB: stm.DB, TB: what.TB, ID: guid.New().String()}
		i.processThing(ctx, key)

	case *sql.Ident:
		key := &keys.Thing{KV: stm.KV, NS: stm.NS, DB: stm.DB, TB: what.ID, ID: guid.New().String()}
		i.processThing(ctx, key)

	// Result of subquery
	case []interface{}:
		key := &keys.Thing{KV: stm.KV, NS: stm.NS, DB: stm.DB}
		i.processOther(ctx, key, what)

	// Result of subquery with LIMIT 1
	case map[string]interface{}:
		key := &keys.Thing{KV: stm.KV, NS: stm.NS, DB: stm.DB}
		i.processOther(ctx, key, []interface{}{what})

	}

	return i.Yield(ctx)

}

func (e *executor) fetchRelate(ctx context.Context, stm *sql.RelateStatement, doc *data.Doc) (interface{}, error) {

	ctx = dive(ctx)

	if doc != nil {
		vars := data.New()
		vars.Set(doc.Data(), varKeyParent)
		ctx = context.WithValue(ctx, ctxKeySubs, vars)
	}

	out, err := e.executeRelate(ctx, stm)
	if err != nil {
		return nil, err
	}

	switch len(out) {
	case 1:
		return data.Consume(out).Get(docKeyOne, docKeyId).Data(), nil
	default:
		return data.Consume(out).Get(docKeyAll, docKeyId).Data(), nil
	}

}

func (d *document) runRelate(ctx context.Context, stm *sql.RelateStatement) (interface{}, error) {

	var ok bool
	var err error
	var met = _CREATE

	if err = d.init(ctx); err != nil {
		return nil, err
	}

	if err = d.lock(ctx); err != nil {
		return nil, err
	}

	if err = d.setup(ctx); err != nil {
		return nil, err
	}

	if d.val.Exi() == false {
		met = _UPDATE
	}

	if err = d.merge(ctx, met, stm.Data); err != nil {
		return nil, err
	}

	if ok, err = d.allow(ctx, met); err != nil {
		return nil, err
	} else if ok == false {
		return nil, nil
	}

	if err = d.storeIndex(ctx); err != nil {
		return nil, err
	}

	if err = d.storeThing(ctx); err != nil {
		return nil, err
	}

	if err = d.table(ctx, met); err != nil {
		return nil, err
	}

	if err = d.event(ctx, met); err != nil {
		return nil, err
	}

	if err = d.lives(ctx, met); err != nil {
		return nil, err
	}

	return d.yield(ctx, stm, stm.Echo)

}
