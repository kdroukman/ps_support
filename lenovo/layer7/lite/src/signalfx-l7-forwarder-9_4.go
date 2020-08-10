package main

import (
	"context"
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/signalfx/golib/v3/datapoint"
	"github.com/signalfx/golib/v3/sfxclient"
	yaml "gopkg.in/yaml.v2"
)

/////////////////////////////////////////////////////
//             Metric Names
/////////////////////////////////////////////////////

const FrontEndReponseTime = "l7.avg_resp_time.ms"
const BackEndReponseTime = "l7.avg_resp_time.ms"
const SuccessCount = "l7.request.success_count"
const TotalCount = "l7.request.count"
const PolicyViolationCount = "l7.request.policy_violation_count"
const RoutingFailureCount = "l7.request.routing_failure_count"

/////////////////////////////////////////////////////
//             Load Configuration
/////////////////////////////////////////////////////

type conf struct {
	ListenAddress       string    `yaml:"listenAddress"`
	SignalFxAccessToken string    `yaml:"signalFxAccessToken"`
	SignalFxRealm       string    `yaml:"signalFxRealm" default:"us1"`
	AppName             string    `yaml:"appName"`
	AppVersion          string    `yaml:"appVersion"`
	AppEnvironment      string    `yaml:"appEnvironment"`
	IntervalSeconds     string    `yaml:"intervalSeconds"`
	Logging             LogConfig `yaml:"logging" default:"{}"`
}

type LogConfig struct {
	// Valid levels include `debug`, `info`, `warn`, `error`.  Note that
	// `debug` logging may leak sensitive configuration (e.g. passwords)
	Level string `yaml:"level" default:"info"`
}

func (c *conf) getConf(f string) *conf {
	yamlFile, err := ioutil.ReadFile(f)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

var configFile = flag.String("config", "/etc/signalfx-l7-forwarder/config.yaml", "Specify the location of SignalFx Layer 7 Configuration file")
var config conf

/////////////////////////////////////////////////////
//             Setup Logging
/////////////////////////////////////////////////////

var (
	Debug   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func initLogging(
	debugHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Debug = log.New(debugHandle,
		"DEBUG: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

var debug = false

/////////////////////////////////////////////////////
//             SignalFx Datapoint Collection
/////////////////////////////////////////////////////

var scheduler = sfxclient.NewScheduler()
var sfx SignalFxCollector

type SignalFxCollector struct {
	dps []*datapoint.Datapoint
}

func (c *SignalFxCollector) Datapoints() []*datapoint.Datapoint {
	if debug {

		Debug.Println("Scheduler sending ", len(c.dps), " datapoints")
	}
	datapoints := c.dps
	c.dps = nil
	return datapoints
}

func ServerHandler(w http.ResponseWriter, r *http.Request) {

	var result map[string]interface{}
	json.NewDecoder(r.Body).Decode(&result)

	sfx.dps = append(sfx.dps, SignalFxProcess(result)...)

}

func makeDatapoint(metric_name string, host string, service_uri string, value string) *datapoint.Datapoint {
	dimensions := map[string]string{
		"host":        host,
		"service_uri": service_uri,
		"environment": config.AppEnvironment,
	}
	var v int64
	switch metric_name {
	case "totalFrontendLatency":
		v, _ := strconv.Atoi(value)
		dimensions["type"] = "frontend"
		return sfxclient.Gauge(FrontEndReponseTime, dimensions, int64(v))
	case "totalBackendLatency":
		v, _ := strconv.Atoi(value)
		dimensions["type"] = "backend"
		return sfxclient.Gauge(BackEndReponseTime, dimensions, int64(v))
	case "isPolicySuccessful":
		if value == "true" {
			v = 1
		} else {
			v = 0
		}
		return sfxclient.Counter(SuccessCount, dimensions, v)
	case "totalRequests":
		v, _ := strconv.Atoi(value)
		return sfxclient.Counter(TotalCount, dimensions, int64(v))
	case "isPolicyViolation":
		if value == "true" {
			v = 1
		} else {
			v = 0
		}
		return sfxclient.Counter(PolicyViolationCount, dimensions, v)
	case "isRoutingFailure":
		if value == "true" {
			v = 1
		} else {
			v = 0
		}
		return sfxclient.Counter(RoutingFailureCount, dimensions, v)
	default:
		return nil
	}
}

func SignalFxProcess(r map[string](interface{})) []*datapoint.Datapoint {

	datapoints := make([]*datapoint.Datapoint, 0)
	metrics := r["request"].(map[string]interface{})

	// Add total requests counter
	dp := makeDatapoint("totalRequests", metrics["nodeName"].(string), metrics["serviceUri"].(string), "1")
	dp.SetProperty("namespace", config.AppName)
	dp.SetProperty("version", config.AppVersion)
	datapoints = append(datapoints, dp)

	for k, v := range metrics {

		dp := makeDatapoint(k, metrics["nodeName"].(string), metrics["serviceUri"].(string), v.(string))
		// Check that the key/value pair was a valid datapoint
		if dp != nil {
			dp.SetProperty("namespace", config.AppName)
			dp.SetProperty("version", config.AppVersion)

			datapoints = append(datapoints, dp)

			if debug {

				Debug.Println("Datapoint ", dp.String(), " with Properties ", dp.GetProperties())
			}

		}
	}
	if debug {

		Debug.Println(len(datapoints), "datapoints generated from request")
	}
	return datapoints
}

/////////////////////////////////////////////////////
//             Main Runtime
/////////////////////////////////////////////////////

func init() {
	// add extra options
	flag.StringVar(configFile, "c", "/etc/signalfx-l7-forwarder/config.yaml", "Specify the location of SignalFx Layer 7 Configuration file")
	flag.Parse()

	config.getConf(*configFile)

	// set up Logging
	initLogging(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
	Info.Println("Config loaded - ", *configFile)
	Info.Println("listenAddress: ", config.ListenAddress)
	Info.Println("realm: ", config.SignalFxRealm)
	Info.Println("service: ", config.AppName)
	Info.Println("version: ", config.AppVersion)
	Info.Println("environment: ", config.AppEnvironment)
	Info.Println("intervalSeconds: ", config.IntervalSeconds)
	if config.Logging.Level == "debug" {
		Info.Println("Debug logging configured")
		debug = true
	}

	// set up scheduler
	// endpoint configuration examples
	scheduler.Sink.(*sfxclient.HTTPSink).DatapointEndpoint = "https://ingest." + config.SignalFxRealm + ".signalfx.com/v2/datapoint"
	scheduler.Sink.(*sfxclient.HTTPSink).EventEndpoint = "https://ingest." + config.SignalFxRealm + ".signalfx.com/v2/event"
	scheduler.Sink.(*sfxclient.HTTPSink).TraceEndpoint = "https://ingest." + config.SignalFxRealm + ".signalfx.com/v1/trace"
	scheduler.Sink.(*sfxclient.HTTPSink).AuthToken = config.SignalFxAccessToken
	interval, _ := time.ParseDuration(config.IntervalSeconds)
	scheduler.ReportingDelay(interval)
	if debug {
		scheduler.Debug(true)
		Debug.Println("Reporting interval set to ", interval)
	}
}

func main() {
	scheduler.AddCallback(&sfx)
	go scheduler.Schedule(context.Background())

	http.HandleFunc("/", ServerHandler)
	http.ListenAndServe(config.ListenAddress, nil)

}
