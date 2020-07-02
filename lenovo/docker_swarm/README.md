# Config for Docker host

*Pre-requisites*:
Clone the image to your local repository where required.

*Step 1*: 

Create configuration directory on the host:
```
$ sudo sh signalfx-agent.sh \
    --action config \
    --realm us1 \
    --env <mandatory environment> \
    --hostname <optional hostname>  \
    --monitors <comma separate list of additional monitors to add. e.g redis,trace> \
    --config-path <web or fs location where .yaml files are located> \
    --trace-endpoint <SignalFx or OTel Collector trace endpoints if APM is on host>
    SIGNALFX_ACCESS_TOKEN
```

*Step 2*:
Run the docker container with appropriate volume mappings. 

We provide a Docker image at <a href=quay.io/signalfx/signalfx-agent>quay.io/signalfx/signalfx-agent</a>. The image is tagged using the same agent version scheme.

If you are using Docker outside of Kubernetes, you can run the agent in a Docker container and still gather metrics on the underlying host by running it with the following flags:

```
$ docker run \
    --name signalfx-agent \ 
    --pid host \
    --net host \
    -v /:/hostfs:ro \ 
    -v /var/run/docker.sock:/var/run/docker.sock:ro \ 
    -v /etc/signalfx/:/etc/signalfx/:ro \ 
    -v /etc/passwd:/etc/passwd:ro \
    quay.io/signalfx/signalfx-agent:<version>
```
This assumes you have the agent config in the conventional directory (/etc/signalfx) on the root mount namespace. This also assumed you are using the provided image. Change the defaults as necessary.
