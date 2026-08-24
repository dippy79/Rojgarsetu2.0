package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/rojgarsetu/backend/internal/db"
)

type InterviewService struct {
	database *db.PostgresDB
}

func NewInterviewService(database *db.PostgresDB) *InterviewService {
	return &InterviewService{database: database}
}

type DailyRoomResponse struct {
	URL string `json:"url"`
	Name string `json:"name"`
}

func (s *InterviewService) CreateRoom(interviewID string) (string, error) {
	apiKey := os.Getenv("DAILY_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("DAILY_API_KEY not configured")
	}

	exp := time.Now().Add(24 * time.Hour).Unix()
	body, _ := json.Marshal(map[string]interface{}{
		"name": "interview-" + interviewID,
		"privacy": "private",
		"properties": map[string]interface{}{
			"exp": exp,
		},
	})

	req, _ := http.NewRequest("POST", "https://api.daily.co/v1/rooms", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("daily.co api error: %d", resp.StatusCode)
	}

	var room DailyRoomResponse
	if err := json.NewDecoder(resp.Body).Decode(&room); err != nil {
		return "", err
	}

	return room.URL, nil
}

func (s *InterviewService) ScheduleInterview(ctx context.Context, appID, candID, compID string, scheduledAt time.Time) (*db.Interview, error) {
	applicationID, _ := uuid.Parse(appID)
	candidateID, _ := uuid.Parse(candID)
	companyID, _ := uuid.Parse(compID)

	interviewID := uuid.New().String()
	roomURL, err := s.CreateRoom(interviewID)
	if err != nil {
		return nil, err
	}

	interview, err := s.database.Queries.CreateInterview(ctx, db.CreateInterviewParams{
		ApplicationID: applicationID,
		CandidateID:   candidateID,
		CompanyID:     companyID,
		ScheduledAt:   scheduledAt,
		RoomUrl:       db.NullStringPtr(roomURL),
	})
	if err != nil {
		return nil, err
	}

	// Queue emails and notifications (to be implemented)
	return &interview, nil
}

func (s *InterviewService) GetInterviewByID(ctx context.Context, id string) (*db.Interview, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	interview, err := s.database.Queries.GetInterviewByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return &interview, nil
}
