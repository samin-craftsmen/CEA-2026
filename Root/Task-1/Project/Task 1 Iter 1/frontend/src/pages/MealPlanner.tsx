import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";

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
    <div style={{ padding: "2rem" }}>
      <h2>Meal Planner 🍽️</h2>

      <p><strong>User:</strong> {username}</p>
      <p><strong>Role:</strong> {role}</p>

      {/* ================= MEAL SELECTION ================= */}
      {(role === "employee" || role === "teamlead" || role === "admin") && (
        <>
          <h3>Tomorrow's Meals</h3>

          {allMeals.map(meal => (
            <div key={meal} style={{ marginBottom: "10px" }}>
              <label>
                <input
                  type="checkbox"
                  checked={selectedMeals.includes(meal)}
                  onChange={() => toggleMeal(meal)}
                />
                {meal}
              </label>

              {/* Show items */}
              {mealItems[meal] && mealItems[meal].length > 0 && (
                <ul style={{ marginLeft: "20px" }}>
                  {mealItems[meal].map((item: string, index: number) => (
                    <li key={index}>{item}</li>
                  ))}
                </ul>
              )}

              {/* Admin item editor */}
              {role === "admin" && (
                <input
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
                  style={{ marginLeft: "20px", width: "300px" }}
                />
              )}
            </div>
          ))}

          <br />
          <button onClick={saveOwnMeals}>Save</button>

          {role === "admin" && (
            <>
              <br /><br />
              <button onClick={saveMealItems}>
                Save Meal Items (Tomorrow)
              </button>
            </>
          )}
        </>
      )}

      {/* ================= OVERRIDE ================= */}
      {(role === "teamlead" || role === "admin") && (
        <>
          <hr />
          <h3>Modify Employee Meals</h3>

          <input
            placeholder="Employee Username"
            value={searchUser}
            onChange={e => setSearchUser(e.target.value)}
          />

          <br /><br />

          <input
            type="date"
            value={searchDate}
            onChange={e => setSearchDate(e.target.value)}
          />

          <br /><br />

          <button onClick={overrideEmployee}>
            Update Employee Meals
          </button>
        </>
      )}

      {/* ================= HEADCOUNT ================= */}
      {role === "admin" && (
        <>
          <hr />
          <h3>Headcount</h3>

          <input
            type="date"
            value={searchDate}
            onChange={e => setSearchDate(e.target.value)}
          />

          <br /><br />

          <button onClick={fetchHeadcount}>
            Get Headcount
          </button>

          {headcount && (
            <div style={{ marginTop: "1rem" }}>
              {Object.entries(headcount).map(([meal, count]) => (
                <p key={meal}>
                  {meal}: {String(count)}
                </p>
              ))}
            </div>
          )}
        </>
      )}

      <hr />
      <button onClick={logout}>Logout</button>
    </div>
  );
}