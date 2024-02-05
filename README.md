# demo-temporal
Demo application built with Temporal in Golang

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

## Testing

### Image Upload

```
> curl -X POST -F "file=@path/to/image" http://localhost:8081/upload
```

## Notes

1. Using a managed service like S3 to handle image uploading would move the
   network bandwidth for images to that service so that the api does not
   need to bear the burden. The api could generate a signed S3 URL for to
   POST the image and add an upload to the database. The frontend would take
   the signed URL from the backend to upload the image directly to S3. The 
   backend could have some kind of process that handles S3 notifications and
   start the image processing workflow once the image is fully uploaded to S3.
