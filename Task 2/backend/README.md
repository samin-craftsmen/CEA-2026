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
powershell Compress-Archive -Path bootstrap -DestinationPath function.zip 

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

### 2. View Meal Participation

```cmd

curl -X POST https://pr807w8a23.execute-api.ap-south-1.amazonaws.com/default/meal/participation/view ^
More?   -H "Content-Type: application/json" ^
More?   -d "{\"discord_id\": \"your_discord_id\"}"

```

Expected Response:

```
{"date":"2026-03-10","meals":{"lunch":"YES","snacks":"YES"},"user_id":"your_discord_id"}

```