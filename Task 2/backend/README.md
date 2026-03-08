# Meal Headcount Planner Backend (Serverless)

This backend is deployed as a serverless API using:
- AWS API Gateway
- AWS Lambda (Go custom runtime)
- AWS DynamoDB

---

## Prerequisites
- AWS CLI configured (`aws configure`)
- Go installed
- PowerShell or shell environment for building and zipping Go binaries

---

## Deploy to AWS

Each Lambda function is built and zipped separately. Example for the `health-lambda`:

```cmd
cd health-lambda
set GOOS=linux
set GOARCH=amd64
go build -o bootstrap
powershell Compress-Archive -Path bootstrap -DestinationPath function.zip -

```

## Test

### 1. Health Check

```cmd

curl https://pr807w8a23.execute-api.ap-south-1.amazonaws.com/default/health

```

Expected Response:
```
OK
```