# Technical Design Document

# Meal Headcount Planner

**UI:** Discord   
**Backend:** AWS Lambda  
**Database:** Amazon DynamoDB  

---

# 1. System-Level Design

## 1.1 High-Level Architecture
Discord User  
↓  
Discord Bot (Slash Commands / Interactions)  
↓   
API Gateway  
↓  
AWS Lambda (Business Logic)  
↓  
DynamoDB (Single Table Design)


### Components

#### Discord Bot
- Handles slash commands
- Sends interaction payloads to API Gateway
- Displays formatted responses

#### API Gateway
- Public endpoint for Discord interactions
- Verifies Discord signatures
- Routes to Lambda

#### AWS Lambda
- Stateless business logic
- Role validation
- Cutoff validation
- Team-based access validation
- Reads/writes to DynamoDB

#### DynamoDB
- Stores users, teams, meals, work locations, special days
- Uses single-table design for scalability

---

## 1.2 DynamoDB Design (Single Table)

**Table Name:** `MHP`

### Primary Key
- `PK` (Partition Key)
- `SK` (Sort Key)

### Entity Patterns

| Entity        | PK            | SK                   | Notes                |
|--------------|--------------|----------------------|----------------------|
| User         | USER#<id>     | META                 | role, teamId         |
| Meal         | MEAL#<date>   | USER#<id>#<mealType> | participation        |
| Work Location| WORK#<date>   | USER#<id>            | Office/WFH           |
| Day Config   | DAY#<date>    | META                 | holiday/closed/event |
| Team         | TEAM#<teamId> | META                 | team metadata        |


---

## 1.3 Role Model

- EMPLOYEE
- TEAM_LEAD
- ADMIN

Role is stored in the USER entity.

---

## 1.4 Cross-Cutting Rules

All features share:

- Date validation (no office-closed updates)
- Cutoff validation
- Role-based access control
- Last-write-wins strategy
- Strongly consistent reads when required

---

# 2. Feature-Level Designs

---

# Feature 1: Employee Meal Management

## Overview

Allows employees to manage meal participation. Limited to employee actions only.
- Employees can:
  - Select a date.
  - Select meal type.
  - Mark meal participation (Yes/No).
  - View their current meal status for a selected date.
- Updates should reflect immediately in the system.
- The system should validate:
  - Date selection (no invalid dates such as office closed).
  - Duplicate updates (latest update overwrites previous one).
  - Update restrictions (such as no update after cutoff time).
  - Role based access.

## Data Model
## DynamoDB Item Structure

**PK:** `MEAL#<date>`  
**SK:** `USER#<userId>#<mealType>`

### Attributes

- `participation`: YES / NO  
- `updatedAt`: <timestamp>


## User Flows

### Opt-In / Opt-Out

1. User runs `/meal set`
2. Bot collects:
   - Date
   - Meal Type
   - YES/NO
3. Lambda:
   - Validate role = EMPLOYEE
   - Check DAY table for closure
   - Check cutoff time
   - Update record
4. Bot confirms action

### View Status

1. User runs `/meal view`
2. Lambda queries:

**PK:** `MEAL#<date>`  
**SK:** begins_with `USER#<userId>`


3. Bot returns status

## Validations

- Closed day → reject
- After cutoff → reject
- Duplicate → overwrite

---

