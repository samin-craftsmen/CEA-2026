import { Routes, Route } from "react-router-dom";
import Login from "./pages/Login";
import MealPlanner from "./pages/MealPlanner";

function App() {
  return (
    <Routes>
      <Route path="/" element={<Login />} />
      <Route path="/meal-planner" element={<MealPlanner />} />
    </Routes>
  );
}

export default App;
