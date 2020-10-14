# APM Error Detector

[<<BACK TO MAIN](https://github.com/kdroukman/ps_support/blob/master/lenovo/workshop/README.md)

View example APM Error Rate detector [HERE](https://app.us1.signalfx.com/#/detector/v2/EkNH-eOA4AA/edit).

Note that while you can use the available APM Metrics to create Infrastructure or Customer Metrics Alert type, there is a simplified APM Detector type available that automatically selects the right metrics for you.
![APM Detector Type](https://github.com/kdroukman/ps_support/blob/master/lenovo/workshop/img/APMDetector-AlertRule.png?raw=true)

The times you may require to use an Infrastructure Detector type instead:
- A single detector for all APM environments. APM Detector types requires that you select an environment.
- An alert condition other than Threshold or Sudden Change. APM Detector type only provides two alert conditions to choose from. For example, you may want to use Historical Anomaly for latency detector.
- A detector managed by Terraform. Terraform creates a V2 type of detectors that require you to provide a SignalFlow program.

Our example detector was created using APM Detector Type. 

When we define Alert Signal, we must choose and environment and signal type. For Error monitoring, we choose Error Rate signal type.
Note that we also have to choose a service and and endpoint. 

To detect on All Services and All Endpoints you can supply a Wildcard: `*:*`.  
![Alert Signal](https://github.com/kdroukman/ps_support/blob/master/lenovo/workshop/img/APMDetectorSignal.png?raw=true)

## Threshold Condition

The threshold condition in APM Detector type is based on a time window. The general best practice is to check if an error rate is consistent across a particular period where there is sufficient load on your service - indicating an clear problem.

Here we check that the error rate is above 25% over a 5 minute period, where we have at least 10 requests in that period. We would get 2 alerts based on some error spikes in our data.
![APM Threshold Alert](https://github.com/kdroukman/ps_support/blob/master/lenovo/workshop/img/APMErrDetector-Threshold.png?raw=true)

## Sudden Change Condition

This condition would be used when you generally capture errors such as HTTP Status 404, but don't need to be alerted on them. Therefore a service will have a consistent ongoing percentage of errors, 
and you want to be alerted when those errors are higher than usual. 

The following Sudden Change condition detects a 30% increase in errors in the current time window of 5 minutes, compared to previous 1 hour. Note that the condition here is met later in the timeline compared to previous alert due as we measure a 30% increase from previous period. We still get alerted on 2 services with an error. 
![APM Sudden Change](https://github.com/kdroukman/ps_support/blob/master/lenovo/workshop/img/APMErrDetector-SuddenChange.png?raw=true)
