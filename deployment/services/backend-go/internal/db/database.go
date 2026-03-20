package db

import (
    "context"
    "database/sql"
    "fmt"
    "time"

    "github.com/google/uuid"
)

type PostgresDB struct {
    *Queries
    DB *sql.DB
}

func NewPostgresDB(connStr string) (*PostgresDB, error) {
    dbConn, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, fmt.Errorf("failed to open database (check DATABASE_URL env): %w", err)
    }
    dbConn.SetMaxOpenConns(25)
    dbConn.SetMaxIdleConns(5)
    dbConn.SetConnMaxLifetime(5 * time.Minute)
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := dbConn.PingContext(ctx); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }
    return &PostgresDB{Queries: New(dbConn), DB: dbConn}, nil
}

func (p *PostgresDB) Close() error { return p.DB.Close() }
func (p *PostgresDB) GetDB() *sql.DB { return p.DB }

func (p *PostgresDB) WithTx(tx *sql.Tx) *PostgresDB {
    return &PostgresDB{Queries: p.Queries.WithTx(tx), DB: p.DB}
}

func (p *PostgresDB) ExecWithRetry(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
    var result sql.Result
    var err error
    for attempt := 0; attempt < 3; attempt++ {
        result, err = p.DB.ExecContext(ctx, query, args...)
        if err == nil {
            return result, nil
        }
        if pgErr, ok := err.(interface{ Code() string }); ok && pgErr.Code() == "40P01" {
            time.Sleep(time.Duration(attempt+1) * 100 * time.Millisecond)
            continue
        }
        return nil, err
    }
    return nil, err
}

func (p *PostgresDB) QueryWithRetry(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
    var rows *sql.Rows
    var err error
    for attempt := 0; attempt < 3; attempt++ {
        rows, err = p.DB.QueryContext(ctx, query, args...)
        if err == nil {
            return rows, nil
        }
        if pgErr, ok := err.(interface{ Code() string }); ok && pgErr.Code() == "40P01" {
            if rows != nil {
                rows.Close()
            }
            time.Sleep(time.Duration(attempt+1) * 100 * time.Millisecond)
            continue
        }
        return nil, err
    }
    return nil, err
}

func (p *PostgresDB) GetGovJobs(filter GovJobFilter, page, limit int) ([]GetGovJobsRow, int, error) {
    ctx := context.Background()
    offset := (page - 1) * limit
    jobs, err := p.Queries.GetGovJobs(ctx, GetGovJobsParams{
        Column1: filter.Department,
        Column2: filter.Location,
        Column3: filter.Source,
        Limit:   int32(limit),
        Offset:  int32(offset),
    })
    if err != nil {
        return nil, 0, err
    }
    count, err := p.Queries.GetGovJobsCount(ctx, GetGovJobsCountParams{
        Column1: filter.Department,
        Column2: filter.Location,
        Column3: filter.Source,
    })
    if err != nil {
        return jobs, 0, err
    }
    return jobs, int(count), nil
}

func (p *PostgresDB) GetGovJobByID(id string) (*GetGovJobByIDRow, error) {
    parsedID, err := uuid.Parse(id)
    if err != nil {
        return nil, fmt.Errorf("invalid uuid: %w", err)
    }
    result, err := p.Queries.GetGovJobByID(context.Background(), parsedID)
    if err != nil {
        return nil, err
    }
    return &result, nil
}

func (p *PostgresDB) GetPrivJobs(filter PrivJobFilter, page, limit int) ([]GetPrivJobsRow, int, error) {
    ctx := context.Background()
    offset := (page - 1) * limit
    jobs, err := p.Queries.GetPrivJobs(ctx, GetPrivJobsParams{
        Column1: filter.Company,
        Column2: filter.Location,
        Column3: filter.Source,
        Column4: filter.JobType,
        Limit:   int32(limit),
        Offset:  int32(offset),
    })
    if err != nil {
        return nil, 0, err
    }
    count, err := p.Queries.GetPrivJobsCount(ctx, GetPrivJobsCountParams{
        Column1: filter.Company,
        Column2: filter.Location,
        Column3: filter.Source,
        Column4: filter.JobType,
    })
    if err != nil {
        return jobs, 0, err
    }
    return jobs, int(count), nil
}

func (p *PostgresDB) GetPrivJobByID(id string) (*GetPrivJobByIDRow, error) {
    parsedID, err := uuid.Parse(id)
    if err != nil {
        return nil, fmt.Errorf("invalid uuid: %w", err)
    }
    result, err := p.Queries.GetPrivJobByID(context.Background(), parsedID)
    if err != nil {
        return nil, err
    }
    return &result, nil
}

func (p *PostgresDB) GetCourses(filter CourseFilter, page, limit int) ([]GetCoursesRow, int, error) {
    ctx := context.Background()
    offset := (page - 1) * limit
    courses, err := p.Queries.GetCourses(ctx, GetCoursesParams{
        Column1: filter.Provider,
        Column2: filter.Mode,
        Column3: filter.Level,
        Limit:   int32(limit),
        Offset:  int32(offset),
    })
    if err != nil {
        return nil, 0, err
    }
    count, err := p.Queries.GetCoursesCount(ctx, GetCoursesCountParams{
        Column1: filter.Provider,
        Column2: filter.Mode,
        Column3: filter.Level,
    })
    if err != nil {
        return courses, 0, err
    }
    return courses, int(count), nil
}

func (p *PostgresDB) GetCourseByID(id string) (*GetCourseByIDRow, error) {
    parsedID, err := uuid.Parse(id)
    if err != nil {
        return nil, fmt.Errorf("invalid uuid: %w", err)
    }
    result, err := p.Queries.GetCourseByID(context.Background(), parsedID)
    if err != nil {
        return nil, err
    }
    return &result, nil
}

func (p *PostgresDB) GetVideos(filter VideoFilter, page, limit int) ([]GetVideosRow, int, error) {
    ctx := context.Background()
    offset := (page - 1) * limit
    videos, err := p.Queries.GetVideos(ctx, GetVideosParams{
        Column1: filter.Channel,
        Column2: filter.Category,
        Limit:   int32(limit),
        Offset:  int32(offset),
    })
    if err != nil {
        return nil, 0, err
    }
    count, err := p.Queries.GetVideosCount(ctx, GetVideosCountParams{
        Column1: filter.Channel,
        Column2: filter.Category,
    })
    if err != nil {
        return videos, 0, err
    }
    return videos, int(count), nil
}

func (p *PostgresDB) GetVideoByID(id string) (*GetVideoByIDRow, error) {
    parsedID, err := uuid.Parse(id)
    if err != nil {
        return nil, fmt.Errorf("invalid uuid: %w", err)
    }
    result, err := p.Queries.GetVideoByID(context.Background(), parsedID)
    if err != nil {
        return nil, err
    }
    return &result, nil
}
