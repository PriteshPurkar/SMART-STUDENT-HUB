import { useQuery } from "@tanstack/react-query";
import { fetchCourses, fetchMyCourses, fetchStudentDashboard } from "../../services/api";
import { Link } from "react-router-dom";

function StudentDashboard() {
  const { data, isLoading } = useQuery({
    queryKey: ["student-dashboard"],
    queryFn: fetchStudentDashboard
  });

  const { data: myCourses } = useQuery({
    queryKey: ["student-my-courses"],
    queryFn: fetchMyCourses
  });

  const { data: allCourses } = useQuery({
    queryKey: ["all-courses"],
    queryFn: fetchCourses
  });

  if (isLoading) {
    return <div className="page">Loading dashboard...</div>;
  }

  const myCourseIds = new Set((myCourses || []).map((c: any) => c.id));
  const availableCourses = (allCourses || []).filter(
    (c: any) => !myCourseIds.has(c.id)
  );

  return (
    <div className="page">
      <h1>Student Dashboard</h1>

      {data?.analytics && (
        <section>
          <h2>My Analytics</h2>
          <div className="card-list">
            <div className="card">
              <h3>Courses Enrolled</h3>
              <p>{data.analytics.totalCourses}</p>
            </div>
            <div className="card">
              <h3>Enrollments</h3>
              <p>{data.analytics.totalEnrollments}</p>
            </div>
            <div className="card">
              <h3>Submissions</h3>
              <p>{data.analytics.totalSubmissions}</p>
            </div>
          </div>
        </section>
      )}

      <section>
        <h2>My Courses</h2>
        <ul className="card-list">
          {myCourses?.map((c: any) => (
            <li key={c.id} className="card">
              <h3>{c.title}</h3>
              <p>{c.description}</p>
              <Link to={`/student/courses/${c.id}`}>Open course</Link>
            </li>
          ))}
        </ul>
      </section>

      <section>
        <h2>Available Courses</h2>
        <ul className="card-list">
          {availableCourses.map((c: any) => (
            <li key={c.id} className="card">
              <h3>{c.title}</h3>
              <p>{c.description}</p>
              <Link to={`/student/courses/${c.id}`}>View & enroll</Link>
            </li>
          ))}
        </ul>
      </section>

      <section>
        <h2>Upcoming Sessions</h2>
        <ul className="card-list">
          {data?.upcoming_sessions?.map((s: any) => (
            <li key={s.id} className="card">
              <h3>{s.title}</h3>
              <p>{s.description}</p>
              <p>Status: {s.status}</p>
            </li>
          ))}
        </ul>
      </section>

      <section>
        <h2>Notifications</h2>
        <ul>
          {data?.notifications?.map((n: any) => (
            <li key={n.id}>{n.message}</li>
          ))}
        </ul>
      </section>
    </div>
  );
}

export default StudentDashboard;

