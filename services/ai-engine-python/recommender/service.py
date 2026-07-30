from fastapi import FastAPI
from pydantic import BaseModel
from typing import List, Optional

app = FastAPI(
    title="RojgarSetu AI Engine",
    description="Recommendation and matching service for RojgarSetu",
    version="1.0.0"
)


class RecommendationRequest(BaseModel):
    user_skills: List[str] = []
    preferred_locations: Optional[List[str]] = []
    jobs_pool: List[dict] = []


@app.get("/")
def read_root():
    return {"service": "RojgarSetu AI Engine", "status": "active"}


@app.get("/health")
def health_check():
    return {"status": "healthy", "service": "AI Engine Python"}


@app.post("/recommend/jobs")
def recommend_jobs(payload: RecommendationRequest):
    user_skills_set = set(skill.lower() for skill in payload.user_skills)
    scored_jobs = []

    for job in payload.jobs_pool:
        job_skills = set(skill.lower() for skill in job.get("skills", []))

        if not user_skills_set or not job_skills:
            score = 0.1
        else:
            # Calculate Jaccard Similarity Score
            intersection = user_skills_set.intersection(job_skills)
            union = user_skills_set.union(job_skills)
            score = round(len(intersection) / len(union), 2) if union else 0.0

        # Location boost factor
        if job.get("location", "").lower() in [loc.lower() for loc in payload.preferred_locations]:
            score += 0.2

        scored_jobs.append({
            "job_id": job.get("id"),
            "title": job.get("title"),
            "match_score": min(score, 1.0)
        })

    # Sort descending by match score
    scored_jobs.sort(key=lambda x: x["match_score"], reverse=True)

    return {
        "status": "success",
        "total": len(scored_jobs),
        "recommendations": scored_jobs
    }
