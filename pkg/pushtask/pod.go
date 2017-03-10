package pushtask

import (
	kmeta "k8s.io/apimachinery/pkg/apis/meta/v1"
	kapi "k8s.io/client-go/pkg/api"
)

func templatePod() *kapi.Pod {
	tru := true // just to be grabbable as ref
	p := &kapi.Pod{
		Spec: kapi.PodSpec{
			Containers: []kapi.Container{
				ObjectMeta: kmeta.ObjectMeta{
					GenerateName: "raceway-",
				},
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
	return p
}
