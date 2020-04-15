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
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertMergeMaps(t *testing.T, cr, instance, expected map[string]string) {
	merged := MergeMapsWithoutOverwrite(cr, instance)
	assert.Equal(t, expected, merged)
}

func TestMergeMapsWithNoInstanceLabels(t *testing.T) {
	cr := map[string]string{
		"foo": "bar",
	}
	instance := map[string]string{}
	expected := map[string]string{
		"foo": "bar",
	}

	assertMergeMaps(t, cr, instance, expected)
}

func TestMergeMapsWithDifferentInstanceLabels(t *testing.T) {
	cr := map[string]string{
		"foo": "bar",
	}
	instance := map[string]string{
		"more": "options",
	}
	expected := map[string]string{
		"foo":  "bar",
		"more": "options",
	}

	assertMergeMaps(t, cr, instance, expected)
}

func TestMergeMapsWithCustomResourceLabelOverride(t *testing.T) {
	cr := map[string]string{
		"foo": "bar",
	}
	instance := map[string]string{
		"foo":  "bar2",
		"more": "options",
	}
	expected := map[string]string{
		"foo":  "bar",
		"more": "options",
	}

	assertMergeMaps(t, cr, instance, expected)
}
