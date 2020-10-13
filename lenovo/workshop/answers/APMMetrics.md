# APM Monitoring MetricSet

APM produces two types of MetricSets as part of processing APM data: Monitoring MetricSet and Troubleshooting MetricSet.

Documentation: [APM MetricSets](https://docs.signalfx.com/en/latest/apm/apm-concepts/apm-metricsets.html)

Monitoring MetricSets get converted to Metric Times Series that can be used in Dashboards and Detectors.

You can find these metrics in your Org if you search on any of:
- service.request.*, eg `service.request.count` 
- spans.*, eg `spans.count`
- traces.*, eg `traces.count`

`service.request.*` metrics have 4 dimensions: sf_environment, sf_service, sf_error, sf_metric

`spans.*` and `trace.*` metrics have the above 4 dimensions plus: sf_operation, sf_httpMethod, sf_kind

When charting the APM Metrics, you can split and filter by any of the above dimensions. 

`sf_error` dimension has values `true` and `false`, therefore we need to obtain any of `service.request.count`, `traces.count`, or `spans.count`, and Filter on `sf_error:true` to get the error count.
We can divide that number of the unfiltered count to obtain Error % Rate.
