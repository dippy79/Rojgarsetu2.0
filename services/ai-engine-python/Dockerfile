FROM python:3.11-slim
WORKDIR /app

RUN groupadd -r appgroup && useradd -r -g appgroup appuser
ENV PYTHONPATH=/app

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY . .
RUN chown -R appuser:appgroup /app
USER appuser

EXPOSE 8000
HEALTHCHECK --interval=30s --timeout=10s --start-period=15s --retries=3 \
  CMD python -c "import urllib.request; urllib.request.urlopen('http://127.0.0.1:8000/health')"

CMD ["uvicorn", "recommender.service:app", "--host", "0.0.0.0", "--port", "8000"]
