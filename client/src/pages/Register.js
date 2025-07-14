import { useState } from "react";
import { useNavigate } from "react-router-dom";
import API from "../api";

export default function Register() {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const navigate = useNavigate();

    const register = async () => {
        try {
            const res = await API.post("/register", { email, password });
            localStorage.setItem("token", res.data.token);
            navigate("/");
        } catch (err) {
            alert("Registration failed")
        }
    };

    return (
        <div className="container">
            <h2 className="center">Register</h2>
            <div className="form-group">
                <input placeholder="Email" onChange={(e) => setEmail(e.target.value)} />
                <input placeholder="Password" type="password" onChange={(e) => setPassword(e.target.value)} />
                <button onClick={register}>Register</button>
            </div>
        </div>
    );
}