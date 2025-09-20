import React, { useState } from "react";
import { useAuth } from "../../contexts/AuthContext";
import { useNavigate } from "react-router-dom";

export const SignupForm: React.FC = () => {
  const { signup } = useAuth();
  const navigate = useNavigate();

  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [err, setErr] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setErr(null);
    setLoading(true);
    try {
      await signup(username, email, password);
      navigate("/home"); // redirect after signup
    } catch (e: any) {
      setErr(e.message || "Signup failed");
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={onSubmit}>
      {/* <h2>Signup</h2> */}
      {err && <div style={{ color: "red" }}>{err}</div>}

      <div>
        <label>Username</label>
        <input
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          autoComplete="username"
          required
        />
      </div>

      <div>
        <label>Email</label>
        <input
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          autoComplete="email"
          required
        />
      </div>

      <div>
        <label>Password</label>
        <input
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
        />
      </div>

      <button type="submit" disabled={loading}>
        {loading ? "Signing up..." : "Signup"}
      </button>
    </form>
  );
};
