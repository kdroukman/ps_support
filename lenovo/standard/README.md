# Configurations for Standard Host

A lot of the monitors here use discovery rules (discoveryRule) which allows you to create different monitors and different permutations of the same monitor for different hosts. If a discoveryRule finds a matching service, it will switch on the monitor. Other the monitor will not turn on, even if it is available on the host. 

See [Service Discovery](https://github.com/signalfx/signalfx-agent/blob/master/docs/auto-discovery.md) for more details.

To assist with configuration and management, each monitor has been split into a separate YAML file. 

Once the installation/config script is run, the following folder structure is created on host:
```
/etc/signalfx/agent.yaml
/etc/signalfx/monitors/<monitor1>.yaml
/etc/signalfx/monitors/<monitor2>.yaml
...
/etc/signalfx/monitors/<monitorN>.yaml
```

For sensitive data, such as passwords, Vault or Zookeeper can be used. 
See [Remote Configuration](https://docs.signalfx.com/en/latest/integrations/agent/remote-config.html) for examples on how to set this.

### caution:
See [Dimensions Names and Values](https://developers.signalfx.com/metrics/data_ingest_overview.html#_dimension_names_and_values) documentation for guidance on naming your environments and hosts. 

Example command:
```
$ export SIGNALFX_ACCESS_TOKEN=abcdefg12345677
$ sudo sh signalfx-agent.sh --action install --realm us1 --env liecomm-nonprod --hostname my-test-host --monitors mysql --config-path https://github.com/kdroukman/ps_support/releases/download/standard $SIGNALFX_ACCESS_TOKEN
```

To install on new host:
```
  $sudo sh signalfx-agent.sh \
  --action install \
  --realm us1 \
  --env <environment name> \
  --hostname <override hostname> \
  --monitors <extra monitors to add> \
  --config-path <path to yaml templates> \
  SIGNALFX_ACCESS_TOKEN
```

To update a version:
```
  $sudo rm -Rf /etc/signalfx
  $sudo sh signalfx-agent.sh \
  --action install \
  --realm us1 \
  --env <environment name> \
  --hostname <override hostname> \
  --monitors <extra monitors to add> \
  --config-path <path to yaml templates> \
  --package-version <version> \
  SIGNALFX_ACCESS_TOKEN
```
  
To update configuration:
```
  $sudo sh signalfx-agent.sh \
  --action config \
  --realm us1 \
  --env <environment name> \
  --hostname <override hostname> \
  --monitors <extra monitors to add> \
  --config-path <path to yaml templates> \
  SIGNALFX_ACCESS_TOKEN
```
