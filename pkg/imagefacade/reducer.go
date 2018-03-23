/*
Copyright (C) 2018 Synopsys, Inc.

Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements. See the NOTICE file
distributed with this work for additional information
regarding copyright ownership. The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied. See the License for the
specific language governing permissions and limitations
under the License.
*/

package imagefacade

import (
	"fmt"
	"reflect"
	"time"

	log "github.com/sirupsen/logrus"
)

type reducer struct{}

func newReducer(initialModel *Model, actions <-chan Action) *reducer {
	stop := time.Now()
	model := initialModel
	go func() {
		for {
			select {
			case nextAction := <-actions:
				// metrics: log message type
				actionString := fmt.Sprintf("%s", reflect.TypeOf(nextAction))
				log.Debug("processing action of type %s", actionString)
				recordActionType(actionString)

				// metrics: how long idling since the last action finished processing?
				start := time.Now()
				recordReducerActivity(false, start.Sub(stop))

				nextAction.apply(model)

				// metrics: how long did the work take?
				stop = time.Now()
				recordReducerActivity(true, stop.Sub(start))
			}
		}
	}()
	return &reducer{}
}
