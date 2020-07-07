# Configurations for a Docker/Docker Swarm host

These instructions are for running on Docker Swarm. See K8s folder for Kubernetes specific instructions. 

On a Docker host you will be setting up SignalFx Smart Agent to run in a container, and corresponding configuration to be placed on host. The configuration is provided by master file - `agent.yaml` - and supplementary files for each of the required monitors to be switched on the host - eg: `rabbitmq.yaml`, `redis.yaml`. 

As your Docker hosts require µAPM - there is an additional configuration file called `trace.yaml` that will need to be placed on each of the µAPM hosts. 

The provided installation and configuration script - `signalfx-agent.sh` - does most of the work for you by importing correct configuration files to a specific directory on your hosts - `/etc/signalfx`

You are required to host the configuration files centrally, either on a network file system, or GitLab. If you are unable to set up the central repository at the time of installation, download the .yaml files to a local files system, or alternatively you can fetch them from here:
[https://github.com/kdroukman/ps_support/releases/download/docker](https://github.com/kdroukman/ps_support/releases/download/docker)

All the YAML files that the installer script relies on must be in the `--config-path` directory or URL passed to the installation script. The script will place them in the appropriate directories. 

The configurations provided are based on POC setup, and default assumptions. You are required to review the configurations set up, and make necessary changes for your environment. 

### How Discovery Rules work

A lot of the monitor configurations here use discovery rules (discoveryRule) which allows you to create different monitors and different permutations of the same monitor for different hosts. If a discoveryRule finds a matching service, it will switch on the monitor. Otherwise the monitor will not turn on.

See [Service Discovery](https://github.com/signalfx/signalfx-agent/blob/master/docs/auto-discovery.md) for more details.

### How monitor configuration has been set up
To assist with configuration and management, each monitor has been split into a separate YAML file. 

Once the installation/config script is run, the following folder structure is created on host:
```
/etc/signalfx/agent.yaml
/etc/signalfx/monitors/<monitor1>.yaml
/etc/signalfx/monitors/<monitor2>.yaml
...
/etc/signalfx/monitors/<monitorN>.yaml
```

### How to reference sensitive data in configuration files
For sensitive data, such as passwords, Vault or Zookeeper can be used. 
See [Remote Configuration](https://docs.signalfx.com/en/latest/integrations/agent/remote-config.html) for examples on how to set this.

### caution:
See [Dimensions Names and Values](https://developers.signalfx.com/metrics/data_ingest_overview.html#_dimension_names_and_values) documentation for guidance on naming your environments and hosts. 

### The importance of Environment variable
It has been agreed that the environment variable will be used to correlate and filter all data across the prod and non-prod infrastructure and application environments. It accepts 8 values:
- licomm-prod, licomm-nonprod
- eservice-prod, eservices-nonprod
- necpc-prod, necpc-nonprod
- accounts-prod, accounts-nonprod

All the configuration settings are being designed to include this. Therefore, you must use configurations provided here, and not the ones installed by RPM or packaged into the provided Docker image. 

### parameters:
The script accepts the following parameters:
<code>
  Option | Decription | Optional/Mandatory
---------|------------|-------------------
**--action <install\|config>**  | Specify whether to run installation or update configuration only. | Optional. Default is 'install'
**--package-version <version>** | The agent package version to install. | Optional.
**--realm <us0\|us1\|eu0\|...>** | SignalFx realm to use (used to set ingest-url and api-url automatically). | Mandatory.
**--trace-endpoint <url>**     | Path to SignalFx trace endpoint or on-premise OpenTelemetry Collector for send µAPM trace to. | Optional for non-APM, mandatory for APM.
**--cluster <custer name>**    | The user-defined environment/cluster to use (corresponds to 'cluster' option in agent). | Optional - not necessary outside of K8s 
**--test**                   | Use the test package repo instead of the primary. | Optional
**--beta**                    | Use the beta package repo instead of the primary. | Optional
**--env <environment name>**   | The name of Lenovo environment/application (liecomm, eservices, accounts, necpc -prod/-nonprod). | Mandatory.
**--hostname <hostname>**      | Override default hostname. | Optional
**--config_path <url of path>** | Location of agent.yaml and corresponding monitors. | Mandatory to use custom config files. Otherwise the bare-bones default one will be used..
**--monitors <list>**          | Comma (,) seperated list of monitors to load. | Mandatory for any hosts that require any extra monitors enabled. Otherwise only host metrics will be collected.
  </code>

### Before you begin: 
Download and host the YAML files in a central location on your local Intranet. 
Alternatively - you can refernce the following Release location - [https://github.com/kdroukman/ps_support/releases/download/docker](https://github.com/kdroukman/ps_support/releases/download/docker)

### SignalFx Smart Agent Docker Image
Docker image is available here [quay.io/signalfx/signalfx-agent](quay.io/signalfx/signalfx-agent). You will need to supply the version tag. Inspect the latest version. At the time of this document - the latest version is 5.3.3

If you are using an internal container registry, pull the provided image, and push it to your internal registry.

Here we provide instructions on setting up docker and host monitoring Smart Agent and setting up APM instrumentation. 

### You will need:
1) SignalFx Smart Agent Docker image hosting the agent
2) signalfx-agent.sh script for setting up the configuration
3) For APM - Language Specific Instrumenation Library. eg: [https://github.com/signalfx/signalfx-java-tracing/releases](https://github.com/signalfx/signalfx-java-tracing/releases) for microservices written in Java.
In these Examples APM traces are being sent directly to SignalFx SaaS ingest URL. Additional Data Masking will not be applied.

## Step 1: 

Run signalfx-agent.sh script to set up configuration files in your /etc/signalfx directory on the host. 
Ensure "--action config" is set so that only config files are created:

Example command:
```
$ export SIGNALFX_ACCESS_TOKEN=abcdefghij123445656
$ sh signalfx_agent.sh --realm us0 --action config --realm us1 --env eservice-prod --hostname kh-docker-1 --monitors trace,redis --config-path https://github.com/kdroukman/ps_support/releases/download/docker --trace-endpoint https://ingest.us1.signalfx.com/v2/trace $SIGNALFX_ACCESS_TOKEN
```

_note: There will be an error at the end of script output stating that signalfx-agent.service is not found. Ignore it. I will update the script to remove it._
```
$ export SIGNALFX_ACCESS_TOKEN=<Change to your token value>
$ sudo sh signalfx-agent.sh \
    --action config \
    --realm us1 \
    --env <mandatory environment> \
    --hostname <optional hostname>  \
    --monitors <Comma separated list, no spaces. For APM include trace. eg: --monitors trace,redis will add trace and redis config> \
    --config-path <web or fs location where .yaml files are located eg: --config-path https://github.com/kdroukman/ps_support/releases/download/docker> \
    --trace-endpoint https://ingest.us1.signalfx.com/v2/trace \
    $SIGNALFX_ACCESS_TOKEN
```
Verify the contents of /etc/signalfx folder after running the script. 

## Step 2:
Run the docker container with appropriate volume mappings. 

We provide a Docker image at [quay.io/signalfx/signalfx-agent](http://quay.io/signalfx/signalfx-agent). 
As instructed in pre-requisites, you can pull this image and push it to your local registry if necessary. 

Start SignalFx Smart Agent with the following options. These will map the Smart Agent to host's network and bind APM trace listener to http://localhost:9080/ on the host. Additionally, the pre-build configuration files will be overwritten with those in you /etc/signalfx directory:


```
$ docker run \
    --pid host \
    --net host \
    -v /:/hostfs:ro \ 
    -v /var/run/docker.sock:/var/run/docker.sock:ro \ 
    -v /etc/signalfx/:/etc/signalfx/:ro \ 
    -v /etc/passwd:/etc/passwd:ro \
    quay.io/signalfx/signalfx-agent:<version - use the number. eg. 5.3.3, there appears to be no latest tag>
```

Verify that the agent is running with docker ps command, or check the logs with docker logs <container name>.
    
## Adding APM

The SignalFx Smart Agent configured as per above will already expose endpoint on http://localhost:9080 and add the necessary environment tag to all the traces.
You will need to apply the necessary instrumentation to your application microservices. 

Available auto-instrumentation options are listed here: [https://docs.signalfx.com/en/latest/apm/apm-instrument/apm-instr-overview.html#automatically-instrument-an-application](https://docs.signalfx.com/en/latest/apm/apm-instrument/apm-instr-overview.html#automatically-instrument-an-application)

_SignalFx can also accept Zipkin v1 or b2 JSON or Jaeger Thrift or gRPC format traces produced by any other libraries, such as OpenTelemetry. This is an option if provided libraries cannot be used._

This following example illustrates setting up APM for a Java microservice.

_note: There are various ways to add libraries and environment variables to containers. This illustrates one such method._

**1)** Download Java Trace agent .jar file from: https://github.com/signalfx/signalfx-java-tracing/releases
**2)** Add the .jar file to a location in your container by packaging it into an image:

Dockerfile (for a basic Spring Pet Clinic application):
```
FROM java:8
COPY ./spring-petclinic/target/ /var/www/java
RUN mkdir -p /opt/signalfx-tracing
COPY signalfx-tracing.jar /opt/signalfx-tracing
WORKDIR /var/www/java

CMD java -javaagent:/opt/signalfx-tracing/signalfx-tracing.jar -jar *.jar
```

**3)** Run the docker container and pass in SignalFx environment variables that are requierd by the Tracing library. Here, only the environment varilables are what is required to setup service name and direct traces to Smart Agent listener. The other options are specific to the application service itself.
Notice that the trace listener at `http://localhost:9080` endpoint can only be accessed on the host, not within the container. Therefore, the $HOSTNAME value is passed to the container, so that it can send traces to the necessary endpoint.

```
docker run -p 8080:8080 --env SIGNALFX_SERVICE_NAME=kh-pet-clinic --env SIGNALFX_ENDPOINT_URL=http://$HOSTNAME:9080/v1/trace -d  pet-clinic
```

You can specify additional optional SignalFx variables are per documentation: [https://github.com/signalfx/signalfx-java-tracing](https://github.com/signalfx/signalfx-java-tracing)

For example you can use SIGNALFX_SPAN_TAGS to tag your traces with additional custom details. eg: `SIGNALFX_SPAN_TAGS="release:canary,version:2.1"`
This will be viewable and searchable when examining traces in SignalFx.
