# r2k8s

```
r2k8s reproducible_reliable.yaml > ready_to_kubectl_apply.yaml
```

r2k8s is a standalone tool for hoisting
[Repeatr](https://github.com/polydawn/repeatr/) formulas to execute in a
[Kubernetes](kubernetes.io/) cluster.

Key reasons to do this include:

- *Truthful, solid versioning*: Repeatr formulas are a reliable way to version assets.  Highly cacheable; zero invalidation problems.
- *Decentralized, offline-first operations*: Using Repeatr means decentralized image and data storage is already the modus operandi.  Deploy anywhere.  Even offline.
- *Efficient images*: Repeatr composes segments of filesystems more efficiently than other container systems, allowing mix-n-match filesystems that contain exactly what you need, and update independently.

r2k8s operates by templating.
This is the easiest way to interface with the huge k8s ecosystem, and plays well with other tools.
Write a Repeatr formulas as normal (perhaps using [Reppl](https://github.com/polydawn/reppl/ to manage pipelines and automated updates).
When writing your k8s config (whether simple PodSpec, Job, Deployment -- whatever), instead of using a URL to a container registry, refer to the Repeatr formula file instead.
r2k8s will combine them, letting you run Repeatr formulas in any k8s cluster -- no additional services, and *no container registry* necessary.

r2k8s does not communicate with kubernetes directly.  r2k8s produces config files.
Feed them into `kubectl apply`, or whatever tools you prefer to manage your kubernetes process flow.
Similarly, standard kubernetes monitoring, service restarting -- everything is the same.
r2k8s *just* frees you from centralized container registries.  No more, no less.
