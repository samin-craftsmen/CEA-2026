import json
import os

USERS_FILE = "users.json"

if os.path.exists(USERS_FILE):
    with open(USERS_FILE, "r") as f:
        users = json.load(f)
else:
    users = []

existing_usernames = {u["username"] for u in users}

print("⚡ User Creation Script")
print("Type 'exit' as username to quit.\n")

while True:
    username = input("Enter username: ").strip()
    if username.lower() == "exit":
        break

    if username in existing_usernames:
        print("Username already exists. Try a different one.\n")
        continue

    password = input("Enter password: ").strip()
    if not password:
        print("Password cannot be empty.\n")
        continue

    role = input("Enter role (admin/teamLead/employee): ").strip()
    if role.lower() not in {"admin", "teamlead", "employee"}:
        print("Invalid role. Must be 'admin', 'teamLead' or 'employee'.\n")
        continue

    users.append({
        "username": username,
        "password": password,
        "role": role
    })
    existing_usernames.add(username)
    print(f"User '{username}' added!\n")

with open(USERS_FILE, "w") as f:
    json.dump(users, f, indent=2)

print("All users saved to users.json!")
