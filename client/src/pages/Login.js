import { useState } from "react";
import { useNavigate } from "react-router-dom";
import API from "../api";

export default function Login() {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const navigate = useNavigate();

    const login = async () => {
        try {
            const res = await API.post("/login", { email, password });
            localStorage.setItem("token", res.data.token);
            navigate("/");
        } catch (err) {
            alert("Login failed bro");
        }
    };

    return (
        <div className="container">
            <h2 className="center">Login</h2>
            <div className="form-group">
                <input placeholder="Email" onChange={(e) => setEmail(e.target.value)} />
                <input placeholder="Password" type="password" onChange={(e) => setPassword(e.target.value)} />
                <button onClick={login}>Login</button>
            </div>
        </div>
    );
}
