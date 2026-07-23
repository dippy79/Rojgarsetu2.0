from fastapi import FastAPI

app = FastAPI(
    title="RojgarSetu AI Engine",
    description="Recommendation and matching service for RojgarSetu",
    version="1.0.0"
)

@app.get("/")
def read_root():
    return {"service": "RojgarSetu AI Engine", "status": "active"}

@app.get("/health")
def health_check():
    return {"status": "healthy"}

@app.post("/recommend/jobs")
def recommend_jobs(payload: dict):
    # Placeholder matching logic
    return {"status": "success", "recommendations": []}

