package kuddle

import (
	"testing"
)

func TestLoadFile(t *testing.T) {
	InterpolateFile(
		"./../../fixtures/example.yaml",
		nil, /* TODO */
		"./../../fixtures/example.r2k",
	)
}
