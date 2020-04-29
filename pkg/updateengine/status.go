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

func NewStatus(body []interface{}) *StatusResult {
	s := &StatusResult{}

	err := proto.Unmarshal(body[0].([]byte), s)
	if err != nil {
		log.Println("Error unmarshalling message: ", err.Error())
		return s
	}

	return s
}
