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

  const [specialDate, setSpecialDate] = useState("");
  const [specialType, setSpecialType] = useState("office_closed");
  const [specialNote, setSpecialNote] = useState("");
  const [currentDayStatus, setCurrentDayStatus] = useState<any>(null);

  // ================= COMPANY-WIDE WFH =================
const [companyWFHStart, setCompanyWFHStart] = useState("");
const [companyWFHEnd, setCompanyWFHEnd] = useState("");
const [companyWFHNote, setCompanyWFHNote] = useState("");
const [companyWFHLoading, setCompanyWFHLoading] = useState(false);
const [companyWFHMessage, setCompanyWFHMessage] = useState("");

const applyCompanyWFH = async () => {
  if (!companyWFHStart || !companyWFHEnd) {
    alert("Select start and end date");
    return;
  }

  setCompanyWFHLoading(true);
  setCompanyWFHMessage("");

  const res = await fetch("http://localhost:8080/admin/company-wfh", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({
      start_date: companyWFHStart,
      end_date: companyWFHEnd,
      note: companyWFHNote ? companyWFHNote : null,
    }),
  });

  const data = await res.json();
  setCompanyWFHLoading(false);

  if (!res.ok) {
    alert(data.error || "Failed to apply company WFH");
    return;
  }

  setCompanyWFHMessage(
    `Company-wide WFH applied for ${data.days_affected} days`
  );
};

  // ================= TEAM LEAD WORK LOCATION =================
  const [leadWorkDate, setLeadWorkDate] = useState("");
  const [leadWorkLocation, setLeadWorkLocation] = useState<"Office" | "WFH">("Office");
  const [leadWorkLoading, setLeadWorkLoading] = useState(false);
  const [leadWorkMessage, setLeadWorkMessage] = useState("");

  // Team member view/edit
  const [memberUsername, setMemberUsername] = useState("");
  const [memberDate, setMemberDate] = useState("");
  const [memberLocation, setMemberLocation] = useState<"Office" | "WFH">("Office");
  const [memberLoading, setMemberLoading] = useState(false);
  const [memberMessage, setMemberMessage] = useState("");
  const [memberLoadedText, setMemberLoadedText] = useState("");

  // Admin own & member work location
const [adminWorkDate, setAdminWorkDate] = useState("");
const [adminWorkLocation, setAdminWorkLocation] = useState<"Office" | "WFH">("Office");
const [adminWorkLoading, setAdminWorkLoading] = useState(false);
const [adminWorkMessage, setAdminWorkMessage] = useState("");

const [adminMemberUsername, setAdminMemberUsername] = useState("");
const [adminMemberDate, setAdminMemberDate] = useState("");
const [adminMemberLocation, setAdminMemberLocation] = useState<"Office" | "WFH">("Office");
const [adminMemberLoading, setAdminMemberLoading] = useState(false);
const [adminMemberMessage, setAdminMemberMessage] = useState("");
const [adminMemberLoadedText, setAdminMemberLoadedText] = useState("");

// ================= ADMIN ANNOUNCEMENT ================= //
const [announcementDate, setAnnouncementDate] = useState("");
const [announcementMsg, setAnnouncementMsg] = useState("");

const fetchAnnouncement = async () => {
  if (!announcementDate) return alert("Select a date");

  const res = await fetch(`http://localhost:8080/admin/headcount/summary/${announcementDate}`, {
    headers: { Authorization: `Bearer ${token}` },
  });

  if (!res.ok) {
    const err = await res.json();
    return alert(err.error || "Failed to fetch headcount");
  }

  const data = await res.json();
  const { total_participants, office, wfh, opted_out, by_meal, day_status, day_note } = data;

  let dayText = "";
  switch(day_status) {
    case "office_closed": dayText = "🏢 Office Closed"; break;
    case "government_holiday": dayText = "🎉 Government Holiday"; break;
    case "special_celebration": dayText = `🎊 Special Celebration${day_note ? `: ${day_note}` : ""}`; break;
    default: dayText = "Normal Day"; break;
  }

  let msg = `📅 Announcement for ${announcementDate}\n`;
  msg += `Status: ${dayText}\n`;
  msg += `Total Participants: ${total_participants} (Office: ${office}, WFH: ${wfh}, Opted Out: ${opted_out})\n\n`;
  msg += `🍽️ Meals Summary:\n`;

  for (const meal of Object.keys(by_meal)) {
    msg += `- ${meal}: ${by_meal[meal]}\n`;
  }

  setAnnouncementMsg(msg);
};

const copyAnnouncement = () => {
  navigator.clipboard.writeText(announcementMsg);
  alert("Announcement copied to clipboard!");
};

const fetchAdminWorkLocation = async (date: string) => {
  if (!date) return;
  
  const res = await fetch(
  `http://localhost:8080/admin/work-location?username=${adminMemberUsername}&date=${adminMemberDate}`,
  {
    method: "GET",
    headers: {
      Authorization: `Bearer ${token}`,
    },
  }
);

  const data = await res.json();
  if (res.ok) setAdminWorkLocation(data.location || "Office");
};

const saveAdminWorkLocation = async () => {
  if (!adminWorkDate) return alert("Select a date");

  if (isPastCutoff(adminWorkDate)) return alert("Cutoff passed.");

  setAdminWorkLoading(true);

  const res = await fetch("http://localhost:8080/admin/work-location/update", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({
      username,
      date: adminWorkDate,
      location: adminWorkLocation,
    }),
  });

  const data = await res.json();
  setAdminWorkLoading(false);

  if (!res.ok) return alert(data.error || "Failed");

  setAdminWorkMessage("Your work location updated");
};

const fetchAdminMemberWorkLocation = async () => {
  if (!adminMemberUsername || !adminMemberDate) return alert("Enter username and date");

  setAdminMemberLoading(true);
  setAdminMemberMessage("");
  setAdminMemberLoadedText("");

  const res = await fetch(`http://localhost:8080/admin/work-location?username=${adminMemberUsername}&date=${adminMemberDate}`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
  });

  const data = await res.json();
  setAdminMemberLoading(false);

  if (!res.ok) return alert(data.error || "Failed");

  setAdminMemberLocation(data.location || "Office");
  setAdminMemberLoadedText(`${adminMemberUsername} is working from ${data.location || "Office"} on ${adminMemberDate}`);
};

const updateAdminMemberWorkLocation = async () => {
  if (!adminMemberUsername || !adminMemberDate) return alert("Enter username and date");

  setAdminMemberLoading(true);

  const res = await fetch("http://localhost:8080/admin/work-location/update", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({
      username: adminMemberUsername,
      date: adminMemberDate,
      location: adminMemberLocation,
    }),
  });

  const data = await res.json();
  setAdminMemberLoading(false);

  if (!res.ok) return alert(data.error || "Failed");

  setAdminMemberMessage("Member work location updated");
};


  const _isPastCutoff = (selectedDate: string) => {
    if (!selectedDate) return true;

    const selected = new Date(selectedDate);
    const cutoff = new Date(selected);
    cutoff.setDate(cutoff.getDate() - 1);
    cutoff.setHours(21, 0, 0, 0); // 9 PM previous day

    return new Date() > cutoff;
  };

  const fetchLeadWorkLocation = async (date: string) => {
    if (!date) return;

    const res = await fetch("http://localhost:8080/me/work-location", {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({
        username,
        date,
      }),
    });

    const data = await res.json();
    if (res.ok) {
      setLeadWorkLocation(data.location || "Office");
    }
  };

  const saveLeadWorkLocation = async () => {
    if (!leadWorkDate) {
      alert("Select a date");
      return;
    }

    if (_isPastCutoff(leadWorkDate)) {
      alert("Cutoff time passed.");
      return;
    }

    setLeadWorkLoading(true);

    const res = await fetch("http://localhost:8080/me/work-location", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({
        username,
        date: leadWorkDate,
        location: leadWorkLocation,
      }),
    });

    const data = await res.json();
    setLeadWorkLoading(false);

    if (!res.ok) {
      alert(data.error || "Failed");
      return;
    }

    setLeadWorkMessage("Work location updated");
  };
  useEffect(() => {
    if (role === "teamlead" && leadWorkDate) {
      fetchLeadWorkLocation(leadWorkDate);
    }
  }, [leadWorkDate]);

  const fetchMemberWorkLocation = async () => {
    if (!memberUsername || !memberDate) {
      alert("Enter username and date");
      return;
    }

    setMemberLoading(true);
    setMemberMessage("");
    setMemberLoadedText(""); // clear previous text

    const res = await fetch("http://localhost:8080/teams/work-location", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({
        username: memberUsername,
        date: memberDate,
      }),
    });

    const data = await res.json();
    setMemberLoading(false);

    if (!res.ok) {
      alert(data.error || "Failed");
      return;
    }

    const location = data.location || "Office";
    setMemberLocation(location);

    setMemberLoadedText(
      `${memberUsername} is working from ${location} on ${memberDate}`
    );
  };
  const updateMemberWorkLocation = async () => {
    if (!memberUsername || !memberDate) {
      alert("Enter username and date");
      return;
    }

    setMemberLoading(true);

    const res = await fetch(
      "http://localhost:8080/teams/work-location/update",
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          username: memberUsername,
          date: memberDate,
          location: memberLocation,
        }),
      }
    );

    const data = await res.json();
    setMemberLoading(false);

    if (!res.ok) {
      alert(data.error || "Failed");
      return;
    }

    setMemberMessage("Member work location updated");
  };
  // ================= WORK LOCATION (EMPLOYEE) =================
  const [workDate, setWorkDate] = useState("");
  const [workLocation, setWorkLocation] = useState<"Office" | "WFH">("Office");
  const [workLoading, setWorkLoading] = useState(false);
  const [workMessage, setWorkMessage] = useState("");

  const isPastCutoff = (selectedDate: string) => {
    if (!selectedDate) return true;

    const selected = new Date(selectedDate);
    const cutoff = new Date(selected);
    cutoff.setDate(cutoff.getDate() - 1);
    cutoff.setHours(21, 0, 0, 0); // 9 PM previous day

    return new Date() > cutoff;
  };

  const fetchWorkLocation = async (date: string) => {
    if (!date) return;

    const res = await fetch("http://localhost:8080/me/work-location", {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({
        username,
        date,
      }),
    });

    const data = await res.json();

    if (res.ok) {
      setWorkLocation(data.location || "Office");
    }
  };

  const saveWorkLocation = async () => {
    if (!workDate) {
      alert("Select a date");
      return;
    }

    if (isPastCutoff(workDate)) {
      alert("Cutoff time passed. You cannot modify this date.");
      return;
    }

    setWorkLoading(true);

    const res = await fetch("http://localhost:8080/me/work-location", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({
        username,
        date: workDate,
        location: workLocation,
      }),
    });

    const data = await res.json();
    setWorkLoading(false);

    if (!res.ok) {
      alert(data.error || "Failed to save");
      return;
    }

    setWorkMessage("Work location updated successfully");

    // Optional: if WFH → reload tomorrow meals because backend opts out
    if (workLocation === "WFH") {
      setSelectedMeals([]);
    }
  };

  useEffect(() => {
    if (role === "employee" && workDate) {
      fetchWorkLocation(workDate);
    }
  }, [workDate]);

  // ========================= Special Day controls =========================//
  const setSpecialDay = async () => {
    if (!specialDate) {
      alert("Select a date");
      return;
    }

    const res = await fetch("http://localhost:8080/admin/day-controls", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({
        date: specialDate,
        type: specialType,
        note:
          specialType === "special_celebration" && specialNote
            ? specialNote
            : null,
      }),
    });

    const data = await res.json();

    if (!res.ok) {
      alert(data.error || "Failed");
      return;
    }

    alert("Special day saved!");
    fetchDayStatus();
  };
  const fetchDayStatus = async () => {
    if (!specialDate) return;

    const res = await fetch(
      `http://localhost:8080/admin/day-controls/${specialDate}`,
      {
        headers: { Authorization: `Bearer ${token}` },
      }
    );

    const data = await res.json();
    setCurrentDayStatus(data);
  };
  const removeSpecialDay = async () => {
    if (!specialDate) {
      alert("Select a date");
      return;
    }

    const res = await fetch(
      `http://localhost:8080/admin/day-controls/${specialDate}`,
      {
        method: "DELETE",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }
    );

    const data = await res.json();

    if (!res.ok) {
      alert(data.error || "Failed");
      return;
    }

    alert("Special day removed!");
    setCurrentDayStatus(null);
  };
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
        {/*================== Special Day controls ========================== */}
        {role === "admin" && (
          <div className="section">
            <h3>Special Day Controls</h3>

            <input
              className="input"
              type="date"
              value={specialDate}
              onChange={e => setSpecialDate(e.target.value)}
            />

            <select
              className="input"
              value={specialType}
              onChange={e => setSpecialType(e.target.value)}
            >
              <option value="office_closed">Office Closed</option>
              <option value="government_holiday">Government Holiday</option>
              <option value="special_celebration">
                Special Celebration Day
              </option>
            </select>

            {specialType === "special_celebration" && (
              <input
                className="input"
                placeholder="Celebration note"
                value={specialNote}
                onChange={e => setSpecialNote(e.target.value)}
              />
            )}

            <div style={{ marginTop: "10px" }}>
              <button className="btn primary" onClick={setSpecialDay}>
                Save
              </button>

              <button
                className="btn secondary"
                onClick={fetchDayStatus}
                style={{ marginLeft: "10px" }}
              >
                Check Status
              </button>

              <button
                className="btn danger"
                onClick={removeSpecialDay}
                style={{ marginLeft: "10px" }}
              >
                Remove
              </button>
            </div>

            {currentDayStatus && (
              <div style={{ marginTop: "15px" }}>
                <strong>Current Status:</strong>
                <p>Type: {currentDayStatus.type}</p>
                {currentDayStatus.note && (
                  <p>Note: {currentDayStatus.note}</p>
                )}
              </div>
            )}
          </div>
        )}

        {role === "employee" && (
          <div className="section">
            <h3>Work Location Per Date</h3>

            <input
              className="input"
              type="date"
              value={workDate}
              onChange={e => {
                setWorkDate(e.target.value);
                setWorkMessage("");
              }}
            />

            {workDate && (
              <>
                <div style={{ marginTop: "10px" }}>
                  <label className="radio">
                    <input
                      type="radio"
                      value="Office"
                      checked={workLocation === "Office"}
                      onChange={() => setWorkLocation("Office")}
                      disabled={isPastCutoff(workDate)}
                    />
                    Office
                  </label>

                  <label className="radio" style={{ marginLeft: "20px" }}>
                    <input
                      type="radio"
                      value="WFH"
                      checked={workLocation === "WFH"}
                      onChange={() => setWorkLocation("WFH")}
                      disabled={isPastCutoff(workDate)}
                    />
                    WFH
                  </label>
                </div>

                {isPastCutoff(workDate) && (
                  <p style={{ color: "red", marginTop: "8px" }}>
                    Cutoff time (9 PM previous day) has passed.
                  </p>
                )}

                <button
                  className="btn primary"
                  style={{ marginTop: "12px" }}
                  onClick={saveWorkLocation}
                  disabled={isPastCutoff(workDate) || workLoading}
                >
                  {workLoading ? "Saving..." : "Save Work Location"}
                </button>

                {workMessage && (
                  <p style={{ marginTop: "8px", color: "green" }}>
                    {workMessage}
                  </p>
                )}
              </>
            )}
          </div>
        )}
        {/* Team Lead work location ops */}
        {role === "teamlead" && (
          <div className="section">
            <h3>My Work Location</h3>

            <input
              className="input"
              type="date"
              value={leadWorkDate}
              onChange={e => {
                setLeadWorkDate(e.target.value);
                setLeadWorkMessage("");
              }}
            />

            {leadWorkDate && (
              <>
                <div style={{ marginTop: "10px" }}>
                  <label>
                    <input
                      type="radio"
                      checked={leadWorkLocation === "Office"}
                      onChange={() => setLeadWorkLocation("Office")}
                      disabled={isPastCutoff(leadWorkDate)}
                    />
                    Office
                  </label>

                  <label style={{ marginLeft: "20px" }}>
                    <input
                      type="radio"
                      checked={leadWorkLocation === "WFH"}
                      onChange={() => setLeadWorkLocation("WFH")}
                      disabled={isPastCutoff(leadWorkDate)}
                    />
                    WFH
                  </label>
                </div>

                {isPastCutoff(leadWorkDate) && (
                  <p style={{ color: "red" }}>
                    Cutoff (9 PM previous day) passed.
                  </p>
                )}

                <button
                  className="btn primary"
                  onClick={saveLeadWorkLocation}
                  disabled={isPastCutoff(leadWorkDate) || leadWorkLoading}
                  style={{ marginTop: "10px" }}
                >
                  Save
                </button>

                {leadWorkMessage && (
                  <p style={{ color: "green" }}>{leadWorkMessage}</p>
                )}
              </>
            )}

            <hr style={{ margin: "25px 0" }} />

            <h3>Team Member Work Location</h3>

            <input
              className="input"
              placeholder="Member Username"
              value={memberUsername}
              onChange={e => setMemberUsername(e.target.value)}
            />

            <input
              className="input"
              type="date"
              value={memberDate}
              onChange={e => setMemberDate(e.target.value)}
            />

            <button
              className="btn secondary"
              onClick={fetchMemberWorkLocation}
            >
              Load
            </button>
            {memberLoadedText && (
              <p style={{ marginTop: "10px", color: "blue" }}>{memberLoadedText}</p>
            )}
            {memberDate && memberUsername && (
              <>
                <div style={{ marginTop: "15px" }}>
                  <label>
                    <input
                      type="radio"
                      checked={memberLocation === "Office"}
                      onChange={() => setMemberLocation("Office")}
                    />
                    Office
                  </label>

                  <label style={{ marginLeft: "20px" }}>
                    <input
                      type="radio"
                      checked={memberLocation === "WFH"}
                      onChange={() => setMemberLocation("WFH")}
                    />
                    WFH
                  </label>
                </div>

                <button
                  className="btn warning"
                  onClick={updateMemberWorkLocation}
                  disabled={memberLoading}
                  style={{ marginTop: "10px" }}
                >
                  Update Member
                </button>

                {memberMessage && (
                  <p style={{ color: "green" }}>{memberMessage}</p>
                )}
              </>
            )}
          </div>
        )}
        {role === "admin" && (
  <div className="section">
    <h3>My Work Location</h3>
    <input
      type="date"
      className="input"
      value={adminWorkDate}
      onChange={e => {
        setAdminWorkDate(e.target.value);
        setAdminWorkMessage("");
        fetchAdminWorkLocation(e.target.value);
      }}
    />

    {adminWorkDate && (
      <>
        <div style={{ marginTop: "10px" }}>
          <label>
            <input
              type="radio"
              checked={adminWorkLocation === "Office"}
              onChange={() => setAdminWorkLocation("Office")}
              disabled={isPastCutoff(adminWorkDate)}
            /> Office
          </label>
          <label style={{ marginLeft: "20px" }}>
            <input
              type="radio"
              checked={adminWorkLocation === "WFH"}
              onChange={() => setAdminWorkLocation("WFH")}
              disabled={isPastCutoff(adminWorkDate)}
            /> WFH
          </label>
        </div>

        <button
          className="btn primary"
          style={{ marginTop: "10px" }}
          onClick={saveAdminWorkLocation}
          disabled={isPastCutoff(adminWorkDate) || adminWorkLoading}
        >
          {adminWorkLoading ? "Saving..." : "Save"}
        </button>

        {adminWorkMessage && <p style={{ color: "green" }}>{adminWorkMessage}</p>}
      </>
    )}

    <hr style={{ margin: "25px 0" }} />

    <h3>Member Work Location</h3>
    <input
      placeholder="Member Username"
      className="input"
      value={adminMemberUsername}
      onChange={e => setAdminMemberUsername(e.target.value)}
    />
    <input
      type="date"
      className="input"
      value={adminMemberDate}
      onChange={e => setAdminMemberDate(e.target.value)}
    />

    <button className="btn secondary" onClick={fetchAdminMemberWorkLocation}>
      Load
    </button>

    {adminMemberLoadedText && <p style={{ marginTop: "10px", color: "blue" }}>{adminMemberLoadedText}</p>}

    {adminMemberDate && adminMemberUsername && (
      <>
        <div style={{ marginTop: "15px" }}>
          <label>
            <input
              type="radio"
              checked={adminMemberLocation === "Office"}
              onChange={() => setAdminMemberLocation("Office")}
            /> Office
          </label>
          <label style={{ marginLeft: "20px" }}>
            <input
              type="radio"
              checked={adminMemberLocation === "WFH"}
              onChange={() => setAdminMemberLocation("WFH")}
            /> WFH
          </label>
        </div>

        <button
          className="btn warning"
          onClick={updateAdminMemberWorkLocation}
          disabled={adminMemberLoading}
          style={{ marginTop: "10px" }}
        >
          Update Member
        </button>

        {adminMemberMessage && <p style={{ color: "green" }}>{adminMemberMessage}</p>}
      </>
    )}
  </div>
)}
<div className="section">
  <h3>Generate Announcement</h3>

  <input
    type="date"
    className="input"
    value={announcementDate}
    onChange={e => setAnnouncementDate(e.target.value)}
  />

  <button className="btn primary" onClick={fetchAnnouncement} style={{ marginTop: "10px" }}>
    Generate Announcement
  </button>

  {announcementMsg && (
    <div style={{ marginTop: "15px" }}>
      <textarea
        className="input"
        value={announcementMsg}
        readOnly
        rows={10}
      />
      <button className="btn secondary" onClick={copyAnnouncement} style={{ marginTop: "10px" }}>
        Copy to Clipboard
      </button>
    </div>
  )}
</div>
{/* ================= COMPANY-WIDE WFH PERIOD ================= */}
{role === "admin" && (
  <div className="section">
    <h3>Company-wide WFH Period</h3>

    <input
      type="date"
      className="input"
      value={companyWFHStart}
      onChange={e => setCompanyWFHStart(e.target.value)}
    />

    <input
      type="date"
      className="input"
      value={companyWFHEnd}
      onChange={e => setCompanyWFHEnd(e.target.value)}
      style={{ marginTop: "10px" }}
    />

    <input
      className="input"
      placeholder="Optional note"
      value={companyWFHNote}
      onChange={e => setCompanyWFHNote(e.target.value)}
      style={{ marginTop: "10px" }}
    />

    <button
      className="btn primary"
      style={{ marginTop: "12px" }}
      onClick={applyCompanyWFH}
      disabled={companyWFHLoading}
    >
      {companyWFHLoading ? "Applying..." : "Apply Company WFH"}
    </button>

    {companyWFHMessage && (
      <p style={{ marginTop: "10px", color: "green" }}>
        {companyWFHMessage}
      </p>
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