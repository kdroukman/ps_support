# Configuration for Layer 7 Hosts

## Setup SignalFx Smart Agent on each Host

### Step 1:

Create a temporary directory on your host and download or copy the necessary files to it.

- Download the agent tar.gz bundle from here: [https://github.com/signalfx/signalfx-agent/releases/download/v5.3.3/signalfx-agent-5.3.3.tar.gz](https://github.com/signalfx/signalfx-agent/releases/download/v5.3.3/signalfx-agent-5.3.3.tar.gz)

- Download the configuration settup script and make it executable:
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

## Setup APM Data Forwarder
Instructions coming soon...
