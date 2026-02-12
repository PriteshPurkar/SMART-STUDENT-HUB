import { useQuery } from "@tanstack/react-query";
import { fetchStudentDashboard } from "../../services/api";
import { Link } from "react-router-dom";

function StudentDashboard() {
  const { data, isLoading } = useQuery({
    queryKey: ["student-dashboard"],
    queryFn: fetchStudentDashboard
  });

  if (isLoading) {
    return <div className="page">Loading dashboard...</div>;
  }

  return (
    <div className="page">
      <h1>Student Dashboard</h1>
      <section>
        <h2>Upcoming Sessions</h2>
        <ul className="card-list">
          {data?.upcoming_sessions?.map((s: any) => (
            <li key={s.id} className="card">
              <h3>{s.title}</h3>
              <p>{s.description}</p>
              <p>Status: {s.status}</p>
              <Link to={`/student/sessions/${s.id}`}>View details</Link>
            </li>
          ))}
        </ul>
      </section>

      <section>
        <h2>Past Sessions</h2>
        <ul className="card-list">
          {data?.past_sessions?.map((s: any) => (
            <li key={s.id} className="card">
              <h3>{s.title}</h3>
              <p>{s.description}</p>
              <a href={s.video_url} target="_blank" rel="noreferrer">
                Watch recording
              </a>
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

