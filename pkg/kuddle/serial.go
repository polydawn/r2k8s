package kuddle

import (
	"k8s.io/apimachinery/pkg/util/yaml"
)

func loadFile() {
	_ = yaml.NewYAMLOrJSONDecoder(nil, 0).Decode
}
