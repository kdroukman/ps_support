# Detectors

## CPU Detectors

1) Identify metrics to measure CPU utilization on the host. [View Answer](https://github.com/kdroukman/ps_support/blob/master/lenovo/workshop/answers/ContCPUMetric.md)
2) Identify metrics to measure CPU utilization of a Docker Container. [View Answer](https://github.com/kdroukman/ps_support/blob/master/lenovo/workshop/answers/DockerCPUMetric.md)
3) Identify metrics to measure CPU utilization of a container running in Kubernetes. [View Answer](https://github.com/kdroukman/ps_support/blob/master/lenovo/workshop/answers/k8sCPUMetric.md)
4) Create custom charts displaying the CPU utilisation of Host, Docker Container, and a container running in Kubernetes. [Explore the Charts Here](https://app.us1.signalfx.com/#/dashboard/Ef5FEgsA0cw?groupId=Ef5FEgsA0cs&configId=Ef5FEgsA0c0)

Review [Alert Conditions](https://docs.signalfx.com/en/latest/detect-alert/set-up-detectors.html#alert-condition) documentation.

Which condition do you think is best for detecting higher than normal CPU activity? [View Answer and Discussion](https://github.com/kdroukman/ps_support/blob/master/lenovo/workshop/answers/CPUDetector.md)


## Memory Detectors

1) Identify metrics to measure Memory Usage on the host. [View Answer](https://github.com/kdroukman/ps_support/blob/master/lenovo/workshop/answers/MemCPUMetric.md)
2) Identify metrics to measure Memory Usage of a Docker Container. [View Answer](https://github.com/kdroukman/ps_support/blob/master/lenovo/workshop/answers/DockerMemMetric.md)
3) Identify metrics to measure Memory Usage of a container running in Kubernetes. [View Answer](https://github.com/kdroukman/ps_support/blob/master/lenovo/workshop/answers/K8sMemMetric.md)
4) Create custom charts displaying the Memory Usage of Host, Docker Container, and a container running in Kubernetes. [Explore the Charts Here](https://app.us1.signalfx.com/#/dashboard/Ef5FEgsA0cw?groupId=Ef5FEgsA0cs&configId=Ef5FEgsA0c0)

Review [Alert Conditions](https://docs.signalfx.com/en/latest/detect-alert/set-up-detectors.html#alert-condition) documentation.

Which condition do you think is best for detecting higher than normal Memory Utilization?

When would you want to activate an alert?
When would you clear an alert?
[View Answer and Discussion](https://github.com/kdroukman/ps_support/blob/master/lenovo/workshop/answers/MemoryDetector.md)

## Error Rate Detectors

1) Identify metrics that measure Errors on for a service. [View Answer](https://github.com/kdroukman/ps_support/blob/master/lenovo/workshop/answers/APMMetrics.md)
2) Create a custom chart to display the number of errors for all services.[View Answer](https://github.com/kdroukman/ps_support/blob/master/lenovo/workshop/answers/APMDashboard.md)

Which alert condition would you use to alert on an increase in errors?
[View Answer and Discussion](https://github.com/kdroukman/ps_support/blob/master/lenovo/workshop/answers/APMDetector.md)

## Detector Tips and Tricks

### How do I set up a detector for ALL hosts, services, operations, etc?

When you select a metric, such as `cpu.utilization` or `spans.count`, that metric reports Metric Time Series from all entites. The number of metric time series can be seen as `# ts` in your signal definition. 

In an APM Detector type, you can use `*:*` wildcard to report on all services and all operations.

### How do I set up a detector for SOME hosts, services, operations, etc?

When you select a metric, such as `cpu.utilization` or `spans.count`, that metric reports Metric Time Series from all entites. You can then use the **Filter** to report only on entities that match your list of dimensions. You can also use wildcard `*` to do a partial match. 

Eg: `kubernetes_workload_name:*agent*` will report on any container that matches agent in its workload name. 

### How do I search for a detector by name or dimension/property value?

The Detectors and Alerts view in SignalFx UI container Filter, Group By, and Search options. To search by name, use Search text box. To search by dimension value use the Filter.

### How do I know which service or host the alert has fired against?

When you recieve and view an alert, all the values that are being captured by that alert are in the text line starting with "Signal details:" .
For metrics that have host or service name dimensions, these will be reported.

Generally, Signal details lists all the dimension values. To preview which dimension values will be reported - you can view the **Data Table** in your Alert Signal definition. 

### What are Alert Rules?

Alert Rules allow you to create several versions of an alert on the **same metric signal definition**. Eg: Alert condition, Alert message, Alert notification. 
All allert rules will use the same reporting metric, if you try to modify metric per alert rule - you will see unexpected behaviour. 

### My alert doesn't fire as expected.

Refer to the follow [Troubleshooting documentation](https://docs.signalfx.com/en/latest/detect-alert/when-detector-not-triggering.html) to look into possible causes of detector not firing as expected. Usually it's due to sensitivity settings on the alert detector not accomodating for data behaviour (sparse data, jittery data, etc).

### What should I detect on for Black Friday?

The following are recommended detectors to setup prior to black friday. These will also be important at any time of the year to report on system and application health. 

1) Sudden Spike in load for all APM services/endpoints.
2) Sudden Spike in Error Rate for all APM services/endpoints.
3) Increase in Latency (Service Performance) for all APM services/endpoints.
4) CPU, Memory, Disk, and Network Usage and Utilization on all hosts and containers.
5) Disk capacity on all hosts using "Resource Running Out" condition. 
6) JVM Memory Utilization ang garbage collection spike for Java-based services, as well as Thread count spike.
7) HAProxy spikes in Backend failures, latency and Front end connection rate for all proxy names (services)
8) Layer 7 spikes in unsuccessful requests and latency for all services.

If you have identified service in a critical path, you may add Critical level detectors for 
1) Sudden Spike in load for APM services in the critical path.
2) Sudden Spike in Error Rate for APM services in the critical path.
3) Increase in Latency (Service Performance) for APM services in the critical path.

While setting all other service alerts to a lower level of severity, eg: High or Medium.

When an APM alert is fired, navigate to **Troubleshoot** link to take straight to the filtered view where you can inspect the issue further.

When a non-APM alert is fired, navigate to the corresponding Dashboard(s) to view where issues have spikes, and cross-reference with Splunk Logs. 

## Metrics
1) Click on the **Metrics** menu in the top navigation bar and inspect the metrics. Search for a metric. Inspect returned results.
2) Navigate to built-in Dashboards. Open a chart of interest and inspect which metrics are used.

## Terraform

ssh into your dedicated EC2 host:
| Name | IP | Command | Password | 
|--|--|--|--|
|vchoo | 3.26.10.62 | ssh ubuntu@3.26.10.62 | Observability2020! |
|bleong | 13.210.133.110 | ssh ubuntu@13.210.133.110 | Observability2020! |
|waijunwo | 3.25.229.14 | ssh ubuntu@3.25.229.14 | Observability2020! |
|thongsim | 3.106.125.182 | ssh ubuntu@3.106.125.182 | Observability2020! |
|kcheng6 | 3.26.40.8 | ssh ubuntu@3.26.40.8 | Observability2020! |
|soomkc | 3.26.30.13 | ssh ubuntu@3.26.30.13 | Observability2020! |
| wkong|3.26.14.0| ssh ubuntu@3.26.14.0 | Observability2020!|
| waisl|3.25.255.196 | ssh ubuntu@3.25.255.196 | Observability2020!| 
| tlim| 3.26.32.32| ssh ubuntu@3.26.32.32 | Observability2020!| 
| gmalayao| 54.206.100.244| ssh ubuntu@54.206.100.244| Observability2020!| 
| htan4| 3.26.42.238| ssh ubuntu@3.26.42.238|Observability2020! | 
| zhaomin10| 3.26.1.73| ssh ubuntu@3.26.1.73|Observability2020! | 
| cteoh1| 3.26.39.163| ssh ubuntu@3.26.39.163| Observability2020!|

1) Navigate to signalfx-jumpstart directory. 
2) Open main.tf and inspect its contents.
3) Open variables.tf and inspects contents.
4) Navigate to signalfx-jumpstart/modules/host directory.
5) Open cpu.tf and inspect it's contents.

Read the program text line. Describe what it does?

6) Navigate back to signalfx-jumpstart directory: `cd ~/signalfx-jumpstart`
7) Run - `terraform init -upgrade` 
8) Set the environment variables:
```
export ACCESS_TOKEN=ACCESS TOKEN, from organisation page
export REALM=REALM e.g. us1
```
9) Run Terraform plan - `terraform plan -var="access_token=$ACCESS_TOKEN" -var="realm=$REALM" -var="sfx_prefix=[$(hostname)]"`
10) Apply the execution plan = `terraform apply -var="access_token=$ACCESS_TOKEN" -var="realm=$REALM" -var="sfx_prefix=[$(hostname)]"`

Check you hostname with `echo $(hostname)` and locate your assets in the SignalFx UI.

11) Run destroy command - `terraform destroy -var="access_token=$ACCESS_TOKEN" -var="realm=$REALM"`

What happens after you run this this command?


## Further Reading

- [What is Time Series Data?](https://blog.timescale.com/blog/what-the-heck-is-time-series-data-and-why-do-i-need-a-time-series-database-dcf3b1b18563/)
- [SignalFx Data Model](https://docs.signalfx.com/en/latest/getting-started/concepts/data-model.html)
- [Detailed Overview of Sudden Change Condition](https://docs.signalfx.com/en/latest/detect-alert/alert-condition-reference/sudden-change.html)
- [Detailed Overview of Static Threshold Condition](https://docs.signalfx.com/en/latest/detect-alert/alert-condition-reference/static-threshold.html)
- [Detailed Overview of Historical Anomaly Condition](https://docs.signalfx.com/en/latest/detect-alert/alert-condition-reference/hist-anomaly.html)
- [Detailed Overview of Resource Running Out Condition](https://docs.signalfx.com/en/latest/detect-alert/alert-condition-reference/resource-running-out.html)
- [Detailed Overview of Outlier Detection Condition](https://docs.signalfx.com/en/latest/detect-alert/alert-condition-reference/outlier-detection.html)
- [SignalFx Detector Overview](https://docs.signalfx.com/en/latest/detect-alert/index.html#detect-alert)
