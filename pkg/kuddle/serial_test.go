package kuddle

import (
	"testing"

	rdef "go.polydawn.net/repeatr/api/def"
)

func TestLoadFile(t *testing.T) {
	frm := &rdef.Formula{
		Inputs: rdef.InputGroup{
			"rootfs": &rdef.Input{
				Type:       "tar",
				Hash:       "aLMH4qK1EdlPDavdhErOs0BPxqO0i6lUaeRE4DuUmnNMxhHtF56gkoeSulvwWNqT",
				Warehouses: rdef.WarehouseCoords{"http+ca://repeatr.s3.amazonaws.com/assets/"},
				MountPath:  "/",
			},
		},
		Action: rdef.Action{
			Entrypoint: []string{"bash", "-c", "echo hello && sleep 5000 && echo yes"},
		},
	}
	InterpolateFile(
		"./../../fixtures/example.yaml",
		FormulaLoaderConstant(frm),
		"./../../fixtures/example.r2k",
	)
}
