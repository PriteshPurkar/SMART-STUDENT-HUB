import { FormEvent, useState } from "react";
import { useNavigate } from "react-router-dom";
import { login, register } from "../services/api";

function LoginPage() {
  const [isRegister, setIsRegister] = useState(false);
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setError(null);
    try {
      if (isRegister) {
        await register(name, email, password);
      }
      const res = await login(email, password);
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

