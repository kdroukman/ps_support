# Host Memory Detector

[<<BACK TO MAIN](https://github.com/kdroukman/ps_support/blob/master/lenovo/workshop/README.md)

Explore an example detector [HERE](https://app.us1.signalfx.com/#/detector/v2/EkNAeJaA0AA/edit)

`memory.utilization` metric has been provided to the detector. The detector uses 5 Metric Time Series corresponding to 5 hosts with SignalFx Smart Agent installed to check if any of the hosts meet the Alert Conditions.

Bellow explore how different alert conditions affect the firing and clearing of the alert.

## Threshold Condition

When an alert is setup to fire immediately as soon as memory reaches 60% utilization we get 27 alerts in 12 hours:
!(Memory Immediate Threshold Detector)[https://github.com/kdroukman/ps_support/blob/master/lenovo/workshop/img/HostMemDetector-Threshold-Immediate.png?raw=true]

When an alert is setup to fire when memory utilization is at or above 60% for 80% of the time in a 5 minute period, we get 1 alert in 12 hours. This removes occastional alert noise:
![Memory Percentage Threshold Detector)[https://github.com/kdroukman/ps_support/blob/master/lenovo/workshop/img/HostMemDetector-Threshold-Percentage.png?raw=true]

## Sudden Change Condition

We set up a sudden change alert condition to look at a 10 minute time window and compare that to the previous 3 hour time window. If the memory utilization in the last 10 minutes has increased by at least 30% above previous 3 hours we fire an alert. In the example bellow, we have 0 alerts, as the memory didn't incrase above 30% in a 10 minute period. We had a few short increases in a 5 minute window and then it settled.
![Memory Sudden Change](https://github.com/kdroukman/ps_support/blob/master/lenovo/workshop/img/HostMemDetector-SuddenChange-Percent.png?raw=true)

You can see that once the time window to analyze is reduced to 1 minute for the same condition, the short bursts in memory usage are captured by alerts.
![Memory Sudden Change - 1 minute](https://github.com/kdroukman/ps_support/blob/master/lenovo/workshop/img/MemoryDetector-SuddenChange-1m.png?raw=true)
