import React from "react";
import { LoginForm } from "../components/forms/LoginForm";
import { Link } from "react-router-dom";

const LoginPage: React.FC = () => {
  return (
    <div style={{ padding: "2rem" }}>
      <h1>Log In</h1>
      <LoginForm />
      <p>
        Don’t have an account? <Link to="/signup">Sign up</Link>
      </p>
    </div>
  );
};

export default LoginPage;
