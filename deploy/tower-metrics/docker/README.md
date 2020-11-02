# tower-metrics

Use the attached [Makefile](Makefile) for assisting on building this artifact. Configuration changes and/or environment settings updates should be done at the files listed below:

- [deploy.env](deploy.env)
- [config.env](config.env)

### Building Local Docker Image

For building local images, execute `make dev-build` or `make prod-build` from the deployment directory as shown in the example below:

For building a local image, execute `make dev-build` from the deployment directory for building a **development** image, and `make prod-push` for buiding a **production** image, as shown in the example below:

```bash
$ git clone https:/github.com/silveiralexandre/logstash-exporter.git
$ cd logstash-exporter/deploy/tower-metrics/docker
$ make dev-build
docker build -t tower-metrics \
                         -f "/git/github.com/silveiralexandre/logstash-exporter/deploy/tower-metrics/docker/Dockerfile" \
                            "/git/github.com/silveiralexandre/logstash-exporter/deploy/tower-metrics/docker" \
                         --tag "gts-cacf-global-team-dev-docker-local.artifactory.swg-devops.com/logstash-custom/tower-metrics:1.0" \
                         --no-cache \
                         --compress
Sending build context to Docker daemon  9.165MB
Step 1/14 : FROM docker.elastic.co/logstash/logstash:7.6.2
 ---> fa5b3b1e9757
Step 2/14 : LABEL name="Logstash pipeline for extracting Ansible Tower inventory data"  maintainer=silveiralexandre@protonmail.com       build-date=20200929
 ---> Running in 122c7e25893d
Removing intermediate container 122c7e25893d
 ---> 04e8b0373af4
Step 3/14 : WORKDIR /usr/share/logstash
 ---> Running in 1ef88b786cd5
Removing intermediate container 1ef88b786cd5
 ---> cf8b567b7c87
Step 4/14 : USER root
 ---> Running in e82ee3347c69
Removing intermediate container e82ee3347c69
 ---> a8131d792f52
Step 5/14 : ADD ./bin/logstash-exporter /usr/bin/logstash-exporter
 ---> f6734ed19d2c
Step 6/14 : RUN chown logstash:root /usr/bin/logstash-exporter && chmod +x /usr/bin/logstash-exporter
 ---> Running in a928ff87849c
Removing intermediate container a928ff87849c
 ---> d18108583246
Step 7/14 : USER logstash
 ---> Running in fefbf9bc44c0
Removing intermediate container fefbf9bc44c0
 ---> 61198c591b5a
Step 8/14 : RUN rm -f /usr/share/logstash/pipeline/logstash.conf
 ---> Running in 1b65b6bec8ef
Removing intermediate container 1b65b6bec8ef
 ---> b3190b689143
Step 9/14 : ADD ./config/ ./config/
 ---> b419adc97686
Step 10/14 : ADD ./pipeline/ ./pipeline/
 ---> e8e0e391cb29
Step 11/14 : ENV LANG=en_US.UTF-8 LC_ALL=en_US.UTF-8
 ---> Running in a3f617b273fb
Removing intermediate container a3f617b273fb
 ---> 268d759ad4bf
Step 12/14 : ENV LS_JAVA_OPTS="-Xmx256m -Xms256m"
 ---> Running in ceeb6f3dc9ff
Removing intermediate container ceeb6f3dc9ff
 ---> 8a6667369ade
Step 13/14 : EXPOSE 5044 9600
 ---> Running in 5b6d72a54758
Removing intermediate container 5b6d72a54758
 ---> b19da90a0de6
Step 14/14 : ENTRYPOINT ["/usr/local/bin/docker-entrypoint"]
 ---> Running in 6e36ce5dfe1a
Removing intermediate container 6e36ce5dfe1a
 ---> 4a75f8a624a4
Successfully built 4a75f8a624a4
Successfully tagged tower-metrics:latest
Successfully tagged gts-cacf-global-team-dev-docker-local.artifactory.swg-devops.com/logstash-custom/tower-metrics:1.0
```

### Pushing to Remote Docker Repository (Artifactory)

For pushing your local image to the target repository, execute `make dev-push` from the deployment directory for pushing to the **development** image repository, and `make prod-push` for pushing the image to the **production** image repository,
 as shown in the example below:

```bash
$ git clone https:/github.com/silveiralexandre/logstash-exporter.git
$ cd logstash-exporter/deploy/tower-metrics/docker
$ make dev-push
publish 1.0 to gts-cacf-global-team-dev-docker-local.artifactory.swg-devops.com
docker push "gts-cacf-global-team-dev-docker-local.artifactory.swg-devops.com/logstash-custom/tower-metrics:1.0"
The push refers to repository [gts-cacf-global-team-dev-docker-local.artifactory.swg-devops.com/logstash-custom/tower-metrics]
8c7d2ca52d86: Pushed
2805b8c2a5ab: Pushed
1ab8fb7f43f8: Pushed
2c61f826eef0: Pushed
bb214dbaa374: Pushed
a2d2943f8c9a: Layer already exists
bfc1b8ecc314: Layer already exists
b729863034fb: Layer already exists
6c9bcb75c5a8: Layer already exists
bbbcfcffed95: Layer already exists
d5479313c55a: Layer already exists
e17d61db2256: Layer already exists
61bbffa29f88: Layer already exists
f1d15f79533e: Layer already exists
aef5ae6629a1: Layer already exists
77b174a6a187: Layer already exists
1.0: digest: sha256:92027209fc830ae53783c8472df4ba7a89a41bf7bdd969f18063b0b13716faa9 size: 3868
```

### Help

Run `make help` for the full reference of commands supported:

```bash
$ make help
help                           This help.
version                        Output the tagged version as defined at `config.env`
stop                           Stop and remove a running container
clean                          Removes all local images of APP_NAME
purge                          Removes all dangling images
dev-login                      Performs docker login to the CACF DEV image repository
dev-build                      Build the container image
dev-run                        Run container on port configured in `config.env`
dev-push                       Publish the `{version}` as defined at `config.env`
prod-login                     Performs docker login to the CACF production image repository
prod-build                     Build the container image
prod-run                       Run container on port configured in `config.env`
prod-push                      Publish the `{version}` as defined at `config.env`
```