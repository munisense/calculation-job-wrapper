# Munisense Calculation Job Wrapper

This tool can be used by Munisense partners that integrate on our queues to handle calculation jobs.
It can read calculation jobs from either the message queue or static files, and presents them via a HTTP endpoint.
Responses on the calculation jobs can be posted back via HTTP.

## Installation

- Download the most recent release for your platform
- Save the config file that you got from your Munisense contact in the same directory
- Open the terminal
- Make sure that the calculation-job-wrapper has executable permissions
```sh
chmod ug+x calculation-job-wrapper
```

## Running the wrapper
There are two modi in which you can run the tool. Either on live jobs, or on static testfiles. 

Run the application from the terminal:
```sh
./calculation-job-wrapper
```

When started it opens a HTTP endpoint and remains running until you close the application. You use that HTTP endpoint to request new jobs and to post responses.

### Startup options

| Option | Description | Example |
| ------ | ------ | ------ |
| --file | Reads job input from file instead of from queue. It looks for the input files in the static/ directory and it will store anything sent to output to the /output directory. | example.input |
| --config | Supply alternative config file. Default is config.json | config_alt.json |

### Process jobs from your application
When the wrapper is active a HTTP endpoint is enabled on which you can request new jobs.
```HTTP
GET http://localhost:8765/input
```

You will get a JSON response which will contain a **correlation_id**. Use this correlation_id to POST back your JSON response:
```HTTP
POST http://localhost:8765/output/{correlationId}
```
The exact format of the input and output structures are partner and solution specific.

Added in v0.1.5:
- Alternatively you can supply the response as a base64 encoded query parameter named *response*.
```HTTP
GET http://localhost:8765/output/{correlationId}?response={base64-encoded-response}
```