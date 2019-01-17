# OPI API
* *This is a WIP.*
* Needs verification on the Eirini side

# Desire an App

## List apps

```
GET /apps
```

Request body:

`<empty>`

Response codes:

* `200`: All right
* `500`: Panic

Response body:

```json
{
    "error": "",
    "desired_lrp_scheduling_infos": [
        {
            "desired_lrp_key": {
                "process_guid": "App guid"
            }
		    "annotation": "timestamp"
        }
    ]
}
```

TODO: `process_guid` parameter is ignored

## Desire an app

```
PUT /apps/:process_guid
```

Request parameters:

* `process_guid`: Uniquely identifies the app to be desired

Request body:

```json
{
    "guid": "Process guid",
    "version": "App version provided by CC",
    "process_guid": "<guid>-<version>",
    "ports": [11, 42],
    "routes": {
        "TODO": "some way of describing a route"
    },
    "docker_image": "reserved for future use",
    "droplet_hash": "hash of the droplet",
    "droplet_guid": "GUID of the droplet",
    "start_command": "/bin/bash -c 'echo Hello World'",
    "environment": {
        "foo": "bar",
        "some": "thing"
    },
    "instances": 23,
    "last_updated": "timestamp",
    "health_check_type": "one of port, http, or process",
    "health_check_http_endpoint": "endpoint to use for health check",
    "health_check_timeout_ms": 4711,
    "memory_mb": 640
}
```

Response codes:

* `201`: App was successfully desired
* `400`: Could not desire app (JSON decoding error or scheduler problem)

Response body:

`<empty>`

TODO: If desiring the app failed on our side, it should not be 4xx, but a 5xx

## Update an app

```
POST /apps/:process_guid
```

Request parameters:

* `process_guid`: Uniquely identifies the app to be desired

Request body:

```json
{
    "process_guid": "<guid>-<version>",
    "instances": 42,
    "routes": {
        "string": "bytes" # TODO What exactly?
    },
    "annotation": "something",
    "guid": "Process guid",
    "version": "App version by CC"
}
```

Response codes:

* `200`: OK
* `400`: Could decode request body
* `500`: Updating the app failed

Response body:

`<empty>`

TODO: Check error after response writing https://github.com/cloudfoundry-incubator/eirini/blob/2adf2e6c59447747c9b6b9254f47d55c8b84530f/handler/app_handler.go#L109

## Stop an app

```
PUT /apps/:process_guid/:version_guid/stop
```

Request parameters:

* `process_guid`: Uniquely identifies the app to be desired
* `version_guid`: Version of the app

Request body:

`<empty>`

Response codes:

* `200`: OK
* `500`: Stopping the app failed

Response body:

`<empty>`

## Get instances of an app

```
GET /apps/:process_guid/:version_guid/instances
```

Request parameters:

* `process_guid`: Uniquely identifies the app to be desired
* `version_guid`: Version of the app

Request body:

`<empty>`

Response codes:

* `200`: OK
* `500`: Encoding the response failed or getting the instances failed

Response body:

```json
{
    "process_guid": "<guid>-<version>",
	"error": "some error message",
	"instances": [{
	       "index": 11,
	       "since": 42,
	       "state": "running"
	   }
	]
}
```

TODO: Check design of error response body

## Get an app

```
GET /apps/:process_guid/:version_guid
```

Request parameters:

* `process_guid`: Uniquely identifies the app to be desired
* `version_guid`: Version of the app

Request body:

`<empty>`

Response codes:

* `200`: OK
* `404`: App is not found
* `500`: Encoding the response failed or getting the instances failed

Response body:

```json
{
    "error": "",
        "desired_lrp": {
            "process_guid": "<guid>-<version>",
            "instances": 5
        }
}
```

TODO: check if process_guid is not garbage. I think that we currently return k8s-specific hash, while we should be returning <guid>-<version>

# Staging

## Stage an app

```
POST /stage/:staging_guid
```

Request parameters:

* `staging_guid`: Unique identifier for staging

Request body:

```json
{
    "app_id": "maybe",
    "log_guid": "",
    "lifecycle_data": "{ \"app_bits_download_uri\": \"uri\", \"droplet_upload_uri\": \"uri\", \"buildpacks\": [{ \"name\": \"some-name\", \"key\": \"some-key\", \"url\": \"some-url\", \"skip_detect\": true  }] }",
    "completion_callback": "call-me-maybe",
    "environment": [{
            "name": "foo",
            "value": "bar"
        }
    ]
}
```

Response codes:

* ``: 

Response body:

`<empty>`

TODO: 

## Mark staging an app complete

```
PUT /stage/:staging_guid/completed
```
Request parameters:

* ``: 

Request body:

```json
{
}
```

Response codes:

* ``: 

Response body:

`<empty>`

TODO: 

