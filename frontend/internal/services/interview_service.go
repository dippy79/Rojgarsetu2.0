package services

import (
    "bytes"
    "encoding/json"
    "net/http"
    "os"
    "time"
)

type InterviewService struct{}

func (s *InterviewService) CreateRoom(interviewID string) (string, error) {
    apiKey := os.Getenv("DAILY_API_KEY")
    url := "https://api.daily.co/v1/rooms"

    body, _ := json.Marshal(map[string]interface{}{
        "name":    "interview-" + interviewID,
        "privacy": "private",
        "properties": map[string]interface{}{
            "exp": time.Now().Add(1 * time.Hour).Unix(),
        },
    })

    req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
    req.Header.Set("Authorization", "Bearer "+apiKey)
    req.Header.Set("Content-Type", "application/json")

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    var result map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&result)
    if roomURL, ok := result["url"].(string); ok {
        return roomURL, nil
    }
    return "", nil
}
