// Global TypeScript type definitions

export interface User {
  id: string;
  email: string;
  role: 'admin' | 'candidate' | 'company';
  is_active: boolean;
  email_verified: boolean;
  created_at: string;
}

export interface Job {
  id: string;
  company_id?: string;
  title: string;
  slug: string;
  description: string;
  category: string;
  type: 'full-time' | 'part-time' | 'contract' | 'internship' | 'remote';
  location: string;
  salary_min?: number;
  salary_max?: number;
  salary_currency: string;
  experience_required?: string;
  education_required?: string;
  skills_required: string[];
  eligibility_criteria?: string;
  fees_structure?: string;
  benefits?: string;
  apply_link: string;
  application_process?: string;
  last_date?: string;
  vacancies: number;
  views: number;
  applications_count: number;
  is_active: boolean;
  is_featured: boolean;
  source?: string;
  source_url?: string;
  created_at: string;
  updated_at: string;
  is_expired?: boolean;
  company_name?: string;
  company_slug?: string;
  company_logo?: string;
  company_description?: string;
  company_website?: string;
  company_location?: string;
  company_verified?: boolean;
}

export interface Course {
  id: string;
  provider_id?: string;
  name: string;
  slug: string;
  description: string;
  category: string;
  subcategory?: string;
  duration: string;
  duration_weeks?: number;
  mode?: 'online' | 'offline' | 'hybrid';
  fees_structure?: string;
  fees_amount?: number;
  fees_currency: string;
  eligibility?: string;
  syllabus?: string;
  certification?: string;
  apply_link: string;
  start_date?: string;
  batch_size?: number;
  enrolled_count: number;
  rating: number;
  reviews_count: number;
  is_active: boolean;
  is_featured: boolean;
  views: number;
  created_at: string;
  updated_at: string;
  provider_name?: string;
  provider_slug?: string;
  provider_logo?: string;
  provider_description?: string;
  provider_website?: string;
  provider_location?: string;
  provider_verified?: boolean;
}

export interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
  message?: string;
  pagination?: Pagination;
}

export interface Pagination {
  page: number;
  limit: number;
  totalCount: number;
  totalPages: number;
  hasNextPage: boolean;
  hasPrevPage: boolean;
}

export interface GamificationState {
  streak: number;
  level: number;
  xp: number;
  xpToNextLevel: number;
  badges: Badge[];
  achievements: Achievement[];
}

export interface Badge {
  id: string;
  name: string;
  description: string;
  icon: string;
  earnedAt?: string;
}

export interface Achievement {
  id: string;
  title: string;
  description: string;
  progress: number;
  target: number;
  xpReward: number;
  completed: boolean;
}

export interface UserProfile {
  id: string;
  email: string;
  role: string;
  first_name?: string;
  last_name?: string;
  phone?: string;
  location?: string;
  education?: string;
  experience?: string;
  skills?: string[];
  bio?: string;
  linkedin_url?: string;
  portfolio_url?: string;
  resume_url?: string;
  expected_salary_min?: number;
  expected_salary_max?: number;
  company_name?: string;
  company_slug?: string;
  description?: string;
  website?: string;
  logo_url?: string;
  industry?: string;
  company_size?: string;
  verified?: boolean;
}