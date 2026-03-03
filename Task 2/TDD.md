# Technical Design Document

# Meal Headcount Planner

**UI:** Discord, React   
**Backend:** AWS Lambda  
**Database:** Amazon DynamoDB  

---

# 1. System-Level Design
## 1. Summary

The Meal Headcount Planner is a system designed to manage daily employee meal participation. Employees can indicate whether they will participate in meals and whether they will work from the office or remotely. The system integrates with Discord for employee interactions and provides a web dashboard for administrators to manage configurations, view reports, and perform bulk operations.The backend is implemented using AWS serverless architecture to ensure scalability, high availability, and minimal operational overhead.

Technologies used:
- Discord Bot for employee interaction
- React Web Dashboard for administrative management
- AWS Lambda for business logic
- Amazon DynamoDB using a single-table design for scalable data access
- API Gateway for HTTP endpoints

Note: The first iteration of the project is developed locally. Iteration 2 will be built on the cloud.

## 2. Problem Statement

Currently, meal participation is managed manually(Excel Based System). This leads to:

- Inaccurate meal headcount
- Food wastage or shortages
- Lack of centralized tracking
- Difficulty managing holidays and office day status
- No clear visibility for administrators

The system aims to provide a centralized platform to manage meal participation and work locations while ensuring role-based access control and enforcing operational rules such as cutoff times and holiday restrictions.

## 3. High Level Architecture

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


