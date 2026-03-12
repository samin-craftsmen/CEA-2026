# DynamoDB Single Table Design -- Access Patterns & Schema

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