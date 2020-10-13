# Host CPU Detector

[<<BACK TO MAIN](https://github.com/kdroukman/ps_support/blob/master/lenovo/workshop/README.md)

Explore the example Host CPU Detector [HERE](https://app.us1.signalfx.com/#/detector/v2/EkM9g-bA4AA/edit).

Explore Alert Signal. Notice that the Alert Signal uses the same `cpu.utilization` metric from the charts.

A CPU metric time series is generally very jittery if you observe it in a 1-minute window at the highest resolution. However there is a normal trend and an abnormal trend.

Usually, we are conserned if CPU Utilization % goes above a certain Threshold - eg. 80%

However, if you have a signal reporting a CPU value every 10 seconds, and it jumps between 75% and 85%, you will get a lot of alerts indicating the same thing. This is where we can extend the time window to look over and quieten down our alerting.

## Threshold Conditions
-> Notice the number of alerts for two different options:

With an immediate threshold we have 7 alerts in 12 hours.
![Immediate Threshold](https://github.com/kdroukman/ps_support/blob/master/lenovo/workshop/img/HostCPUDetector-Threshold-Immediate.png?raw=true)

With adding duration, we check if the high CPU Utilization persists for 1 minute. This reduces our number of alerts to 4 in 12 hours.
![Duration Threshold](https://github.com/kdroukman/ps_support/blob/master/lenovo/workshop/img/HostCPUDetector-Threshold-Duration2m.png?raw=true)

You may also have situations where 80% CPU Utilization is normal for some hosts and services. This where you can use Sudden Change to see if there is a spike in utilization compared to a previous period.

## Sudden Change Conditions
-> Notice the number of alerts for three different Sudden Change optoins:

Sudden Change using Mean plus Standard Deviation - 4 alerts in 12 hours:
![Mean+StdDev Sudden Change](https://github.com/kdroukman/ps_support/blob/master/lenovo/workshop/img/HostCPUDetector%20-%20SuddenChange-MeanStdDev.png?raw=true)
_You can read more about Standard Deviation [here](https://www.mathsisfun.com/data/standard-deviation.html)_

Sudden Change using Percentile - 3 alerts in 12 hours:
![Percentile Sudden Change](https://github.com/kdroukman/ps_support/blob/master/lenovo/workshop/img/HostCPUDetector-SuddenChange-Percentile.png?raw=true)
_Percentiles are less sensitive to extreme values than mean-based calculation._

Sudden Change using Mean plus Percentage - 4 alerts in 12 hours. Notice when the alerts are cleared:
![Mean+Percent Sudden Change](https://github.com/kdroukman/ps_support/blob/master/lenovo/workshop/img/HostCPUDetector-SuddenChange-MeanPerct.png?raw=true)
