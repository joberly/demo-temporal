# demo-temporal

Demo application built with Temporal in Golang

## Description

The demo uses Temporal to manage an image processing workflow.

The image processing workflow takes an uploaded image, ensures it's of a
valid image type, converts it to grayscale, and makes the resulting image
available for download.

To access this workflow, an API is provided on localhost port 8081. This
API has an /upload path to which an image may be uploaded. This starts the
workflow running. The /status/:workflowID/run:runID queries the status of
the workflow given information returned from the upload. Once complete,
the status will read "processing complete". At that point, the result image
can be downloaded from the /download/:imageID API path.

## Running

1. Clone [github.com/temporalio/docker-compose](https://github.com/temporalio/docker-compose) in a separate directory.
2. Change to the directory with the above repository and run `docker compose up -d`.
3. Change back to this repository's directory and run `docker compose up -d`.
4. Check that the API is up and running by curling its health check.
   It should return successfully with a status of "ok".
   ```
   $ curl http://localhost:8081/health
   {"status":"ok"}
   ```

## Usage

The following contains examples. The imageId, workflowId, and runId is 
randomly generated for each image uploaded.

### Upload Image

```
$ curl -X POST -F "file=@test1.webp" http://localhost:8081/upload
{"imageId":"79839d04-5dd1-47a9-a2c6-ba91bb7edbb1","message":"file uploaded","runId":"6c2a3179-6dc8-4ddc-919a-3eb1fa6c58a6","workflowId":"79839d04-5dd1-47a9-a2c6-ba91bb7edbb1"}
```

### Get Image Processing Status

```
$ curl http://localhost:8081/status/79839d04-5dd1-47a9-a2c6-ba91bb7edbb1/run/6c2a3179-6dc8-4ddc-919a-3eb1fa6c58a6
{"error":"","runId":"6c2a3179-6dc8-4ddc-919a-3eb1fa6c58a6","status":"converting image to grayscale","workflowId":"79839d04-5dd1-47a9-a2c6-ba91bb7edbb1"}
```

### Download Processed Image

Open `http://localhost:8081/download/<imageId>` with your browser, replacing the `<imageId>` with your imageId returned from the upload.

## Notes

1. Using a managed service like S3 to handle image uploading would move the
   network bandwidth for images to that service so that the api does not
   need to bear the burden. The api could generate a signed S3 URL for to
   POST the image and add an upload to the database. The frontend would take
   the signed URL from the backend to upload the image directly to S3. The 
   backend could have some kind of process that handles S3 notifications and
   start the image processing workflow once the image is fully uploaded to S3.
2. Everything needs unit tests badly. Definitely needs error path testing. :)
3. API needs standardization. Same for the workflow status that gets returned
   from the API status endpoint.
4. There's no auth so please don't run this publicly. The image ID probably 
   isn't even really large enough to make it hard to guess.
5. There are no neat Grafana dashboards for service status.
6. The worker process does not provide an endpoint for Prometheus to scrape.
