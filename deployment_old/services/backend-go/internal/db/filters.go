package db

type GovJob = JobsGovernment
type PrivJob = JobsPrivate
type Video = YoutubeVideo

type GovJobFilter struct {
    Department string
    Location   string
    Source     string
}

type PrivJobFilter struct {
    Company  string
    Location string
    Source   string
    JobType  string
}

type CourseFilter struct {
    Provider string
    Mode     string
    Level    string
}

type VideoFilter struct {
    Channel  string
    Category string
}
