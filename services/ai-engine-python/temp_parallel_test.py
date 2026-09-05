import os
import threading
import time
from unittest.mock import patch

import recommender.service as service


def fake_fetch(table_name: str, query: str):
    start = time.perf_counter()
    print(f"start:{table_name}:{start:.3f}")
    time.sleep(0.5)
    end = time.perf_counter()
    print(f"end:{table_name}:{end:.3f}")
    return [
        ("1", "Test Job", "Delhi", "{Python}", "desc", "full-time")
    ]


# Force a cache miss and isolate the concurrency check.
service.fetch_jobs_from_db_cached.cache_clear()
with patch.object(service, "get_db_connection", return_value=None):
    with patch.object(service, "_fetch_jobs_table", side_effect=fake_fetch) as mocked:
        start = time.perf_counter()
        jobs = service.fetch_jobs_from_db_cached(999999)
        elapsed = time.perf_counter() - start
        print(f"jobs={len(jobs)} elapsed={elapsed:.3f}s futures={mocked.call_count}")
        if mocked.call_count != 3:
            raise SystemExit(f"expected 3 database queries, got {mocked.call_count}")
        if elapsed >= 1.4:
            raise SystemExit(f"queries ran serially or slower than expected: {elapsed:.3f}s")
        print("PARALLEL_QUERY_CHECK=PASS")
