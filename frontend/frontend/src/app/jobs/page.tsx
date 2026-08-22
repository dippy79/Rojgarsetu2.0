"use client";

import React, { useState } from "react";
import { 
  Search, MapPin, Flame, Zap, Bookmark, Briefcase, 
  CheckCircle2, ArrowUpRight, Filter, Sparkles, Clock, DollarSign 
} from "lucide-react";
import { jobsAPI } from "@/lib/api";
import { useQuery } from "@tanstack/react-query";
import type { Job } from "@/types";
import { formatSalary, getMatchScoreColor, getMatchScoreBg } from "@/lib/utils";

// Mock Data - Fully expanded with realistic job data
const MOCK_JOBS: Job[] = [
  {
    id: "job-1",
    title: "Senior Go Microservices Engineer",
    company_name: "TechCorp Global",
    location: "Delhi-NCR (Remote)",
    salary_min: 1800000,
    salary_max: 2400000,
    salary_currency: "INR",
    type: "full-time",
    skills_required: ["Go", "PostgreSQL", "Docker", "Redis"],
    description: "We are looking for a Senior Go Engineer to lead our microservices architecture...",
    category: "Technology",
    experience_required: "5+ years",
    education_required: "B.Tech/B.E. in CS",
    apply_link: "https://example.com/apply",
    last_date: "2024-09-15",
    vacancies: 2,
    views: 0,
    applications_count: 0,
    is_active: true,
    is_featured: true,
    created_at: "2024-08-21T10:00:00Z",
    updated_at: "2024-08-21T10:00:00Z"
  },
  {
    id: "job-2",
    title: "Full-Stack React + Python Engineer",
    company_name: "DataScale Systems",
    location: "Ludhiana (Hybrid)",
    salary_min: 1200000,
    salary_max: 1600000,
    salary_currency: "INR",
    type: "full-time",
    skills_required: ["React", "FastAPI", "Tailwind", "Playwright"],
    description: "Join our team to build scalable web applications using React and Python...",
    category: "Technology",
    experience_required: "3+ years",
    education_required: "B.Tech/B.E. in CS",
    apply_link: "https://example.com/apply",
    last_date: "2024-09-20",
    vacancies: 3,
    views: 0,
    applications_count: 0,
    is_active: true,
    is_featured: false,
    created_at: "2024-08-21T09:00:00Z",
    updated_at: "2024-08-21T09:00:00Z"
  },
  {
    id: "job-3",
    title: "DevOps & Cloud Infrastructure Lead",
    company_name: "Nexus Solutions",
    location: "Chandigarh (On-site)",
    salary_min: 2000000,
    salary_max: 2800000,
    salary_currency: "INR",
    type: "contract",
    skills_required: ["Kubernetes", "AWS", "CI/CD", "Docker"],
    description: "Lead our cloud infrastructure team and implement DevOps best practices...",
    category: "Technology",
    experience_required: "7+ years",
    education_required: "B.Tech/B.E. in CS",
    apply_link: "https://example.com/apply",
    last_date: "2024-09-10",
    vacancies: 1,
    views: 0,
    applications_count: 0,
    is_active: true,
    is_featured: true,
    created_at: "2024-08-20T15:00:00Z",
    updated_at: "2024-08-20T15:00:00Z"
  }
];

export default function JobsBentoPage() {
  const [searchQuery, setSearchQuery] = useState("");
  const [selectedTag, setSelectedTag] = useState("All");
  const [savedJobs, setSavedJobs] = useState<string[]>([]);
  const [streakCount] = useState(4); // Gamified state

  // Fetch real jobs data (will use mock data as fallback)
  const { data: jobsData, isLoading, error } = useQuery({
    queryKey: ['jobs'],
    queryFn: () => jobsAPI.getJobs().then(res => res.data),
    // For now, we'll use mock data since the API might not be fully set up
    enabled: false, // Disable for now, use mock data
  });

  const jobs = MOCK_JOBS; // Using mock data for demonstration

  const toggleSave = (id: string) => {
    setSavedJobs(prev => 
      prev.includes(id) ? prev.filter(jId => jId !== id) : [...prev, id]
    );
  };

  const tagsList = ["All", "Go", "React", "Python", "Docker", "Remote"];

  const filteredJobs = jobs.filter(job => {
    const matchesSearch = job.title.toLowerCase().includes(searchQuery.toLowerCase()) || 
                          job.company_name?.toLowerCase().includes(searchQuery.toLowerCase());
    const matchesTag = selectedTag === "All" || job.skills_required?.includes(selectedTag);
    return matchesSearch && matchesTag;
  });

  // Calculate mock match scores
  const getMatchScore = (job: Job) => {
    // In a real app, this would be calculated based on user profile
    const scores: Record<string, number> = {
      "job-1": 94,
      "job-2": 88,
      "job-3": 79
    };
    return scores[job.id] || Math.floor(Math.random() * 30) + 70;
  };

  return (
    <div className="min-h-screen bg-slate-navy text-off-white p-4 md:p-8 font-sans">
      <div className="max-w-7xl mx-auto space-y-6">
        
        {/* Bento Top Cluster: Search Bar + Streak Meter */}
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-4">
          
          {/* Main Search Bento Box */}
          <div className="lg:col-span-8 bg-slate-900/80 border border-slate-800 rounded-bento-lg p-6 backdrop-blur-md flex flex-col justify-between shadow-xl">
            <div className="space-y-2">
              <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-cyan-950/60 border border-cyan-800/50 text-cyan-400 text-xs font-semibold">
                <Sparkles className="w-3.5 h-3.5" />
                <span>AI-Filtered Jobs Feed</span>
              </div>
              <h1 className="text-2xl md:text-3xl font-bold tracking-tight text-white">
                Explore Aggregated Opportunities
              </h1>
            </div>

            {/* Interactive Search Field */}
            <div className="mt-4 relative">
              <Search className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-slate-400" />
              <input 
                type="text"
                placeholder="Search by role, skill, or company name..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="w-full bg-slate-950 border border-slate-800 text-slate-100 rounded-bento-sm pl-12 pr-4 py-3.5 text-sm focus:outline-none focus:border-cyan-500 transition-colors"
              />
            </div>

            {/* Tag Filter Pills */}
            <div className="mt-4 flex flex-wrap gap-2 items-center">
              <Filter className="w-4 h-4 text-slate-400 mr-1" />
              {tagsList.map(tag => (
                <button
                  key={tag}
                  onClick={() => setSelectedTag(tag)}
                  className={`px-3.5 py-1.5 rounded-lg text-xs font-medium transition-all ${
                    selectedTag === tag 
                      ? "bg-cyan-500 text-slate-950 font-semibold shadow-md shadow-cyan-500/20" 
                      : "bg-slate-800/80 text-slate-300 hover:bg-slate-700 hover:text-white"
                  }`}
                >
                  {tag}
                </button>
              ))}
            </div>
          </div>

          {/* Gamified Application Streak Widget */}
          <div className="lg:col-span-4 bg-gradient-to-br from-slate-900 to-amber-950/20 border border-amber-500/20 rounded-bento-lg p-6 flex flex-col justify-between relative overflow-hidden shadow-xl">
            <div className="flex justify-between items-start">
              <div>
                <span className="text-xs font-bold text-amber-500 uppercase tracking-wider">Candidate Streak</span>
                <h3 className="text-2xl font-black text-white mt-0.5">{streakCount} Days Active 🔥</h3>
              </div>
              <div className="p-3 bg-amber-500/10 border border-amber-500/30 rounded-bento-sm">
                <Flame className="w-6 h-6 text-amber-500 animate-pulse" />
              </div>
            </div>

            <p className="text-xs text-slate-400 mt-2">
              Apply to 1 more job today to extend your streak and earn <span className="text-amber-400 font-semibold">+50 XP</span> towards profile boost!
            </p>

            {/* Streak Progress Bar */}
            <div className="mt-4 space-y-1.5">
              <div className="flex justify-between text-[11px] text-slate-400 font-medium">
                <span>Daily Goal</span>
                <span className="text-amber-400">1 / 2 Applied</span>
              </div>
              <div className="w-full h-2 bg-slate-950 rounded-full overflow-hidden border border-slate-800">
                <div className="h-full bg-gradient-to-r from-amber-500 to-emerald-400 w-1/2 rounded-full"></div>
              </div>
            </div>
          </div>
        </div>

        {/* Bento Jobs Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {filteredJobs.length > 0 ? (
            filteredJobs.map((job) => {
              const isSaved = savedJobs.includes(job.id);
              const matchScore = getMatchScore(job);
              
              return (
                <div 
                  key={job.id} 
                  className="bento-card-dark group relative overflow-hidden"
                >
                  {/* AI Match Score Badge */}
                  <div className={`absolute top-4 right-4 px-2.5 py-1 rounded-lg text-xs font-bold flex items-center gap-1.5 ${getMatchScoreBg(matchScore)} ${getMatchScoreColor(matchScore)} border`}>
                    <Zap className="w-3 h-3" />
                    {matchScore}% Match
                  </div>

                  {/* Job Header */}
                  <div className="space-y-3">
                    <div className="flex items-start justify-between pr-16">
                      <div>
                        <h3 className="text-lg font-bold text-white group-hover:text-cyan-400 transition-colors">
                          {job.title}
                        </h3>
                        <p className="text-sm text-slate-400 mt-1">{job.company_name}</p>
                      </div>
                    </div>

                    {/* Meta Info */}
                    <div className="flex flex-wrap gap-2 text-xs text-slate-400">
                      <span className="flex items-center gap-1">
                        <MapPin className="w-3.5 h-3.5" />
                        {job.location}
                      </span>
                      <span className="flex items-center gap-1">
                        <Briefcase className="w-3.5 h-3.5" />
                        {job.type}
                      </span>
                      <span className="flex items-center gap-1">
                        <DollarSign className="w-3.5 h-3.5" />
                        {formatSalary(job.salary_min, job.salary_max)}
                      </span>
                    </div>

                    {/* Skills Tags */}
                    <div className="flex flex-wrap gap-1.5">
                      {job.skills_required?.slice(0, 3).map((skill, idx) => (
                        <span key={idx} className="px-2 py-0.5 bg-slate-800 text-slate-300 rounded text-[10px]">
                          {skill}
                        </span>
                      ))}
                    </div>

                    {/* XP Reward Badge */}
                    <div className="flex items-center gap-1.5 text-xs text-amber-400">
                      <Sparkles className="w-3.5 h-3.5" />
                      <span className="font-semibold">+{Math.floor(matchScore / 2)} XP</span>
                    </div>
                  </div>

                  {/* Action Buttons */}
                  <div className="flex items-center justify-between mt-4 pt-4 border-t border-slate-800">
                    <div className="flex items-center gap-2 text-[10px] text-slate-500">
                      <Clock className="w-3 h-3" />
                      <span>Posted 2h ago</span>
                    </div>
                    <div className="flex items-center gap-2">
                      <button
                        onClick={() => toggleSave(job.id)}
                        className={`p-2 rounded-lg transition-colors ${
                          isSaved 
                            ? "text-amber-500 bg-amber-500/10" 
                            : "text-slate-400 hover:text-amber-500 hover:bg-amber-500/10"
                        }`}
                      >
                        <Bookmark className={`w-4 h-4 ${isSaved ? "fill-current" : ""}`} />
                      </button>
                      <button className="btn-primary px-3 py-1.5 text-xs flex items-center gap-1.5">
                        Apply Now
                        <ArrowUpRight className="w-3.5 h-3.5" />
                      </button>
                    </div>
                  </div>
                </div>
              );
            })
          ) : (
            <div className="col-span-full text-center py-12">
              <p className="text-slate-400">No jobs found matching your criteria.</p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}