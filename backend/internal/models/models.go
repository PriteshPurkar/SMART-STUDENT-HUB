package models

import "time"

type Role string

const (
	RoleStudent   Role = "STUDENT"
	RoleInstructor     = "INSTRUCTOR"
	RoleAdmin          = "ADMIN"
)

type User struct {
	ID        int64     `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password_hash" json:"-"`
	Role      Role      `db:"role" json:"role"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type SessionStatus string

const (
	SessionScheduled SessionStatus = "SCHEDULED"
	SessionActive    SessionStatus = "ACTIVE"
	SessionCompleted SessionStatus = "COMPLETED"
)

type Session struct {
	ID          int64         `db:"id" json:"id"`
	Title       string        `db:"title" json:"title"`
	Description string        `db:"description" json:"description"`
	StartTime   time.Time     `db:"start_time" json:"start_time"`
	EndTime     time.Time     `db:"end_time" json:"end_time"`
	Status      SessionStatus `db:"status" json:"status"`
	VideoURL    string        `db:"video_url" json:"video_url"`
	CreatedBy   int64         `db:"created_by" json:"created_by"`
}

type MaterialType string

const (
	MaterialPDF  MaterialType = "PDF"
	MaterialPPT  MaterialType = "PPT"
	MaterialLink MaterialType = "LINK"
)

type StudyMaterial struct {
	ID        int64        `db:"id" json:"id"`
	SessionID *int64       `db:"session_id" json:"session_id,omitempty"`
	Title     string       `db:"title" json:"title"`
	Type      MaterialType `db:"type" json:"type"`
	S3Key     string       `db:"s3_key" json:"s3_key"`
	URL       string       `db:"url" json:"url"`
	UploadedBy int64       `db:"uploaded_by" json:"uploaded_by"`
	CreatedAt time.Time    `db:"created_at" json:"created_at"`
}

type ExamType string

const (
	ExamTypeExam       ExamType = "EXAM"
	ExamTypeAssignment ExamType = "ASSIGNMENT"
)

type Exam struct {
	ID        int64     `db:"id" json:"id"`
	SessionID int64     `db:"session_id" json:"session_id"`
	Title     string    `db:"title" json:"title"`
	OpenTime  time.Time `db:"open_time" json:"open_time"`
	CloseTime time.Time `db:"close_time" json:"close_time"`
	MaxScore  int       `db:"max_score" json:"max_score"`
	Type      ExamType  `db:"type" json:"type"`
}

type SubmissionStatus string

const (
	SubmissionSubmitted SubmissionStatus = "SUBMITTED"
	SubmissionGraded    SubmissionStatus = "GRADED"
)

type Submission struct {
	ID         int64            `db:"id" json:"id"`
	ExamID     int64            `db:"exam_id" json:"exam_id"`
	StudentID  int64            `db:"student_id" json:"student_id"`
	SubmittedAt time.Time       `db:"submitted_at" json:"submitted_at"`
	FileS3Key  string           `db:"file_s3_key" json:"file_s3_key"`
	Score      *int             `db:"score" json:"score,omitempty"`
	Status     SubmissionStatus `db:"status" json:"status"`
}

type Notification struct {
	ID        int64     `db:"id" json:"id"`
	UserID    int64     `db:"user_id" json:"user_id"`
	Type      string    `db:"type" json:"type"`
	Message   string    `db:"message" json:"message"`
	IsRead    bool      `db:"is_read" json:"is_read"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type ActivityLog struct {
	ID           int64     `db:"id" json:"id"`
	UserID       int64     `db:"user_id" json:"user_id"`
	Action       string    `db:"action" json:"action"`
	ResourceType string    `db:"resource_type" json:"resource_type"`
	ResourceID   int64     `db:"resource_id" json:"resource_id"`
	Metadata     string    `db:"metadata" json:"metadata"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

