import { useEffect } from "react";
import { useParams } from "react-router-dom";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { fetchSession, fetchSessionMaterials } from "../../services/api";

function SessionDetail() {
  const { id } = useParams<{ id: string }>();
  const queryClient = useQueryClient();

  const { data: session } = useQuery({
    queryKey: ["session", id],
    queryFn: () => fetchSession(id || "")
  });

  const { data: materials } = useQuery({
    queryKey: ["session-materials", id],
    queryFn: () => fetchSessionMaterials(id || ""),
    enabled: !!id
  });

  useEffect(() => {
    if (!id) return;

    // Subscribe to realtime session status updates via SSE.
    const source = new EventSource("/api/v1/realtime/stream");
    source.onmessage = event => {
      try {
        const parsed = JSON.parse(event.data);
        if (
          parsed.type === "session_status_updated" &&
          String(parsed.data.session_id) === String(id)
        ) {
          queryClient.invalidateQueries({ queryKey: ["session", id] });
        }
      } catch {
        // ignore malformed events
      }
    };

    return () => {
      source.close();
    };
  }, [id, queryClient]);

  if (!session) {
    return <div className="page">Loading session...</div>;
  }

  return (
    <div className="page">
      <h1>{session.title}</h1>
      <p>{session.description}</p>
      <p>Status: {session.status}</p>
      <a href={session.video_url} target="_blank" rel="noreferrer" className="btn-primary">
        Join / Watch Session
      </a>

      <section>
        <h2>Study Materials</h2>
        <ul>
          {materials?.map((m: any) => (
            <li key={m.id}>
              {m.title} -{" "}
              <a href={m.url} target="_blank" rel="noreferrer">
                Open
              </a>
            </li>
          ))}
        </ul>
      </section>
    </div>
  );
}

export default SessionDetail;

