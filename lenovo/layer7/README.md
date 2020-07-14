# Configuration for Layer 7 Hosts

## Setup SignalFx Smart Agent on each Host

### Step 1:

Create a temporary directory on your host and download or copy the necessary files to it.

- Download the agent tar.gz bundle from here: [https://github.com/signalfx/signalfx-agent/releases/download/v5.3.3/signalfx-agent-5.3.3.tar.gz](https://github.com/signalfx/signalfx-agent/releases/download/v5.3.3/signalfx-agent-5.3.3.tar.gz)

- Download the configuration setup script and make it executable:
```
$ sudo wget https://raw.githubusercontent.com/kdroukman/ps_support/master/lenovo/signalfx_agent.sh
$ sudo chmod 744 signalfx_agent.sh
```

- Download the base agent.yaml file:
```
$ sudo wget https://raw.githubusercontent.com/kdroukman/ps_support/master/lenovo/standard/agent.yaml
```

### Step 2:

Extract the file to a directory on your host - /usr/lib
```
$ sudo tar xzfv signalfx-agent-5.3.3.tar.gz -C /usr/lib
```

Create a symbolic link to the agent executable:
```
$ sudo ln -s /usr/lib/signalfx-agent/bin/signalfx-agent /usr/bin/signalfx-agent
```

### Step 4:

Navigate to the newly created Signalfx Agent directory, and set the right loader for the libraries in the bundle:
```
$ cd /usr/lib/signalfx-agent/
$ sudo bin/patch-interpreter $(pwd)
```

### Step 5:
Create a signalfx-agent user:
```
$ getent passwd signalfx-agent >/dev/null || \
          useradd --system --home-dir /usr/lib/signalfx-agent --no-create-home --shell /sbin/nologin signalfx-agent
```

### Step 6:
Add signalfx-agent as a service on your RHEL or CentOS host. 

Download the following file to your `init.d` directory and make it executable:

```
$ wget https://raw.githubusercontent.com/signalfx/signalfx-agent/master/packaging/etc/init.d/signalfx-agent.rhel -O /etc/init.d/signalfx-agent
$ sudo chmod 755 /etc/init.d/signalfx-agent
```

### Step 7:

**Go back to your temporary directory** and run the following commands to set up the configuration. Make sure you replace the options with your values:
```
$ export SIGNALFX_ACCESS_TOKEN=<replace with your access token>
$ sudo sh signalfx_agent.sh \
  --action config \
  --realm us1 \
  --env <replace with the right environment name: liecomm-prod, liecomm-nonprod, eservice-prod, eservice-nonprod> \
  --config-path . \
  $SIGNALFX_ACCESS_TOKEN
```

Confirm that configuration directory was created with content:
```
$ ls /etc/signalfx
```

### Step 8:

If the previous step did not start the agent, start it with the following command:
```
$ sudo service signalfx-agent start
```

Check the status of the service:
```
$ sudo service signalfx-agent status
```

### Step 9:

Enable the service to be automatically started on boot:
```
$ sudo chkconfig signalfx-agent on
```

## Setup APM Metrics Forwarder

### Pre-requsites
You must have a SignalFx Smart Agent installed on your host. APM Metrics Forwarders provided here will be using `signalfx-agent` user and Access Token created as part of the SignalFx Smart Agent setup.

## Setting up the Metrics Forwarder on Layer 7 API Gateway version 9.3
_Disclaimer: While a metrics forwarder is provided for version 9.3, this is not the recommended or supported method. It is assumed that Lenovo will make the necessary upgrades to utilized recommended and supported method. In version 9.3 we incercept APM metrics published by the Layer 7 vendor APM solution provided alongside Layer 7 API Gateway._

**Pre-requisites**
1) APM Metrics Forwarder for Layer 7 API Gateway version 9.3 must be deployed only on nodes that run Layer 7 Vendor's APM solution. 
2) Prepare the services for APM Metric collection: For each service, Enable Metrics as per product documentation -  [Configure Trace and Metrics Collection](https://techdocs.broadcom.com/content/broadcom/techdocs/us/en/ca-enterprise-software/layer7-api-management/precision-api-monitoring/3-4/configuring/configure-trace-and-metrics-collection.html#concept.dita_c24702cd46b050b1320a1b1527e8095af6de0e1d_OptionalConfigureLatencyMetricsAssertion_)
3) Turn off Layer 7 Vendor's APM and make sure Port 9080 is free.

### Step 1:

Download the forwarder script and configuration file to a specific directory on your host:
```
$ sudo mkdir /etc/signalfx-l7-forwarder
$ sudo wget https://raw.githubusercontent.com/kdroukman/ps_support/master/lenovo/layer7/signalfx-l7-forwarder-9_3.py -0 /etc/signalfx-l7-forwarder/signalfx-l7-forwarder-9_3.py
$ sudo wget https://raw.githubusercontent.com/kdroukman/ps_support/master/lenovo/layer7/config.cfg -O /etc/signalfx-l7-forwarder/config.cfg
```

### Step 2:
Edit the /etc/signalfx-l7-forwarder/config.cfg file to ensure the following values are set:
```
[Server]
port=9080

[SignalFx]
realm=us1
dopost=1
service=layer7
version=9.3
env=**<Enter relevant application here: liecomm-nonprod, eservice-nonprod, liecomm-prod, eservice-prod>**

[Logging]
file=stdout
level=INFO
```

### Step 3:

Download the APM Metrics Forwarder service to your RHEL or CentOS host: 

Download the following file to your `init.d` directory and make it executable:

```
$ sudo wget https://github.com/kdroukman/ps_support/blob/master/lenovo/layer7/signalfx-l7-forwarder-9_3.init -O /etc/init.d/signalfx-l7-forwarder
$ sudo chmod 755 /etc/init.d/signalfx-l7-forwarder
```

### Step 4:

Start the service with the following command:
```
$ sudo service signalfx-l7-forwarder start
```

Check the status of the service:
```
$ sudo service signalfx-l7-forwarder status
```

### Step 5

Enable the service to be automatically started on boot:
```
$ sudo chkconfig signalfx-l7-forwarder on
```

### Troubleshooting:

Should you have any issues with starting the service, execute the following steps to collect Debug logs - 
1) Amend `/etc/signalfx-l7-forwarder/config.cfg` and set Logging level to `DEBUG`
2) Restart the service. 
3) Collect logs at `/var/log/signalfx-l7-forwarder.log` and send to our team for troubleshooting. 

