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
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/godbus/dbus"
)

//
// method call time=1586301148.023978 sender=:1.2974643 -> destination=org.chromium.UpdateEngine serial=3 path=/org/chromium/UpdateEngine; interface=org.chromium.UpdateEngineInterface; member=AttemptUpdateWithFlags
// string "ForcedUpdate"
// string "type='signal',interface='org.freedesktop.DBus',member='NameOwnerChanged',path='/org/freedesktop/DBus',sender='org.freedesktop.DBus',arg0='org.chromium.UpdateEngine'"
// string "org.chromium.UpdateEngine"
// string "type='signal', sender='org.chromium.UpdateEngine', interface='org.chromium.UpdateEngineInterface', path='/org/chromium/UpdateEngine'"
// method call time=1586301148.041487 sender=:1.2974643 -> destination=org.chromium.UpdateEngine serial=7 path=/org/chromium/UpdateEngine; interface=org.chromium.UpdateEngineInterface; member=GetStatusAdvanced

const (
	dbusObject    = "org.chromium.UpdateEngine"
	dbusPath      = "/org/chromium/UpdateEngine"
	dbusInterface = "org.chromium.UpdateEngineInterface"
	dbusMember    = "StatusUpdateAdvanced"
	signalBuffer  = 32 // TODO(bp): What is a reasonable value here?
)

type Client struct {
	conn   *dbus.Conn
	object dbus.BusObject
	ch     chan *dbus.Signal
}

func New() (*Client, error) {
	c := new(Client)
	var err error

	c.conn, err = dbus.SystemBusPrivate()
	if err != nil {
		return nil, err
	}

	methods := []dbus.Auth{dbus.AuthExternal(strconv.Itoa(os.Getuid()))}
	err = c.conn.Auth(methods)
	if err != nil {
		c.conn.Close()
		return nil, err
	}

	err = c.conn.Hello()
	if err != nil {
		c.conn.Close()
		return nil, err
	}

	c.object = c.conn.Object(dbusObject, dbus.ObjectPath(dbusPath))

	// Setup the filter for the StatusUpdate signals
	match := fmt.Sprintf("type='signal',interface='%s',member='%s'", dbusInterface, dbusMember)

	call := c.conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, match)
	if call.Err != nil {
		return nil, call.Err
	}

	c.ch = make(chan *dbus.Signal, signalBuffer)
	c.conn.Signal(c.ch)

	return c, nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// ReceiveStatuses receives signal messages from dbus and sends them as Statues
// on the rcvr channel, until the stop channel is closed. An attempt is made to
// get the initial status and send it on the rcvr channel before receiving
// starts.
func (c *Client) ReceiveStatuses(rcvr chan *StatusResult, stop <-chan struct{}) {
	// if there is an error getting the current status, ignore it and just
	// move onto the main loop.
	st, err := c.GetStatus()
	if err != nil {
		log.Println("Got error: ", err.Error())
	}

	rcvr <- st

	for {
		select {
		case <-stop:
			return
		case signal := <-c.ch:
			if signal == nil {
				return
			}
			rcvr <- NewStatus(signal.Body)
		}
	}
}

func (c *Client) RebootNeededSignal(rcvr chan *StatusResult, stop <-chan struct{}) {
	for {
		select {
		case <-stop:
			return
		case signal := <-c.ch:
			s := NewStatus(signal.Body)
			if s.CurrentOperation == Operation_UPDATED_NEED_REBOOT {
				rcvr <- s
			}
		}
	}
}

// GetStatus gets the current status from update_engine
func (c *Client) GetStatus() (*StatusResult, error) {
	call := c.object.Call(dbusInterface+".GetStatusAdvanced", 0)
	if call.Err != nil {
		return &StatusResult{}, call.Err
	}

	return NewStatus(call.Body), nil
}

// AttemptUpdate will trigger an update if available. This is an asynchronous
// call - it returns immediately.
func (c *Client) AttemptUpdate() error {
	call := c.object.Call(dbusInterface+".AttemptUpdate", 0)
	return call.Err
}
