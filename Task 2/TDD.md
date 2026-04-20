# Technical Design Document

# Meal Headcount Planner

**UI:** Discord, React   
**Backend:** AWS Lambda  
**Database:** Amazon DynamoDB  

---

# 1. System-Level Design
## 1.1 Summary

The Meal Headcount Planner is a system designed to manage daily employee meal participation. Employees can indicate whether they will participate in meals and whether they will work from the office or remotely. The system integrates with Discord for employee interactions and provides a web dashboard for administrators to manage configurations, view reports, and perform bulk operations.The backend is implemented using AWS serverless architecture to ensure scalability, high availability, and minimal operational overhead.

Technologies used:
- Discord Bot for employee interaction
- React Web Dashboard for administrative management
- AWS Lambda for business logic
- Amazon DynamoDB using a single-table design for scalable data access
- API Gateway for HTTP endpoints

Note: The first iteration of the project is developed locally. Iteration 2 will be built on the cloud.

## 1.2 Problem Statement

Currently, meal participation is managed manually(Excel Based System). This leads to:

- Inaccurate meal headcount
- Food wastage or shortages
- Lack of centralized tracking
- Difficulty managing holidays and office day status
- No clear visibility for administrators

The system aims to provide a centralized platform to manage meal participation and work locations while ensuring role-based access control and enforcing operational rules such as cutoff times and holiday restrictions.

## 1.3 High Level Architecture

#### Employee Interaction Flow

Discord User  
↓  
Discord Bot (Slash Commands / Interactions)  
↓   
API Gateway  
↓  
AWS Lambda (Business Logic)  
↓  
DynamoDB (Single Table Design)

#### Administrator Interaction Flow  
↓  
Admin User  
↓   
React Web Dashboard  
↓  
API Gateway  
↓  
Lambda Functions
↓
DynamoDB


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

#### React Admin Dashboard

The web dashboard provides administrative capabilities.

- Day management (holidays/events/office closed)
- Meal reports
- Employee management
- Bulk updates

The dashboard communicates with backend APIs exposed via API Gateway.

## 1.4 DynamoDB Design (Single Table)

**Table Name:** `MHP`

### Primary Key
- `PK` (Partition Key)
- `SK` (Sort Key)

The system stores multiple logical entities within a single DynamoDB table.

#### Entities include:

User  
Team  
Meal Participation  
Work Location  
Day Configuration  
Meal Configuration

### Entity Patterns

| Entity                                              | PK                | SK                     | Purpose                                                                                                                                     |
| --------------------------------------------------- | ----------------- | ---------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **User**                                            | `USER#<id>`       | `META`                 | Stores employee metadata: role, team. Essential for role-based access, team validation, and identifying users across meals and work locations.       |
| **Meal Participation**                              | `MEAL#<date>`     | `USER#<id>#<mealType>` | Represents whether a user participates in a meal (Yes/No). Needed to track daily meal headcount, allow employees to opt-in/out, and generate reports.     |
| **Work Location**                                   | `WORK#<date>`     | `USER#<id>`            | Stores where a user will work (Office/WFH) on a specific date. Necessary for linking attendance to meals and automating meal logic. |
| **Day Configuration**                               | `DAY#<date>`      | `META`                 | Stores special days like holidays, office closed, and events. Needed to enforce date validation and cutoff rules.                                       |
| **Team**                                            | `TEAM#<teamId>`   | `META`                 | Stores team metadata (name, lead, members). Will help in team based features like Team Lead meal management, and bulk operations and report generation.                           |
| **Meal Configuration**  | `CONFIG#MEALTYPE` | `<mealType>`           | Stores allowed meal types (Lunch, Snacks, Event wise options) and settings. Supports admin operations to manage meals and enforce constraints.                     |

### Note: How they are queried will be explained feature by feature. For example, for feature 1 we will need to use GSI between user and team entity.

#### DynamoDB Access Patterns

The DynamoDB schema is designed based on the following primary access patterns:

1. Get a user profile
2. Get meal participation for a user on a specific date
3. Get all meal participants for a specific date
4. Get work location for a user on a specific date
5. Get all users for a team
6. Get day configuration for a specific date
7. Update meal participation for a user

### Note: Access pattern can be modified based on new features based on which this will be updated. Also how these access patterns are executed is explained on feature level design. For feature one we will only need to deal with employees(only) meal participation.

---

## 1.5 Role Implementation 

### Role-Based Access Control

The system implements role-based access control using the USER entity.

#### Roles

EMPLOYEE  
TEAM_LEAD  
ADMIN  

#### Role Storage

User role information is stored in DynamoDB.

Example Item

PK: USER#123  
SK: META  

Attributes:

role: EMPLOYEE  
teamId: TEAM#10

#### Authorization Enforcement

All Lambda handlers perform role validation before executing operations.

Examples:

Employee operations:
- Only allow access to their own meal or location records.

Team Lead operations:
- Can modify records for users belonging to the same team.

Admin operations:
- Can modify any record in the system.

Authorization checks occur in Lambda before database operations.

#### Note: For iteration 1 lambda is not used. It will be built in interation 2.

## 1.6 Cross-Cutting Rules

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

## Entities

| Entity                 | Attributes (for this feature)                                            | Purpose                                                     |
| ---------------------- | ------------------------------------------------------------------------ | ----------------------------------------------------------- |
| **User**               | `role`, `teamId`                                                         | Validate role (EMPLOYEE) and team for internal Lambda logic |
| **Meal Participation** | `participation` (YES/NO)                                                 | Core data for marking and tracking meals                    |
| **Day Configuration**  | `type` (GOV HOLIDAY/CLOSED/EVENT)                                        | Used for read-only validation of date before update         |
| **Team (via GSI)**     | `teamId`, `name`, `leadId`, `description`, `additionalMetadata`          | Supports team-level queries for Team Leads/Admins           |

### Note: Team entity is not needed for feature 1 yet. Only shown for documentation as to show how User and Team entities will be connected. It is not needed as only admin/team lead can do team wise operations. 

Stored in **single DynamoDB table (`MHP`)** with PK/SK prefixes: `USER#<id>`, `MEAL#<date>`, `DAY#<date>`

## Access Patterns & Queries

1. **View own meals**
   - `Query PK=MEAL#<date> SK begins_with USER#<userId>`

2. **Update own meals**
   - `PutItem PK=MEAL#<date> SK=USER#<userId>#<mealType>`
   - Lambda validates role, cutoff, closed day

3. **Day validation**
   - `Query PK=DAY#<date> SK=META`

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

**PK:** `MEAL#<date>`  
**SK:** begins_with `USER#<userId>#<mealType>`

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