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

### 3. Team Lead - Set Meal Status

#### via discord bot ->

- Use the `/team-meal set` Discord slash command to update a team member's meal participation status.

**Command format:**
```
/team-meal set employee:@TeamMember date:YYYY-MM-DD meal_type:Lunch/Snacks status:Yes/No
```
**Example:**
```
/team-meal set employee:@JohnDoe date:2026-03-19 meal_type:Lunch status:No
```

**Expected Response:**
```
❌ @JohnDoe has been opted out of Lunch on 2026-03-19.
```

#### via direct api call ->

```cmd
curl -X POST https://pr807w8a23.execute-api.ap-south-1.amazonaws.com/default/meal/team/set -H "Content-Type: application/json" -d "{\"team_lead_discord_id\": \"your_team_lead_discord_id\", \"target_discord_id\": \"employee_discord_id\", \"date\": \"2026-03-19\", \"meal_type\": \"lunch\", \"status\": \"NO\"}"
```

Expected Response:

```
{"message":"meal status updated successfully"}
```

### 4. Team Lead - View Meal Status

#### via discord bot ->

- Use the `/team-meal view` Discord slash command to view meal participation for all members of your team on a specific date.

**Command format:**
```
/team-meal view date:YYYY-MM-DD
```
**Example:**
```
/team-meal view date:2026-03-19
```

**Expected Response:**
```
Team Meal Status — 2026-03-19
@JohnDoe    ✅ Lunch: YES  |  ✅ Snacks: YES
@JaneDoe    ❌ Lunch: NO   |  ✅ Snacks: YES
```

#### via direct api call ->

```cmd
curl -X POST https://pr807w8a23.execute-api.ap-south-1.amazonaws.com/default/meal/team/view -H "Content-Type: application/json" -d "{\"team_lead_discord_id\": \"your_team_lead_discord_id\", \"date\": \"2026-03-19\"}"
```

Expected Response:

```
{"date":"2026-03-19","team_id":"TEAM#engineering","members":{"employee_discord_id":{"lunch":"YES","snacks":"YES"}}}
```

