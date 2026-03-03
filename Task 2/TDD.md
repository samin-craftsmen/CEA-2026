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

### Note: Access pattern can be modified based on new features based on which this will be updated. Also how these access patterns are executed is explained on feature level design.
---



