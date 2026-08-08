# Task Group D - mechanical transform script
# Copies original sources/*.go files to new package dirs, rewrites package
# clause, adds shared import, and prefixes the 16 helper calls with `shared.`
# and the data types with `shared.` where they cross package boundaries.

$ErrorActionPreference = "Stop"
$root = "f:\Rojgarsetu2.0\rojgarsetu2\services\crawler-go\internal"
$src = "$root\sources"

# type -> package mapping for the files being moved
$typeMap = @{
    "GovJobSource"     = "shared"
    "PrivJobSource"    = "shared"
    "CourseSource"     = "shared"
    "YouTubeVideoSource" = "shared"
    "GovJobFetcher"    = "shared"
    "PrivJobFetcher"   = "shared"
    "CourseFetcher"    = "shared"
    "VideoFetcher"     = "shared"
    "BaseSource"       = "shared"
    "JobSource"        = "shared"
    "RSSDocument"      = "shared"
    "RSSChannel"       = "shared"
    "RSSItem"          = "shared"
}

# 16 renamed helpers
$helperRenames = @{
    "parseRSSXML"          = "shared.ParseRSSXML"
    "cleanString"          = "shared.CleanString"
    "extractURL"           = "shared.ExtractURL"
    "extractField"         = "shared.ExtractField"
    "parseDateString"      = "shared.ParseDateString"
    "isValidJob"           = "shared.IsValidJob"
    "isValidPrivJob"       = "shared.IsValidPrivJob"
    "isValidCourse"        = "shared.IsValidCourse"
    "isValidVideo"         = "shared.IsValidVideo"
    "extractMatches"       = "shared.ExtractMatches"
    "extractYouTubeVideoID" = "shared.ExtractYouTubeVideoID"
    "parseYouTubeDuration" = "shared.ParseYouTubeDuration"
    "parseViewCount"       = "shared.ParseViewCount"
    "normalizeJobType"     = "shared.NormalizeJobType"
    "normalizeCourseMode"  = "shared.NormalizeCourseMode"
    "normalizeCourseLevel" = "shared.NormalizeCourseLevel"
    # base.go exported helpers
    "SetUserAgentAndCheck" = "shared.SetUserAgentAndCheck"
    "CheckRobotsTxt"       = "shared.CheckRobotsTxt"
    "CheckStatusAndPause"  = "shared.CheckStatusAndPause"
    "NewDomainLimiter"     = "shared.NewDomainLimiter"
    "NewUserAgentRotator"  = "shared.NewUserAgentRotator"
    "SanitizeString"       = "shared.SanitizeString"
    "DoRequest"            = "shared.DoRequest"
}

# ordered word-boundary replacement to avoid partial matches
function Replace-Word([string]$text, [string]$word, [string]$repl) {
    return [regex]::Replace($text, "(?<![\w.])" + [regex]::Escape($word) + "(?![\w])", $repl)
}

function Transform-File([string]$srcFile, [string]$dstFile, [string]$pkg) {
    $content = Get-Content -Raw -Encoding UTF8 $srcFile

    # package clause
    $content = [regex]::Replace($content, "^package\s+\w+", "package $pkg")

    # Append shared import is tricky; we rewrite the import block. Simpler:
    # replace the import block by inserting the shared path.
    # Find the import opening paren
    $content = Add-SharedImport $content

    # First replace long helper names (longest first to avoid substring issues)
    $ordered = $helperRenames.GetEnumerator() | Sort-Object { $_.Key.Length } -Descending
    foreach ($h in $ordered) {
        $content = Replace-Word $content $h.Key $h.Value
    }

    # Replace type references (boundary-aware)
    $orderedTypes = $typeMap.GetEnumerator() | Sort-Object { $_.Key.Length } -Descending
    foreach ($t in $orderedTypes) {
        $content = Replace-Word $content $t.Key "$($t.Value).$($t.Key)"
    }

    # Remove the local BaseSource literal usage that now becomes shared.BaseSource (already handled by typeMap)
    Set-Content -Path $dstFile -Value $content -Encoding UTF8
    Write-Host "Wrote $dstFile"
}

function Add-SharedImport($content) {
    # We need to add "github.com/rojgarsetu/crawler/internal/shared" to the import block.
    # Replace the import block content.
    $impPattern = '(?s)import \(\s*\n(.*?)\n\s*\)'
    $sharedImp = "`\t`"github.com/rojgarsetu/crawler/internal/shared`"`n"
    $m = [regex]::Match($content, $impPattern)
    if ($m.Success) {
        $inner = $m.Groups[1].Value
        # add shared import at top of inner (after any leading blank). We'll insert before first line.
        $inner = "`t`"github.com/rojgarsetu/crawler/internal/shared`"`n" + $inner
        $newBlock = "import (`n" + $inner + "`n)"
        return $content.Substring(0, $m.Index) + $newBlock + $content.Substring($m.Index + $m.Length)
    }
    return $content
}

# ---- gov files already done manually; process priv, courses, videos ----
$jobs = @{
    "priv" = @("indeed.go","linkedin.go","google_jobs.go","company_pages.go","company_pages_test.go","greenhouse.go","lever.go","naukri.go")
    "courses" = @("nptel.go","swayam.go","nsdc.go","coursera.go","udemy.go","geeksforgeeks.go")
    "videos" = @("youtube.go")
}

foreach ($pkg in $jobs.Keys) {
    $dstDir = "$root\$pkg"
    if ($pkg -eq "priv") { $dstDir = "$root\jobs\priv" }
    elseif ($pkg -eq "gov") { $dstDir = "$root\jobs\gov" }
    New-Item -ItemType Directory -Force -Path $dstDir | Out-Null
    foreach ($f in $jobs[$pkg]) {
        $srcFile = "$src\$f"
        $dstFile = "$dstDir\$f"
        if (Test-Path $srcFile) {
            Transform-File $srcFile $dstFile $pkg
        } else {
            Write-Host "SKIP missing: $srcFile"
        }
    }
}

Write-Host "DONE"
