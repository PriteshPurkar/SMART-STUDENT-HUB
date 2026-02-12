import { FormEvent, useState } from "react";
import { useParams } from "react-router-dom";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import {
  createTest,
  fetchCourseDetails,
  scheduleLecture,
  uploadMaterialToCourse
} from "../../services/api";

function FacultyCourseDetail() {
  const { id } = useParams<{ id: string }>();
  const queryClient = useQueryClient();

  const { data } = useQuery({
    queryKey: ["faculty-course-detail", id],
    queryFn: () => fetchCourseDetails(id || "")
  });

  const [lectureTitle, setLectureTitle] = useState("");
  const [lectureDesc, setLectureDesc] = useState("");
  const [materialTitle, setMaterialTitle] = useState("");
  const [materialUrl, setMaterialUrl] = useState("");
  const [testTitle, setTestTitle] = useState("");

  const lectureMutation = useMutation({
    mutationFn: () =>
      scheduleLecture(id || "", lectureTitle, lectureDesc),
    onSuccess: () => {
      setLectureTitle("");
      setLectureDesc("");
      queryClient.invalidateQueries({
        queryKey: ["faculty-course-detail", id]
      });
    }
  });

  const materialMutation = useMutation({
    mutationFn: () =>
      uploadMaterialToCourse(id || "", materialTitle, materialUrl),
    onSuccess: () => {
      setMaterialTitle("");
      setMaterialUrl("");
      queryClient.invalidateQueries({
        queryKey: ["faculty-course-detail", id]
      });
    }
  });

  const testMutation = useMutation({
    mutationFn: () => createTest(id || "", testTitle),
    onSuccess: () => {
      setTestTitle("");
      queryClient.invalidateQueries({
        queryKey: ["faculty-course-detail", id]
      });
    }
  });

  if (!data) {
    return <div className="page">Loading course...</div>;
  }

  return (
    <div className="page">
      <h1>{data.course.title}</h1>
      <p>{data.course.description}</p>

      <section>
        <h2>Plan / Schedule Lecture</h2>
        <form
          onSubmit={(e: FormEvent) => {
            e.preventDefault();
            lectureMutation.mutate();
          }}
        >
          <div className="form-group">
            <label>Title</label>
            <input
              value={lectureTitle}
              onChange={e => setLectureTitle(e.target.value)}
              required
            />
          </div>
          <div className="form-group">
            <label>Description</label>
            <textarea
              value={lectureDesc}
              onChange={e => setLectureDesc(e.target.value)}
              required
            />
          </div>
          <button type="submit" className="btn-primary">
            Schedule Lecture
          </button>
        </form>
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
        <h2>Upload Study Materials</h2>
        <form
          onSubmit={(e: FormEvent) => {
            e.preventDefault();
            materialMutation.mutate();
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
        <h2>Create Test / Assignment</h2>
        <form
          onSubmit={(e: FormEvent) => {
            e.preventDefault();
            testMutation.mutate();
          }}
        >
          <div className="form-group">
            <label>Title</label>
            <input
              value={testTitle}
              onChange={e => setTestTitle(e.target.value)}
              required
            />
          </div>
          <button type="submit" className="btn-primary">
            Create Test
          </button>
        </form>
        <ul>
          {data.tests.map((t: any) => (
            <li key={t.id}>{t.title}</li>
          ))}
        </ul>
      </section>
    </div>
  );
}

export default FacultyCourseDetail;

