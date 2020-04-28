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
	"log"

	"google.golang.org/protobuf/proto"
)

// The possible update statuses returned from the update engine
//
// These correspond to current operation values exposed over DBus and defined by `update_engine`:
// https://github.com/coreos/update_engine/blob/v0.4.3/src/update_engine/update_attempter.h#L34-L43
const (
	UpdateStatusIdle                = "UPDATE_STATUS_IDLE"
	UpdateStatusCheckingForUpdate   = "UPDATE_STATUS_CHECKING_FOR_UPDATE"
	UpdateStatusUpdateAvailable     = "UPDATE_STATUS_UPDATE_AVAILABLE"
	UpdateStatusDownloading         = "UPDATE_STATUS_DOWNLOADING"
	UpdateStatusVerifying           = "UPDATE_STATUS_VERIFYING"
	UpdateStatusFinalizing          = "UPDATE_STATUS_FINALIZING"
	UpdateStatusUpdatedNeedReboot   = "UPDATE_STATUS_UPDATED_NEED_REBOOT"
	UpdateStatusReportingErrorEvent = "UPDATE_STATUS_REPORTING_ERROR_EVENT"
)

/*
type Status struct {
	LastCheckedTime  int64
	Progress         float64
	CurrentOperation string
	NewVersion       string
	NewSize          int64
}
*/

func NewStatus(body []interface{}) *StatusResult {
	s := &StatusResult{}

	err := proto.Unmarshal(body[0].([]byte), s)
	if err != nil {
		log.Println("Error unmarshalling message: ", err.Error())
		return s
	}

	return s
}
