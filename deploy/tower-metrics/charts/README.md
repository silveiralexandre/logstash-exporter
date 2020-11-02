# tower-metrics

This Helm Chart deploys a Logstash pipeline for extracting Inventory metadata from a target Ansible Tower API

## Requirements

1. Connectivity from the Kubernetes worker nodes where Logstash pipeline is deployed to the target Ansible Tower and Elasticsearch instance.
2. Valid credentials at the target Ansible Tower and Elasticsearch instances (username, password, SSL Certificate).
3. Customize **values.yaml** file with appropriate credentials and requirements as instructed at the file itself.
4. Deploy the this chart on the same project/namespace where Ansible Tower is currently deployed or from where it can access the target Ansible Tower instance.

## Deployment

1. Clone [logstash-exporter](https://github.com/silveiralexandre/logstash-exporter) repository:

```shell
$ git clone https:/github.com/silveiralexandre/logstash-exporter.git
$ cd logstash-exporter/deploy/tower-metrics/charts
```

1. Customize **values.yaml** as instructions provided on the file itself.
2. Perform Helm install, as shown in the example below:

```shell
$ helm install tower-metrics -f values.yaml .
Release "tower-metrics" has been upgraded. Happy Helming!
NAME: tower-metrics
LAST DEPLOYED: Fri Sep 29 13:57:06 2020
NAMESPACE: tower
STATUS: deployed
REVISION: 2
TEST SUITE: None
```

## Post Deployment Verification

Confirm that the pod is running as expected as shown in the example below:

```shell
$ oc get deployment | grep -E 'NAME|tower-metrics'
NAME               DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
tower-metrics      1         1         1            1           1d

$ oc get pod | grep -E 'NAME|tower-metrics'
NAME                                READY     STATUS      RESTARTS   AGE
tower-metrics-5675d47555-m4j9l      1/1       Running     0          1d
```

Confirm data that is being properly extracted, by checking the pod logs as shown in the example below:

```shell
$ oc logs tower-metrics-5675d47555-m4j9l
[2020-09-28T22:15:52,606][INFO ][logstash.runner          ] Starting Logstash {"logstash.version"=>"7.6.2"}
[2020-09-28T22:15:52,622][INFO ][logstash.agent           ] No persistent UUID file found. Generating new UUID {:uuid=>"65191f5e-5e6c-4d48-88d9-b6a878191162", :path=>"/usr/share/logstash/data/uuid"}
[2020-09-28T22:15:54,005][INFO ][org.reflections.Reflections] Reflections took 32 ms to scan 1 urls, producing 20 keys and 40 values
[2020-09-28T22:15:57,403][WARN ][logstash.outputs.elasticsearch] You are using a deprecated config setting "document_type" set in elasticsearch. Deprecated settings will continue to work, but are scheduled for removal from logstash in the future. Document types are being deprecated in Elasticsearch 6.0, and removed entirely in 7.0. You should avoid this feature If you have any questions about this, please visit the #logstash channel on freenode irc. {:name=>"document_type", :plugin=><LogStash::Outputs::ElasticSearch
ssl_certificate_verification=>true, proxy=>http://myproxyhost.local:8080, password=><password>, hosts=>[//elasticsearch:9200], doc_as_upsert=>true, cacert=>"/usr/share/logstash/config/ca/ca.crt", ilm_enabled=>"false", index=>"ibm.cacf.tower.reportmetrics-%{+YYYY.ww}", id=>"7f869696a0326ef2963e940e4ccf29b43bf3c21d92e6c7a0dd7245651650c7c4", user=>"silveiralexandre@protonmail.com", ssl=>true, document_type=>"_doc", enable_metric=>true, codec=><LogStash::Codecs::Plain id=>"plain_9fc1c5ea-b713-41f7-965a-f64ef5b88dfe", enable_metric=>true, charset=>"UTF-8">, workers=>1, manage_template=>true, template_name=>"logstash", template_overwrite=>false, script_type=>"inline", script_lang=>"painless", script_var_name=>"event", scripted_upsert=>false, retry_initial_interval=>2, retry_max_interval=>64, retry_on_conflict=>1, ilm_rollover_alias=>"logstash", ilm_pattern=>"{now/d}-000001", ilm_policy=>"logstash-policy", action=>"index", sniffing=>false, sniffing_delay=>5, timeout=>60, pool_max=>1000, pool_max_per_route=>100, resurrect_delay=>5, validate_after_inactivity=>10000, http_compression=>false>}
[2020-09-28T22:15:57,801][INFO ][logstash.outputs.elasticsearch] Elasticsearch pool URLs updated {:changes=>{:removed=>[], :added=>[https://silveiralexandre%2540protonmail.com:xxxxxx@elasticsearch:9200/]}}
[2020-09-28T22:15:58,524][WARN ][logstash.outputs.elasticsearch] Restored connection to ES instance {:url=>"https://silveiralexandre%2540protonmail.com:xxxxxx@elasticsearch:9200/"}
[2020-09-28T22:15:58,652][INFO ][logstash.outputs.elasticsearch] ES Output version determined {:es_version=>7}
[2020-09-28T22:15:58,654][WARN ][logstash.outputs.elasticsearch] Detected a 6.x and above cluster: the `type` event field won't be used to determine the document _type {:es_version=>7}
[2020-09-28T22:15:58,785][INFO ][logstash.outputs.elasticsearch] New Elasticsearch output {:class=>"LogStash::Outputs::ElasticSearch", :hosts=>["//elasticsearch:9200"]}
[2020-09-28T22:15:58,886][WARN ][org.logstash.instrument.metrics.gauge.LazyDelegatingGauge] A gauge metric of an unknown type (org.jruby.RubyArray) has been created for key:
cluster_uuids. This may result in invalid serialization.  It is recommended to log an issue to the responsible developer/development team.
[2020-09-28T22:15:58,888][INFO ][logstash.javapipeline    ] Starting pipeline {:pipeline_id=>"main", "pipeline.workers"=>1, "pipeline.batch.size"=>125, "pipeline.batch.delay"=>50, "pipeline.max_inflight"=>125, "pipeline.sources"=>["/usr/share/logstash/pipeline/logstash.conf"], :thread=>"#<Thread:0x44d4d959 run>"}
[2020-09-28T22:15:58,915][INFO ][logstash.outputs.elasticsearch] Using default mapping template
[2020-09-28T22:15:58,924][INFO ][logstash.outputs.elasticsearch] Attempting to install template {:manage_template=>{"index_patterns"=>"logstash-*", "version"=>60001, "settings"=>{"index.refresh_interval"=>"5s", "number_of_shards"=>1}, "mappings"=>{"dynamic_templates"=>[{"message_field"=>{"path_match"=>"message", "match_mapping_type"=>"string", "mapping"=>{"type"=>"text", "norms"=>false}}}, {"string_fields"=>{"match"=>"*", "match_mapping_type"=>"string", "mapping"=>{"type"=>"text", "norms"=>false, "fields"=>{"keyword"=>{"type"=>"keyword", "ignore_above"=>256}}}}}], "properties"=>{"@timestamp"=>{"type"=>"date"}, "@version"=>{"type"=>"keyword"}, "geoip"=>{"dynamic"=>true, "properties"=>{"ip"=>{"type"=>"ip"}, "location"=>{"type"=>"geo_point"}, "latitude"=>{"type"=>"half_float"}, "longitude"=>{"type"=>"half_float"}}}}}}}
[2020-09-28T22:15:59,518][INFO ][logstash.inputs.exec     ] Registering Exec Input {:type=>nil, :command=>"/usr/bin/reportmetrics", :interval=>nil, :schedule=>"*/1 * * * *"}
[2020-09-28T22:15:59,524][INFO ][logstash.javapipeline    ] Pipeline started {"pipeline.id"=>"main"}
[2020-09-28T22:15:59,600][INFO ][logstash.agent           ] Pipelines running {:count=>1, :running_pipelines=>[:main], :non_running_pipelines=>[]}
[2020-09-28T22:15:59,825][INFO ][logstash.agent           ] Successfully started Logstash API endpoint {:port=>9600}
/usr/share/logstash/vendor/bundle/jruby/2.5.0/gems/rufus-scheduler-3.0.9/lib/rufus/scheduler/cronline.rb:77: warning: constant ::Fixnum is deprecated
```

You may also confirm on your target Elasticsearch if the indices were properly created as expected, as shown in the example below:

```shell
$ curl https://your-target-elasticsearch-host:9200/_cat/indices?v
health status index                                         uuid                   pri rep docs.count docs.deleted store.size pri.store.size
green  open   ibm.cacf.tower.reportmetrics-2020.40          M-KgnYtfS1emVuXAT30LiQ   1   1       3249            0    584.2kb        305.5kb
```

And on Kibana as shown in the example below:

![kibana](../docs/../../../docs/kibana.png)