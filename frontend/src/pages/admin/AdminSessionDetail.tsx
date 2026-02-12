import { useParams } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import { fetchSession } from "../../services/api";

function AdminSessionDetail() {
  const { id } = useParams<{ id: string }>();
  const { data: session } = useQuery({
    queryKey: ["admin-session", id],
    queryFn: () => fetchSession(id || "")
  });

  if (!session) {
    return <div className="page">Loading session...</div>;
  }

  return (
    <div className="page">
      <h1>Session Control: {session.title}</h1>
      <p>{session.description}</p>
      <p>Status: {session.status}</p>
      <p>Video URL: {session.video_url}</p>
      <p>
        This is a scaffolded control panel. In a full implementation, you would
        start/stop the session, upload materials, and view real-time student
        counts from here.
      </p>
    </div>
  );
}

export default AdminSessionDetail;

