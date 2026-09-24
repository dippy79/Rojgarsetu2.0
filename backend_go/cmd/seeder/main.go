package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening DB: %v", err)
	}
	defer db.Close()

	// 1. Seed Admin User
	adminEmail := "admin@rojgarsetu.in"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("Admin@12345678"), bcrypt.DefaultCost)

	var adminID string
	err = db.QueryRow(`
		INSERT INTO users (email, password_hash, role, name, is_active)
		VALUES ($1, $2, 'admin', 'System Administrator', true)
		ON CONFLICT (email) DO UPDATE SET is_active = true RETURNING id::text`,
		adminEmail, string(hashedPassword)).Scan(&adminID)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Admin seed note: %v", err)
	}

	// 2. Seed Sample Candidate User & Profile
	candEmail := "candidate@rojgarsetu.in"
	candHash, _ := bcrypt.GenerateFromPassword([]byte("Candidate@123456"), bcrypt.DefaultCost)
	var candID string
	err = db.QueryRow(`
		INSERT INTO users (email, password_hash, role, name, is_active)
		VALUES ($1, $2, 'candidate', 'Rahul Sharma', true)
		ON CONFLICT (email) DO UPDATE SET is_active = true RETURNING id::text`,
		candEmail, string(candHash)).Scan(&candID)

	if candID != "" {
		_, _ = db.Exec(`
			INSERT INTO candidates (user_id, phone, location, bio, skills, experience_years, resume_parsed)
			VALUES ($1, '+919876543210', 'New Delhi', 'Full-stack software engineer', ARRAY['Go', 'React', 'Python', 'SQL'], 3, '{"skills": ["Go", "React"], "experience": "3 years"}'::jsonb)
			ON CONFLICT DO NOTHING`, candID)
	}

	// 3. Seed Sample Employer User & Company
	empEmail := "employer@techcorp.in"
	empHash, _ := bcrypt.GenerateFromPassword([]byte("Employer@123456"), bcrypt.DefaultCost)
	var empID string
	err = db.QueryRow(`
		INSERT INTO users (email, password_hash, role, name, is_active)
		VALUES ($1, $2, 'company', 'TechCorp Hiring', true)
		ON CONFLICT (email) DO UPDATE SET is_active = true RETURNING id::text`,
		empEmail, string(empHash)).Scan(&empID)

	if empID != "" {
		_, errCompany := db.Exec(`
			INSERT INTO companies (user_id, name, industry, location, website)
			VALUES ($1, 'TechCorp India', 'Information Technology', 'Bangalore', 'https://techcorp.example.com')
			ON CONFLICT DO NOTHING`, empID)
		if errCompany != nil {
			log.Printf("Company seed error: %v", errCompany)
		}
	}

	// 4. Seed Government Jobs
	_, _ = db.Exec(`
		INSERT INTO jobs_government (title, department, location, eligibility, apply_url, source, is_active, job_hash)
		VALUES
			('Assistant Section Officer', 'Ministry of External Affairs', 'New Delhi', 'Bachelor Degree in Any Stream', 'https://ssc.gov.in', 'SSC', true, md5('ASO MEA')),
			('Railway Junior Engineer', 'Indian Railways', 'Pan India', 'Diploma / B.Tech in Civil/Mechanical', 'https://rrbcdg.gov.in', 'RRB', true, md5('RRB JE'))
		ON CONFLICT (job_hash) DO NOTHING`)

	// 5. Seed Private Jobs
	_, _ = db.Exec(`
		INSERT INTO jobs_private (title, company, location, skills, description, job_type, apply_url, source, is_active, job_hash)
		VALUES
			('Senior Go Developer', 'TechCorp India', 'Bangalore', ARRAY['Go', 'Docker', 'PostgreSQL'], 'Build scalable backend microservices.', 'Full-time', 'https://techcorp.example.com/jobs/1', 'Direct', true, md5('Senior Go Dev')),
			('Frontend React Engineer', 'DesignTech Solutions', 'Remote', ARRAY['React', 'TypeScript', 'Tailwind'], 'Develop modern web interfaces.', 'Full-time', 'https://designtech.example.com/jobs/2', 'Direct', true, md5('Frontend React Eng'))
		ON CONFLICT (job_hash) DO NOTHING`)

	// 6. Seed Courses
	_, _ = db.Exec(`
		INSERT INTO courses (title, provider, category, url, is_active)
		VALUES
			('Complete Go Programming Masterclass', 'NPTEL', 'Software Development', 'https://nptel.ac.in/courses/go', true),
			('Web Development with React and Node', 'SWAYAM', 'Web Development', 'https://swayam.gov.in/courses/webdev', true)
		ON CONFLICT DO NOTHING`)

	// 7. Seed YouTube Educational Videos
	_, _ = db.Exec(`
		INSERT INTO youtube_videos (title, channel, video_id, url, thumbnail, category, is_active)
		VALUES
			('Go Backend Engineering Course', 'Tech Edu', 'v_go_123', 'https://youtube.com/watch?v=v_go_123', 'https://img.youtube.com/vi/v_go_123/0.jpg', 'Engineering', true),
			('Full Stack Web Development Overview', 'Code Academy', 'v_web_456', 'https://youtube.com/watch?v=v_web_456', 'https://img.youtube.com/vi/v_web_456/0.jpg', 'Software', true)
		ON CONFLICT DO NOTHING`)

	log.Println("Database seeding completed successfully.")
}
