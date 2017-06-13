## Summary

The influxdb-firehose-nozzle is a CF component which forwards metrics from the Loggregator Firehose to [influxdb](http://github.com/influxdata/influxdb)

### Configure CloudFoundry UAA for Firehose Nozzle

The InfluxDB firehose nozzle requires a UAA user who is authorized to access the loggregator firehose. You can add a user by editing your CloudFoundry manifest to include the details about this user under the properties.uaa.clients section. For example to add a user `influxdb-firehose-nozzle`:

```yaml
properties:
  uaa:
    clients:
      influxdb-firehose-nozzle:
        access-token-validity: 1209600
        authorized-grant-types: authorization_code,client_credentials,refresh_token
        override: true
        secret: influxdb-firehose-nozzle
        scope: openid,oauth.approvals,doppler.firehose,scim.write
        authorities: oauth.login,doppler.firehose
```

If you are using the `uaac` command line tool to create the user it might look something like:

```
uaac client \
  add influxdb-firehose-nozzle \
  --scope openid,oauth.approvals,doppler.firehose \
  --authorized_grant_types authorization_code,client_credentials,refresh_token \
  --authorities oauth.login,doppler.firehose,scim.write \
  -s influxdb-firehose-nozzle
```

### Running

The InfluxDB nozzle uses a configuration file to obtain the firehose URL, InfluxDB API credentials and other configuration parameters. The firehose and the InfluxDB servers both require authentication -- the firehose requires a valid CloudFoundry username/client-secret and InfluxDB requires authentication credentials.

You can start the firehose nozzle by executing:
```
go run main.go -config config/influxdb-firehose-nozzle.json"
```

### Firehose Event Type Filtering

The firehose exports four main event types: ContainerMetric, CounterEvent, HttpStartStop, ValueMetric.  The event types the nozzle will export can be selected via the json configuration file or via enviornment variables.  

If you plan to export HTttpStartStop or ContainerMetrics, you will need to specify the path to a secondary service which exports information about the applications in cloud foundry which will be included as tags to the influxdb data points.  Specifically the service is expected to export json in the following format:
```
[{
  "name": "app1",
  "guid": "d5697f98-1a94-4d92-a93b-6ba812c9f67a",
  "space": "testing",
  "org": "myorg"
}, {
  "name": "app2",
  "guid": "164686ca-0c4f-42b7-95d3-3b1ee0fb095c",
  "space": "testing",
  "org": "myorg"
}]
```
If you do not have a centralized service for providing this information in your environment already, a sample application is included in this repository [for you to use](app-api-example).  

### Batching

The configuration file specifies the interval at which the nozzle will flush metrics to influxdb. By default this is set to 15 seconds.

### `slowConsumerAlert`
For the most part, the influxdb-firehose-nozzle forwards metrics from the loggregator firehose to influxdb without too much processing. A notable exception is the `influxdb.nozzle.slowConsumerAlert` metric. The metric is a binary value (0 or 1) indicating whether or not the nozzle is forwarding metrics to influxdb at the same rate that it is receiving them from the firehose: `0` means the the nozzle is keeping up with the firehose, and `1` means that the nozzle is falling behind.

The nozzle determines the value of `influxdb.nozzle.slowConsumerAlert` with the following rules:

1. **When the nozzle receives a `TruncatingBuffer.DroppedMessages` metric, it publishes the value `1`.** The metric indicates that Doppler determined that the client (in this case, the nozzle) could not consume messages as quickly as the firehose was sending them, so it dropped messages from its queue of messages to send.

2. **When the nozzle receives a websocket Close frame with status `1008`, it publishes the value `1`.** Traffic Controller pings clients to determine if the connections are still alive. If it does not receive a Pong response before the KeepAlive deadline, it decides that the connection is too slow (or even dead) and sends the Close frame.

3. **Otherwise, the nozzle publishes `0`.**

### Tests

You need [ginkgo](http://onsi.github.io/ginkgo/) to run the tests. The tests can be executed by:
```
ginkgo -r

```

## Deploying

### [Bosh](http://bosh.io)

There is a bosh release that will configure, start and monitor the influxdb nozzle:
[https://github.com/cloudfoundry-incubator/influxdb-firehose-nozzle-release](https://github.com/cloudfoundry-incubator/influxdb-firehose-nozzle-release
)

### Configuration 

Any of the configuration parameters can be overloaded by using environment variables. The following
parameters are supported

| Environment variable          | Description            |
|-------------------------------|------------------------|
| NOZZLE_UAAURL                 | UAA URL which the nozzle uses to get an authentication token for the firehose |
| NOZZLE_CLIENT                 | Client who has access to the firehose |
| NOZZLE_CLIENT_SECRET          | Secret for the client |
| NOZZLE_TRAFFICCONTROLLERURL   | Loggregator's traffic controller URL |
| NOZZLE_FIREHOSESUBSCRIPTIONID | Subscription ID used when connecting to the firehose. Nozzles with the same subscription ID get a proportional share of the firehose |
| NOZZLE_INFLUXDB_URL           | The influxdb API URL |
| NOZZLE_INFLUXDB_DATABASE      | The database name used when publishing metrics to influxdb |
| NOZZLE_INFLUXDB_USER          | The username name used when publishing metrics to influxdb |
| NOZZLE_INFLUXDB_PASSWORD      | The password name used when publishing metrics to influxdb |
| NOZZLE_METRICPREFIX           | The metric prefix is prepended to all metrics flowing through the nozzle |
| NOZZLE_DEPLOYMENT             | The deployment name for the nozzle. Used for tagging metrics internal to the nozzle |
| NOZZLE_DEPLOYMENT_FILTER      | If set, the nozzle will only send metrics with this deployment name |
| NOZZLE_EVENT_FILTER           | If set, the nozzle will only send metrics from these event types (Ex. (ContainerMetric, CounterEvent, HttpStartStop, ValueMetric) |
| NOZZLE_APP_API_URL            | Url of the service which provides a list of cf apps and their information. Only used if ContainerMetric or HttpStartStop events are enabled |
| NOZZLE_FLUSHDURATIONSECONDS   | Number of seconds to buffer data before publishing to influxdb |
| NOZZLE_INSECURESSLSKIPVERIFY  | If true, allows insecure connections to the UAA and the Trafficcontroller |
| NOZZLE_DISABLEACCESSCONTROL   | If true, disables authentication with the UAA. Used in lattice deployments |