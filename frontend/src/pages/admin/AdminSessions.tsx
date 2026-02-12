import { FormEvent, useState } from "react";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { createAdminSession, fetchAdminSessions } from "../../services/api";
import { Link } from "react-router-dom";

function AdminSessions() {
  const queryClient = useQueryClient();
  const { data: sessions } = useQuery({
    queryKey: ["admin-sessions"],
    queryFn: fetchAdminSessions
  });

  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [videoUrl, setVideoUrl] = useState("");

  const handleCreate = async (e: FormEvent) => {
    e.preventDefault();
    const now = new Date();
    await createAdminSession({
      title,
      description,
      video_url: videoUrl,
      start_time: now.toISOString(),
      end_time: new Date(now.getTime() + 60 * 60 * 1000).toISOString()
    });
    setTitle("");
    setDescription("");
    setVideoUrl("");
    queryClient.invalidateQueries({ queryKey: ["admin-sessions"] });
  };

  return (
    <div className="page">
      <h1>Admin Sessions</h1>
      <section>
        <h2>Create Live Session</h2>
        <form onSubmit={handleCreate}>
          <div className="form-group">
            <label>Title</label>
            <input value={title} onChange={e => setTitle(e.target.value)} required />
          </div>
          <div className="form-group">
            <label>Description</label>
            <textarea value={description} onChange={e => setDescription(e.target.value)} required />
          </div>
          <div className="form-group">
            <label>Video URL</label>
            <input value={videoUrl} onChange={e => setVideoUrl(e.target.value)} required />
          </div>
          <button type="submit" className="btn-primary">
            Create Session
          </button>
        </form>
      </section>

      <section>
        <h2>Existing Sessions</h2>
        <ul className="card-list">
          {sessions?.map((s: any) => (
            <li key={s.id} className="card">
              <h3>{s.title}</h3>
              <p>{s.description}</p>
              <p>Status: {s.status}</p>
              <Link to={`/admin/sessions/${s.id}`}>Open control panel</Link>
            </li>
          ))}
        </ul>
      </section>
    </div>
  );
}

export default AdminSessions;

