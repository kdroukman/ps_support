# Configurations for a Kubernetes

These instructions are for running on Kubernetes. For Docker Swarm see docker-swarm directory.  

In a Kubernetes environment you will be deploying everything via the provided YAML files. The SignalFx Smart Agent will run as a DaemonSet, which means exactly one copy will deployed to every host in your Kubernetes cluster.
Corresponding Configuration will be created from ConfigMapAgent.yaml file. 
This configuration has been designed to discovery additional services such as Redis, RabbitMQ, etc. 

As your Kubernetes hosts require µAPM - ConfigMapAgent.yaml provides additional settings to enable this. 

The provided YAML do most of the work for you by setting up the Smart Agent and required configuration. 

You can download and host these configuration files centrally, together with your existing Kubernetes configurations. 
The listed files are also available for download here:
[https://github.com/kdroukman/ps_support/releases/download/k8s](https://github.com/kdroukman/ps_support/releases/download/k8s)

You will be using command line tool provided by Kubernetes - `kubectl` - to perform the deployments. 

### Do you need to use Helm?

As not all your hosts have Helm installed, the provided YAML files can be installed using standard Kubernetes command line tool. 
You will need to update some configurations to your environment, cluster, and application names. Instructions on what needs to be updated are listed in the upcoming sections. 

### How Discovery Rules work

The monitor configurations you see in ConfigMapAgent.yaml use discovery rules (discoveryRule) which allows you to create different monitors and different permutations of the same monitor for different hosts. If a discoveryRule finds a matching service, it will switch on the monitor. Otherwise the monitor will not turn on.

See [Service Discovery](https://github.com/signalfx/signalfx-agent/blob/master/docs/auto-discovery.md) for more details.

### How monitor configuration has been set up
The various services to be discovered are listed in ConfigMapAgent.yaml. This ConfigMap generates a file called agent.yaml that will be planced in a SignalFx Smart Agent container created by the DaemonSet. You are required to review the listed monitors and adjust any discover rules as necessary. The rules listed assume standard port numbers are used. 

### How to reference sensitive data in configuration files
For sensitive data, such as passwords, Vault or Zookeeper can be used. 
See [Remote Configuration](https://docs.signalfx.com/en/latest/integrations/agent/remote-config.html) for examples on how to set this.

### caution:
See [Dimensions Names and Values](https://developers.signalfx.com/metrics/data_ingest_overview.html#_dimension_names_and_values) documentation for guidance on naming your environments. 

### The importance of Environment variable
It has been agreed that the environment variable will be used to correlate and filter all data across the prod and non-prod infrastructure and application environments. It accepts 8 values:
- licomm-prod, licomm-nonprod
- eservice-prod, eservices-nonprod
- necpc-prod, necpc-nonprod
- accounts-prod, accounts-nonprod

All the configuration settings are being designed to include this. Therefore, you must use configurations provided here, and not the ones hosted in main SignalFx repository. 

### Before you begin: 
Download the YAML files where you can easily deploy them to your Kubernetes clusters.. 
[https://github.com/kdroukman/ps_support/releases/download/k8s](https://github.com/kdroukman/ps_support/releases/download/k8s)

### SignalFx Smart Agent Docker Image
The K8s DaemonSet uses the Docker image available here [quay.io/signalfx/signalfx-agent](quay.io/signalfx/signalfx-agent).
Inspect the latest version. At the time of this document - the latest version is 5.3.3

If you are using an internal container registry, pull the provided image, and push it to your internal registry.
You will need to update DaemonSet.yaml to reference the image in your internal registry. 



Here we provide instructions on setting up Kuberentes Smart Agent and setting up APM instrumentation. 

### You will need:
1) SignalFx Smart Agent Kuberenetes YAML files
2) For APM - Language Specific Instrumenation Library. eg: [https://github.com/signalfx/signalfx-java-tracing/releases](https://github.com/signalfx/signalfx-java-tracing/releases) for microservices written in Java.

In these Examples APM traces are being sent directly to SignalFx SaaS ingest URL. Additional Data Masking will not be applied.

## Step 1: 

Download all the YAML files listed here to a directory where you will be deploying them from. 
```
$ mkdir signalfx
$ cd signalfx
$ wget https://github.com/kdroukman/ps_support/releases/download/k8s/*.yaml
```

Update the following files:

File Name | What to update
----------|---------------
ConfigMapEnv.yaml | This file sets up all the necessary enviroment variables. Update `CLUSTER_NAME` to your Kubernetes cluster name. Update `ENV_NAME` to one of 8 environment values. `REALM` and `TRACE_ENDPOINT_URL` should already be set to us1. Don't forget to update this for each environment and cluster.
ClusterRoleBinding.yaml | Open this file and update `namespace` to your namespace.
DaemonSet.yaml | If you are hosting the SignalFx Smart Agent image in your Docker registry, update the image location and associate tag. Otherwise leave as is.
ConfigMapAgent.yaml | This file sets up the agent.yaml text to output. You need to review this, specifically the discoveryRules for monitors. Note, that there is a long list of monitors, not all will be in your environment. Only review the ones you require, ignore the rest.


## Step 2:
After the YAML files have been reviewed and updated, deploy the Smart Agent.

1. Create a secret in K8s with your org's access token:
```
$ kubectl create secret generic --from-literal access-token=<YOUR_ACCESS_TOKEN> signalfx-agent
```
2. Create the DaemonSet and associated Objects in your Kubernetes cluster. In the directory where all your YAML files are located, run:
```
$ kubectl apply -R -f .
```
_This command assumes the directory contains only the YAML files listed here. If there are any other YAML files, it will create everything._ 


SignalFx Smart Agent DaemonSet will be deployed to every node in your Kubernetes cluster. 
For APM it will expose `http://<NODE_IP>:9080/`, the configuration is already set up to derive NODE_IP, you don't need to do anything extra.

Verify that the agent is running with `kubectl get ds` command, or check events with `kubectl get events` and logs with `kubectl logs <agent pod name>`.
    
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
Notice that the trace listener at `http://localhost:9080` endpoint can only be accessed on the host, not within the container. Therefore, the $HOSTNAME value is passed to the container, so that it can send traces to the necessary endpoint.

```
docker run -p 8080:8080 --env SIGNALFX_SERVICE_NAME=kh-pet-clinic --env SIGNALFX_ENDPOINT_URL=http://$HOSTNAME:9080/v1/trace -d  pet-clinic
```

You can specify additional optional SignalFx variables are per documentation: [https://github.com/signalfx/signalfx-java-tracing](https://github.com/signalfx/signalfx-java-tracing)

For example you can use SIGNALFX_SPAN_TAGS to tag your traces with additional custom details. eg: `SIGNALFX_SPAN_TAGS="release:canary,version:2.1"`
This will be viewable and searchable when examining traces in SignalFx.
