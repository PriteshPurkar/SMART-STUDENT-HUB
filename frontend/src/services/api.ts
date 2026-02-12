// Full mock API layer with role-based users (STUDENT / FACULTY / ADMIN),
// courses, enrollments, sessions, materials, tests, submissions, and
// simple analytics for student, faculty, and admin dashboards.

export type Role = "STUDENT" | "FACULTY" | "ADMIN";

type User = {
  id: number;
  name: string;
  email: string;
  password: string;
  role: Role;
};

type Course = {
  id: number;
  title: string;
  description: string;
  facultyId: number;
};

type Enrollment = {
  id: number;
  courseId: number;
  studentId: number;
  status: "ENROLLED" | "PAID";
};

type Session = {
  id: number;
  courseId: number;
  title: string;
  description: string;
  status: "SCHEDULED" | "ACTIVE" | "COMPLETED";
  videoUrl: string;
};

type Material = {
  id: number;
  courseId: number;
  title: string;
  url: string;
};

type Test = {
  id: number;
  courseId: number;
  title: string;
};

type Submission = {
  id: number;
  testId: number;
  studentId: number;
  score: number | null;
  status: "SUBMITTED" | "GRADED";
};

type Notification = {
  id: number;
  userId: number;
  message: string;
  type: string;
};

let users: User[] = [
  {
    id: 1,
    name: "Demo Student",
    email: "student@example.com",
    password: "123",
    role: "STUDENT"
  },
  {
    id: 2,
    name: "Demo Faculty",
    email: "faculty@example.com",
    password: "123",
    role: "FACULTY"
  },
  {
    id: 3,
    name: "Demo Admin",
    email: "admin@example.com",
    password: "123",
    role: "ADMIN"
  }
];

let courses: Course[] = [
  {
    id: 1,
    title: "Scalable System Design",
    description: "Designing scalable learning platforms for high-traffic events.",
    facultyId: 2
  },
  {
    id: 2,
    title: "Distributed Systems Basics",
    description: "Fundamentals of distributed systems and consistency models.",
    facultyId: 2
  }
];

let enrollments: Enrollment[] = [];
let sessions: Session[] = [
  {
    id: 1,
    courseId: 1,
    title: "Live Kickoff Lecture",
    description: "Introduction to the course and architecture overview.",
    status: "SCHEDULED",
    videoUrl: "https://video.example.com/course/1/session/1"
  }
];

let materials: Material[] = [
  {
    id: 1,
    courseId: 1,
    title: "Course Syllabus (PDF)",
    url: "https://cdn.example.com/materials/syllabus.pdf"
  }
];

let tests: Test[] = [
  {
    id: 1,
    courseId: 1,
    title: "Midterm Assignment"
  }
];

let submissions: Submission[] = [];
let notifications: Notification[] = [];

let currentUser: User | null = null;

const USER_KEY = "slp_user";

// Hydrate currentUser from localStorage so route guards work after reload.
if (typeof window !== "undefined") {
  const stored = window.localStorage.getItem(USER_KEY);
  if (stored) {
    try {
      currentUser = JSON.parse(stored) as User;
    } catch {
      currentUser = null;
    }
  }
}

function persistCurrentUser() {
  if (typeof window === "undefined") return;
  if (!currentUser) {
    window.localStorage.removeItem(USER_KEY);
    return;
  }
  window.localStorage.setItem(USER_KEY, JSON.stringify(currentUser));
}

// Auth -----------------------------------------------------------------------

export async function register(
  name: string,
  email: string,
  password: string,
  role: Role
) {
  const existing = users.find(u => u.email === email);
  if (existing) {
    throw new Error("User already exists");
  }
  const id = users.length ? Math.max(...users.map(u => u.id)) + 1 : 1;
  const user: User = { id, name, email, password, role };
  users.push(user);
  currentUser = user;
  persistCurrentUser();
  return user;
}

export async function login(email: string, password: string) {
  const user = users.find(u => u.email === email && u.password === password);
  if (!user) {
    throw new Error("Invalid credentials");
  }
  currentUser = user;
  persistCurrentUser();
  return { token: "mock-token", user };
}

export async function getCurrentUser() {
  return currentUser;
}

// Student APIs ---------------------------------------------------------------

export async function fetchStudentDashboard() {
  if (!currentUser || currentUser.role !== "STUDENT") {
    return null;
  }
  const myEnrollments = enrollments.filter(e => e.studentId === currentUser!.id);
  const myCourseIds = new Set(myEnrollments.map(e => e.courseId));
  const myCourses = courses.filter(c => myCourseIds.has(c.id));

  const mySessions = sessions.filter(s => myCourseIds.has(s.courseId));
  const upcoming_sessions = mySessions.filter(s => s.status !== "COMPLETED");
  const past_sessions = mySessions.filter(s => s.status === "COMPLETED");

  const myNotifications = notifications.filter(n => n.userId === currentUser!.id);

  const analytics = {
    totalCourses: myCourses.length,
    totalEnrollments: myEnrollments.length,
    totalSubmissions: submissions.filter(
      sub => sub.studentId === currentUser!.id
    ).length
  };

  return {
    analytics,
    my_courses: myCourses,
    upcoming_sessions,
    past_sessions,
    notifications: myNotifications
  };
}

export async function fetchCourses() {
  return courses;
}

export async function enrollInCourse(courseId: string) {
  if (!currentUser || currentUser.role !== "STUDENT") {
    throw new Error("Only students can enroll");
  }
  const cid = Number(courseId);
  const existing = enrollments.find(
    e => e.courseId === cid && e.studentId === currentUser!.id
  );
  if (existing) return existing;

  const id = enrollments.length ? Math.max(...enrollments.map(e => e.id)) + 1 : 1;
  const enrollment: Enrollment = {
    id,
    courseId: cid,
    studentId: currentUser.id,
    status: "PAID"
  };
  enrollments.push(enrollment);

  notifications.push({
    id: Date.now(),
    userId: currentUser.id,
    message: "Payment successful and enrolled in course.",
    type: "ENROLLED"
  });

  return enrollment;
}

export async function fetchMyCourses() {
  if (!currentUser || currentUser.role !== "STUDENT") return [];
  const myEnrollments = enrollments.filter(e => e.studentId === currentUser.id);
  const ids = new Set(myEnrollments.map(e => e.courseId));
  return courses.filter(c => ids.has(c.id));
}

export async function fetchCourseDetails(courseId: string) {
  const cid = Number(courseId);
  const course = courses.find(c => c.id === cid);
  if (!course) return null;
  const courseMaterials = materials.filter(m => m.courseId === cid);
  const courseSessions = sessions.filter(s => s.courseId === cid);
  const courseTests = tests.filter(t => t.courseId === cid);
  return {
    course,
    materials: courseMaterials,
    sessions: courseSessions,
    tests: courseTests
  };
}

export async function submitExam(testId: string, answers: string) {
  if (!currentUser || currentUser.role !== "STUDENT") {
    throw new Error("Only students can submit");
  }
  const tid = Number(testId);
  const id = submissions.length ? Math.max(...submissions.map(s => s.id)) + 1 : 1;
  const submission: Submission = {
    id,
    testId: tid,
    studentId: currentUser.id,
    score: null,
    status: "SUBMITTED"
  };
  submissions.push(submission);

  notifications.push({
    id: Date.now(),
    userId: currentUser.id,
    message: `Submission received for test ${tid}.`,
    type: "SUBMISSION"
  });

  return {
    submission_id: submission.id,
    message: "Submission received successfully"
  };
}

export async function fetchMySubmission(testId: string) {
  if (!currentUser) return null;
  const tid = Number(testId);
  return (
    submissions.find(
      s => s.testId === tid && s.studentId === currentUser!.id
    ) ?? null
  );
}

// Faculty APIs ---------------------------------------------------------------

export async function fetchFacultyDashboard() {
  if (!currentUser || currentUser.role !== "FACULTY") return null;
  const myCourses = courses.filter(c => c.facultyId === currentUser!.id);
  const myCourseIds = new Set(myCourses.map(c => c.id));

  const myEnrollments = enrollments.filter(e => myCourseIds.has(e.courseId));
  const mySessions = sessions.filter(s => myCourseIds.has(s.courseId));

  const analytics = {
    totalCourses: myCourses.length,
    totalStudents: new Set(myEnrollments.map(e => e.studentId)).size,
    totalSessions: mySessions.length
  };

  const myNotifications = notifications.filter(n => n.userId === currentUser!.id);

  return {
    analytics,
    my_courses: myCourses,
    notifications: myNotifications
  };
}

export async function fetchFacultyCourses() {
  if (!currentUser || currentUser.role !== "FACULTY") return [];
  return courses.filter(c => c.facultyId === currentUser.id);
}

export async function createCourse(title: string, description: string) {
  if (!currentUser || currentUser.role !== "FACULTY") {
    throw new Error("Only faculty can create courses");
  }
  const id = courses.length ? Math.max(...courses.map(c => c.id)) + 1 : 1;
  const course: Course = {
    id,
    title,
    description,
    facultyId: currentUser.id
  };
  courses.push(course);
  return course;
}

export async function scheduleLecture(
  courseId: string,
  title: string,
  description: string
) {
  if (!currentUser || currentUser.role !== "FACULTY") {
    throw new Error("Only faculty can schedule lectures");
  }
  const cid = Number(courseId);
  const id = sessions.length ? Math.max(...sessions.map(s => s.id)) + 1 : 1;
  const session: Session = {
    id,
    courseId: cid,
    title,
    description,
    status: "SCHEDULED",
    videoUrl: `https://video.example.com/course/${cid}/session/${id}`
  };
  sessions.push(session);

  const enrolledStudents = enrollments
    .filter(e => e.courseId === cid)
    .map(e => e.studentId);
  enrolledStudents.forEach(studentId => {
    notifications.push({
      id: Date.now() + studentId,
      userId: studentId,
      message: `New lecture scheduled: ${title}`,
      type: "LECTURE_SCHEDULED"
    });
  });

  return session;
}

export async function uploadMaterialToCourse(
  courseId: string,
  title: string,
  url: string
) {
  if (!currentUser || currentUser.role !== "FACULTY") {
    throw new Error("Only faculty can upload materials");
  }
  const cid = Number(courseId);
  const id = materials.length ? Math.max(...materials.map(m => m.id)) + 1 : 1;
  const material: Material = { id, courseId: cid, title, url };
  materials.push(material);

  const enrolledStudents = enrollments
    .filter(e => e.courseId === cid)
    .map(e => e.studentId);
  enrolledStudents.forEach(studentId => {
    notifications.push({
      id: Date.now() + studentId,
      userId: studentId,
      message: `New material uploaded in ${title}`,
      type: "MATERIAL_UPLOADED"
    });
  });

  return material;
}

export async function createTest(courseId: string, title: string) {
  if (!currentUser || currentUser.role !== "FACULTY") {
    throw new Error("Only faculty can create tests");
  }
  const cid = Number(courseId);
  const id = tests.length ? Math.max(...tests.map(t => t.id)) + 1 : 1;
  const test: Test = { id, courseId: cid, title };
  tests.push(test);
  return test;
}

export async function gradeSubmission(
  testId: string,
  studentId: number,
  score: number
) {
  if (!currentUser || currentUser.role !== "FACULTY") {
    throw new Error("Only faculty can grade");
  }
  const tid = Number(testId);
  const sub =
    submissions.find(
      s => s.testId === tid && s.studentId === studentId
    ) ?? null;
  if (!sub) return null;
  sub.score = score;
  sub.status = "GRADED";
  return sub;
}

// Admin APIs -----------------------------------------------------------------

export async function fetchAdminAnalytics() {
  const totalUsers = users.length;
  const totalStudents = users.filter(u => u.role === "STUDENT").length;
  const totalFaculty = users.filter(u => u.role === "FACULTY").length;
  const totalAdmins = users.filter(u => u.role === "ADMIN").length;
  const totalCourses = courses.length;
  const totalEnrollments = enrollments.length;

  return {
    totalUsers,
    totalStudents,
    totalFaculty,
    totalAdmins,
    totalCourses,
    totalEnrollments
  };
}

// Shared helpers used by multiple roles -------------------------------------

// Sessions list for admin view (all sessions across courses).
export async function fetchAdminSessions() {
  return sessions;
}

// Create a new live session not tied to a specific course (admin shortcut).
export async function createAdminSession(payload: {
  title: string;
  description: string;
  start_time: string;
  end_time: string;
  video_url: string;
}) {
  const id = sessions.length ? Math.max(...sessions.map(s => s.id)) + 1 : 1;
  const session: Session = {
    id,
    courseId: 1,
    title: payload.title,
    description: payload.description,
    status: "SCHEDULED",
    videoUrl: payload.video_url
  };
  sessions.push(session);
  return session;
}

// Fetch a single session by id (used by student/admin session detail pages).
export async function fetchSession(id: string) {
  const sid = Number(id);
  return sessions.find(s => s.id === sid) ?? null;
}

// For compatibility with existing UI, return materials for the session's course.
export async function fetchSessionMaterials(id: string) {
  const sid = Number(id);
  const session = sessions.find(s => s.id === sid);
  if (!session) return [];
  return materials.filter(m => m.courseId === session.courseId);
}

// Simple session stats mock used in admin session detail.
export async function fetchSessionStats(id: string) {
  const sid = Number(id);
  const session = sessions.find(s => s.id === sid);
  if (!session) {
    return { active_students: 0 };
  }
  const enrolledCount = enrollments.filter(
    e => e.courseId === session.courseId
  ).length;
  return {
    active_students: enrolledCount
  };
}

// Update session status (used by admin controls).
export async function updateSessionStatus(
  id: string,
  status: "SCHEDULED" | "ACTIVE" | "COMPLETED"
) {
  const sid = Number(id);
  const session = sessions.find(s => s.id === sid);
  if (!session) return null;
  session.status = status;
  return session;
}

// Upload material for a specific session (admin view convenience).
export async function uploadMaterial(
  sessionId: string,
  title: string,
  url: string
) {
  const sid = Number(sessionId);
  const session = sessions.find(s => s.id === sid);
  if (!session) return null;
  return uploadMaterialToCourse(String(session.courseId), title, url);
}




