import os
import sqlite3
import requests
from fastapi import FastAPI
from pydantic import BaseModel
from typing import List, Optional

app = FastAPI()

class MatchRequest(BaseModel):
    skills: List[str]
    experience_years: int
    preferred_location: Optional[str] = ""
    expected_salary: Optional[float] = 0.0

class ResumeParseRequest(BaseModel):
    text: str

@app.get("/health")
def health():
    return {"status": "healthy"}

@app.post("/recommend/jobs")
def recommend_jobs(req: MatchRequest):
    conn = sqlite3.connect(os.Getenv("DATABASE_URL", "rojgarsetu.db"))
    cursor = conn.cursor()
    cursor.execute("SELECT id, title, skills, location, salary_max FROM jobs")
    jobs = cursor.fetchall()

    recommendations = []
    user_skills = set([s.lower() for s in req.skills])

    for j in jobs:
        j_id, title, j_skills_raw, location, salary = j
        j_skills = set([s.strip().lower() for s in (j_skills_raw or "").split(",") if s.strip()])
        
        overlap = len(user_skills.intersection(j_skills))
        total = len(j_skills) if j_skills else 1
        skill_score = (overlap / total) * 60

        loc_score = 20 if req.preferred_location and location and req.preferred_location.lower() in location.lower() else 0
        salary_score = 20 if salary and req.expected_salary and salary >= req.expected_salary else 0

        final_score = round(skill_score + loc_score + salary_score, 2)

        recommendations.append({
            "job_id": j_id,
            "title": title,
            "score": final_score,
            "matching_skills": list(user_skills.intersection(j_skills))
        })

    recommendations.sort(key=lambda x: x["score"], reverse=True)
    return recommendations[:10]

@app.post("/parse-resume")
def parse_resume(req: ResumeParseRequest):
    api_key = os.Getenv("GEMINI_API_KEY")
    url = f"https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key={api_key}"
    
    prompt = f"Extract JSON only: {{skills:[], experience:[], education:[], name, email}} from: {req.text}"
    payload = {"contents": [{"parts": [{"text": prompt}]}]}

    res = requests.post(url, json=payload)
    if res.status_code == 200:
        try:
            content = res.json()['candidates'][0]['content']['parts'][0]['text']
            return {"data": content}
        except Exception:
            return {"error": "parse_failed"}
    return {"error": "parse_failed"}
