# Detectors

## CPU Detectors

1) Identify metrics to measure CPU utilization on the host.
2) Identify metrics to measure CPU utilization of a Docker Container.
3) Identify metrics to measure CPU utilization of a container running in Kubernetes.
4) Create custom charts displaying the CPU utilisation of Host, Docker Container, and a container running in Kubernetes. 

Review [Alert Conditions](https://docs.signalfx.com/en/latest/detect-alert/set-up-detectors.html#alert-condition) documentation.

Which condition do you think is best for detecting higher than normal CPU activity?


## Memory Detectors

1) Identify metrics to measure Memory Usage on the host.
2) Identify metrics to measure Memory Usage of a Docker Container.
3) Identify metrics to measure Memory Usage of a container running in Kubernetes.
4) Create custom charts displaying the Memory Usage of Host, Docker Container, and a container running in Kubernetes. 

Review [Alert Conditions](https://docs.signalfx.com/en/latest/detect-alert/set-up-detectors.html#alert-condition) documentation.

Which condition do you think is best for detecting higher than normal Memory Utilization?

When would you want to activate an alert?
When would you clear an alert?


## Error Rate Detectors

1) Identify metrics that measure Errors on for a service.
2) Create a custom chart to display the number of errors for all services.

Which alert condition would you use to alert on an increase in errors?

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
