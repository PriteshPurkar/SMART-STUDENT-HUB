import { useParams } from "react-router-dom";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import {
  enrollInCourse,
  fetchCourseDetails,
  fetchMyCourses
} from "../../services/api";

function StudentCourseDetail() {
  const { id } = useParams<{ id: string }>();
  const queryClient = useQueryClient();

  const { data } = useQuery({
    queryKey: ["course-detail", id],
    queryFn: () => fetchCourseDetails(id || "")
  });

  const enrollMutation = useMutation({
    mutationFn: () => enrollInCourse(id || ""),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["student-dashboard"] });
      queryClient.invalidateQueries({ queryKey: ["student-my-courses"] });
    }
  });

  const { data: myCourses } = useQuery({
    queryKey: ["student-my-courses"],
    queryFn: fetchMyCourses
  });

  if (!data) {
    return <div className="page">Loading course...</div>;
  }

  const alreadyEnrolled = myCourses?.some((c: any) => c.id === data.course.id);

  return (
    <div className="page">
      <h1>{data.course.title}</h1>
      <p>{data.course.description}</p>

      {!alreadyEnrolled && (
        <button
          className="btn-primary"
          onClick={() => enrollMutation.mutate()}
          style={{ marginBottom: "1rem" }}
        >
          Pay & Enroll
        </button>
      )}

      <section>
        <h2>Study Materials</h2>
        <ul>
          {data.materials.map((m: any) => (
            <li key={m.id}>
              {m.title} -{" "}
              <a href={m.url} target="_blank" rel="noreferrer">
                Open
              </a>
            </li>
          ))}
        </ul>
      </section>

      <section>
        <h2>Live Lectures</h2>
        <ul>
          {data.sessions.map((s: any) => (
            <li key={s.id}>
              {s.title} ({s.status}) -{" "}
              <a href={s.videoUrl} target="_blank" rel="noreferrer">
                Join
              </a>
            </li>
          ))}
        </ul>
      </section>

      <section>
        <h2>Tests & Assignments</h2>
        <ul>
          {data.tests.map((t: any) => (
            <li key={t.id}>{t.title}</li>
          ))}
        </ul>
      </section>
    </div>
  );
}

export default StudentCourseDetail;

