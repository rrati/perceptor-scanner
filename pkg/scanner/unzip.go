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

package scanner

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

// This code is from https://golangcode.com/unzip-files-in-go/ and
// https://stackoverflow.com/a/24792688/894284
func unzip(source string, destination string) error {
	r, err := zip.OpenReader(source)
	if err != nil {
		log.Errorf("unable to open reader: %s", err.Error())
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			log.Errorf("unable to open file: %s", err.Error())
			return err
		}
		defer rc.Close()

		fpath := filepath.Join(destination, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			var fdir string
			if lastIndex := strings.LastIndex(fpath, string(os.PathSeparator)); lastIndex > -1 {
				fdir = fpath[:lastIndex]
			}

			err = os.MkdirAll(fdir, os.ModePerm)
			if err != nil {
				log.Errorf("unable to make directory: %s", err.Error())
				return err
			}
			f, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				log.Errorf("unable to open file: %s", err.Error())
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				log.Errorf("unable to copy file: %s", err.Error())
				return err
			}
		}
	}
	return nil
}
