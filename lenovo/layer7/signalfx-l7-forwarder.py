######################################################################
# Broadcom Layer 7 API Gateway v9.4+ metrics forwarder to SignalFx
# Custom forwarder developed for Lenovo
# khymers@splunk.com
# v1 - initial release
#######################################################################

from SimpleHTTPServer import SimpleHTTPRequestHandler
import SocketServer
import argparse
import logging
from logging.handlers import RotatingFileHandler
import ConfigParser
import json
import requests
import threading
import re
import sys


def process_args():
    parser = argparse.ArgumentParser(description='Splunk SignalFx Layer7 Metrics Forwarder')
    parser.add_argument('-c','--config', help='forwarder configuration file', required=True)
    parser.add_argument('-t','--token', help='signalfx access token', required=True)
    return parser.parse_args()


# Apply Configuration
args = process_args()
config = ConfigParser.RawConfigParser()
config.read(args.config)

port = config.get('Server','port')

realm = config.get('SignalFx','realm')
sfx_token = args.token
service = config.get('SignalFx', 'service')
version = config.get('SignalFx', 'version')
env = config.get('SignalFx', 'env')

log_file = config.get('Logging','file')
log_level = getattr(logging, config.get('Logging', 'level').upper())
dopost = bool(config.get('SignalFx','dopost'))

if log_file == "stdout":
    log_handler = logging.StreamHandler(sys.stdout)
else:
    log_handler = RotatingFileHandler(log_file, maxBytes=5*1024*1024, backupCount=2)
log_formatter = logging.Formatter('%(asctime)s - %(levelname)s - %(message)s')
log_handler.setFormatter(log_formatter)

logger = logging.getLogger(__name__)
logger.setLevel(log_level)
logger.addHandler(log_handler)

logger.info("Loaded config(" + "service=" + service + ",version=" + version + ",realm=" + realm + ",env=" + env + ")")


# Post metric data to SignalFx
class SFxHandler():
    def __init__(self, r,t):
        self.realm = r
        self.token = t
        self.endpoint_dp = 'https://ingest.'+r+'.signalfx.com/v2/datapoint'
        self.endpoint_tags = 'https://api.'+r+'.signalfx.com/v2/dimension'
        self.headers = {
            'Content-Type': 'application/json',
            'X-SF-TOKEN': t
         }


    def post_dp(self,dp):
        logger.debug("Sending datapoint - " + str(dp))
        rs = requests.post(self.endpoint_dp,headers=self.headers, json=dp)
        logger.debug("Responses:" + str(rs))
        if rs.status_code != 200:
            logger.error("HTTP Error:" + str(rs.status_code) + " " + str(rs.reason) + "; Request:" + str(dp))


    def put_tags(self, tags, dk, dv):
        logger.debug("Sending datapoint - " + str(tags))
        rs = requests.put(self.endpoint_tags + '/' + dk + '/' + dv, headers=self.headers, json=tags)
        logger.debug("Responses:" + str(rs))
        if rs.status_code != 200:
            logger.error("HTTP Error:" + str(rs.status_code) + " " + str(rs.reason) + "; Request:" + str(tags))


# Convert Layer7 v 9.4 APM Metrics to SignalFx datapoints
def front_end_avg_response_time(host, service_uri, value):
    return {"gauge": [{
        "metric": "l7.avg_resp_time.ms",
        "dimensions": {"host": host, "service_uri": service_uri, "type": "frontend", "environment": env},
        "value": value}
    ]}


def back_end_avg_response_time(host, service_uri, value):
    return {"gauge": [{
        "metric": "l7.avg_resp_time.ms",
        "dimensions": {"host": host, "service_uri": service_uri, "type": "backend", "environment": env},
        "value": value}
    ]}


def success_count(host, service_uri, value):
    return {"counter": [{
        "metric": "l7.request.success_count",
        "dimensions": {"host": host, "service_uri": service_uri, "environment": env},
        "value": 1 if value == "true" else 0}
    ]}


def policy_violation_count(host, service_uri, value):
    return {"counter": [{
        "metric": "l7.request.policy_violation_count",
        "dimensions": {"host": host, "service_uri": service_uri, "environment": env},
        "value": 1 if value == "true" else 0}
    ]}


def routing_failure_count(host, service_uri, value):
    return {"counter": [{
        "metric": "l7.request.routing_failure_count",
        "dimensions": {"host": host, "service_uri": service_uri, "environment": env},
        "value": 1 if value == "true" else 0}
    ]}

def total_requests(host, service_uri, value):
    return {"counter": [{
        "metric": "l7.request.count",
        "dimensions": {"host": host, "service_uri": service_uri,"environment": env},
        "value": value}
    ]}


def get_sfx_json(metric_name, host, service_uri, value):
    return {
       "Total Frontend Latency" : front_end_avg_response_time(host, service_uri, value),
       "Total Backend Latency" : back_end_avg_response_time(host, service_uri, value),
       "Success Count" : success_count(host, service_uri, value),
       "Total Requests" : total_requests(host, service_uri, value),
       "Policy Violation Count" : policy_violation_count(host, service_uri, value),
       "Routing Failure Count" : routing_failure_count(host, service_uri, value)
     }[metric_name]


def get_sfx_tags(dv,dk):
    return {
        "key": dv,
        "value": dk,
        "tags": [],
        "customProperties": {
            "namespace": service,
            "version": version
        }
    }


def l72sfx(metric_data):
    data = json.loads(metric_data)
    sfx_datapoints = []
    host = ''
    for k,v in data.items():
        host = v["nodeName"]
        sfx_datapoints.append(get_sfx_json("Total Frontend Latency", v["nodeName"], v["serviceUri"], v["totalFrontendLatency"]))
        sfx_datapoints.append(get_sfx_json("Total Backend Latency", v["nodeName"], v["serviceUri"], v["totalBackendLatency"]))
        sfx_datapoints.append(get_sfx_json("Success Count", v["nodeName"], v["serviceUri"], v["isPolicySuccessful"]))
        sfx_datapoints.append(get_sfx_json("Total Requests", v["nodeName"], v["serviceUri"], 1))
        sfx_datapoints.append(get_sfx_json("Policy Violation Count", v["nodeName"], v["serviceUri"], v["isPolicyViolation"]))
        sfx_datapoints.append(get_sfx_json("Routing Failure Count", v["nodeName"], v["serviceUri"], v["isRoutingFailure"]))
    return sfx_datapoints, host


sfxh = SFxHandler(realm,sfx_token)


class ServerHandler(SimpleHTTPRequestHandler):
    def do_POST(self):
        content_len = int(self.headers.getheader('content-length', 0))
        post_body = self.rfile.read(content_len)
        self.send_response(200)
        self.end_headers()

        dps, host = l72sfx(post_body)
        if not dps:
            logger.debug("Request received - no APM data to forward.")
        else:
            if dopost:
                for dp in dps:
                    thread = threading.Thread(target=sfxh.post_dp, args=[dp])
                    thread.start()
                sfxh.put_tags(get_sfx_tags("host",host),"host",host)


def run(port):
    logger.info("Serving at port " + str(port))
    httpd = SocketServer.TCPServer(("", port), ServerHandler)
    httpd.serve_forever()


if __name__ == "__main__":
    run(int(port))

