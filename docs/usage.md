# logstash-exporter

## Usage

1. Download the latest compiled binary from [releases page](https://github.com/silveiralexandre/logstash-exporter/releases) -- for the full instructions on how to execute the program, execute it whithout any flags or using the `-h` flag as shown below:

```shell
$ ./logstash-exporter -h
Generic exporter to feed JSON data from input of a source API into a target Logstash instance

Usage:
        logstash-exporter -n [ <report_name> ] [ -t <timeout> ] [ -l <limit> ] [ -r <retries> ] [ -ls ]

Examples:
        $ logstash-exporter -n tower-inventorydata -l 50 -t 90 -r 3
        $ logstash-exporter -ls

This will execute the informed script/command 50 times, limiting concurrency to 10 at the time.
Options:
   -n        Name of report to be executed
   -ls       List of currently supported reports
   -l        Limit of concurrent requests to be executed (default is 50)
   -t        Timeout in seconds for requests (default is 90)
   -r        Number of retries for HTTP requests (default is 3)
   -v        Prints version information

```
1. For running locally, export the environment variables and execute the script with the appropriate flags as demonstrated below:

```shell
./logstash-exporter  -n tower-inventorydata -t 10 -l 100
2020/10/06 13:02:50 Required environment variables are not set: $TOWER_HOST, $TOWER_USER, $TOWER_PASSWORD

$ export TOWER_HOST="YOUR-ANSIBLE-TOWER-HOST"
$ export TOWER_USER="your-service-id@somemail.com"
$ export TOWER_PASSWORD="YOUR-SECRET-PASSWORD"

$ ./logstash-exporter -n tower-inventorydata -t 10 -l 100
[
    {
        "account_name": "TestCustomer1",
        "country": "Netherlands",
        "geo": "EU",
        "market": "To be defined",
        "sector": "SampleSector",
        "type": "SampleType"
    },
    {
        "account_name": "TestCustomer2",
        "country": "Brazil",
        "geo": "EU",
        "market": "To be defined",
        "sector": "Retail",
        "type": "SampleType"
    },
    {
        "account_name": "TestCustomer3",
        "country": "France",
        "geo": "EU",
        "market": "France",
        "sector": "Oil&Gas",
        "type": "SampleType"
    }
]

```

3. Confirm that the data presented by the exporter matches the inventory metadata on your Ansible Tower instance.