# r2k8s

```
r2k8s reproducible_reliable.yaml > ready_to_kubectl_apply.yaml
```

r2k8s is a standalone tool for hoisting
[Repeatr](https://github.com/polydawn/repeatr/) formulas to execute in a
[Kubernetes](kubernetes.io/) cluster.

Coming from vanilla Kubernetes, this offers several advantages:

- *Truthful, solid versioning*: Repeatr formulas are a reliable way to version assets.  Highly cacheable; zero invalidation problems.
- *Decentralized, offline-first operations*: Using Repeatr means decentralized image and data storage is already the modus operandi.  Deploy anywhere.  Even offline.
- *Efficient images*: Repeatr composes segments of filesystems more efficiently than other container systems, allowing mix-n-match filesystems that contain exactly what you need, and update independently.

Coming from Repeatr:

- Everything is the same
- Plus now you can do it all on a cluster!  Farm your builds trivially...
- Or ship your final products to run as services, k8s-style!

r2k8s operates by templating.
This is the easiest way to interface with the huge k8s ecosystem, and plays well with other tools.

The steps are simple:

1. Write a Repeatr formulas as normal (perhaps using [Reppl](https://github.com/polydawn/reppl/) to manage pipelines and automated updates).
2. When writing your k8s config (whether simple PodSpec, Job, Deployment -- whatever), instead of using a URL to a container registry, refer to the Repeatr formula file instead.
3. r2k8s will combine them, letting you run Repeatr formulas in any k8s cluster -- no additional services, and *no container registry* necessary.

r2k8s runs in two modes.  You choose which suits your needs, depending on whether you want to run a long-lived service, or run a formula that produces results.

For services
------------

In the "service" path, r2k8s generates a k8s config file.

In this mode, r2k8s does not communicate with kubernetes directly -- instead, reuse all the standard kubernetes interfaces.  The config files produced are ready to feed directly `kubectl apply`, or whatever tools you prefer to manage your kubernetes process flow.
Standard kubernetes monitoring, lifecycles, service config: everything is business as usual.  Only the image transport changes.

This is great for freeing you from centralized container registries, and connecting directly to Repeatr-ized build processes and image pipelines.

For builds
----------

In the "build" path, r2k8s acts exactly like `repeatr run`, while farming out the real heavy lifting to a kubernetes cluster.

In this mode, r2k8s communicates with the cluster in the background: spawning a pod to handle the work, monitoring it directly.
Autorestarts and other overly-"helpful" k8s features are turned off, so in the event r2k8s crashes or you simply decide to kill it, the job in the cluster will evaporate when done.
No tempfiles and template management are necessary; just give r2k8s the config to contact your cluster.
