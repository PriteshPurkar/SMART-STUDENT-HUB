import { Navigate, Route, Routes } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import { getCurrentUser } from "./services/api";
import LoginPage from "./pages/LoginPage";
import StudentDashboard from "./pages/student/StudentDashboard";
import SessionDetail from "./pages/student/SessionDetail";
import ExamPage from "./pages/student/ExamPage";
import StudentCourseDetail from "./pages/student/StudentCourseDetail";
import FacultyDashboard from "./pages/faculty/FacultyDashboard";
import FacultyCourseDetail from "./pages/faculty/FacultyCourseDetail";
import AdminDashboard from "./pages/admin/AdminDashboard";
import AdminSessions from "./pages/admin/AdminSessions";
import AdminSessionDetail from "./pages/admin/AdminSessionDetail";

function App() {
  const { data: user, isLoading } = useQuery({
    queryKey: ["me"],
    queryFn: getCurrentUser,
    retry: false
  });

  if (isLoading) {
    return <div className="page-center">Loading...</div>;
  }

  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route
        path="/student/dashboard"
        element={
          <RequireRole user={user} allowed={["STUDENT"]}>
            <StudentDashboard />
          </RequireRole>
        }
      />
      <Route
        path="/student/sessions/:id"
        element={
          <RequireRole user={user} allowed={["STUDENT"]}>
            <SessionDetail />
          </RequireRole>
        }
      />
      <Route
        path="/student/courses/:id"
        element={
          <RequireRole user={user} allowed={["STUDENT"]}>
            <StudentCourseDetail />
          </RequireRole>
        }
      />
      <Route
        path="/student/exams/:id"
        element={
          <RequireRole user={user} allowed={["STUDENT"]}>
            <ExamPage />
          </RequireRole>
        }
      />
      <Route
        path="/faculty/dashboard"
        element={
          <RequireRole user={user} allowed={["FACULTY"]}>
            <FacultyDashboard />
          </RequireRole>
        }
      />
      <Route
        path="/faculty/courses/:id"
        element={
          <RequireRole user={user} allowed={["FACULTY"]}>
            <FacultyCourseDetail />
          </RequireRole>
        }
      />
      <Route
        path="/admin/dashboard"
        element={
          <RequireRole user={user} allowed={["ADMIN"]}>
            <AdminDashboard />
          </RequireRole>
        }
      />
      <Route
        path="/admin/sessions"
        element={
          <RequireRole user={user} allowed={["ADMIN", "INSTRUCTOR"]}>
            <AdminSessions />
          </RequireRole>
        }
      />
      <Route
        path="/admin/sessions/:id"
        element={
          <RequireRole user={user} allowed={["ADMIN", "INSTRUCTOR"]}>
            <AdminSessionDetail />
          </RequireRole>
        }
      />
      <Route path="*" element={<Navigate to="/login" />} />
    </Routes>
  );
}

type RequireRoleProps = {
  user: any;
  allowed: string[];
  children: JSX.Element;
};

function RequireRole({ user, allowed, children }: RequireRoleProps) {
  if (!user) {
    return <Navigate to="/login" replace />;
  }
  if (!allowed.includes(user.role)) {
    return <Navigate to="/login" replace />;
  }
  return children;
}

export default App;

