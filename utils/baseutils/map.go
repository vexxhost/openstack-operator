// Copyright 2020 VEXXHOST, Inc.
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

package baseutils

// MergeMaps merges all maps in the list
func MergeMaps(MapList ...map[string]string) map[string]string {
	var baseMap = make(map[string]string)
	for _, imap := range MapList {
		for k, v := range imap {
			baseMap[k] = v
		}
	}
	return baseMap
}

// MergeMapsWithoutOverwrite merges all maps in the list without overwriting. The priority is the same as the sequence of the list.
func MergeMapsWithoutOverwrite(MapList ...map[string]string) map[string]string {
	var baseMap = make(map[string]string)
	for _, imap := range MapList {
		for k, v := range imap {
			if _, ok := baseMap[k]; !ok {
				baseMap[k] = v
			}

		}
	}
	return baseMap
}
