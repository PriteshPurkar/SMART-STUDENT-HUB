import { FormEvent, useState } from "react";
import { useNavigate } from "react-router-dom";
import { useQueryClient } from "@tanstack/react-query";
import { login, register, type Role } from "../services/api";

function LoginPage() {
  const [isRegister, setIsRegister] = useState(false);
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [role, setRole] = useState<Role>("STUDENT");
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setError(null);
    try {
      if (isRegister) {
        await register(name, email, password, role);
      }
      const res = await login(email, password);

      // Keep React Query auth state in sync so protected routes see the user.
      queryClient.setQueryData(["me"], res.user);

      if (res.user.role === "STUDENT") {
        navigate("/student/dashboard");
      } else {
        navigate("/admin/sessions");
      }
    } catch (err) {
      console.error(err);
      setError("Authentication failed");
    }
  };

  return (
    <div className="auth-page">
      <div className="auth-card">
        <h1>{isRegister ? "Student Registration" : "Login"}</h1>
        <form onSubmit={handleSubmit}>
          {isRegister && (
            <div className="form-group">
              <label>Name</label>
              <input value={name} onChange={e => setName(e.target.value)} required />
            </div>
          )}
          {isRegister && (
            <div className="form-group">
              <label>Role</label>
              <select value={role} onChange={e => setRole(e.target.value as Role)}>
                <option value="STUDENT">Student</option>
                <option value="FACULTY">Faculty</option>
                <option value="ADMIN">Admin</option>
              </select>
            </div>
          )}
          <div className="form-group">
            <label>Email</label>
            <input type="email" value={email} onChange={e => setEmail(e.target.value)} required />
          </div>
          <div className="form-group">
            <label>Password</label>
            <input type="password" value={password} onChange={e => setPassword(e.target.value)} required />
          </div>
          {error && <div className="error">{error}</div>}
          <button type="submit" className="btn-primary">
            {isRegister ? "Register" : "Login"}
          </button>
        </form>
        <button className="link-btn" onClick={() => setIsRegister(!isRegister)}>
          {isRegister ? "Already have an account? Login" : "New student? Register"}
        </button>
      </div>
    </div>
  );
}

export default LoginPage;

