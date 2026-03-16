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

### 2. Set Meal Status


#### via discord bot ->

- Use the `/meal view {date}` Discord slash command to check meal participation status for a specific date.

**Command format:**
```
 /meal set date:YYYY-MM-DD meal_type:Lunch/Snacks status:Yes/No
```
**Example:**
```
/meal view date:2026-03-10 meal_type:Lunch status:No
```

**Expected Response:**
```
{"You have opted out of Lunch on 2026-03-13."}
```

#### via direct api call ->

```cmd

curl -X POST https://pr807w8a23.execute-api.ap-south-1.amazonaws.com/default/meal/participation/set ^
  -H "Content-Type: application/json" ^
  -d "{\"discord_id\": \"123456789\", \"date\": \"2026-03-13\", \"meal_type\": \"lunch\", \"status\": \"NO\"}"

```

Expected Response:

```
{"message":"meal status updated successfully"}

```



