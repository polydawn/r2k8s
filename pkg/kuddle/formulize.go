package kuddle

import (
	"fmt"
	"strings"
)

/*
	Examines the `podSpec` containers; if any of them refer to formulas,
	the `podSpec` will be altered to run the repeatr/radd container with
	those formulas injected via the environment.

	Specifically:

	  - the container name is unchanged;
	  - the container image *is* changed (it becomes 'radd');
	  - the container command *is* changed (it becomes 'repeatr');
	  - any volumeMounts are currently unchanged (so make sure you pass them through
	    appropriately in your formula; we might add more automatic help here later);
	  - except we will add one volumeMount, which we need to put COW filesystems in
	    (this is due to obscure behavior in AUFS, which we unfortunately presume is still default);
	  - the container security context is set to Privileged=true (but don't worry;
	    the repeatr container allows you to lock it down again).

*/
func formulize(podSpec map[string]interface{}, getFrm FormulaLoader) error {
	fmt.Printf("podspec found.\n")
	containerSpecs := podSpec["containers"].([]interface{})
	for _, v := range containerSpecs {
		containerSpec, ok := v.(map[string]interface{})
		if !ok { // Weird.  K8s will error at this, but ok.
			continue
		}
		imageName, ok := containerSpec["image"].(string)
		if !ok {
			continue
		}
		if !strings.HasPrefix(imageName, "./") {
			fmt.Printf("image %q looks like a public name; leaving it\n", imageName)
			continue
		}
		fmt.Printf("image %q looks like a local file; looking for a formula\n", imageName)
		// TODO
	}
	return nil
}
