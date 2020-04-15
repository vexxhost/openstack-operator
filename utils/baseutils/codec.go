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

import (
	"encoding/base64"
)

// Base64DecodeByte2Str returns plain text as string from the encrypted text as byte array
func Base64DecodeByte2Str(enc []byte) string {
	encStr := string(enc)
	decStr, err := base64.StdEncoding.DecodeString(encStr)
	if err != nil {
		return ""
	}
	return string(decStr)
}
