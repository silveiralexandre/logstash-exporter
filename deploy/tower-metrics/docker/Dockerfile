FROM docker.elastic.co/logstash/logstash:7.6.2

LABEL name="Logstash pipeline for extracting Ansible Tower inventory data" \
      maintainer=silveiralexandre@protonmail.com \
      build-date=20200929

WORKDIR /usr/share/logstash

USER root
ADD ./bin/logstash-exporter /usr/bin/logstash-exporter
RUN chown logstash:root /usr/bin/logstash-exporter && chmod +x /usr/bin/logstash-exporter

# Replace Logstash defaults with custom configs and pipeline
USER logstash
RUN rm -f /usr/share/logstash/pipeline/logstash.conf
ADD ./config/ ./config/
ADD ./pipeline/ ./pipeline/

ENV LANG=en_US.UTF-8 LC_ALL=en_US.UTF-8
ENV LS_JAVA_OPTS="-Xmx256m -Xms256m"

EXPOSE 5044 9600
ENTRYPOINT ["/usr/local/bin/docker-entrypoint"]