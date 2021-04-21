# env-replace

A simple program to process a file with an environment file to replace the variables inline. Outputs to standard output to allow easy piping to other commands.

## Motivation

While beginning to move from Docker Compose centered application deploys to Kubernetes deployment centered ones, it became frustrating that environment variables did not seem as easy to utilize anymore. The process to use environment variables with deployment files felt cumbersome (at least from how I tried to do it) and I did not find any solution that fit what I was trying to do, so here we are.

## Usage

It is assumed a the file you want to substitute value into and a `.env` file exists where you are running the command.

`env-replace <file_name>`

### Example file to have values replaced
```
apiVersion: v1
kind: Service
metadata:
  name: ${SERVICE_NAME}
spec:
  ports:
    - protocol: TCP
      port: ${SERVICE_PORT}
```

### Example .env file
```
SERVICE_NAME=test-service
SERVICE_PORT=8080
```