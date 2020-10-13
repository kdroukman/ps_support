# container_cpu_utilization

[<<BACK TO MAIN](https://github.com/kdroukman/ps_support/blob/master/lenovo/workshop/README.md)

`container_cpu_utilization` is collected by SignalFx Smart Agent through three different monitor types `cadvisor`, `kubelet-stats` and `kubelet-metrics`.

The SmartAgent you deployed in your Kubernetes environment via pre-configured agent.yaml will use `kubelet-stats` monitor for versions less than 5.3.0, and `kubelet-metrics` monitor otherwise.

Documentation: [Kubernetes Container CPU Metrics](https://docs.signalfx.com/en/latest/integrations/integrations-reference/integrations.kubernetes.html#container-cpu-utilization)

Kubernetes monitor already publishes this as a percentage, and no further formulas are required.

The following charts demonstrates `container_cpu_utilization` emitted by numerous contianers running in Kubernetes.

![Kubernetes CPU Utilization Chart](https://github.com/kdroukman/ps_support/blob/master/lenovo/workshop/img/K8sCPUMetric.png?raw=true)

_Click on the image to enlarge in a new tab_
