import { useEffect } from "react";
import { useNavigate } from "react-router-dom";

export default function MealPlanner() {
    const navigate = useNavigate();

    useEffect(() => {
        const token = localStorage.getItem("token");

        fetch("http://localhost:8080/me", {
            headers: {
                Authorization: `Bearer ${token}`,
            },
        }).then(res => {
            if (!res.ok) navigate("/");
        });
    }, []);


    const handleLogout = () => {
        localStorage.removeItem("token");
        navigate("/");
    };


    return (
        <div style={{ padding: "2rem" }}>
            <h2>Meal Planner Page 🍽️</h2>
            <p>Welcome! You are logged in.</p>
            <button onClick={handleLogout}>Logout</button>
        </div>
    );
}
