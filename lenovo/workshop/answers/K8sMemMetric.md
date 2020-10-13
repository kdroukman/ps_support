# container_memory_usage_bytes & kubernetes.container_memory_limit

[<<BACK TO MAIN](https://github.com/kdroukman/ps_support/blob/master/lenovo/workshop/README.md)

`container_memory_usage_bytes` and `kubernetes.container_memory_limit` are collected by SignalFx Smart Agent through thre different monitor types `kubernetes-cluster`,`kubelet-stats` and `kubelet-metrics`.

The SmartAgent you deployed in your Kubernetes environment via pre-configured agent.yaml will use `kubelet-stats` monitor for versions less than 5.3.0, and `kubelet-metrics` monitor otherwise. 
`kubernetes-cluster` monitor is also preconfigured to be enabled.

Documentation: [Kubernetes Memory Limit](https://docs.signalfx.com/en/latest/integrations/agent/monitors/kubernetes-cluster.html)

Documentation: [Kubernetes Memory Usage](https://docs.signalfx.com/en/latest/integrations/agent/monitors/kubelet-metrics.html)

In order to get a Memory Utilization as a percentage, a formula must be applied. A percentage is a useful measure to detect against a consistent Threshold.

**Note:** If a memory limit is not set on a container, the result of the formula will be null, and therefore you cannot chart utilization as a percentage, and need to rely on Sudden Change detector condition to monitor increases in memory. 

The following charts demonstrates `container_memory_usage_bytes` and `kubernetes.container_memory_limit` emitted by contianers running in Kubernetes. A formula is used to calculate the resulting utilization %.

![Kubernetes Memory Utilization Chart](https://github.com/kdroukman/ps_support/blob/master/lenovo/workshop/img/K8sMemMetric.png?raw=true)

_Click on the image to enlarge in a new tab_
