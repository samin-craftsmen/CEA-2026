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
powershell Compress-Archive -Path bootstrap -DestinationPath function.zip -Force

```

## Test


### 1. View Meal Participation

#### via discord bot ->

- Use the `/meal view {date}` Discord slash command to check meal participation status for a specific date.

**Command format:**
```
 /meal view date:YYYY-MM-DD
```
**Example:**
```
/meal view date:2026-03-10
```

**Expected Response:**
```
{"date":"2026-03-10","meals":{"lunch":"YES","snacks":"YES"},"user_id":"your_discord_id"}
```

#### via direct api call ->

```cmd

curl -X POST https://pr807w8a23.execute-api.ap-south-1.amazonaws.com/default/meal/participation/view ^
  -H "Content-Type: application/json" ^
  -d "{\"discord_id\": \"your_discord_id\", \"date\": \"2026-03-10\"}"

```

Expected Response:

```
{"date":"2026-03-10","meals":{"lunch":"YES","snacks":"YES"},"user_id":"your_discord_id"}

```