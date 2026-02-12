import { FormEvent, useState } from "react";
import { useParams } from "react-router-dom";
import { submitExam } from "../../services/api";

function ExamPage() {
  const { id } = useParams<{ id: string }>();
  const [answers, setAnswers] = useState("");
  const [message, setMessage] = useState<string | null>(null);

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    if (!id) return;
    const res = await submitExam(id, answers);
    setMessage(`Submission ID ${res.submission_id}: ${res.message}`);
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
    </div>
  );
}

export default ExamPage;

