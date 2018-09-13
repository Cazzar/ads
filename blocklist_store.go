// Copyright 2018 Christian Müller <cmueller.dev@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ads

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

type StoredBlocklistConfiguration struct {
	UpdateTimestamp int      `json:"update_timestamp"`
	Blocklists      []string `json:"blocklists"`
	BlockedNames    BlockMap `json:"blocked_names"`
}

func ReadBlocklistConfiguration(path string) (*StoredBlocklistConfiguration, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	compressedData, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	data, err := gunzip(compressedData)

	var config StoredBlocklistConfiguration

	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (s *StoredBlocklistConfiguration) Persist(path string) error {
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}

	compressed, err := gzip(data)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(compressed)

	return err
}

func (s *StoredBlocklistConfiguration) NeedsUpdate(updateDuration time.Duration) bool {
	log.Info(time.Now(),time.Unix(int64(s.UpdateTimestamp),0), updateDuration,time.Now().After(time.Unix(int64(s.UpdateTimestamp), 0).Add(updateDuration)))
	return time.Now().After(time.Unix(int64(s.UpdateTimestamp), 0).Add(updateDuration))
}
