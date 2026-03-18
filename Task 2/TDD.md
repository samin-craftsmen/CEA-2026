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

**Table Name:** trainee-2026-samin-dynamoDB-mhp

## Primary Keys
- **PK** – Partition Key
- **SK** – Sort Key

## Conventions / Prefixes
- `USER#<userId>` – user entity  
- `TEAM#<teamId>` – team entity  
- `DAY#<YYYY-MM-DD>` – date-based entity  
- `MEAL#<mealType>` – meal type (lunch/dinner)  
- `LOC#<locationType>` – work location (office/WFH)  
- `SPECIAL#<YYYY-MM-DD>` – special day  


---

## Users

### Access Patterns
- Get user by `userId`  
- Get user by external identity (e.g., `discordId`)  
- Get all users  

### DB Schema
- **PK:** USERS  
- **SK:** `<userId>`  

### Notes
- User metadata stored in the item (role, team)  
- Map user to team when adding a user to the Users entity  

---

## Team

### Access Patterns
- Get team by `teamId`  
- Get all teams  
- Get all users in a team  

### DB Schema
- **PK:** TEAM#<teamId>  
- **SK:** USER#<userId>  

### Notes
- Scanning to get all teams is acceptable as it is rare and the dataset is small  

---

## Meal Participation

### Access Patterns
- Get all participation records for a date (daily totals)  
- Get a specific user's meals for a date  
- Get a specific meal record (user + date + mealType)  
- Get a user's participation history across dates (reporting)  
- Admin quick daily totals  
- Team Lead quick team-level totals  

### DB Schema
- **PK:** DAY#<YYYY-MM-DD>  
- **SK:** TEAM#<teamId>MEAL#<mealType>#USER#<userId>  

### Notes
- Admin daily totals: Query `PK = DAY#<date>` → all `MEAL#*` items  
- Team Lead totals: Include `TEAM#<teamId>` as an attribute  

---

## Work Location

### Access Patterns
- Get all location records for a date (daily headcount)  
- Get a specific user's location for a date  
- Get a user's location records for a month  

### DB Schema
- **Date View:**  
  - **PK:** DAY#<YYYY-MM-DD>  
  - **SK:** USER#<userId>#LOC#<locationType>  

- **User View:**  
  - **PK:** USER#<userId>  
  - **SK:** DAY#<YYYY-MM-DD>#LOC#<locationType>  

### Notes
- Records are stored twice since a single view cannot satisfy all access patterns  
- Range search on PKs is not feasible  

---

## Special Day

### Access Patterns
- Get special day for a specific date  
- Get all special days in a month  

### DB Schema
- **PK:** SPECIAL  
- **SK:** Date#<YYYY-MM-DD>

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


## Access Patterns & Queries

1. **View own meals**
Retrieve all meal records for a user on a specific date.
   - `Query PK=DAY#<date> Filter: contains(SK, USER#<userId>)`

2. **Update own meals**
The employee can opt in or opt out of a meal.
   - `PutItem PK=DAY#<date> SK = TEAM#<teamId>#MEAL#<mealType>#USER#<userId>`
   - Attributes: YES/NO
   - Lambda validates role, cutoff, closed day

3. **Day validation**
   - Before allowing updates, the system verifies whether the selected date is a special day.
   - `Query PK=SPECIAL SK = DATE#<YYYY-MM-DD>`
   - If result exists and type indicates closure, updates are rejected.

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

**PK:** `DAY#<date>`  
**SK:** `TEAM#<teamId>#MEAL#<mealType>#USER#<userId>`

### View Status

1. User runs `/meal view`
2. Lambda queries:

**PK:** `DAY#<date>`  
**Filter:** `contains(SK, USER#<userId>)`


3. Bot returns status

## Validations

- Closed day → reject
- After cutoff → reject
- Duplicate → overwrite
---

---

# Feature 2: Team Lead Meal Management

## Overview

Allows Team Leads to manage meal participation for employees within their assigned team.

Team Leads can:

- Select a date
- Select meal type
- Mark meal participation (Yes/No) for employees within their team
- View meal status for employees in their team for a selected date

The feature improves operational efficiency by allowing delegated management of meal participation while enforcing strict role and team-based access restrictions.

Updates performed by Team Leads must follow the same operational policies as employee updates, including date validation, cutoff rules, and data integrity checks.

---

## Access Patterns & Queries

### 1. View team meal status for a date

Retrieve all meal records for a specific team and date.

**Query**

PK: `DAY#<date>`

Filter condition:
`contains(SK, TEAM#<teamId>)`

Result:
Returns all team member meal records for that date.

---

### 2. Update meal participation for a team member

Team Lead updates meal participation for a user belonging to their team.

**PutItem**

PK: `DAY#<date>`  
SK: `TEAM#<teamId>#MEAL#<mealType>#USER#<userId>`

Attributes:

```
participation: YES | NO
updatedBy: TEAM_LEAD
updatedAt: timestamp
```

Lambda validations performed before update:

- Role must be `TEAM_LEAD`
- Target user must belong to the same team
- Date must not be a closed/special day
- Update must be before cutoff time

Latest update overwrites previous value.

---

### 3. Validate team membership

Before allowing updates, Lambda verifies that the target user belongs to the Team Lead’s team.

**Query**

PK: `TEAM#<teamId>`  
SK: `USER#<userId>`

If no record exists, the update is rejected.

---


## User Flows

### Update Team Member Meal

1. Team Lead runs `/team-meal set`

2. Bot collects:
   - Date
   - Meal Type
   - Target Employee
   - YES/NO participation

3. Lambda performs validations:
   - Validate role = `TEAM_LEAD`
   - Validate employee belongs to Team Lead's team
   - Check for closed day
   - Check cutoff time

4. If validation passes:

```
PK: DAY#<date>
SK: TEAM#<teamId>#MEAL#<mealType>#USER#<userId>
```

5. Record is written to DynamoDB.

6. Bot confirms update.

---

### View Team Meal Status

1. Team Lead runs `/team-meal view`

2. Bot collects:
   - Date
   - Meal Type (optional)

3. Lambda queries DynamoDB:

```
PK: DAY#<date>
Filter: contains(SK, TEAM#<teamId>)
```

4. Lambda formats response to display:

```
Employee Name | Lunch | Dinner
```

5. Bot returns formatted team status.

---

## Validations

Team Lead operations must satisfy the following conditions:

- **Role Validation**
  - User role must be `TEAM_LEAD`.

- **Team Boundary Enforcement**
  - Team Leads can only modify records of employees in their own team.

- **Closed Day Restriction**
  - Updates are rejected if the selected date is marked as office closed or holiday.

- **Cutoff Enforcement**
  - Updates cannot occur after the configured meal cutoff time.

- **Duplicate Handling**
  - Latest update overwrites previous participation value.

- **Data Integrity**
  - Invalid user IDs or meal types are rejected.

---

# Feature 3: Admin Meal Management

## Overview

Allows Administrators to fully manage and control meal participation across the system.

Admins can:

- Select a date  
- Select meal type  
- Mark meal participation (Yes/No) for any employee  
- View meal status for any employee or across the system  
- Perform bulk updates for multiple employees  
- Add or modify meal types  
- Add or modify meal items  

All admin operations must comply with system-wide rules such as cutoff time, special day restrictions, and data integrity constraints.

---

## Access Patterns & Queries

### 1. View meal status for a date (global view)

Retrieve all meal participation records for a specific date.

**Query**

PK: `DAY#<date>`

Result:  
Returns all meal records for all teams and users.

---

### 2. View meal status for a specific employee

Retrieve meal participation for a specific user on a given date.

**Query**

PK: `DAY#<date>`  
Filter: `contains(SK, USER#<userId>)`

---

### 3. Update meal participation for any employee

Admin updates meal participation for any user in the system.

**PutItem**

PK: `DAY#<date>`  
SK: `TEAM#<teamId>#MEAL#<mealType>#USER#<userId>`

Attributes:
participation: YES | NO
updatedBy: ADMIN
updatedAt: timestamp


Validations before update:

- Role must be `ADMIN`
- Date must not be a closed/special day
- Update must be before cutoff time

---

### 4. Bulk update meal participation

Admins can update multiple users in a single operation.

**Approach**

- Batch processing using multiple `PutItem` operations (or BatchWrite if optimized)

For each user:

PK: `DAY#<date>`  
SK: `TEAM#<teamId>#MEAL#<mealType>#USER#<userId>`

Attributes:
participation: YES | NO
updatedBy: ADMIN
updatedAt: timestamp


---

### 5. Manage meal types

Admins can create or update meal types (e.g., lunch, dinner, snacks).

**DB Schema**

PK: `CONFIG#MEALTYPE`  
SK: `<mealType>`

Attributes:
name: string
isActive: boolean
createdAt: timestamp
updatedAt: timestamp


---

### 6. Manage meal items

Admins can define items under each meal type.

**DB Schema**

PK: `MEALTYPE#<mealType>`  
SK: `ITEM#<itemId>`

Attributes:
itemName: string
isActive: boolean
createdAt: timestamp
updatedAt: timestamp


---

## User Flows

### Update Employee Meal

1. Admin runs `/admin-meal set`

2. Bot collects:
   - Date  
   - Meal Type  
   - Target Employee  
   - YES/NO participation  

3. Lambda validations:
   - Validate role = `ADMIN`  
   - Check special/closed day  
   - Check cutoff time  

4. Write to DynamoDB:
PK: DAY#<date>
SK: TEAM#<teamId>#MEAL#<mealType>#USER#<userId>


5. Bot confirms update

---

### Bulk Update Meals

1. Admin runs `/admin-meal bulk-set`

2. Bot collects:
   - Date  
   - Meal Type  
   - List of Employees  
   - YES/NO  

3. Lambda:
   - Validate role = `ADMIN`  
   - Validate cutoff and special day  
   - Iterate and write records  

4. Bot confirms bulk update

---

### View Meal Status

1. Admin runs `/admin-meal view`

2. Bot collects:
   - Date  
   - (Optional) Team / Employee filter  

3. Lambda queries:
PK: DAY#<date>


4. Returns aggregated or filtered results

---

### Manage Meal Types

1. Admin runs `/meal-type add` or `/meal-type update`

2. Lambda writes:
PK: CONFIG#MEALTYPE
SK: <mealType>


---

### Manage Meal Items

1. Admin runs `/meal-item add` or `/meal-item update`

2. Lambda writes:
PK: MEALTYPE#<mealType>
SK: ITEM#<itemId>


---

## Validations

Admin operations must satisfy:

- **Role Validation**
  - User role must be `ADMIN`

- **Closed Day Restriction**
  - No updates allowed if the date is marked as office closed

- **Cutoff Enforcement**
  - Updates must occur before configured cutoff time

- **Duplicate Handling**
  - Latest update overwrites previous value (last-write-wins)

- **Data Integrity**
  - Invalid user IDs, meal types, or item IDs are rejected

---