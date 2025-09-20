import React, { useState } from "react";
import { useAuth } from "../../contexts/AuthContext";
import { useNavigate } from "react-router-dom";

export const LoginForm: React.FC = () => {
  const { login } = useAuth();
  const navigate = useNavigate();

  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [err, setErr] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setErr(null);
    setLoading(true);
    try {
      await login(username, password);
      navigate("/home");                                 // redirect after successful login
    } catch (e: any) {
      setErr(e.message || "Login failed");
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={onSubmit}>
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
        <label>Password</label>
        <input
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
        />
      </div>

      <button type="submit" disabled={loading}>
        {loading ? "Logging in..." : "Login"}
      </button>
    </form>
  );
};
