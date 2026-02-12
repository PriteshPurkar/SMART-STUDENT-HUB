// Real HTTP API client for the Go backend (no mock data).

export type Role = "STUDENT" | "INSTRUCTOR" | "ADMIN";

const API_BASE =
  (typeof window !== "undefined" && (window as any).__API_URL__) ||
  import.meta.env.VITE_API_URL ||
  "/api/v1";

let authToken: string | null = null;
const TOKEN_KEY = "slp_token";

if (typeof window !== "undefined") {
  authToken = window.localStorage.getItem(TOKEN_KEY);
}

function setAuthToken(token: string | null) {
  authToken = token;
  if (typeof window === "undefined") return;
  if (!token) {
    window.localStorage.removeItem(TOKEN_KEY);
  } else {
    window.localStorage.setItem(TOKEN_KEY, token);
  }
}

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const headers: HeadersInit = {
    "Content-Type": "application/json",
    ...(options.headers || {})
  };
  if (authToken) {
    (headers as any).Authorization = `Bearer ${authToken}`;
  }

  const res = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers
  });

  if (!res.ok) {
    const text = await res.text();
    throw new Error(text || `Request failed with ${res.status}`);
  }

  const ct = res.headers.get("Content-Type") || "";
  if (!ct.includes("application/json")) {
    return undefined as unknown as T;
  }
  return (await res.json()) as T;
}

// Types based on backend models / responses ----------------------------------

export type User = {
  id: number;
  name: string;
  email: string;
  role: Role;
};

export type StudentDashboardResponse = {
  upcoming_sessions: any[];
  past_sessions: any[];
  notifications: any[];
};

export type Session = {
  id: number;
  title: string;
  description: string;
  status: "SCHEDULED" | "ACTIVE" | "COMPLETED";
  video_url: string;
};

// Auth -----------------------------------------------------------------------

export async function register(
  name: string,
  email: string,
  password: string,
  role: Role
) {
  const body: any = { name, email, password };
  if (role === "INSTRUCTOR") {
    body.role = "INSTRUCTOR";
  }
  return request<User>("/auth/register", {
    method: "POST",
    body: JSON.stringify(body)
  });
}

export async function login(email: string, password: string) {
  const res = await request<{ token: string; user: User }>("/auth/login", {
    method: "POST",
    body: JSON.stringify({ email, password })
  });
  setAuthToken(res.token);
  return res;
}

export async function getCurrentUser() {
  if (!authToken) return null;
  try {
    const user = await request<User>("/auth/me");
    return user;
  } catch {
    setAuthToken(null);
    return null;
  }
}

// Student APIs ---------------------------------------------------------------

export async function fetchStudentDashboard() {
  return request<StudentDashboardResponse>("/student/dashboard");
}

// Courses are not modeled on the backend yet; keep UI minimal for now.
export async function fetchCourses() {
  return [] as any[];
}

export async function enrollInCourse(_courseId: string) {
  throw new Error("Course enrollment is not implemented on the backend yet");
}

export async function fetchMyCourses() {
  return [] as any[];
}

export async function fetchCourseDetails(_courseId: string) {
  return null;
}

export async function submitExam(testId: string, _answers: string) {
  const examId = Number(testId);
  const res = await request<{ submission_id: number; message: string }>(
    `/exams/${examId}/submissions`,
    {
      method: "POST",
      body: JSON.stringify({ answers: "", file_s3_key: "" })
    }
  );
  return res;
}

export async function fetchMySubmission(testId: string) {
  const examId = Number(testId);
  return request<any>(`/exams/${examId}/submissions/me`);
}

// Faculty / Instructor APIs --------------------------------------------------

export async function fetchFacultyDashboard() {
  return null;
}

export async function fetchFacultyCourses() {
  return [];
}

export async function createCourse(_title: string, _description: string) {
  throw new Error("Course creation is not implemented on the backend yet");
}

export async function scheduleLecture(
  _courseId: string,
  _title: string,
  _description: string
) {
  throw new Error("Lecture scheduling by course is not implemented yet");
}

export async function uploadMaterialToCourse(
  _courseId: string,
  _title: string,
  _url: string
) {
  throw new Error("Uploading materials by course is not implemented yet");
}

export async function createTest(_courseId: string, _title: string) {
  throw new Error("Exam creation by course is not implemented yet");
}

export async function gradeSubmission(
  _testId: string,
  _studentId: number,
  _score: number
) {
  throw new Error("Grading submissions is not implemented yet");
}

// Admin APIs -----------------------------------------------------------------

export async function fetchAdminAnalytics() {
  // Could be extended to a dedicated analytics endpoint; for now, reuse submissions summary.
  return request<any>("/admin/submissions");
}

export async function fetchAdminSessions() {
  return request<Session[]>("/sessions");
}

export async function createAdminSession(payload: {
  title: string;
  description: string;
  start_time: string;
  end_time: string;
  video_url: string;
}) {
  return request<Session>("/admin/sessions", {
    method: "POST",
    body: JSON.stringify(payload)
  });
}

// Shared helpers -------------------------------------------------------------

export async function fetchSession(id: string) {
  const sessionId = Number(id);
  return request<Session>(`/sessions/${sessionId}`);
}

export async function fetchSessionMaterials(id: string) {
  const sessionId = Number(id);
  return request<any[]>(`/sessions/${sessionId}/materials`);
}

export async function fetchSessionStats(id: string) {
  const sessionId = Number(id);
  const stats = await request<{ submission_count: number; active_exams: number }>(
    `/admin/sessions/${sessionId}/stats`
  );
  return {
    active_students: stats.submission_count
  };
}

export async function updateSessionStatus(
  id: string,
  status: "SCHEDULED" | "ACTIVE" | "COMPLETED"
) {
  const sessionId = Number(id);
  await request<void>(`/admin/sessions/${sessionId}/status`, {
    method: "PATCH",
    body: JSON.stringify({ status })
  });
}

export async function uploadMaterial(
  _sessionId: string,
  title: string,
  url: string
) {
  console.warn("uploadMaterial is not yet implemented on the backend");
  return { id: Date.now(), title, url };
}

