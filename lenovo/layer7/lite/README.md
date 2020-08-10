## Setup APM Metrics Forwarder (Lite Version)

### Pre-requsites
You must have a SignalFx Smart Agent installed on your hosts. APM Metrics Forwarders provided here will be using `signalfx-agent` user.

## Setting up the Metrics Forwarder on Layer 7 API Gateway version 9.3
_Disclaimer: While a metrics forwarder is provided for version 9.3, this is not the recommended or supported method. It is assumed that Lenovo will make the necessary upgrades to utilize the recommended and supported method. In version 9.3 we incercept APM metrics published by the Layer 7 vendor APM solution provided alongside Layer 7 API Gateway._

**Pre-requisites**
1) APM Metrics Forwarder for Layer 7 API Gateway version 9.3 must be deployed only on nodes that run Layer 7 Vendor's APM solution. 

2) Prepare the services for APM Metric collection: For each service, Enable Metrics as per product documentation -  [Configure Trace and Metrics Collection](https://techdocs.broadcom.com/content/broadcom/techdocs/us/en/ca-enterprise-software/layer7-api-management/precision-api-monitoring/3-4/configuring/configure-trace-and-metrics-collection.html#concept.dita_c24702cd46b050b1320a1b1527e8095af6de0e1d_OptionalConfigureLatencyMetricsAssertion_)

3) Turn off Layer 7 Vendor's APM and make sure Port 9080 is free.

4) If you haven't already done so create the `/etc/signalfx-l7-forwarder` directory:
```
$ sudo mkdir /etc/signalfx-l7-forwarder
```

### Downloads:

- Download the SignalFx L7 forwarder 9.3 .tar.gz bundle from here - [https://github.com/kdroukman/ps_support/releases/download/layer7-9.3/sfx-l7-fwd-linux.tar.gz](https://github.com/kdroukman/ps_support/releases/download/layer7-9.3/sfx-l7-fwd-linux.tar.gz)

- Download the configuration file:
```
$ sudo wget https://raw.githubusercontent.com/kdroukman/ps_support/master/lenovo/layer7/lite/config.yaml -O /etc/signalfx-l7-forwarder/config.yaml
```

- Stop current `signalfx-l7-forwarder` service from running and download the service script to your `init.d` directory:
```
$ sudo service signalfx-l7-forwarder stop
$ sudo wget https://raw.githubusercontent.com/kdroukman/ps_support/master/lenovo/layer7/lite/signalfx-l7-forwarder-9_3.init -O /etc/init.d/signalfx-l7-forwarder
$ sudo chmod 755 /etc/init.d/signalfx-l7-forwarder
```

### Step 1:

Extract the forwarder to the `/etc/signalfx-l7-forwarder` directory:
```
sudo tar xzfv sfx-l7-fwd-linux.tar.gz -C /etc/signalfx-l7-forwarder
```

### Step 2:

Edit the /etc/signalfx-l7-forwarder/config.yaml file to ensure the following values are set:

```
listenAddress: 127.0.0.1:9080
signalFxAccessToken: <Replace with your SFX token>
signalFxRealm: us1
appName: layer7
appVersion: < specify version 9.3 or 9.4>
appEnvironment: <Replace with eservice-prod, eservice-nonprod, liecomm-prod, liecomm-nonprod>
intervalSeconds: 10s
logging:
  level: info
```


### Step 3:

Start the service with the following command:
```
$ sudo service signalfx-l7-forwarder start
```

Check the status of the service:
```
$ sudo service signalfx-l7-forwarder status
```

### Step 4

Enable the service to be automatically started on boot:
```
$ sudo chkconfig signalfx-l7-forwarder on
```

### Troubleshooting:

Should you have any issues with starting the service, execute the following steps to collect Debug logs - 
1) Amend `/etc/signalfx-l7-forwarder/config.yaml` and set Logging level to `debug`
2) Restart the service. 
3) Collect logs at `/var/log/signalfx-l7-forwarder.log` and send to our team for troubleshooting. 


### Metrics

SignalFx APM Metrics Forwarder for Layer 7 version 9.3 collects the following metrics in SignalFx:
Metric Name | Type | Dimensions
------------|------|-----------
l7.avg_resp_time.ms | gauge | host, service_uri, type (frontend or backend), environment
l7.req_size.bytes | gauge | host, service_uri, environment
l7.res_size.bytes | gauge | host, service_uri, environment
l7.request.success_count | counter | host, service_uri, environment
l7.request.count | counter | host, service_uri, environment



## Setting up the Metrics Forwarder on Layer 7 API Gateway version 9.4+ (or later)

**Pre-requisites**
1) APM Metrics Forwarder for Layer 7 API Gateway version 9.4 or later can be deployed either on each of the nodes, or centrally on a separate server. 

2) **Prepare the services for APM Metric collection by Configuring Layer 7 Gateway for External Metrics Collection** as documented in [Configure Gateway for External Metrics Collection](https://techdocs.broadcom.com/content/broadcom/techdocs/us/en/ca-enterprise-software/layer7-api-management/api-gateway/9-4/learning-center/overview-of-the-policy-manager/gateway-dashboard/configure-gateway-for-external-service-metrics.html)

3) When creating Service Metrics Event Listener Backing Policy, make sure that the output format is in the exact JSON format as per bellow. This is the format that the SignalFx APM Metrics Forwarder expects as input:
```
{"request":
{"id":"00000171720ae706-41ade0",
 "nodeId":"28db3239014749319d1e6c7276e79a58",
 "nodeName":"Gateway1",
 "nodeIp":"100.11.111.193",
 "serviceId":"e001cfd0c1c1ffaa18e187b5e72fdd38",
 "serviceName":"service name",
 "serviceUri":"service/name",
 "isPolicySuccessful":"true",
 "isPolicyViolation":"false",
 "isRoutingFailure":"false",
 "totalFrontendLatency":"21",
 "totalBackendLatency":"0",
 "time":"1587550389361000000"}
 }
```
Do not deviate from the above format. Do not use an escape sequence when printing out quotes - eg. make sure the strings in the output look like this: `"request"`, not this: `\"request\"`.

In the Backing Policy you are creating, set the metrics to be Routed to the Appropriate HTTP Server. 
  1) If you are deploying the Forwarder to each of the nodes - it should be `http://127.0.0.1:9080` (localhost)
  2) If you are deploying the Forwarder to a central server - it should be `http://<CENTRAL SERVER IP OR DNS>:9080`
  
  _Note: If port 9080 is in use by another application, you can change it to another suitable port number_

**Test the JSON Output**

If you need to test the output of your Backing Policy and ensure it adheres to the above, you can use the Simple Server created in POC: [server.py] (https://raw.githubusercontent.com/kdroukman/poc_support/master/server.py)

Download the Simple server to the same host where you plan to run the actual Forwarder on.
If you are not using Port value 9080, change it to the value you will be using for the SignalFx Forwarder. 

Run the server with a Python command. 
```
$ sudo python server.py
```
The script will write received output to a file called `layer7output_v3.txt` in you working directory. You can inspect the output format that is being received from Layer 7 API Gateway Backing Policy that you created, and make tweaks as necessary to ensure that it matches the required format. 

### Downloads:

- Download the SignalFx L7 forwarder 9.4 .tar.gz bundle from here - [https://github.com/kdroukman/ps_support/releases/download/layer7-9.4/sfx-l7-fwd-linux_9_4-v2.tar.gz](https://github.com/kdroukman/ps_support/releases/download/layer7-9.4/sfx-l7-fwd-linux_9_4-v2.tar.gz)

- Download the configuration file:
```
$ sudo wget https://raw.githubusercontent.com/kdroukman/ps_support/master/lenovo/layer7/lite/config.yaml -O /etc/signalfx-l7-forwarder/config.yaml
```

- Stop current `signalfx-l7-forwarder` service from running and download the service script to your `init.d` directory:
```
$ sudo service signalfx-l7-forwarder stop
$ sudo wget https://raw.githubusercontent.com/kdroukman/ps_support/master/lenovo/layer7/lite/signalfx-l7-forwarder.init -O /etc/init.d/signalfx-l7-forwarder
$ sudo chmod 755 /etc/init.d/signalfx-l7-forwarder
```

### Step 1:

Extract the forwarder to the `/etc/signalfx-l7-forwarder` directory:
```
sudo tar xzfv sfx-l7-fwd-linux_9_4.tar.gz -C /etc/signalfx-l7-forwarder
```

### Step 2:

Edit the /etc/signalfx-l7-forwarder/config.yaml file to ensure the following values are set:
_Note: If you are deploying to a central server, replace 127.0.0.1 with appropriate listner address_

```
listenAddress: 127.0.0.1:9080
signalFxAccessToken: <Replace with your SFX token>
signalFxRealm: us1
appName: layer7
appVersion: < specify version 9.3 or 9.4>
appEnvironment: <Replace with eservice-prod, eservice-nonprod, liecomm-prod, liecomm-nonprod>
intervalSeconds: 10s
logging:
  level: info
```


### Step 3:

Start the service with the following command:
```
$ sudo service signalfx-l7-forwarder start
```

Check the status of the service:
```
$ sudo service signalfx-l7-forwarder status
```

### Step 4

Enable the service to be automatically started on boot:
```
$ sudo chkconfig signalfx-l7-forwarder on
```

### Troubleshooting:

Should you have any issues with starting the service, execute the following steps to collect Debug logs - 
1) Amend `/etc/signalfx-l7-forwarder/config.yaml` and set Logging level to `debug`
2) Restart the service. 
3) Collect logs at `/var/log/signalfx-l7-forwarder.log` and send to our team for troubleshooting. 

### Metrics

SignalFx APM Metrics Forwarder for Layer 7 version 9.4+ collects the following metrics in SignalFx:
Metric Name | Type | Dimensions
------------|------|-----------
l7.avg_resp_time.ms | gauge | host, service_uri, type (frontend or backend), environment
l7.request.success_count | counter | host, service_uri, environment
l7.request.count | counter | host, service_uri, environment
l7.request.policy_violation_count | counter | host, service_uri, environment
l7.request.routing_failure_count | counter | host, service_uri, environment
