"use client";

import React, { useState } from "react";
import { 
  Search, BookOpen, Clock, Award, Users, Star, 
  ArrowRight, Filter, Sparkles, CheckCircle2, DollarSign 
} from "lucide-react";
import { coursesAPI } from "@/lib/api";
import { useQuery } from "@tanstack/react-query";
import type { Course } from "@/types";
import { formatFees } from "@/lib/utils";

// Mock Data - Fully expanded with realistic course data
const MOCK_COURSES: Course[] = [
  {
    id: "course-1",
    name: "Full-Stack Web Development Bootcamp",
    provider_name: "TechAcademy Pro",
    category: "Web Development",
    subcategory: "Full-Stack",
    duration: "12 weeks",
    duration_weeks: 12,
    mode: "online",
    fees_structure: "₹15,000",
    fees_amount: 15000,
    fees_currency: "INR",
    eligibility: "Basic programming knowledge",
    certification: "Industry-recognized certificate",
    apply_link: "https://example.com/enroll",
    start_date: "2024-09-01",
    batch_size: 30,
    enrolled_count: 1250,
    rating: 4.8,
    reviews_count: 342,
    is_active: true,
    is_featured: true,
    views: 0,
    created_at: "2024-08-21T10:00:00Z",
    updated_at: "2024-08-21T10:00:00Z",
    description: "Master modern web development with React, Node.js, and databases through hands-on projects and real-world applications."
  },
  {
    id: "course-2",
    name: "Data Science & Machine Learning Specialization",
    provider_name: "DataMasters Institute",
    category: "Data Science",
    subcategory: "Machine Learning",
    duration: "16 weeks",
    duration_weeks: 16,
    mode: "hybrid",
    fees_structure: "₹25,000",
    fees_amount: 25000,
    fees_currency: "INR",
    eligibility: "Python programming, basic math",
    certification: "Professional Data Scientist Certificate",
    apply_link: "https://example.com/enroll",
    start_date: "2024-09-15",
    batch_size: 25,
    enrolled_count: 890,
    rating: 4.9,
    reviews_count: 567,
    is_active: true,
    is_featured: true,
    views: 0,
    created_at: "2024-08-20T14:00:00Z",
    updated_at: "2024-08-20T14:00:00Z",
    description: "Comprehensive data science program covering ML algorithms, deep learning, and practical data analysis with industry projects."
  },
  {
    id: "course-3",
    name: "Cloud Computing & DevOps Engineering",
    provider_name: "CloudSkills Academy",
    category: "Cloud Computing",
    subcategory: "DevOps",
    duration: "10 weeks",
    duration_weeks: 10,
    mode: "online",
    fees_structure: "₹18,000",
    fees_amount: 18000,
    fees_currency: "INR",
    eligibility: "Basic Linux knowledge",
    certification: "AWS & DevOps Certified",
    apply_link: "https://example.com/enroll",
    start_date: "2024-09-05",
    batch_size: 40,
    enrolled_count: 2100,
    rating: 4.7,
    reviews_count: 892,
    is_active: true,
    is_featured: false,
    views: 0,
    created_at: "2024-08-19T11:00:00Z",
    updated_at: "2024-08-19T11:00:00Z",
    description: "Master cloud platforms (AWS, Azure, GCP), containerization, CI/CD pipelines, and infrastructure as code."
  }
];

export default function CoursesBentoPage() {
  const [searchQuery, setSearchQuery] = useState("");
  const [selectedCategory, setSelectedCategory] = useState("All");
  const [selectedMode, setSelectedMode] = useState("All");

  // Fetch real courses data (will use mock data as fallback)
  const { data: coursesData, isLoading, error } = useQuery({
    queryKey: ['courses'],
    queryFn: () => coursesAPI.getCourses().then(res => res.data),
    enabled: false, // Disable for now, use mock data
  });

  const courses = MOCK_COURSES; // Using mock data for demonstration

  const categories = ["All", "Web Development", "Data Science", "Cloud Computing", "Mobile Development"];
  const modes = ["All", "online", "offline", "hybrid"];

  const filteredCourses = courses.filter(course => {
    const matchesSearch = course.name.toLowerCase().includes(searchQuery.toLowerCase()) || 
                          course.provider_name?.toLowerCase().includes(searchQuery.toLowerCase());
    const matchesCategory = selectedCategory === "All" || course.category === selectedCategory;
    const matchesMode = selectedMode === "All" || course.mode === selectedMode;
    return matchesSearch && matchesCategory && matchesMode;
  });

  return (
    <div className="min-h-screen bg-slate-navy text-off-white p-4 md:p-8 font-sans">
      <div className="max-w-7xl mx-auto space-y-6">
        
        {/* Header Section */}
        <div className="bg-slate-900/80 border border-slate-800 rounded-bento-lg p-6 backdrop-blur-md">
          <div className="space-y-2">
            <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-emerald-950/60 border border-emerald-800/50 text-emerald-400 text-xs font-semibold">
              <Sparkles className="w-3.5 h-3.5" />
              <span>Skill & Course Aggregator</span>
            </div>
            <h1 className="text-2xl md:text-3xl font-bold tracking-tight text-white">
              Upskill with Expert-Led Courses
            </h1>
            <p className="text-slate-400 text-sm">
              Discover courses from top providers with industry-recognized certifications
            </p>
          </div>

          {/* Search Field */}
          <div className="mt-4 relative">
            <Search className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-slate-400" />
            <input 
              type="text"
              placeholder="Search courses by name, provider, or skill..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-full bg-slate-950 border border-slate-800 text-slate-100 rounded-bento-sm pl-12 pr-4 py-3.5 text-sm focus:outline-none focus:border-emerald-500 transition-colors"
            />
          </div>

          {/* Filter Pills */}
          <div className="mt-4 flex flex-wrap gap-3">
            <div className="flex items-center gap-2">
              <Filter className="w-4 h-4 text-slate-400" />
              <span className="text-xs text-slate-400">Category:</span>
            </div>
            {categories.map(category => (
              <button
                key={category}
                onClick={() => setSelectedCategory(category)}
                className={`px-3 py-1.5 rounded-lg text-xs font-medium transition-all ${
                  selectedCategory === category 
                    ? "bg-emerald-500 text-slate-950 font-semibold shadow-md shadow-emerald-500/20" 
                    : "bg-slate-800/80 text-slate-300 hover:bg-slate-700 hover:text-white"
                }`}
              >
                {category}
              </button>
            ))}
          </div>

          <div className="mt-3 flex flex-wrap gap-3">
            <div className="flex items-center gap-2">
              <span className="text-xs text-slate-400">Mode:</span>
            </div>
            {modes.map(mode => (
              <button
                key={mode}
                onClick={() => setSelectedMode(mode)}
                className={`px-3 py-1.5 rounded-lg text-xs font-medium transition-all capitalize ${
                  selectedMode === mode 
                    ? "bg-cyan-500 text-slate-950 font-semibold shadow-md shadow-cyan-500/20" 
                    : "bg-slate-800/80 text-slate-300 hover:bg-slate-700 hover:text-white"
                }`}
              >
                {mode}
              </button>
            ))}
          </div>
        </div>

        {/* Courses Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {filteredCourses.length > 0 ? (
            filteredCourses.map((course) => (
              <div 
                key={course.id} 
                className="bento-card-dark group relative overflow-hidden"
              >
                {/* Featured Badge */}
                {course.is_featured && (
                  <div className="absolute top-4 left-4 px-2.5 py-1 rounded-lg text-xs font-bold bg-amber-500/10 text-amber-500 border border-amber-500/20 flex items-center gap-1.5">
                    <Star className="w-3 h-3 fill-current" />
                    Featured
                  </div>
                )}

                {/* Course Header */}
                <div className="space-y-3">
                  <div>
                    <h3 className="text-lg font-bold text-white group-hover:text-emerald-400 transition-colors">
                      {course.name}
                    </h3>
                    <p className="text-sm text-slate-400 mt-1">{course.provider_name}</p>
                  </div>

                  {/* Rating & Students */}
                  <div className="flex items-center gap-3 text-xs">
                    <div className="flex items-center gap-1 text-amber-400">
                      <Star className="w-3.5 h-3.5 fill-current" />
                      <span className="font-semibold">{course.rating}</span>
                    </div>
                    <div className="flex items-center gap-1 text-slate-400">
                      <Users className="w-3.5 h-3.5" />
                      <span>{course.enrolled_count.toLocaleString()} students</span>
                    </div>
                  </div>

                  {/* Meta Info */}
                  <div className="flex flex-wrap gap-2 text-xs text-slate-400">
                    <span className="flex items-center gap-1">
                      <Clock className="w-3.5 h-3.5" />
                      {course.duration}
                    </span>
                    <span className="flex items-center gap-1 capitalize">
                      <BookOpen className="w-3.5 h-3.5" />
                      {course.mode}
                    </span>
                  </div>

                  {/* Description Preview */}
                  <p className="text-xs text-slate-500 line-clamp-2">
                    {course.description}
                  </p>

                  {/* Certification Badge */}
                  {course.certification && (
                    <div className="flex items-center gap-1.5 text-xs text-emerald-400">
                      <Award className="w-3.5 h-3.5" />
                      <span className="font-medium">Certificate Included</span>
                    </div>
                  )}
                </div>

                {/* Price & Action */}
                <div className="flex items-center justify-between mt-4 pt-4 border-t border-slate-800">
                  <div className="flex items-center gap-1.5">
                    <DollarSign className="w-4 h-4 text-emerald-400" />
                    <span className="text-lg font-bold text-emerald-400">
                      {formatFees(course.fees_structure)}
                    </span>
                  </div>
                  <button className="btn-primary px-3 py-1.5 text-xs flex items-center gap-1.5">
                    Enroll Now
                    <ArrowRight className="w-3.5 h-3.5" />
                  </button>
                </div>
              </div>
            ))
          ) : (
            <div className="col-span-full text-center py-12">
              <p className="text-slate-400">No courses found matching your criteria.</p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}