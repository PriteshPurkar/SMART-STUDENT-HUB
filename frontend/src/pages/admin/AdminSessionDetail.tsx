import { FormEvent, useState } from "react";
import { useParams } from "react-router-dom";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  fetchSession,
  fetchSessionMaterials,
  fetchSessionStats,
  updateSessionStatus,
  uploadMaterial
} from "../../services/api";

function AdminSessionDetail() {
  const { id } = useParams<{ id: string }>();
  const queryClient = useQueryClient();

  const { data: session } = useQuery({
    queryKey: ["admin-session", id],
    queryFn: () => fetchSession(id || "")
  });

  const { data: materials } = useQuery({
    queryKey: ["admin-session-materials", id],
    queryFn: () => fetchSessionMaterials(id || ""),
    enabled: !!id
  });

  const { data: stats } = useQuery({
    queryKey: ["admin-session-stats", id],
    queryFn: () => fetchSessionStats(id || ""),
    enabled: !!id
  });

  const statusMutation = useMutation({
    mutationFn: (status: "SCHEDULED" | "ACTIVE" | "COMPLETED") =>
      updateSessionStatus(id || "", status),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin-session", id] });
    }
  });

  const materialMutation = useMutation({
    mutationFn: (payload: { title: string; url: string }) =>
      uploadMaterial(id || "", payload.title, payload.url),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ["admin-session-materials", id]
      });
    }
  });

  const [materialTitle, setMaterialTitle] = useState("");
  const [materialUrl, setMaterialUrl] = useState("");

  if (!session) {
    return <div className="page">Loading session...</div>;
  }

  return (
    <div className="page">
      <h1>Session Control: {session.title}</h1>
      <p>{session.description}</p>
      <p>Status: {session.status}</p>
      <p>Video URL: {session.video_url}</p>
      <section>
        <h2>Live Controls</h2>
        <button
          className="btn-primary"
          onClick={() => statusMutation.mutate("ACTIVE")}
          style={{ marginRight: "0.5rem" }}
        >
          Start Session
        </button>
        <button
          className="btn-primary"
          onClick={() => statusMutation.mutate("COMPLETED")}
        >
          Stop Session
        </button>
      </section>

      <section>
        <h2>Active Students</h2>
        <p>Currently active: {stats?.active_students ?? 0}</p>
      </section>

      <section>
        <h2>Upload Study Materials</h2>
        <form
          onSubmit={(e: FormEvent) => {
            e.preventDefault();
            materialMutation.mutate({ title: materialTitle, url: materialUrl });
            setMaterialTitle("");
            setMaterialUrl("");
          }}
        >
          <div className="form-group">
            <label>Title</label>
            <input
              value={materialTitle}
              onChange={e => setMaterialTitle(e.target.value)}
              required
            />
          </div>
          <div className="form-group">
            <label>URL</label>
            <input
              value={materialUrl}
              onChange={e => setMaterialUrl(e.target.value)}
              required
            />
          </div>
          <button type="submit" className="btn-primary">
            Add Material
          </button>
        </form>
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

export default AdminSessionDetail;

