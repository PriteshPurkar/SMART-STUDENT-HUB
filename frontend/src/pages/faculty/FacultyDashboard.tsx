import { FormEvent, useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { fetchFacultyDashboard, createCourse } from "../../services/api";
import { Link } from "react-router-dom";

function FacultyDashboard() {
  const queryClient = useQueryClient();
  const { data, isLoading } = useQuery({
    queryKey: ["faculty-dashboard"],
    queryFn: fetchFacultyDashboard
  });

  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");

  const courseMutation = useMutation({
    mutationFn: () => createCourse(title, description),
    onSuccess: () => {
      setTitle("");
      setDescription("");
      queryClient.invalidateQueries({ queryKey: ["faculty-dashboard"] });
    }
  });

  if (isLoading) {
    return <div className="page">Loading dashboard...</div>;
  }

  return (
    <div className="page">
      <h1>Faculty Dashboard</h1>

      {data?.analytics && (
        <section>
          <h2>Analytics</h2>
          <div className="card-list">
            <div className="card">
              <h3>Courses</h3>
              <p>{data.analytics.totalCourses}</p>
            </div>
            <div className="card">
              <h3>Students</h3>
              <p>{data.analytics.totalStudents}</p>
            </div>
            <div className="card">
              <h3>Sessions</h3>
              <p>{data.analytics.totalSessions}</p>
            </div>
          </div>
        </section>
      )}

      <section>
        <h2>Create Course</h2>
        <form
          onSubmit={(e: FormEvent) => {
            e.preventDefault();
            courseMutation.mutate();
          }}
        >
          <div className="form-group">
            <label>Title</label>
            <input
              value={title}
              onChange={e => setTitle(e.target.value)}
              required
            />
          </div>
          <div className="form-group">
            <label>Description</label>
            <textarea
              value={description}
              onChange={e => setDescription(e.target.value)}
              required
            />
          </div>
          <button type="submit" className="btn-primary">
            Create Course
          </button>
        </form>
      </section>

      <section>
        <h2>My Courses</h2>
        <ul className="card-list">
          {data?.my_courses?.map((c: any) => (
            <li key={c.id} className="card">
              <h3>{c.title}</h3>
              <p>{c.description}</p>
              <Link to={`/faculty/courses/${c.id}`}>Manage course</Link>
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

export default FacultyDashboard;

