# Technical Documentation  
## Meal Headcount Planner

---

## Scope

- A lightweight internal web app to replace the current Excel-based process for collecting and calculating daily meal headcount for 100+ employees.  
- **Frontend:** React  
- **Backend:** Gin  
- Role-based access:
  - Employee  
  - Team Lead  
  - Admin / Logistics  
- Basic meal management and information display  
- File based JSON storage

---

## Assumptions

- All employees are opted in by default  
- Modifications can only be made for the current working day (unless an admin allows otherwise)  
- Each user has only one role  
- Employees must opt out before a certain daily cutoff time  
- Admins and Team Leads can modify employee meal options  

---

## Key Flows

### User Access

| Step   | API       | Method | Purpose                       |
|--------|----------|--------|-------------------------------|
| Login  | `/login` | POST   | Authenticate user credentials |
| Logout | `/logout`| POST   | Remove user session           |

---

### Meal Participation (Employees)

| Step              | API               | Method | Purpose                         |
|-------------------|------------------|--------|---------------------------------|
| View today’s meal | `/meals/today`   | GET    | Get meal list                   |
| Opt out of meal   | `/meals/opt-out` | POST   | Mark employee as opted out      |
| Opt back in       | `/meals/opt-in`  | POST   | Remove opted-out record         |

---

### Meal Types

| Step                  | API            | Method | Purpose                        |
|-----------------------|---------------|--------|--------------------------------|
| Get meals for the day | `/meal/types` | GET    | Get types of meals for the day |

---

### Employee Management (Admin / Team Lead)

| Step                 | API                     | Method | Purpose                        |
|----------------------|-------------------------|--------|--------------------------------|
| Load employees       | `/users`                | GET    | Get employee list              |
| Update participation | `/participation/update` | POST   | Update opt-in / opt-out status |

---

### Reporting

| Step          | API                | Method | Purpose                                                     |
|---------------|--------------------|--------|-------------------------------------------------------------|
| Get headcount | `/meals/headcount` | GET    | Get number of people opted in per meal type for the day     |

---

## Verification & Definition of Done

- Employees can log in and view meal information  
- Meals are opted in by default for all employees  
- Employees can opt out of any single or all meals for the day  
- Team Leads and Admins can update employee meal participation  
- Admin / Logistics can view meal headcount  
- Access is restricted based on user roles  
