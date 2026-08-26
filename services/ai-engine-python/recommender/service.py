import os
import logging
import json
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from typing import List, Optional
import google.generativeai as genai

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

app = FastAPI(
    title="RojgarSetu AI Engine",
    description="Recommendation and matching service for RojgarSetu",
    version="1.0.0"
)

# Database connection URL from environment
DATABASE_URL = os.getenv(
    "DATABASE_URL",
    "postgresql://postgres:postgres@postgres:5432/rojgarsetu"
)

# Configure Gemini
genai.configure(api_key=os.getenv("GEMINI_API_KEY", ""))

class RecommendationRequest(BaseModel):
    user_skills: List[str] = []
    experience_years: Optional[int] = 0
    preferred_locations: Optional[List[str]] = []
    expected_salary: Optional[int] = 0

class ResumeParseRequest(BaseModel):
    text: str

@app.post("/parse-resume")
def parse_resume(payload: ResumeParseRequest):
    if not os.getenv("GEMINI_API_KEY"):
        return {"error": "Gemini API key not configured"}

    try:
        model = genai.GenerativeModel("gemini-2.0-flash")
        prompt = f"""
        Extract professional information from the following resume text.
        Return ONLY a JSON object with this exact structure:
        {{
            "name": "Full Name",
            "email": "email@example.com",
            "skills": ["skill1", "skill2"],
            "experience": ["job1", "job2"],
            "education": ["degree1", "degree2"]
        }}

        Resume text:
        {payload.text}
        """
        response = model.generate_content(prompt)
        text = response.text.strip()
        # Handle cases where model might wrap JSON in markdown blocks
        if "```json" in text:
            text = text.split("```json")[1].split("```")[0].strip()
        elif "```" in text:
            text = text.split("```")[1].split("```")[0].strip()

        return json.loads(text)
    except Exception as e:
        logger.error(f"Resume parse failed: {e}")
        return {"error": "parse_failed", "details": str(e)}

class JobRecord(BaseModel):
    id: str
    title: str
    location: Optional[str] = ""
    skills: Optional[List[str]] = []
    description: Optional[str] = ""
    job_type: Optional[str] = ""
    source_table: Optional[str] = ""


def get_db_connection():
    """Create a new database connection."""
    try:
        import psycopg2
        conn = psycopg2.connect(DATABASE_URL)
        return conn
    except Exception as e:
        logger.error(f"Database connection failed: {e}")
        return None


def fetch_jobs_from_db() -> List[JobRecord]:
    """Fetch available jobs from the company_jobs table."""
    conn = get_db_connection()
    if not conn:
        return []

    try:
        with conn.cursor() as cur:
            cur.execute("""
                SELECT
                    cj.id::text AS id,
                    COALESCE(cj.title, '') AS title,
                    COALESCE(cj.location, '') AS location,
                    COALESCE(cj.skills, '{}') AS skills,
                    COALESCE(cj.description, '') AS description,
                    COALESCE(cj.job_type, '') AS job_type,
                    'company_jobs' AS source_table
                FROM company_jobs cj
                WHERE cj.is_active = true

                UNION ALL

                SELECT
                    pj.id::text AS id,
                    COALESCE(pj.title, '') AS title,
                    COALESCE(pj.location, '') AS location,
                    COALESCE(pj.skills, '{}') AS skills,
                    COALESCE(pj.description, '') AS description,
                    COALESCE(pj.job_type, '') AS job_type,
                    'jobs_private' AS source_table
                FROM jobs_private pj
                WHERE pj.is_active = true

                UNION ALL

                SELECT
                    gj.id::text AS id,
                    COALESCE(gj.title, '') AS title,
                    COALESCE(gj.location, '') AS location,
                    '{}'::text[] AS skills,
                    TRIM(
                        COALESCE(gj.eligibility, '') || ' ' || COALESCE(gj.department, '')
                    ) AS description,
                    '' AS job_type,
                    'jobs_government' AS source_table
                FROM jobs_government gj
                WHERE gj.is_active = true
                LIMIT 200
            """)
            rows = cur.fetchall()
            jobs = []
            for row in rows:
                job_id, title, location, skills, description, job_type, source_table = row
                # skills comes from PostgreSQL as a list (text[])
                if isinstance(skills, str):
                    # Handle case where it's returned as a string
                    skills_list = [s.strip().strip('"') for s in skills.strip('{}').split(',') if s.strip()]
                else:
                    skills_list = list(skills) if skills else []
                jobs.append(JobRecord(
                    id=job_id,
                    title=title,
                    location=location or "",
                    skills=skills_list,
                    description=description or "",
                    job_type=job_type or "",
                    source_table=source_table or ""
                ))
            return jobs
    except Exception as e:
        logger.error(f"Error fetching jobs from database: {e}")
        return []
    finally:
        conn.close()


def calculate_skill_match_score(user_skills: set, job_skills: set) -> float:
    """Calculate Jaccard similarity between user skills and job skills."""
    if not user_skills or not job_skills:
        return 0.0
    intersection = user_skills.intersection(job_skills)
    union = user_skills.union(job_skills)
    if not union:
        return 0.0
    return round(len(intersection) / len(union), 2)


def calculate_keyword_overlap(user_skills: set, job_text: str) -> float:
    """Calculate simple keyword overlap between user skills and job title/description."""
    if not user_skills or not job_text:
        return 0.0
    job_words = set(job_text.lower().split())
    intersection = user_skills.intersection(job_words)
    if not intersection:
        return 0.0
    # Score based on how many skills matched in the text
    return round(len(intersection) / len(user_skills), 2)


def recommend_jobs_from_db(user_skills: List[str], preferred_locations: List[str]) -> dict:
    """Core recommendation logic: fetch jobs from DB and score them."""
    # Validate input
    if not user_skills:
        return {
            "status": "error",
            "message": "user_skills list is empty. Provide at least one skill to get recommendations.",
            "recommendations": []
        }

    # Normalize user skills
    user_skills_set = set(skill.lower().strip() for skill in user_skills if skill.strip())
    if not user_skills_set:
        return {
            "status": "error",
            "message": "user_skills contains only empty values. Provide at least one valid skill.",
            "recommendations": []
        }

    # Normalize preferred locations
    preferred_locations_lower = [loc.lower().strip() for loc in preferred_locations if loc.strip()]

    # Fetch jobs from database
    jobs = fetch_jobs_from_db()
    if not jobs:
        return {
            "status": "error",
            "message": "No jobs found in the database. Please ensure jobs have been crawled/added.",
            "recommendations": []
        }

    logger.info(f"Fetched {len(jobs)} jobs from DB, scoring against {len(user_skills_set)} skills")

    scored_jobs = []

    for job in jobs:
        # Normalize job skills
        job_skills_set = set(skill.lower().strip() for skill in job.skills if skill.strip())

        # Calculate skill match score (Jaccard similarity)
        skill_score = calculate_skill_match_score(user_skills_set, job_skills_set)

        # Calculate keyword overlap with title and description
        job_text = f"{job.title} {job.description}"
        keyword_score = calculate_keyword_overlap(user_skills_set, job_text)

        # Combined score: 70% skill match, 30% keyword overlap
        score = (skill_score * 0.7) + (keyword_score * 0.3)

        # Location boost: +0.15 if job location matches preferred location
        if preferred_locations_lower and job.location.lower() in preferred_locations_lower:
            score += 0.15

        # Cap score at 1.0
        score = min(score, 1.0)

        if score > 0:
            scored_jobs.append({
                "job_id": job.id,
                "title": job.title,
                "location": job.location,
                "job_type": job.job_type,
                "source_table": job.source_table,
                "match_score": round(score, 2),
                "matched_skills": list(user_skills_set.intersection(job_skills_set))
            })

    # Sort descending by match score
    scored_jobs.sort(key=lambda x: x["match_score"], reverse=True)

    # Return top 20 results
    top_recommendations = scored_jobs[:20]

    return {
        "status": "success",
        "total": len(top_recommendations),
        "total_scored": len(scored_jobs),
        "recommendations": top_recommendations
    }


@app.get("/")
def read_root():
    return {"service": "RojgarSetu AI Engine", "status": "active"}


@app.get("/health")
def health_check():
    # Check database connectivity
    db_healthy = False
    conn = get_db_connection()
    if conn:
        try:
            with conn.cursor() as cur:
                cur.execute("SELECT 1")
                db_healthy = True
        except Exception:
            pass
        finally:
            conn.close()

    return {
        "status": "healthy" if db_healthy else "degraded",
        "service": "AI Engine Python",
        "database": "connected" if db_healthy else "disconnected"
    }


@app.post("/recommend/jobs")
def recommend_jobs(payload: RecommendationRequest):
    """
    Recommend jobs based on candidate skills and preferred locations.
    Queries the company_jobs table from PostgreSQL and scores jobs by
    skill overlap (Jaccard similarity + keyword matching).
    """
    try:
        result = recommend_jobs_from_db(
            user_skills=payload.user_skills,
            preferred_locations=payload.preferred_locations or []
        )
        return result
    except Exception as e:
        logger.error(f"Unexpected error in recommend_jobs: {e}")
        raise HTTPException(
            status_code=500,
            detail=f"Internal server error: {str(e)}"
        )
