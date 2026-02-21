import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import "./MealPlanner.css";

const allMeals = [
  "Lunch",
  "Snacks",
  "Iftar",
  "Event Dinner",
  "Optional Dinner"
];

export default function MealPlanner() {
  const navigate = useNavigate();
  const token = localStorage.getItem("token");

  const [role, setRole] = useState("");
  const [username, setUsername] = useState("");
  const [team, setTeam] = useState("");

  const [selectedMeals, setSelectedMeals] = useState<string[]>([]);
  const [searchUser, setSearchUser] = useState("");
  const [searchDate, setSearchDate] = useState("");
  const [headcount, setHeadcount] = useState<any>(null);

  const [mealItems, setMealItems] = useState<any>({});
  const [itemInputs, setItemInputs] = useState<any>({});

  const tomorrow = new Date();
  tomorrow.setDate(tomorrow.getDate() + 1);
  const tomorrowStr = tomorrow.toISOString().split("T")[0];

  const [teamMeals, setTeamMeals] = useState<any[]>([]);
  const [showTeamMeals, setShowTeamMeals] = useState(false);

  const [allTeamsData, setAllTeamsData] = useState<any>(null);
  const [adminDate, setAdminDate] = useState("");

  const [bulkDate, setBulkDate] = useState("");
  const [bulkMeals, setBulkMeals] = useState<string[]>([]);
  const [bulkLoading, setBulkLoading] = useState(false);

  const [adminBulkDate, setAdminBulkDate] = useState("");
  const [adminBulkMeals, setAdminBulkMeals] = useState<string[]>([]);
  const [adminBulkLoading, setAdminBulkLoading] = useState(false);
  // ============================= Admin Bulk Ops ==========================//
  const toggleAdminBulkMeal = (meal: string) => {
    if (adminBulkMeals.includes(meal)) {
      setAdminBulkMeals(adminBulkMeals.filter(m => m !== meal));
    } else {
      setAdminBulkMeals([...adminBulkMeals, meal]);
    }
  };

  const toggleBulkMeal = (meal: string) => {
    if (bulkMeals.includes(meal)) {
      setBulkMeals(bulkMeals.filter(m => m !== meal));
    } else {
      setBulkMeals([...bulkMeals, meal]);
    }
  };

  const adminBulkOptOut = async () => {
    if (!adminBulkDate || adminBulkMeals.length === 0) {
      alert("Select date and at least one meal");
      return;
    }

    setAdminBulkLoading(true);

    const res = await fetch(
      `http://localhost:8080/admin/meals/opt-out/${adminBulkDate}`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          date: adminBulkDate,
          meals: adminBulkMeals,
        }),
      }
    );

    const data = await res.json();
    setAdminBulkLoading(false);

    if (!res.ok) {
      alert(data.error || "Operation failed");
      return;
    }

    alert(`Opt-out applied to ${data.updated_count} users`);
  };

  const adminBulkOptIn = async () => {
    if (!adminBulkDate || adminBulkMeals.length === 0) {
      alert("Select date and at least one meal");
      return;
    }

    setAdminBulkLoading(true);

    const res = await fetch(
      `http://localhost:8080/admin/meals/opt-in/${adminBulkDate}`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          date: adminBulkDate,
          meals: adminBulkMeals,
        }),
      }
    );

    const data = await res.json();
    setAdminBulkLoading(false);

    if (!res.ok) {
      alert(data.error || "Operation failed");
      return;
    }

    alert(`Opt-in applied to ${data.updated_count} users`);
  };

  //======================== Team Bulk In Option =======================//
  const bulkOptOut = async () => {
    if (!bulkDate || bulkMeals.length === 0) {
      alert("Select date and at least one meal");
      return;
    }

    setBulkLoading(true);

    const res = await fetch("http://localhost:8080/teams/meals/optout", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({
        date: bulkDate,
        meals: bulkMeals,
      }),
    });

    const data = await res.json();
    setBulkLoading(false);

    if (!res.ok) {
      alert(data.error || "Operation failed");
      return;
    }

    alert(`Opt-out successful for ${data.updated_count} members`);
  };

  // ================= Team Bulk Out Option ==============================//
  const bulkOptIn = async () => {
    if (!bulkDate || bulkMeals.length === 0) {
      alert("Select date and at least one meal");
      return;
    }

    setBulkLoading(true);

    const res = await fetch("http://localhost:8080/teams/meals/optin", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({
        date: bulkDate,
        meals: bulkMeals,
      }),
    });

    const data = await res.json();
    setBulkLoading(false);

    if (!res.ok) {
      alert(data.error || "Operation failed");
      return;
    }

    alert(`Opt-in successful for ${data.updated_count} members`);
  };

  // ================= FETCH ALL TEAMS PARTICIPATION (ADMIN ONLY) =================
  const fetchAllTeamsParticipation = async () => {
    if (!adminDate) {
      alert("Select a date first");
      return;
    }

    const res = await fetch(
      `http://localhost:8080/admin/teams/meals/${adminDate}`,
      {
        headers: { Authorization: `Bearer ${token}` },
      }
    );

    const data = await res.json();
    setAllTeamsData(data);
  };


  // ================= LOAD USER INFO =================
  useEffect(() => {
    fetch("http://localhost:8080/me", {
      headers: { Authorization: `Bearer ${token}` },
    })
      .then(res => {
        if (!res.ok) {
          navigate("/");
          return;
        }
        return res.json();
      })
      .then(data => {
        if (!data) return;
        setRole(data.role.toLowerCase());
        setUsername(data.username);
        setTeam(data.team);
      });
  }, []);

  // ================= LOAD TEAM TODAY MEALS (TEAM LEAD) =================
  useEffect(() => {
    if (role !== "teamlead") return;

    fetch("http://localhost:8080/teams/meals/today", {
      headers: { Authorization: `Bearer ${token}` },
    })
      .then(res => res.json())
      .then(data => {
        setTeamMeals(data || []);
      });
  }, [role]);


  // ================= LOAD TOMORROW MEALS =================
  useEffect(() => {
    if (!role) return;

    fetch("http://localhost:8080/meals/tomorrow", {
      headers: { Authorization: `Bearer ${token}` },
    })
      .then(res => res.json())
      .then(data => {
        if (data.meals) {
          setSelectedMeals(data.meals);
        }
      });

    // Load meal items for tomorrow
    fetch(`http://localhost:8080/meals/items/${tomorrowStr}`, {
      headers: { Authorization: `Bearer ${token}` },
    })
      .then(res => res.json())
      .then(data => {
        setMealItems(data.items || {});
        setItemInputs(data.items || {});
      });

  }, [role]);

  const toggleMeal = (meal: string) => {
    if (selectedMeals.includes(meal)) {
      setSelectedMeals(selectedMeals.filter(m => m !== meal));
    } else {
      setSelectedMeals([...selectedMeals, meal]);
    }
  };

  // ================= SAVE OWN MEALS =================
  const saveOwnMeals = async () => {
    await fetch("http://localhost:8080/meals/update", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({ meals: selectedMeals }),
    });

    alert("Saved for tomorrow!");
  };

  // ================= SAVE MEAL ITEMS (ADMIN ONLY) =================
  const saveMealItems = async () => {
    await fetch("http://localhost:8080/meals/items/update", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({
        date: tomorrowStr,
        items: itemInputs,
      }),
    });

    alert("Meal items updated!");
  };

  // ================= OVERRIDE =================
  const overrideEmployee = async () => {
    if (!searchUser || !searchDate) {
      alert("Please enter employee username and date");
      return;
    }

    await fetch("http://localhost:8080/meals/override", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({
        username: searchUser,
        meals: selectedMeals,
        date: searchDate,
      }),
    });

    alert("Employee meals updated!");
  };

  // ================= HEADCOUNT =================
  const fetchHeadcount = async () => {
    if (!searchDate) {
      alert("Select a date first");
      return;
    }

    const res = await fetch(
      `http://localhost:8080/meals/headcount/${searchDate}`,
      {
        headers: { Authorization: `Bearer ${token}` },
      }
    );

    const data = await res.json();
    setHeadcount(data);
  };

  const logout = () => {
    localStorage.removeItem("token");
    navigate("/");
  };

  return (
    <div className="page">
      <div className="card">
        <div className="header">
          <h2>🍽️ Meal Planner</h2>
          <div>
            <span className="badge">{username}</span>
            <span className="badge role">{role}</span>
            <span className="badge role">{team}</span>
          </div>
        </div>

        {/* ================= MEAL SELECTION ================= */}
        {(role === "employee" || role === "teamlead" || role === "admin") && (
          <>
            <h3>Tomorrow's Meals</h3>

            <div className="meal-grid">
              {allMeals.map(meal => (
                <div key={meal} className="meal-card">
                  <label className="checkbox">
                    <input
                      type="checkbox"
                      checked={selectedMeals.includes(meal)}
                      onChange={() => toggleMeal(meal)}
                    />
                    {meal}
                  </label>

                  {mealItems[meal]?.length > 0 && (
                    <ul>
                      {mealItems[meal].map((item: string, i: number) => (
                        <li key={i}>{item}</li>
                      ))}
                    </ul>
                  )}

                  {role === "admin" && (
                    <input
                      className="input"
                      placeholder="Add items (comma separated)"
                      value={(itemInputs[meal] || []).join(", ")}
                      onChange={e =>
                        setItemInputs({
                          ...itemInputs,
                          [meal]: e.target.value
                            .split(",")
                            .map(s => s.trim()),
                        })
                      }
                    />
                  )}
                </div>
              ))}
            </div>

            <button className="btn primary" onClick={saveOwnMeals}>
              Save Selection
            </button>

            {role === "admin" && (
              <button className="btn secondary" onClick={saveMealItems}>
                Save Meal Items
              </button>
            )}
          </>
        )}

        {/* ================= ALL TEAMS PARTICIPATION (ADMIN ONLY) ================= */}
        {role === "admin" && (
          <div className="section">
            <h3>All Teams Participation</h3>

            <input
              className="input"
              type="date"
              value={adminDate}
              onChange={e => setAdminDate(e.target.value)}
            />

            <button
              className="btn primary"
              onClick={fetchAllTeamsParticipation}
            >
              Load Participation
            </button>

            {allTeamsData &&
              Object.entries(allTeamsData).map(([teamName, members]: any) => (
                <div key={teamName} className="team-card">
                  <div className="team-header">
                    <strong>{teamName}</strong>
                  </div>

                  {members.map((member: any, index: number) => (
                    <div key={index} style={{ marginBottom: "8px" }}>
                      <strong>{member.username}:</strong>{" "}
                      {member.meals?.join(", ") || "No meals"}
                    </div>
                  ))}
                </div>
              ))}
          </div>
        )}


        {/* ================= TEAM TODAY MEALS (TEAM LEAD ONLY) ================= */}
        {role === "teamlead" && (
          <div className="section">
            <button
              className="btn secondary"
              onClick={() => setShowTeamMeals(!showTeamMeals)}
            >
              {showTeamMeals ? "Hide Team Meal Status" : "View Team Meal Status"}
            </button>

            {showTeamMeals && (
              <>
                <h3 style={{ marginTop: "20px" }}>Today's Team Meal Status</h3>

                {teamMeals.length === 0 && <p>No data found</p>}

                {teamMeals.map((entry, index) => (
                  <div key={index} className="team-card">
                    <div className="team-header">
                      <strong>{entry.username}</strong>
                    </div>

                    <div className="team-meals">
                      {entry.meals?.length > 0 ? (
                        <p>{entry.meals.join(", ")}</p>
                      ) : (
                        <p>No meals selected</p>
                      )}
                    </div>
                  </div>
                ))}

              </>
            )}
          </div>
        )}



        {/* ================= OVERRIDE ================= */}
        {(role === "teamlead" || role === "admin") && (
          <div className="section">
            <h3>Modify Employee Meals</h3>

            <input
              className="input"
              placeholder="Employee Username"
              value={searchUser}
              onChange={e => setSearchUser(e.target.value)}
            />

            <input
              className="input"
              type="date"
              value={searchDate}
              onChange={e => setSearchDate(e.target.value)}
            />

            <button className="btn warning" onClick={overrideEmployee}>
              Update Employee Meals
            </button>
          </div>
        )}

        {/* ================= HEADCOUNT ================= */}
        {role === "admin" && (
          <div className="section">
            <h3>Headcount</h3>

            <input
              className="input"
              type="date"
              value={searchDate}
              onChange={e => setSearchDate(e.target.value)}
            />

            <button className="btn primary" onClick={fetchHeadcount}>
              Get Headcount
            </button>

            {headcount && (
              <div className="headcount">
                {Object.entries(headcount).map(([meal, count]) => (
                  <div key={meal} className="headcount-item">
                    <span>{meal}</span>
                    <span>{String(count)}</span>
                  </div>
                ))}
              </div>
            )}
          </div>
        )}

        {/* ================= TEAM LEAD BULK HANDLING ================= */}
        {role === "teamlead" && (
          <div className="section">
            <h3>Bulk & Exception Handling (Team)</h3>

            <input
              className="input"
              type="date"
              value={bulkDate}
              onChange={e => setBulkDate(e.target.value)}
            />

            <div className="meal-grid">
              {allMeals.map(meal => (
                <label key={meal} className="checkbox">
                  <input
                    type="checkbox"
                    checked={bulkMeals.includes(meal)}
                    onChange={() => toggleBulkMeal(meal)}
                  />
                  {meal}
                </label>
              ))}
            </div>

            <div style={{ marginTop: "12px" }}>
              <button
                className="btn warning"
                onClick={bulkOptOut}
                disabled={bulkLoading}
              >
                Bulk Opt-Out
              </button>

              <button
                className="btn primary"
                onClick={bulkOptIn}
                disabled={bulkLoading}
                style={{ marginLeft: "10px" }}
              >
                Bulk Opt-In
              </button>
            </div>
          </div>
        )}

        {/* ================= ADMIN BULK HANDLING (EVERYONE) ================= */}
        {role === "admin" && (
          <div className="section">
            <h3>Admin Bulk & Exception Handling (All Employees)</h3>

            <input
              className="input"
              type="date"
              value={adminBulkDate}
              onChange={e => setAdminBulkDate(e.target.value)}
            />

            <div className="meal-grid">
              {allMeals.map(meal => (
                <label key={meal} className="checkbox">
                  <input
                    type="checkbox"
                    checked={adminBulkMeals.includes(meal)}
                    onChange={() => toggleAdminBulkMeal(meal)}
                  />
                  {meal}
                </label>
              ))}
            </div>

            <div style={{ marginTop: "12px" }}>
              <button
                className="btn warning"
                onClick={adminBulkOptOut}
                disabled={adminBulkLoading}
              >
                Opt-Out For Everyone
              </button>

              <button
                className="btn primary"
                onClick={adminBulkOptIn}
                disabled={adminBulkLoading}
                style={{ marginLeft: "10px" }}
              >
                Opt-In For Everyone
              </button>
            </div>
          </div>
        )}

        <button className="btn danger logout" onClick={logout}>
          Logout
        </button>
      </div>
    </div>
  );

}