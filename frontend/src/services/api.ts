// NOTE: Because Go and Docker are not available in this environment,
// these API helpers are implemented as in-memory mocks so that the
// React app runs without hitting a real backend (and without Vite
// proxy errors like ECONNREFUSED). When you have a running Go API,
// replace these with real HTTP calls pointing at /api/v1/*.

let currentUser: any | null = null;
let sessions: any[] = [
  {
    id: 1,
    title: "Live Exam Prep",
    description: "Final preparation for exams",
    status: "SCHEDULED",
    video_url: "https://video.example.com/session/1"
  }
];

export async function login(email: string, password: string) {
  currentUser = {
    id: 1,
    name: "Demo Student",
    email,
    role: "STUDENT"
  };
  return { token: "mock-token", user: currentUser };
}

export async function register(name: string, email: string, password: string) {
  currentUser = {
    id: 1,
    name,
    email,
    role: "STUDENT"
  };
  return currentUser;
}

export async function getCurrentUser() {
  return currentUser;
}

export async function fetchStudentDashboard() {
  return {
    upcoming_sessions: sessions,
    past_sessions: [],
    notifications: [
      {
        id: 1,
        message: "Your exam prep session starts in 30 minutes."
      }
    ]
  };
}

export async function fetchSessions() {
  return sessions;
}

export async function fetchSession(id: string) {
  return sessions.find(s => String(s.id) === String(id));
}

export async function fetchSessionMaterials(id: string) {
  return [
    {
      id: 1,
      title: "Exam Syllabus",
      url: "https://cdn.example.com/materials/exam-syllabus.pdf"
    }
  ];
}

export async function submitExam(id: string, answers: string) {
  return {
    submission_id: Date.now(),
    message: "Submission received successfully"
  };
}

export async function fetchAdminSessions() {
  return sessions;
}

export async function createAdminSession(payload: {
  title: string;
  description: string;
  start_time: string;
  end_time: string;
  video_url: string;
}) {
  const id = Date.now();
  const s = { id, status: "SCHEDULED", ...payload };
  sessions.push(s);
  return s;
}

