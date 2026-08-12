package db

// Type aliases — maps database.go names to SQLC generated model names
type GovJob = JobsGovernment
type PrivJob = JobsPrivate
type Video = YoutubeVideo

// GovJobFilter — fields map to GetGovJobsParams Column1/Column2/Column3
type GovJobFilter struct {
	Department string
	Location   string
	Source     string
}

// PrivJobFilter — fields map to GetPrivJobsParams Column1/Column2/Column3/Column4
type PrivJobFilter struct {
	Company  string
	Location string
	Source   string
	JobType  string
}

// CourseFilter — fields map to GetCoursesParams Column1/Column2/Column3
type CourseFilter struct {
	Provider string
	Mode     string
	Level    string
}

// VideoFilter — fields map to GetVideosParams Column1/Column2
type VideoFilter struct {
	Channel  string
	Category string
}

// JobTradeFilter — fields map to GetJobTradesParams Column1/Column2
type JobTradeFilter struct {
	CategoryID  string
	DemandLevel string
}

// UserEnrollmentFilter — fields map to GetUserEnrollmentsParams Column1/Column2
type UserEnrollmentFilter struct {
	Status string
}

// NotificationLogFilter — fields map to GetUserNotificationLogsParams Column1/Column2
type NotificationLogFilter struct {
	NotificationType string
}
