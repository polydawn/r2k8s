package kuddle

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ugorji/go/codec"
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
		if !ok { // Weird.  K8s will error at lack of image name, but ok.
			continue
		}
		if !strings.HasPrefix(imageName, "./") {
			fmt.Printf("image %q looks like a public name; leaving it\n", imageName)
			continue
		}
		fmt.Printf("image %q looks like a local file; looking for a formula\n", imageName)

		// Ok, we've got the image name, and feel like we're gonna handle it.
		// Let's load the formula.
		frm, err := getFrm(imageName)
		if err != nil {
			fmt.Printf("image %q -- skipping, can't find a formula (%s)\n", imageName, err)
			continue
		}
		// And serialize it.
		frmBuf := bytes.Buffer{}
		if err := codec.NewEncoder(&frmBuf, &codec.JsonHandle{}).Encode(frm); err != nil {
			fmt.Printf("image %q -- skipping, error serializing formula: %s\n", imageName, err)
			continue
		}

		// Start altering.
		containerSpec["image"] = "radd.repeatr.io/radd"
		containerSpec["imagePullPolicy"] = "Never"
		containerSpec["securityContext"] = map[string]interface{}{"privileged": true}
		delete(containerSpec, "workingDir")
		containerSpec["command"] = []string{
			"/bin/bash", "-c",
			"/opt/repeatr/repeatr run -s --ignore-job-exit <(echo \"$FRM\")",
		}
		// Injecting the env is a most complicated part.
		// Note that we haven't countered the substitution system that k8s adds to this.
		// We may be able to sanely get around that by using a EnvVarSource; not yet tested.
		env, ok := containerSpec["env"].([]interface{})
		if !ok {
			env = []interface{}{}
		}
		containerSpec["env"] = append(env, map[string]interface{}{"name": "FRM", "value": frmBuf.String()})
		// TODO if the target environment is using AUFS, we would still need the
		// mounts to escape nest-anything-inside-AUFS problems.
		// Joyously, recent generations of GKE at least are now shipping with overlayfs.
		// TODO those mounts still make sense to let repeatr instances share cache.
		fmt.Printf("image %q -- has now been jibbled\n", imageName)
	}
	return nil
}
