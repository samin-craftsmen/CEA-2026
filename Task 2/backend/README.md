# Meal Headcount Planner Backend (Serverless)

This backend is deployed as a serverless API using:
- AWS API Gateway
- AWS Lambda (Go custom runtime)
- AWS DynamoDB

## Prerequisites
- AWS CLI configured (`aws configure`)
- AWS SAM CLI installed
- Go installed

## Deploy to AWS
From project root:

```bash
sam build
sam deploy --guided
```

Recommended answers for first deploy:
- Stack Name: `meal-headcount-backend`
- AWS Region: your target region (e.g. `us-east-1`)
- Confirm changes before deploy: `Y`
- Allow SAM CLI IAM role creation: `Y`
- Save arguments to configuration file: `Y`

For subsequent deploys:

```bash
sam build
sam deploy
```

## Test after Deployment
### 1) Get deployed endpoint outputs

```bash
aws cloudformation describe-stacks \
  --stack-name meal-headcount-backend \
  --query "Stacks[0].Outputs"
```

Copy `HealthEndpoint` value.

### 2) Verify health endpoint

```bash
curl <HealthEndpoint>
```

Expected response:

```json
{"status":"healthy"}
```

### 3) API Gateway live check (optional)
Open `<HealthEndpoint>` in browser and confirm HTTP 200 response.

## Notes for local development
- `DYNAMODB_ENDPOINT` can be set for local DynamoDB.
- In AWS deployment, `DYNAMODB_ENDPOINT` is left empty by default.
