# Mongo APP for grafana's simpe-json plugin

## Installation

```
go get github.com/m-masataka/grafana-simplejson-mongo
```

## Usage

Import package
```
import "github.com/m-masataka/grafana-simplejson-mongo/api"
```

Set up the parameter to ``api.Config``.

```
	conf := api.Config{
		Port: 8080,              // Server port
		MongoHost: "localhost",  // MongoDB Connection Host
	}
```

Start Http Server.
```
	errs := make(chan error, 2)
	api.StartHTTPServer(conf, errs)
	for {
		err := <-errs
		log.Println(err)
	}
```

### Connect from Grafana

1. Set up Grafana [simple-json-datasource](https://github.com/grafana/simple-json-datasource) plugin.
2. Add DataSource  with Grafana UI, and connect to this application.
3. Add [Dashbord](http://docs.grafana.org/guides/getting_started/), and [Edit Query](#Support)

### Support Query {#Support}

#### TimeSeries Query
Select ``timeseries`` Query, and edit Target.
```
[Database name].[Collection name].{[value column(Y-axis) name],[time column(X-axis) name]}
```
Example

```
fluentd.memory.{memory,time}
```


You can also add a match query on a field
```
[Database name].[Collection name].{[value column(Y-axis) name],[time column(X-axis) name],[field name],[value to match]}
```
Example

```
fluentd.memory.{memory,time,server,JP}
```


#### Table Query
Select ``table`` Query, and edit Target.
```
[Database name].[Collection name]
```

Example

```
fluentd.memory
```
