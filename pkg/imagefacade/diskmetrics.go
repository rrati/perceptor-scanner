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
	"syscall"

	log "github.com/sirupsen/logrus"
)

type DiskMetrics struct {
	FreeBytes      uint64
	AvailableBytes uint64
	TotalBytes     uint64
	UsedBytes      uint64
}

func getDiskMetrics() (*DiskMetrics, error) {
	var stat syscall.Statfs_t
	err := syscall.Statfs("/var/images", &stat)
	if err != nil {
		log.Errorf("unable to get disk stats: %s", err.Error())
		return nil, err
	}
	metrics := &DiskMetrics{
		FreeBytes:      stat.Bfree * uint64(stat.Bsize),
		AvailableBytes: stat.Bavail * uint64(stat.Bsize),
		TotalBytes:     stat.Blocks * uint64(stat.Bsize),
	}
	metrics.UsedBytes = metrics.TotalBytes - metrics.FreeBytes
	return metrics, nil
}
