package kuddle

import (
	"testing"

	rdef "go.polydawn.net/repeatr/api/def"
)

func TestLoadFile(t *testing.T) {
	InterpolateFile(
		"./../../fixtures/example.yaml",
		FormulaLoaderConstant(&rdef.Formula{}),
		"./../../fixtures/example.r2k",
	)
}
