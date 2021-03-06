// Copyright 2018 The Operator-SDK Authors
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

package framework

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"k8s.io/client-go/rest"
)

type TestCtx struct {
	ID         string
	CleanUpFns []finalizerFn
	Namespace  string
	CRClient   *rest.RESTClient
}

type finalizerFn func() error

func (f *Framework) NewTestCtx(t *testing.T) TestCtx {
	// TestCtx is used among others for namespace names where '/' is forbidden
	prefix := strings.TrimPrefix(
		strings.Replace(
			strings.ToLower(t.Name()),
			"/",
			"-",
			-1,
		),
		"test",
	)

	id := prefix + "-" + strconv.FormatInt(time.Now().Unix(), 10)
	return TestCtx{
		ID: id,
	}
}

func (ctx *TestCtx) GetID() string {
	return ctx.ID
}

// GetObjID returns an ascending ID based on the length of cleanUpFns. It is
// based on the premise that every new object also appends a new finalizerFn on
// cleanUpFns.
func (ctx *TestCtx) GetObjID() string {
	return ctx.ID + "-" + strconv.Itoa(len(ctx.CleanUpFns))
}

func (ctx *TestCtx) Cleanup(t *testing.T) {
	for i := len(ctx.CleanUpFns) - 1; i >= 0; i-- {
		err := ctx.CleanUpFns[i]()
		if err != nil {
			t.Errorf("A cleanup function failed with error: %v\n", err)
		}
	}
}

func (ctx *TestCtx) AddFinalizerFn(fn finalizerFn) {
	ctx.CleanUpFns = append(ctx.CleanUpFns, fn)
}
