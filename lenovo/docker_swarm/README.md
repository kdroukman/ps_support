# Config for Docker host

These instructions are for running on Docker Swarm. See K8s folder for Kubernetes specific instructions. 

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

We provide a Docker image at[quay.io/signalfx/signalfx-agent](quay.io/signalfx/signalfx-agent). 
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
    
    
