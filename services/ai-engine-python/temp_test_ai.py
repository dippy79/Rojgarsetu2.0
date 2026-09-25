import os
import unittest
from fastapi.testclient import TestClient

# Set smoke test mode for testing imports
os.environ["SMOKE_TEST"] = "true"
os.environ["AI_ENGINE_API_KEY"] = "mock_ai_key"
os.environ["DATABASE_URL"] = "postgres://postgres:postgres@127.0.0.1:5435/rojgarsetu2?sslmode=disable"

from recommender.service import app

class TestAIEngine(unittest.TestCase):
    def setUp(self):
        self.client = TestClient(app)
        self.headers = {"X-AI-Secret-Key": "mock_ai_key"}

    def test_health_check(self):
        response = self.client.get("/health")
        self.assertEqual(response.status_code, 200)
        data = response.json()
        self.assertIn("status", data)

    def test_parse_resume_fallback(self):
        sample_resume = "Rahul Sharma is a Software Engineer with 5 years experience in Python, Go, React, and SQL."
        response = self.client.post("/parse-resume", json={"text": sample_resume}, headers=self.headers)
        self.assertEqual(response.status_code, 200)
        data = response.json()
        self.assertIn("skills", data)
        self.assertIn("Python", data["skills"])
        self.assertIn("Go", data["skills"])
        print("Parse Resume Test Output:", data)

    def test_recommend_jobs(self):
        payload = {
            "user_skills": ["Go", "React", "Python"],
            "preferred_locations": ["Bangalore", "New Delhi"]
        }
        response = self.client.post("/recommend/jobs", json=payload, headers=self.headers)
        self.assertEqual(response.status_code, 200)
        data = response.json()
        self.assertEqual(data["status"], "success")
        self.assertIn("recommendations", data)
        print("Recommend Jobs Test Output:", len(data["recommendations"]), "recommendations found")

if __name__ == "__main__":
    unittest.main()
