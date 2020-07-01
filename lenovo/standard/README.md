<b>Configurations for Standard Host</b>

To install on new host:
<code>
  sudo sh signalfx-agent.sh --action install --realm us1 --env <environment name> --hostname <override hostname> --monitors <extra monitors to add> --config-path <path to yaml templates> SIGNALFX_ACCESS_TOKEN
  </code>

To update a version:
<code>
  sudo rm -Rf /etc/signalfx
  sudo sh signalfx-agent.sh --action install --realm us1 --env <environment name> --hostname <override hostname> --monitors <extra monitors to add> --config-path <path to yaml templates> --package-version <version> SIGNALFX_ACCESS_TOKEN
</code>
  
To update configuration:
<code>
  sudo sh signalfx-agent.sh --action config --realm us1 --env <environment name> --hostname <override hostname> --monitors <extra monitors to add> --config-path <path to yaml templates> SIGNALFX_ACCESS_TOKEN
</code>
