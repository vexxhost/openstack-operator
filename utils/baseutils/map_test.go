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
