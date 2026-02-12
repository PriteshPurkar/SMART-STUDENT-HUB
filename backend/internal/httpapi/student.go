package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/example/scalable-learning-platform/backend/internal/middleware"
)

func registerStudentRoutes(r chi.Router, services *Services) {
	r.Get("/dashboard", handleStudentDashboard(services))
}

type dashboardResponse struct {
	UpcomingSessions []interface{}  `json:"upcoming_sessions"`
	PastSessions     []interface{}  `json:"past_sessions"`
	Notifications    []interface{}  `json:"notifications"`
}

func handleStudentDashboard(services *Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := middleware.CurrentUser(r)
		if user == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// Get upcoming sessions
		upcomingSessions, err := services.Sessions.GetUpcomingSessions(10)
		if err != nil {
			http.Error(w, "failed to fetch sessions", http.StatusInternalServerError)
			return
		}

		// Get past sessions
		pastSessions, err := services.Sessions.GetPastSessions(10)
		if err != nil {
			http.Error(w, "failed to fetch sessions", http.StatusInternalServerError)
			return
		}

		// Get user notifications
		notifications, err := services.Notifications.GetNotificationsByUserID(user.UserID)
		if err != nil {
			http.Error(w, "failed to fetch notifications", http.StatusInternalServerError)
			return
		}

		// Convert to interface slices for JSON response
		var upcomingInterface []interface{}
		for _, s := range upcomingSessions {
			upcomingInterface = append(upcomingInterface, s)
		}

		var pastInterface []interface{}
		for _, s := range pastSessions {
			pastInterface = append(pastInterface, s)
		}

		var notifInterface []interface{}
		for _, n := range notifications {
			notifInterface = append(notifInterface, n)
		}

		resp := dashboardResponse{
			UpcomingSessions: upcomingInterface,
			PastSessions:     pastInterface,
			Notifications:    notifInterface,
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}
}

