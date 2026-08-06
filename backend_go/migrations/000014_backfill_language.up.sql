-- 000014_backfill_language.up.sql
-- One-time backfill of the `language` column for rows that existed BEFORE the
-- 000013_add_language_columns migration added the column with DEFAULT 'en'.
--
-- Rationale: DEFAULT 'en' would silently mark all pre-existing rows (the ~800
-- already-crawled Greenhouse/Coursera/GeeksforGeeks/Naukri/UPSC rows) as
-- English without running real detection. This backfill re-runs the same
-- Unicode-script heuristic as the crawler's internal/lang.Detect (pure
-- codepoint scoring, no external deps) in pure SQL so pre-existing non-English
-- content is correctly tagged.
--
-- The SQL function `detect_language(title, description)` mirrors lang.Detect:
--   * If no strong script signal (empty/English/URL-only), returns 'en'.
--   * A non-English code is returned only when its script rune count is
--     >= 3 AND >= 20% of total non-space runes (the same conservative gate).
--   * Recognized scripts: Devanagari(hi), Tamil(ta), Telugu(te),
--     Bengali(bn), Gujarati(gu), Kannada(kn), Malayalam(ml), Gurmukhi(pa),
--     Arabic/Urdu(ur), Oriya(or).
--
-- Run AFTER 000013_add_language_columns.up.sql. Idempotent: safe to re-run.
--
-- NOTE: PostgreSQL has no UNICODE() builtin; we use ascii(ch) which returns
-- the Unicode code point of a single-character string (works for UTF-8).

BEGIN;

CREATE OR REPLACE FUNCTION detect_language(title TEXT, description TEXT)
RETURNS TEXT AS $$
DECLARE
    combined TEXT;
    total_runes INT := 0;
    hi_count INT := 0;
    ta_count INT := 0;
    te_count INT := 0;
    bn_count INT := 0;
    gu_count INT := 0;
    kn_count INT := 0;
    ml_count INT := 0;
    pa_count INT := 0;
    ur_count INT := 0;
    or_count INT := 0;
    ch CHAR;
    best TEXT := '';
    best_count INT := 0;
BEGIN
    combined := COALESCE(title, '') || ' ' || COALESCE(description, '');

    FOR ch IN SELECT unnest(string_to_array(combined, NULL)) AS c
    LOOP
        IF ch = ' ' OR ch = E'\n' OR ch = E'\t' OR ch = E'\r' THEN
            CONTINUE;
        END IF;
        total_runes := total_runes + 1;

        -- ascii(ch) returns the Unicode code point of the single char ch.
        -- Devanagari U+0900-U+097F (Hindi/Marathi/Nepali)
        IF ascii(ch) BETWEEN 0x0900 AND 0x097F THEN hi_count := hi_count + 1;
        -- Tamil U+0B80-U+0BFF
        ELSIF ascii(ch) BETWEEN 0x0B80 AND 0x0BFF THEN ta_count := ta_count + 1;
        -- Telugu U+0C00-U+0C7F
        ELSIF ascii(ch) BETWEEN 0x0C00 AND 0x0C7F THEN te_count := te_count + 1;
        -- Bengali U+0980-U+09FF
        ELSIF ascii(ch) BETWEEN 0x0980 AND 0x09FF THEN bn_count := bn_count + 1;
        -- Gujarati U+0A80-U+0AFF
        ELSIF ascii(ch) BETWEEN 0x0A80 AND 0x0AFF THEN gu_count := gu_count + 1;
        -- Kannada U+0C80-U+0CFF
        ELSIF ascii(ch) BETWEEN 0x0C80 AND 0x0CFF THEN kn_count := kn_count + 1;
        -- Malayalam U+0D00-U+0D7F
        ELSIF ascii(ch) BETWEEN 0x0D00 AND 0x0D7F THEN ml_count := ml_count + 1;
        -- Gurmukhi U+0A00-U+0A7F (Punjabi)
        ELSIF ascii(ch) BETWEEN 0x0A00 AND 0x0A7F THEN pa_count := pa_count + 1;
        -- Arabic U+0600-U+06FF (Urdu)
        ELSIF ascii(ch) BETWEEN 0x0600 AND 0x06FF THEN ur_count := ur_count + 1;
        -- Oriya U+0B00-U+0B7F
        ELSIF ascii(ch) BETWEEN 0x0B00 AND 0x0B7F THEN or_count := or_count + 1;
        END IF;
    END LOOP;

    IF total_runes = 0 THEN
        RETURN 'en';
    END IF;

    -- Conservative gate mirroring lang.Detect: count >= 3 AND count >= total/5.
    IF hi_count >= 3 AND hi_count * 5 >= total_runes AND hi_count > best_count THEN
        best := 'hi'; best_count := hi_count;
    END IF;
    IF ta_count >= 3 AND ta_count * 5 >= total_runes AND ta_count > best_count THEN
        best := 'ta'; best_count := ta_count;
    END IF;
    IF te_count >= 3 AND te_count * 5 >= total_runes AND te_count > best_count THEN
        best := 'te'; best_count := te_count;
    END IF;
    IF bn_count >= 3 AND bn_count * 5 >= total_runes AND bn_count > best_count THEN
        best := 'bn'; best_count := bn_count;
    END IF;
    IF gu_count >= 3 AND gu_count * 5 >= total_runes AND gu_count > best_count THEN
        best := 'gu'; best_count := gu_count;
    END IF;
    IF kn_count >= 3 AND kn_count * 5 >= total_runes AND kn_count > best_count THEN
        best := 'kn'; best_count := kn_count;
    END IF;
    IF ml_count >= 3 AND ml_count * 5 >= total_runes AND ml_count > best_count THEN
        best := 'ml'; best_count := ml_count;
    END IF;
    IF pa_count >= 3 AND pa_count * 5 >= total_runes AND pa_count > best_count THEN
        best := 'pa'; best_count := pa_count;
    END IF;
    IF ur_count >= 3 AND ur_count * 5 >= total_runes AND ur_count > best_count THEN
        best := 'ur'; best_count := ur_count;
    END IF;
    IF or_count >= 3 AND or_count * 5 >= total_runes AND or_count > best_count THEN
        best := 'or'; best_count := or_count;
    END IF;

    IF best <> '' THEN
        RETURN best;
    END IF;

    RETURN 'en';
END;
$$ LANGUAGE plpgsql;

-- Apply detection to existing rows only (skip rows the crawler already tagged
-- with a real value — though at first run all rows are DEFAULT 'en', so this
-- is effectively a full-table backfill).
UPDATE jobs_government
SET language = detect_language(title, COALESCE(eligibility, '') || ' ' || COALESCE(department, ''))
WHERE language IS NULL OR language = 'en';

UPDATE jobs_private
SET language = detect_language(title, COALESCE(description, ''))
WHERE language IS NULL OR language = 'en';

UPDATE company_jobs
SET language = detect_language(title, COALESCE(description, ''))
WHERE language IS NULL OR language = 'en';

UPDATE courses
SET language = detect_language(title, COALESCE(description, ''))
WHERE language IS NULL OR language = 'en';

UPDATE youtube_videos
SET language = detect_language(title, COALESCE(description, ''))
WHERE language IS NULL OR language = 'en';

COMMIT;
