// Copyright 2015 CoreOS, Inc.
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

package updateengine

import (
	"testing"
	"time"

	"github.com/godbus/dbus"
	"google.golang.org/protobuf/proto"
)

func makeSig(t *testing.T, curOp Operation) *dbus.Signal {
	status := makeStat(curOp)
	rawStatus, err := proto.Marshal(status)
	if err != nil {
		t.Fatal("Unable to marshal protobuf")
	}

	return &dbus.Signal{
		Body: []interface{}{rawStatus},
	}
}

func makeStat(curOp Operation) *StatusResult {
	return &StatusResult{
		LastCheckedTime:  0,
		Progress:         0.0,
		CurrentOperation: curOp,
		NewVersion:       "newVer",
		NewSize:          1024,
	}
}

func TestRebootNeededSignal(t *testing.T) {
	c := &Client{
		ch: make(chan *dbus.Signal, signalBuffer),
	}
	r := make(chan *StatusResult)
	s := make(chan struct{})
	var done bool
	go func() {
		c.RebootNeededSignal(r, s)
		done = true
	}()

	if done {
		t.Fatal("RebootNeededSignal stopped prematurely")
	}
	c.ch <- makeSig(t, Operation_UPDATED_NEED_REBOOT)
	if done {
		t.Fatal("RebootNeededSignal stopped prematurely")
	}

	time.Sleep(10 * time.Millisecond)

	select {
	case stat := <-r:
		if !proto.Equal(stat, makeStat(Operation_UPDATED_NEED_REBOOT)) {
			t.Fatalf("bad status received: %#v", stat)
		}
	default:
		t.Fatal("RebootNeededSignal did not send expected Status update")
	}

	if done {
		t.Fatal("RebootNeededSignal stopped prematurely")
	}

	c.ch <- makeSig(t, Operation_DOWNLOADING) // some other ignored signal

	time.Sleep(10 * time.Millisecond)

	select {
	case stat := <-r:
		t.Fatalf("unexpected status on unknown signal: %#v", stat)
	default:
	}

	if done {
		t.Fatal("RebootNeededSignal stopped prematurely")
	}

	close(s)

	time.Sleep(10 * time.Millisecond)

	if !done {
		t.Fatal("RebootNeededSignal did not stop as expected")
	}
}
