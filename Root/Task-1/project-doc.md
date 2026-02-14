# Meal Headcount Planner (MHP) -- Technical Design Document

## 1. Header

**Project Name:** Meal Headcount Planner (MHP)\
**Iteration:** 1 --- Daily Meal Opt-In/Out + Basic Visibility\
**Stack:** React (Frontend) + Go (Gin) Backend\
**Target Users:** Employees, Team Leads, Admin, Logistics

------------------------------------------------------------------------

## 2. Summary

Meal Headcount Planner (MHP) is a lightweight internal web application
designed to replace the current Excel based meal tracking process. It
enables employees to manage their daily meal participation while
providing operations teams with real time meal headcounts.

Iteration 1 focuses on daily participation tracking, role-based
access, and basic headcount visibility for logistics planning.

------------------------------------------------------------------------

## 3. Problem Statement

The current Excel-based system for tracking daily meal headcount:

-   Is manual and error-prone
-   Lacks real-time visibility
-   Makes it difficult to track last-minute changes
-   Does not scale well beyond 100+ employees
-   Requires manual aggregation for logistics planning

This results in inaccurate meal counts, food wastage, or shortages, and
high operational overhead.

------------------------------------------------------------------------

## 4. Goals and Non-Goals

### Goals (Iteration 1)

-   Provide a simple web interface for employees to manage daily meal
    participation
-   Default all employees as opted in unless they opt out
-   Allow authorized roles to update entries on behalf of employees
-   Provide real-time headcount per meal type for logistics/admin
-   Replace daily Excel workflow with a centralized system

### Non-Goals (Iteration 1)

-   No advanced reporting 
-   No payroll or billing integrations
-   No notifications
-   No complex scheduling beyond "today"

------------------------------------------------------------------------

## 5. Tech Stack and Rationale

### Frontend: **React**

-   As per training instruction

### Backend: **Go (Gin)**

-   As per training instruction

### Storage: File-based JSON (Iteration 1)

-   Low setup overhead
-   Easy to inspect and debug
-   Sufficient for \~100+ users with daily records
-   Designed to be replaceable by a relational DB later

------------------------------------------------------------------------

## 6. Scope of Changes

### New System

-   New backend service (Gin API)
-   New React frontend
-   JSON-based storage for users and daily participation

### Replaces

-   Excel sheets currently used for daily meal tracking

------------------------------------------------------------------------

## 7. Requirements

### Functional Requirements

#### FR1 --- Authentication

-   Users can log in with username and password
-   The system allows authorized roles(admin) to create new users
-   JWT based authentication
-   Each user has a role:
    -   Employee
    -   Team Lead
    -   Admin
    -   Logistics

#### FR2 --- View Today's Meals

-   Employees see a list of today's meal types:
    -   Lunch
    -   Snacks
    -   Iftar
    -   Event Dinner
    -   Optional Dinner
-   Default status: Participating

#### FR3 --- Employee Opt-Out

-   Employees can opt out of any meal for today
-   Changes allowed until cutoff time 

#### FR4 --- Role-Based Editing

-   Team Leads / Admin can update meal participation for any employee
-   Logistics has read-only access to participation but can view totals

#### FR5 --- Headcount View

-   Logistics/Admin can view:
    -   Total participating count per meal type for today

------------------------------------------------------------------------

### Non-Functional Requirements

-   JWT based security
-   Support at least 100 concurrent users
-   Basic auditability (who changed what)

------------------------------------------------------------------------

## 8. User Flows

### 8.1 Employee Flow

1.  User logs in
2.  Lands on "Today's Meals" page
3.  Sees list of meals with toggle (Participating / Opted Out)
4.  Changes status for one or more meals
5.  Clicks Save
6.  System confirms update

------------------------------------------------------------------------

### 8.2 Team Lead / Admin Flow (Edit on Behalf)

1.  Logs in
2.  Navigates to Employee Participation page
3.  Searches/selects employee
4.  Views their meal participation for today
5.  Edits participation
6.  Saves changes

------------------------------------------------------------------------

### 8.3 Logistics Flow (Headcount)

1.  Logs in
2.  Navigates to Headcount Dashboard
3.  Sees totals per meal type for today

------------------------------------------------------------------------

## 9. Design

### 9.1 High-Level Architecture

    React Frontend  →  Gin REST API  →  JSON File Storage
                          |
                       Auth Middleware
                          |
                    Role-Based Access

------------------------------------------------------------------------


### 9.2 API Endpoints (Gin)

| Method | Endpoint                     | Description                          | Roles            |
|--------|------------------------------|--------------------------------------|------------------|
| POST   | /login                       | Login                                | All              |
| POST   | /logout                      | Logout                               | All              |
| POST   | /register                    | Register                             | Admin            |
| GET    | /meals/today                 | Get today's meal types               | All              |
| GET    | /me/participation            | Get logged-in user's participation   | Employee+        |
| PUT    | /me/participation            | Update own participation             | Employee+        |
| GET    | /users                       | List users                           | Team Lead+       |
| GET    | /users/:id/participation     | Get participation for a user         | Team Lead+       |
| PUT    | /users/:id/participation     | Update participation for a user      | Team Lead+       |
| GET    | /headcount/today             | Get totals per meal                  | Logistics/Admin  |


------------------------------------------------------------------------

## 10. Key Decisions and Trade-offs

 | Decision           | Why                                   | Trade-off                                   |
|--------------------|----------------------------------------|---------------------------------------------|
| Default opt-in     | Matches real-world expectation         | Must ensure opt-out is easy                 |
| JSON storage       | Fast to build, no DB setup required    | Limited scalability and concurrency control |
| Single-day focus   | Simplifies the data model              | Historical reporting is deferred            |


------------------------------------------------------------------------

## 11. Security and Access Control

### Authentication

-   Username/password login
-   JWT based authentication

### Authorization 

| Role       | Permissions                              |
|------------|-------------------------------------------|
| Employee   | View/update own participation             |
| Team Lead  | Edit others                               |
| Admin      | Full access including headcount           |
| Logistics  | View headcount only                       |


------------------------------------------------------------------------

## 12. Testing Plan

### Unit Tests (Backend)

-   Participation logic (default opt-in)
-   Headcount aggregation
-   Role permission checks

### Integration Tests

-   Login flow
-   Employee updating meals
-   Team Lead editing another user
-   Headcount totals correctness

------------------------------------------------------------------------

## 13. Operations

### Running Locally

**Backend**

``` bash
go run main.go
```

**Frontend**

``` bash
npm install
npm run dev
```


------------------------------------------------------------------------

