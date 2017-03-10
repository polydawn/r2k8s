package pushtask

import (
	kapi "k8s.io/client-go/pkg/api"
)

func templatePod() *kapi.Pod {
	// Assign true to a value so we can take an address of it.
	// Pointer-to-bool is used by the k8s api to implement trinary logic.
	tru := true
	// Stamp out a pod object.
	p := &kapi.Pod{
		Spec: kapi.PodSpec{
			Containers: []kapi.Container{
				{
					Name:       "main",
					Image:      "radd.repeatr.io/radd",
					WorkingDir: "/task",
					Command: []string{
						"/bin/bash", "-c",
						"/opt/repeatr/repeatr run -s --ignore-job-exit <(echo \"$FRM\")",
					},
					Env: []kapi.EnvVar{
						{
							Name: "FRM", Value: string("todo"),
						},
					},
					SecurityContext: &kapi.SecurityContext{Privileged: &tru},
					ImagePullPolicy: "Never",
				},
			},
			RestartPolicy: kapi.RestartPolicyNever,
		},
	}
	// Some fields are easier to set by walking the object like this afterwards,
	// rather than using the struct initializers, due to some vagueries of go library
	// imports.  (Specifically, traversing these fields rather than using struct
	// initializers lets us avoid importing another package, which avoids a whole
	// argument with re-importing vendored packages, which... Ufdah.
	// This seems like a silly thing to have manifest impacts in code, but, well,
	// cest la vie.)
	p.ObjectMeta.GenerateName = "r2k8s-"
	return p
}
