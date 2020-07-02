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
    
## Adding APM

The SignalFx Smart Agent configured as per above will already expose endpoint on http://localhost:9080 and add the necessary environment tag to all the traces.
You will need to apply the necessary instrumentation to your application microservices. 

Available auto-instrumentation options are listed here: [https://github.com/kdroukman/ps_support/blob/master/lenovo/standard/README.md](https://github.com/kdroukman/ps_support/blob/master/lenovo/standard/README.md)

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

```
docker run -p 8080:8080 --env SIGNALFX_SERVICE_NAME=kh-pet-clinic --env SIGNALFX_ENDPOINT_URL=http://$HOSTNAME:9080/v1/trace -d  pet-clinic
```

You can specify additional optional SignalFx variables are per documentation: [https://github.com/signalfx/signalfx-java-tracing](https://github.com/signalfx/signalfx-java-tracing)

For example you can use SIGNALFX_SPAN_TAGS to tag your traces with additional custom details. eg: `SIGNALFX_SPAN_TAGS="release:canary,version:2.1"`
This will be viewable and searchable when examining traces in SignalFx.
