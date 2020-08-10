package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
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
const RequestSize = "l7.req_size.bytes"
const ResponseSize = "l7.res_size.bytes"
const SuccessCount = "l7.request.success_count"
const TotalCount = "l7.request.count"

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

	var result = make(map[string](interface{}))
	json.NewDecoder(r.Body).Decode(&result)

	sfx.dps = append(sfx.dps, SignalFxProcess(result)...)

}

func makeDatapoint(metric_name string, dims []string, value int64) *datapoint.Datapoint {
	dimensions := map[string]string{
		"host":        dims[0],
		"service_uri": dims[len(dims)-2],
		"environment": config.AppEnvironment,
	}

	switch metric_name {
	case "Front End Average Response Time (ms)":
		dimensions["type"] = "frontend"
		return sfxclient.Gauge(FrontEndReponseTime, dimensions, value)
	case "Back End Average Response Time (ms)":
		dimensions["type"] = "backend"
		return sfxclient.Gauge(BackEndReponseTime, dimensions, value)
	case "Request size (bytes)":
		return sfxclient.Gauge(RequestSize, dimensions, value)
	case "Response size (bytes)":
		return sfxclient.Gauge(ResponseSize, dimensions, value)
	case "Success Count":
		return sfxclient.Counter(SuccessCount, dimensions, value)
	case "Total Requests":
		return sfxclient.Counter(TotalCount, dimensions, value)
	default:
		return nil
	}
}

func SignalFxProcess(r map[string](interface{})) []*datapoint.Datapoint {

	metrics := r["metrics"].([]interface{})
	datapoints := make([]*datapoint.Datapoint, 0)
	for _, v := range metrics {
		data := v.(map[string]interface{})
		name := fmt.Sprintf("%v", data["name"])

		if strings.Contains(name, "Services") {
			re := regexp.MustCompile(`[\w\-.*/() ]+`)
			dims := re.FindAllString(name, -1)

			metric_name := dims[len(dims)-1]

			v_str, _ := data["value"].(string)
			v_int, _ := strconv.Atoi(v_str)

			dp := makeDatapoint(metric_name, dims, int64(v_int))
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
