package kuddle

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
	kyaml "k8s.io/apimachinery/pkg/util/yaml"
)

/*
	Using the full k8s API to deserialize their typed, versioned objects is...
	somewhat of an interesting errand, and one this code doesn't pursue.

	One of the nearest examples I see to a clear entrypoint to the
	factory/registry/codec/schema shenanigans is this snippet from the k8s main repo:

	```
		switch obj.(type) {
		case runtime.Unstructured, *runtime.Unknown:
			actualObj, err = runtime.Decode(
				api.Codecs.UniversalDecoder(),
				[]byte(/.../)))
		default:
			actualObj = obj
			err = nil
		}
	```

	However, the overall impression of this is:

	  - very complicated
	  - still results in a huge swath of case-switches for the real types
	  - using it is likely to result in a kind of tight coupling which will be
	    unpleasant for our own maintainability;
	  - and does not cause our code to Do The Right Thing for types it doesn't
	    recognize, which doesn't play well with extensions, nor ease upgrade cycles.

	Instead, we proceed from some observations:

	  - PodSpecs are the single most solid and consistent part of the k8s API over time;
	  - PodSpecs are often embedded in other API types;
	  - but the way *we* are interested in PodSpecs never changes,
	    no matter what they're embedded in.

	As a result, it makes sense for us to
	deserialize objects nearly schema-free,
	detect PodSpecs by their (quite clear) structure,
	alter our fields of interest patch-wise (leaving the rest untouched),
	and emit the result without further processing.
*/

func Interpolate(k8sDocuments []byte, getFrm FormulaLoader) (result []byte, err error) {
	decoder := kyaml.NewYAMLOrJSONDecoder(bytes.NewBuffer(k8sDocuments), 2<<6)
	resultBuf := bytes.Buffer{}
	for err == nil {
		var slot interface{}
		err = decoder.Decode(&slot)
		if err == io.EOF {
			return resultBuf.Bytes(), nil
		} else if err != nil {
			return resultBuf.Bytes(), err
		}
		err = interpolateObj(slot, getFrm)
		if err != nil {
			return resultBuf.Bytes(), err
		}
		bs, err := yaml.Marshal(&slot)
		resultBuf.WriteString("---\n")
		resultBuf.Write(bs)
		if err != nil {
			return resultBuf.Bytes(), err
		}
	}
	panic("unreachable")
}

func InterpolateFile(k8sDocumentPath string, getFrm FormulaLoader, writePath string) error {
	f, err := os.Open(k8sDocumentPath)
	if err != nil {
		return err
	}
	defer f.Close()
	bs, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	bs, err = Interpolate(bs, getFrm)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(writePath, bs, 0644)
	return err
}

func interpolateObj(obj interface{}, getFrm FormulaLoader) error {
	switch obj2 := obj.(type) {
	case map[string]interface{}:
		spec, ok := obj2["spec"]
		if ok { // Might be a PodSpec!
			specMap, ok := spec.(map[string]interface{})
			if !ok {
				goto notAPod
			}
			containers, ok := specMap["containers"]
			if !ok {
				goto notAPod
			}
			_, ok = containers.([]interface{})
			if !ok {
				goto notAPod
			}
			// Looks like a PodSpec!
			formulize(specMap, getFrm)
		}
	notAPod:
		// If this object didn't contain any pod spec, recurse; its children might.
		for _, v := range obj2 {
			if err := interpolateObj(v, getFrm); err != nil {
				return err
			}
		}
		return nil
	case []interface{}:
		// Always recurse.
		for _, v := range obj2 {
			if err := interpolateObj(v, getFrm); err != nil {
				return err
			}
		}
		return nil
	case interface{}:
		// All leaf types like string and float64 end up here.
		// None of which can contain a PodSpec, so, ignore.
		return nil
	default:
		panic(fmt.Errorf("unhandled type %T", obj))
	}
}
