#  Pull Request: Team-Based Visibility (Backend)

##  Related Ticket #7 
Feature: Team-Based Visibility (Backend Only)

---

##  Description

This PR implements backend support for role-based team visibility for daily meal participation.

### Access Rules

- **Employee**
  - Can view their assigned team.
  - Can view their own participation only.

- **Team Lead**
  - Can view daily participation of members within their own team.
  - Cannot access data from other teams.

- **Admin / Logistics**
  - Can view participation data across all teams.

---

##  Changes Made

- [ ] Added team field to user model
- [ ] Implemented team-based filtering logic
- [ ] Added team-based grouping for participation response
- [ ] Enforced access restriction for Team Leads
- [ ] Applied default opt-in logic when participation record is missing

---

##  How to Test

1. Login as **Employee**
   - Verify only own participation is visible.

2. Login as **Team Lead**
   - Verify only team members’ participation is visible.
   - Verify other teams’ data is not accessible.

3. Login as **Admin / Logistics**
   - Verify all teams’ participation data is accessible.

---

##  Acceptance Criteria

- [ ] Team information is stored with users
- [ ] Participation data is filtered based on role and team
- [ ] Team Leads cannot access other teams’ data
- [ ] Admin have full visibility across teams
- [ ] No empty team keys in grouped response

