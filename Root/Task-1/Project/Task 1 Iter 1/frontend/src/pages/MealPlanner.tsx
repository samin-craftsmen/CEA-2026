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

      <button className="btn danger logout" onClick={logout}>
        Logout
      </button>
    </div>
  </div>
);

}