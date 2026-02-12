import { useQuery } from "@tanstack/react-query";
import { fetchAdminAnalytics } from "../../services/api";

function AdminDashboard() {
  const { data, isLoading } = useQuery({
    queryKey: ["admin-analytics"],
    queryFn: fetchAdminAnalytics
  });

  if (isLoading) {
    return <div className="page">Loading admin analytics...</div>;
  }

  return (
    <div className="page">
      <h1>Admin Analytics & Management</h1>
      <section>
        <h2>Platform Overview</h2>
        <div className="card-list">
          <div className="card">
            <h3>Total Users</h3>
            <p>{data?.totalUsers}</p>
          </div>
          <div className="card">
            <h3>Students</h3>
            <p>{data?.totalStudents}</p>
          </div>
          <div className="card">
            <h3>Faculty</h3>
            <p>{data?.totalFaculty}</p>
          </div>
          <div className="card">
            <h3>Admins</h3>
            <p>{data?.totalAdmins}</p>
          </div>
          <div className="card">
            <h3>Courses</h3>
            <p>{data?.totalCourses}</p>
          </div>
          <div className="card">
            <h3>Enrollments</h3>
            <p>{data?.totalEnrollments}</p>
          </div>
        </div>
      </section>
      <p>
        This mock dashboard summarizes key metrics. In a real system, you would
        extend this with user management, course approvals, and detailed
        reports.
      </p>
    </div>
  );
}

export default AdminDashboard;

