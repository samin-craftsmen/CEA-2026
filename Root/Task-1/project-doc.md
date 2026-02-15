# Meal Headcount Planner (MHP) -- Technical Design Document

## 1. Header

**Project Name:** Meal Headcount Planner (MHP)\
**Iteration:** 2 - Team/Department Views + Special Day Controls\
**Stack:** React (Frontend) + Go (Gin) Backend\
**Target Users:** Employees, Team Leads, Admin, Logistics

------------------------------------------------------------------------

## 2. Summary

Meal Headcount Planner (MHP) replaces the manual Excel-based process for
tracking daily meal participation.

Iteration 2 expands the system to support:

-   Team-based visibility and access control
-   Bulk updates and exception handling
-   Special day controls (Office Closed / Holiday / Celebration)
-   Work location tracking (Office / WFH)
-   Company-wide WFH periods
-   Improved headcount reporting
-   Live dashboard updates (no page refresh)
-   Auto-generated daily announcement draft

The system evolves from simple daily tracking to an operational planning
tool.

------------------------------------------------------------------------

## 3. Problem Statement

While Iteration 1 solved daily participation tracking, operational gaps
remain:

-   No team-level visibility
-   No bulk handling for offsites or events
-   No way to handle holidays or office closures
-   No work-location distinction (Office vs WFH)
-   No real-time UI updates
-   No standardized daily announcement message

Iteration 2 addresses planning complexity across teams and special
scenarios.

------------------------------------------------------------------------

## 4. Goals and Non-Goals

### Goals (Iteration 2)

-   Introduce team-based visibility and reporting
-   Enable bulk participation changes within role scope
-   Support special day types with system-level effects
-   Track work location per date
-   Support company-wide WFH periods
-   Provide real-time headcount updates
-   Generate formatted daily announcement text
-   Extend reporting breakdowns


------------------------------------------------------------------------

## 5. Scope of Changes

### Backend Enhancements

-   Team association per user
-   Day configuration entity
-   Work-location tracking
-   WFH date-range logic
-   Bulk update logic
-   Real-time update mechanism (WebSocket / SSE)
-   Announcement generation endpoint

### Frontend Enhancements

-   Team view screens
-   Special day configuration UI
-   Work-location selector
-   Live dashboard updates
-   Announcement generator panel

------------------------------------------------------------------------

## 6. Functional Requirements

### Functional Requirements 1 - Team-Based Visibility

-   Each employee belongs to a team
-   Employees can view their team name
-   Team Leads can view participation for their team only
-   Admin/Logistics can view all teams

### Functional Requirements 2 - Bulk Participation Actions

Team Lead/Admin can:

-   Apply bulk opt-out for a team or selected users
-   Apply bulk change for a specific meal
-   Mark group as opted-out due to event/offsite

Bulk changes must respect role scope.

### Functional Requirements 3 - Special Day Controls

Admin/Logistics can configure a date as:

-   Office Closed
-   Government Holiday
-   Special Celebration Day (with note)

| Day Type             | Behavior                          |
|----------------------|-----------------------------------|
| Office Closed        | Meals disabled                    |
| Government Holiday   | Meals disabled                    |
| Special Celebration  | Meals enabled + note displayed    |
| Normal               | Default behavior                  |


### Functional Requirements 4 - Work Location Per Date

-   Employees can select: Office / WFH for a date
-   Team Leads/Admin can correct missing entries
-   Location impacts reporting

### Functional Requirements 5 - Company-Wide WFH Period

-   Admin/Logistics can define a date range as "WFH for everyone"
-   All users default to WFH during that range
-   Reports treat users as WFH unless explicitly overridden

### Functional Requirements 6 - Improved Headcount Reporting

Headcount must be available by:

-   Meal type
-   Team
-   Overall total
-   Office vs WFH split

### Functional Requirements 7 - Live Updates (No Refresh)

When participation or work location changes:

-   Dashboard updates instantly
-   Bulk updates reflect immediately

Implementation: WebSocket or Server-Sent Events (SSE)

### Functional Requirements 8 - Daily Announcement Draft

Admin/Logistics can:

-   Select a date
-   Generate formatted message including:
    -   Meal-wise totals
    -   Office/WFH split
    -   Special day note

Output must be copy/paste friendly.

------------------------------------------------------------------------

## 7. Data Model (JSON Structure --- Iteration 2)

### users.json

``` json
{
  "id": 1,
  "username": "john",
  "role": "employee",
  "team": "Engineering"
}
```

### participation.json

``` json
{
  "date": "2026-02-14",
  "user_id": 1,
  "meals": {
    "lunch": true,
    "snacks": false
  },
  "work_location": "office"
}
```

### day_config.json

``` json
{
  "date": "2026-02-21",
  "type": "office_closed",
  "note": ""
}
```

### wfh_periods.json

``` json
{
  "start_date": "2026-02-10",
  "end_date": "2026-02-15"
}
```

------------------------------------------------------------------------


## 8. User Flows

### 8.1 Employee Flow

1.  Log in
2.  View team name
3.  Set work location
4.  Update meal participation
5.  Save
6.  Dashboard updates live

### 8.2 Team Lead Flow

1.  Log in
2.  View team participation
3.  Apply bulk change (if needed)
4.  Correct work-location entries
5.  See real-time updated totals

### 8.3 Admin/Logistics Flow

1.  Log in
2.  Configure special day (if needed)
3.  Declare WFH period (if required)
4.  View detailed report
5.  Generate daily announcement
6.  Monitor live headcount

------------------------------------------------------------------------

## 9. Architecture

    React Frontend
       │
       ├── REST API (Gin)
       │       ├── Auth Middleware
       │       ├── Role Guard
       │       ├── Day Config Service
       │       ├── Reporting Engine
       │       └── Bulk Action Handler
       │
       ├── WebSocket/SSE Channel
       │
       └── JSON File Storage

------------------------------------------------------------------------

## 10. Authorization Matrix

| Role       | Own Meals | Edit Others | Bulk | View Team | View All | Special Day | WFH Period | Announcement |
|------------|-----------|-------------|------|-----------|----------|------------|-----------|--------------|
| Employee   | ✓         | ✗           | ✗    | ✗         | ✗        | ✗          | ✓         | ✗            |
| Team Lead  | ✓         | ✓ (Team)    | ✓ (Team) | ✓     | ✗        | ✗          | ✓         | ✗            |
| Admin      | ✓         | ✓           | ✓    | ✓         | ✓        | ✓          | ✓         | ✓            |


------------------------------------------------------------------------

## 11. Verification Approach (Definition of Done)

Iteration 2 is considered complete when:

-   Team-based filtering works correctly
-   Bulk updates respect role boundaries
-   Special day types disable meals appropriately
-   Work location is stored and reflected in reports
-   WFH period overrides default behavior
-   Reports correctly calculate:
    -   Meal totals
    -   Team totals
    -   Office vs WFH split
-   Dashboard updates instantly without reload
-   Announcement output matches reporting totals

------------------------------------------------------------------------

## 12. Operations

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
