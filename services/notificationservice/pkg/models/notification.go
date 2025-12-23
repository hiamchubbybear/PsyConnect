package models

// Email Notification Types
type UserCreateNotification struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Code     string `json:"code"`
	Fullname string `json:"fullname"`
}

type UserActivateNotification struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Code     string `json:"code"`
	Fullname string `json:"fullname"`
}

type PasswordResetNotification struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Code     string `json:"code"`
}

type AccountChangeNotification struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

// Therapist Notification Types
type TherapistApproveNotification struct {
	UserID   string `json:"userId"`
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
}

type TherapistRejectNotification struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
}

// Push Notification Types
type NewMessageNotification struct {
	UserID string `json:"userId"`
	From   string `json:"from"`
}

type ConsultationCreatedNotification struct {
	UserID string `json:"userId"`
}

type ConsultationUpdatedNotification struct {
	UserID string `json:"userId"`
}

type ConsultationReminderNotification struct {
	UserID string `json:"userId"`
}

type ConsultationCompletedNotification struct {
	UserID string `json:"userId"`
}

type ProfileUpdateNotification struct {
	UserID string `json:"userId"`
}

type ProfileApprovedNotification struct {
	UserID string `json:"userId"`
}

type ProfileRejectedNotification struct {
	UserID string `json:"userId"`
}

type NewReviewNotification struct {
	UserID   string `json:"userId"`
	Reviewer string `json:"reviewer"`
}

type NewClientNotification struct {
	UserID     string `json:"userId"`
	ClientName string `json:"clientName"`
}

type SystemNotification struct {
	UserID  string `json:"userId"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

// Social Notification Types
type PostUpvoteNotification struct {
	UserID    string `json:"userId"`
	VoterName string `json:"voterName"`
	PostID    string `json:"postId"`
}

type PostBookmarkNotification struct {
	UserID         string `json:"userId"`
	BookmarkerName string `json:"bookmarkerName"`
	PostID         string `json:"postId"`
}

type PostShareNotification struct {
	UserID     string `json:"userId"`
	SharerName string `json:"sharerName"`
	PostID     string `json:"postId"`
}

type PostCommentNotification struct {
	UserID        string `json:"userId"`
	CommenterName string `json:"commenterName"`
	PostID        string `json:"postId"`
	CommentID     string `json:"commentId"`
}

type UserFollowNotification struct {
	UserID       string `json:"userId"`
	FollowerName string `json:"followerName"`
	FollowerID   string `json:"followerId"`
}
