import { FormEvent, useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { fetchMySubmission, submitExam } from "../../services/api";

function ExamPage() {
  const { id } = useParams<{ id: string }>();
  const [answers, setAnswers] = useState("");
  const [message, setMessage] = useState<string | null>(null);
  const [result, setResult] = useState<any | null>(null);

  useEffect(() => {
    if (!id) return;
    fetchMySubmission(id).then(setResult);
  }, [id]);

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    if (!id) return;
    const res = await submitExam(id, answers);
    setMessage(`Submission ID ${res.submission_id}: ${res.message}`);
    const latest = await fetchMySubmission(id);
    setResult(latest);
  };

  return (
    <div className="page">
      <h1>Exam / Assignment</h1>
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label>Your Answers</label>
          <textarea value={answers} onChange={e => setAnswers(e.target.value)} required />
        </div>
        <button type="submit" className="btn-primary">
          Submit
        </button>
      </form>
      {message && <div className="success">{message}</div>}
      {result && (
        <div className="card" style={{ marginTop: "1rem" }}>
          <h2>Your Result</h2>
          <p>Status: {result.status}</p>
          <p>
            Score: {result.score != null ? result.score : "Pending evaluation"}
          </p>
        </div>
      )}
    </div>
  );
}

export default ExamPage;

