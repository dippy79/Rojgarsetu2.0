"""
RojgarSetu 2.0 - AI Recommendation Engine
FastAPI-based microservice for job and course recommendations
"""

import os
from datetime import datetime
from typing import List, Optional

from dotenv import load_dotenv
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import redis
from sqlalchemy import create_engine, text
import numpy as np
from prometheus_client import Counter, Histogram, generate_latest, CONTENT_TYPE_LATEST
from fastapi.responses import Response

load_dotenv()

# Configuration
DATABASE_URL = os.getenv("DATABASE_URL", "postgresql://rojgarsetu:rojgarsetu_secret@localhost:5432/rojgarsetu")
REDIS_URL = os.getenv("REDIS_URL", "redis://localhost:6379")

# Initialize FastAPI
app = FastAPI(
    title="RojgarSetu AI Engine",
    version="2.0.0",
    description="ML-powered job and course recommendations"
)

# Prometheus metrics
request_count = Counter('http_requests_total', 'Total HTTP requests', ['method', 'endpoint', 'status'])
request_duration = Histogram('http_request_duration_seconds', 'HTTP request latency', ['method', 'endpoint'])

# Initialize Redis
redis_client = redis.from_url(REDIS_URL, decode_responses=True)

# Initialize SQLAlchemy
engine = create_engine(DATABASE_URL)


# Pydantic models
class UserProfile(BaseModel):
    user_id: str
    skills: List[str] = []
    preferences: dict = {}


class JobRecommendation(BaseModel):
    job_id: str
    title: str
    company: str
    location: str
    score: float
    reason: str


class CourseRecommendation(BaseModel):
    course_id: str
    title: str
    provider: str
    score: float
    reason: str


class RecommendationRequest(BaseModel):
    userId: str


# Helper functions
def get_user_profile(user_id: str) -> dict:
    """Fetch user profile from database"""
    with engine.connect() as conn:
        result = conn.execute(
            text("SELECT id, name, skills FROM users WHERE id = :user_id"),
            {"user_id": user_id}
        ).fetchone()
        
        if not result:
            return None
        
        return {
            "user_id": str(result.id),
            "name": result.name,
            "skills": result.skills or []
        }


def calculate_skill_match(user_skills: List[str], job_skills: List[str]) -> float:
    """Calculate skill match score"""
    if not user_skills or not job_skills:
        return 0.0
    
    user_skills_set = set(s.lower() for s in user_skills)
    job_skills_set = set(s.lower() for s in job_skills)
    
    if not job_skills_set:
        return 0.0
    
    intersection = user_skills_set.intersection(job_skills_set)
    return len(intersection) / len(job_skills_set)


def get_recommended_jobs(user_id: str, limit: int = 10) -> List[dict]:
    """Get personalized job recommendations"""
    # Check cache first
    cache_key = f"job_recs:{user_id}"
    cached = redis_client.get(cache_key)
    if cached:
        import json
        return json.loads(cached)
    
    # Get user profile
    user_profile = get_user_profile(user_id)
    if not user_profile:
        return []
    
    user_skills = user_profile.get("skills", [])
    
    # Get active jobs
    with engine.connect() as conn:
        jobs = conn.execute(
            text("""
                SELECT j.id, j.title, j.location, j.job_type, j.salary_min, j.salary_max,
                       c.name as company_name
                FROM jobs j
                LEFT JOIN companies c ON j.company_id = c.id
                WHERE j.is_active = true
                ORDER BY j.posted_at DESC
                LIMIT 100
            """)
        ).fetchall()
    
    # Score jobs based on user skills (simplified algorithm)
    scored_jobs = []
    for job in jobs:
        # Simple scoring - in production, use ML model
        score = np.random.uniform(0.5, 1.0)  # Placeholder
        
        if user_skills:
            # Boost score if skills match
            score += 0.1
        
        scored_jobs.append({
            "job_id": str(job.id),
            "title": job.title,
            "company": job.company_name,
            "location": job.location,
            "job_type": job.job_type,
            "salary_min": job.salary_min,
            "salary_max": job.salary_max,
            "score": round(score, 4),
            "reason": "Based on your skills and preferences"
        })
    
    # Sort by score and limit
    scored_jobs.sort(key=lambda x: x["score"], reverse=True)
    recommendations = scored_jobs[:limit]
    
    # Cache for 15 minutes
    import json
    redis_client.setex(cache_key, 900, json.dumps(recommendations))
    
    return recommendations


def get_recommended_courses(user_id: str, limit: int = 10) -> List[dict]:
    """Get personalized course recommendations"""
    cache_key = f"course_recs:{user_id}"
    cached = redis_client.get(cache_key)
    if cached:
        import json
        return json.loads(cached)
    
    user_profile = get_user_profile(user_id)
    if not user_profile:
        return []
    
    user_skills = user_profile.get("skills", [])
    
    with engine.connect() as conn:
        courses = conn.execute(
            text("""
                SELECT id, title, provider, skills, duration, level, is_free
                FROM courses
                ORDER BY created_at DESC
                LIMIT 50
            """)
        ).fetchall()
    
    # Score courses
    scored_courses = []
    for course in courses:
        score = np.random.uniform(0.5, 1.0)
        
        # Boost if skills align
        if user_skills and course.skills:
            skill_match = calculate_skill_match(user_skills, course.skills)
            score += skill_match * 0.3
        
        scored_courses.append({
            "course_id": str(course.id),
            "title": course.title,
            "provider": course.provider,
            "skills": course.skills,
            "duration": course.duration,
            "level": course.level,
            "is_free": course.is_free,
            "score": round(score, 4),
            "reason": "Recommended to enhance your skills"
        })
    
    scored_courses.sort(key=lambda x: x["score"], reverse=True)
    recommendations = scored_courses[:limit]
    
    import json
    redis_client.setex(cache_key, 900, json.dumps(recommendations))
    
    return recommendations


# API Endpoints
@app.get("/")
def root():
    return {
        "service": "RojgarSetu AI Engine",
        "version": "2.0.0",
        "status": "running"
    }


@app.get("/health")
def health_check():
    # Check database
    try:
        with engine.connect() as conn:
            conn.execute(text("SELECT 1"))
        db_status = "healthy"
    except Exception as e:
        db_status = f"unhealthy: {str(e)}"
    
    # Check Redis
    try:
        redis_client.ping()
        redis_status = "healthy"
    except Exception as e:
        redis_status = f"unhealthy: {str(e)}"
    
    return {
        "status": "healthy" if db_status == "healthy" and redis_status == "healthy" else "degraded",
        "database": db_status,
        "redis": redis_status
    }


@app.post("/recommend/jobs")
def recommend_jobs(request: RecommendationRequest):
    """Get personalized job recommendations"""
    try:
        recommendations = get_recommended_jobs(request.userId)
        return {
            "status": "success",
            "data": {
                "recommendations": recommendations,
                "count": len(recommendations)
            }
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/recommend/courses")
def recommend_courses(request: RecommendationRequest):
    """Get personalized course recommendations"""
    try:
        recommendations = get_recommended_courses(request.userId)
        return {
            "status": "success",
            "data": {
                "recommendations": recommendations,
                "count": len(recommendations)
            }
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.delete("/cache/{user_id}")
def clear_cache(user_id: str):
    """Clear recommendation cache for a user"""
    redis_client.delete(f"job_recs:{user_id}")
    redis_client.delete(f"course_recs:{user_id}")
    return {"status": "success", "message": f"Cache cleared for user {user_id}"}


@app.get("/metrics")
def metrics():
    """Prometheus metrics endpoint"""
    return Response(content=generate_latest(), media_type=CONTENT_TYPE_LATEST)


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)

