# APM Metrics Dashboard

Here we use `service.request.count` metric to calculate and view error rate for All Services:
![Service Error Rate](https://github.com/kdroukman/ps_support/blob/master/lenovo/workshop/img/APMErrMetric.png?raw=true)

Similarly we can use `spans.count` to arrive at the same calculation. Note that in spans there are 2 extra Time Series for this total count. This is because spans metric further splits Metric Time Series by Operation name, where service.requests metric doesn't.
![Spans Error Rate](https://github.com/kdroukman/ps_support/blob/master/lenovo/workshop/img/APMErrorRate-Spans.png?raw=true)

_Note: screenshots taken at different time windows, hence the lines are different_
