# Configurations for Standard Host

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
