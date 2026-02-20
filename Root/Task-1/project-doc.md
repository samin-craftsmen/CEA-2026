# Meal Headcount Planner (MHP) — Technical Design Document

## 1. Header

**Project Name:** Meal Headcount Planner (MHP)  
**Current Iteration:** 3 — Scheduling + Events + Operational Readiness  
**Previous Iterations Covered:**  
- Iteration 1 — Daily Meal Opt-In/Out + Basic Visibility  
- Iteration 2 — Team Views + Special Days + Live Updates  
**Stack:** React (Frontend) + Go (Gin) Backend  
**Target Users:** Employees, Team Leads, Admin, Logistics  

---

## 2. Summary

Meal Headcount Planner (MHP) is a lightweight internal web application designed to replace the Excel-based meal tracking process.

The system now supports:

- Daily and future meal participation
- Team-based visibility and bulk handling
- Special day configuration (Holiday / Office Closed / Celebration)
- Event-specific meals
- Work location tracking (Office/WFH)
- Company-wide WFH declarations
- Live headcount updates
- Forecasting
- Audit logging
- Monthly WFH allowance tracking with over-limit reporting
- Operational dashboards for Logistics/Admin

Iteration 3 focuses on operational readiness, forecasting, scheduling flexibility, and accountability.

---

## 3. Problem Statement

The Excel-based system:

- Is manual and error-prone
- Lack real-time visibility
- Cannot scale beyond 100+ employees
- No team-level insights
- Provides no forecasting capability
- Has no traceability of edits
- Does not track WFH allowance
- Requires manual aggregation for logistics reporting

This results in inaccurate counts, food wastage, shortages, and high operational overhead.

---

## 4. Goals (Iteration 1–3)

- Replace Excel workflow with a centralized system
- Default all employees as opted-in unless opted-out
- Provide role-based access and editing
- Enable team-level and company-wide visibility
- Support bulk updates and exception handling
- Support special day controls
- Enable live headcount updates
- Allow future meal scheduling within a limited window
- Provide headcount forecasting
- Track work location (Office/WFH)
- Support company-wide WFH periods
- Track monthly WFH allowance (soft limit: 5 days)
- Highlight over-limit employees in reports
- Provide audit logs for accountability
- Support event-based meals
- Provide operational dashboards for logistics

---

## 5. Tech Stack

### Frontend: React
### Backend: Go (Gin)
### Storage : JSON-based file storage  

---


# 6. Functional Requirements

## Functional Requirement 1 — Authentication & Roles

- Username/password login
- JWT-based authentication
- Roles:
  - Employee
  - Team Lead
  - Admin

---

## Functional Requirement 2 — Daily Meal Participation

- Employees see meals:
  - Lunch
  - Snacks
  - Iftar
  - Event Dinner
  - Optional Dinner
- Default: Participating
- Employees can opt out before cutoff
- Team Lead/Admin can edit on behalf

---

## Functional Requirement 3 — Team-Based Visibility 

- Employees see their assigned team
- Team Leads view participation for their own team
- Admin view across all teams
- Data grouped by team

---

## Functional Requirement 4 — Bulk & Exception Handling

- Team Lead can bulk update within own team
- Admin can bulk update across teams

---

## Functional Requirement 5 — Special Day Controls

Admin can mark a date as:

- Office Closed
- Government Holiday
- Special Celebration Day (with note)

System behavior:

- Office Closed → Meals disabled
- Holiday/Celebration → Adjust messaging and reporting

---

## Functional Requirement 6 — Improved Headcount Reporting

Headcount available by:

- Meal type
- Team
- Overall total
- Office vs WFH split

---

## Functional Requirement 7 — Live Updates

- Participation changes update totals instantly
- No page refresh required
- WebSocket-based broadcasting

---

## Functional Requirement 8 — Work Location per Date

Employees can set:

- Office
- WFH

Team Lead/Admin can correct entries.

---

## Functional Requirement 9 — Company-Wide WFH Period

Admin/Logistics can declare a date range as:

- “WFH for everyone”

System treats all employees as WFH for reporting during that period.

---

## Functional Requirement 10 — Future Meal Scheduling 

- Employees can set meal participation for future dates
- Forward window is configurable
- Admin/Logistics can view forecasted headcount

---

## Functional Requirement 11 — Event Meals

Admin/Logistics can create event meals with:

- Date
- Meal type
- Optional note

Employees can opt in/out for event meals separately.

---

## Functional Requirement 12 — Audit Logging

System records:

- Who changed participation
- What was changed
- When it was changed
- Whether change was self-edit or role-edit

Admin can view audit logs.

---

## Functional Requirement 13 — Operational Dashboard

Dashboard includes:

- Daily snapshot
- Upcoming forecast snapshot
- Special day indicators
- Over-limit WFH indicators
- Team breakdown
- Office vs WFH split

---

## Functional Requirement 14 — Monthly WFH Usage 

- Track WFH days per employee per month
- Standard allowance: 5 days
- System does NOT block excess usage

---

## Functional Requirement 15 — Over-Limit Indicators & Filters

Team Lead/Admin views:

- Highlight employees exceeding 5 days
- Show:
  - Number of employees over limit
  - Total extra WFH days
- Filter:
  - Show only over-limit employees

---

# 7. Non Functional Requirements

- JWT security
- Role-based middleware authorization
- Support 100+ concurrent users
- Real-time updates via WebSocket
- Basic auditability
---

# 8. User Flows

## 8.1 Employee Flow

1. Login  
2. Select today or future date  
3. Set:
   - Meal participation
   - Work location  
4. Save  
5. See live updates  
6. View monthly WFH usage  

---

## 8.2 Team Lead Flow

1. Login  
2. View team participation  
3. Bulk update if required  
4. Filter over-limit employees  
5. View WFH summary  
6. Edit entries when necessary  

---

## 8.3 Admin / Logistics Flow

1. Login  
2. View operational dashboard  
3. Configure:
   - Special days
   - Event meals
   - WFH periods  
4. Generate announcement draft  
5. View forecast  
6. View audit logs  
7. Analyze over-limit reports  

---

# 9. High-Level Architecture

```
React Frontend
        ↓
Gin REST API
        ↓
Auth Middleware
        ↓
Role-Based Access
        ↓
Business Logic Layer
        ↓
JSON Storage
        ↓
WebSocket Broadcast Layer
```

---

# 10. Testing Plan

## Backend Unit Tests

- Participation logic
- Forecast calculations
- Bulk updates
- WFH monthly calculation
- Over-limit detection
- Audit log creation
- Special day behavior

## Integration Tests

- Login flow
- Team visibility restrictions
- Live update broadcasting
- Event meal behavior
- Company-wide WFH override

---

# 11. Operations

## Run Locally

### Backend

```bash
go run main.go

```
### Frontend

```bash

npm install
npm run dev

```
