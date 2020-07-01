Configurations for Docker host


We provide a Docker image at quay.io/signalfx/signalfx-agent. The image is tagged using the same agent version scheme.

If you are using Docker outside of Kubernetes, you can run the agent in a Docker container and still gather metrics on the underlying host by running it with the following flags:

$ docker run \
    --name signalfx-agent \
    --pid host \
    --net host \
    -v /:/hostfs:ro \
    -v /var/run/docker.sock:/var/run/docker.sock:ro \
    -v /etc/signalfx/:/etc/signalfx/:ro \
    -v /etc/passwd:/etc/passwd:ro \
    quay.io/signalfx/signalfx-agent:<version>
This assumes you have the agent config in the conventional directory (/etc/signalfx) on the root mount namespace. 
