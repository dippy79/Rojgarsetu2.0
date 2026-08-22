"use client";

import React, { useState } from "react";
import { 
  User, Flame, Award, Target, TrendingUp, CheckCircle2, 
  Zap, BookOpen, Briefcase, Settings, Bell, LogOut,
  Calendar, MapPin, Mail, Phone, Linkedin, Globe
} from "lucide-react";

// Mock user data
const MOCK_USER = {
  name: "Amit Kumar",
  email: "amit.kumar@example.com",
  role: "candidate",
  location: "Delhi-NCR, India",
  phone: "+91 98765 43210",
  linkedin: "linkedin.com/in/amitkumar",
  portfolio: "amitkumar.dev",
  level: 3,
  xp: 2450,
  xpToNextLevel: 3000,
  streak: 7,
  profileCompletion: 75,
  skills: ["React", "TypeScript", "Node.js", "Python", "PostgreSQL"],
  applications: 12,
  savedJobs: 8,
  enrolledCourses: 3
};

// Mock achievements
const MOCK_ACHIEVEMENTS = [
  {
    id: "first-application",
    title: "First Step",
    description: "Applied to your first job",
    icon: <Briefcase className="w-5 h-5" />,
    progress: 1,
    target: 1,
    xpReward: 50,
    completed: true,
    earnedAt: "2024-08-15"
  },
  {
    id: "streak-warrior",
    title: "Streak Warrior",
    description: "Maintain a 7-day application streak",
    icon: <Flame className="w-5 h-5" />,
    progress: 7,
    target: 7,
    xpReward: 200,
    completed: true,
    earnedAt: "2024-08-21"
  },
  {
    id: "skill-collector",
    title: "Skill Collector",
    description: "Add 5 skills to your profile",
    icon: <Zap className="w-5 h-5" />,
    progress: 5,
    target: 5,
    xpReward: 100,
    completed: true,
    earnedAt: "2024-08-18"
  },
  {
    id: "course-enthusiast",
    title: "Course Enthusiast",
    description: "Enroll in 5 courses",
    icon: <BookOpen className="w-5 h-5" />,
    progress: 3,
    target: 5,
    xpReward: 150,
    completed: false
  },
  {
    id: "network-builder",
    title: "Network Builder",
    description: "Save 10 jobs",
    icon: <Target className="w-5 h-5" />,
    progress: 8,
    target: 10,
    xpReward: 75,
    completed: false
  }
];

export default function ProfilePage() {
  const [user] = useState(MOCK_USER);
  const [achievements] = useState(MOCK_ACHIEVEMENTS);
  const [activeTab, setActiveTab] = useState("overview");

  const completedAchievements = achievements.filter(a => a.completed);
  const inProgressAchievements = achievements.filter(a => !a.completed);

  return (
    <div className="min-h-screen bg-slate-navy text-off-white p-4 md:p-8 font-sans">
      <div className="max-w-7xl mx-auto space-y-6">
        
        {/* Profile Header */}
        <div className="bg-slate-900/80 border border-slate-800 rounded-bento-lg p-6 backdrop-blur-md">
          <div className="flex flex-col md:flex-row items-start md:items-center gap-6">
            {/* Avatar & Level */}
            <div className="relative">
              <div className="w-24 h-24 rounded-full bg-gradient-to-br from-cyan-500 to-emerald-500 flex items-center justify-center text-3xl font-bold text-white">
                {user.name.split(' ').map(n => n[0]).join('')}
              </div>
              <div className="absolute -bottom-2 -right-2 bg-amber-500 text-slate-950 text-xs font-bold px-2 py-1 rounded-full">
                Lvl {user.level}
              </div>
            </div>

            {/* User Info */}
            <div className="flex-1">
              <h1 className="text-2xl font-bold text-white">{user.name}</h1>
              <p className="text-slate-400 text-sm mt-1">{user.email}</p>
              <div className="flex flex-wrap gap-4 mt-3 text-xs text-slate-400">
                <span className="flex items-center gap-1">
                  <MapPin className="w-3.5 h-3.5" />
                  {user.location}
                </span>
                <span className="flex items-center gap-1">
                  <Briefcase className="w-3.5 h-3.5" />
                  {user.applications} Applications
                </span>
                <span className="flex items-center gap-1">
                  <Flame className="w-3.5 h-3.5 text-amber-500" />
                  {user.streak} Day Streak
                </span>
              </div>
            </div>

            {/* XP Progress */}
            <div className="bg-slate-800/50 rounded-xl p-4 min-w-[200px]">
              <div className="flex justify-between items-center mb-2">
                <span className="text-xs text-slate-400">Level {user.level}</span>
                <span className="text-xs font-semibold text-cyan-400">{user.xp} / {user.xpToNextLevel} XP</span>
              </div>
              <div className="w-full h-2 bg-slate-700 rounded-full overflow-hidden">
                <div 
                  className="h-full bg-gradient-to-r from-cyan-500 to-emerald-400 transition-all duration-500"
                  style={{ width: `${(user.xp / user.xpToNextLevel) * 100}%` }}
                />
              </div>
              <p className="text-[10px] text-slate-500 mt-1">
                {user.xpToNextLevel - user.xp} XP to next level
              </p>
            </div>
          </div>
        </div>

        {/* Main Content Grid */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          
          {/* Left Column - Profile Stats & Completion */}
          <div className="space-y-6">
            {/* Profile Completion Ring */}
            <div className="bento-card-dark">
              <h3 className="text-lg font-semibold text-white mb-4 flex items-center gap-2">
                <Target className="w-5 h-5 text-cyan-400" />
                Profile Completion
              </h3>
              
              <div className="flex items-center justify-center mb-4">
                <div className="relative w-32 h-32">
                  <svg className="w-full h-full transform -rotate-90">
                    <circle
                      cx="64"
                      cy="64"
                      r="56"
                      stroke="currentColor"
                      strokeWidth="8"
                      fill="transparent"
                      className="text-slate-800"
                    />
                    <circle
                      cx="64"
                      cy="64"
                      r="56"
                      stroke="currentColor"
                      strokeWidth="8"
                      fill="transparent"
                      strokeDasharray={`${2 * Math.PI * 56}`}
                      strokeDashoffset={`${2 * Math.PI * 56 * (1 - user.profileCompletion / 100)}`}
                      className="text-cyan-500 transition-all duration-500"
                    />
                  </svg>
                  <div className="absolute inset-0 flex items-center justify-center">
                    <span className="text-2xl font-bold text-white">{user.profileCompletion}%</span>
                  </div>
                </div>
              </div>

              <div className="space-y-2">
                <div className="flex items-center gap-2 text-xs">
                  <CheckCircle2 className="w-4 h-4 text-emerald-400" />
                  <span className="text-slate-300">Basic information</span>
                </div>
                <div className="flex items-center gap-2 text-xs">
                  <CheckCircle2 className="w-4 h-4 text-emerald-400" />
                  <span className="text-slate-300">Skills added</span>
                </div>
                <div className="flex items-center gap-2 text-xs">
                  <div className="w-4 h-4 rounded-full border-2 border-slate-600" />
                  <span className="text-slate-500">Resume uploaded</span>
                </div>
                <div className="flex items-center gap-2 text-xs">
                  <div className="w-4 h-4 rounded-full border-2 border-slate-600" />
                  <span className="text-slate-500">Portfolio linked</span>
                </div>
              </div>

              <p className="text-xs text-amber-400 mt-4">
                Add 1 more skill to boost visibility by 2x
              </p>
            </div>

            {/* Quick Stats */}
            <div className="bento-card-dark">
              <h3 className="text-lg font-semibold text-white mb-4">Quick Stats</h3>
              <div className="space-y-3">
                <div className="flex justify-between items-center p-3 bg-slate-800/50 rounded-lg">
                  <div className="flex items-center gap-2">
                    <Briefcase className="w-4 h-4 text-cyan-400" />
                    <span className="text-sm text-slate-300">Applications</span>
                  </div>
                  <span className="text-lg font-bold text-white">{user.applications}</span>
                </div>
                <div className="flex justify-between items-center p-3 bg-slate-800/50 rounded-lg">
                  <div className="flex items-center gap-2">
                    <Target className="w-4 h-4 text-amber-400" />
                    <span className="text-sm text-slate-300">Saved Jobs</span>
                  </div>
                  <span className="text-lg font-bold text-white">{user.savedJobs}</span>
                </div>
                <div className="flex justify-between items-center p-3 bg-slate-800/50 rounded-lg">
                  <div className="flex items-center gap-2">
                    <BookOpen className="w-4 h-4 text-emerald-400" />
                    <span className="text-sm text-slate-300">Courses</span>
                  </div>
                  <span className="text-lg font-bold text-white">{user.enrolledCourses}</span>
                </div>
              </div>
            </div>
          </div>

          {/* Middle Column - Skills & Activity */}
          <div className="space-y-6">
            {/* Skills */}
            <div className="bento-card-dark">
              <h3 className="text-lg font-semibold text-white mb-4 flex items-center gap-2">
                <Zap className="w-5 h-5 text-amber-400" />
                Skills
              </h3>
              <div className="flex flex-wrap gap-2">
                {user.skills.map((skill, index) => (
                  <span 
                    key={index}
                    className="px-3 py-1.5 bg-cyan-500/10 text-cyan-400 border border-cyan-500/20 rounded-lg text-sm font-medium"
                  >
                    {skill}
                  </span>
                ))}
                <button className="px-3 py-1.5 bg-slate-800 text-slate-400 border border-slate-700 rounded-lg text-sm hover:bg-slate-700 transition-colors">
                  + Add Skill
                </button>
              </div>
            </div>

            {/* Recent Activity */}
            <div className="bento-card-dark">
              <h3 className="text-lg font-semibold text-white mb-4 flex items-center gap-2">
                <TrendingUp className="w-5 h-5 text-emerald-400" />
                Recent Activity
              </h3>
              <div className="space-y-3">
                <div className="flex items-start gap-3 p-3 bg-slate-800/50 rounded-lg">
                  <div className="p-2 bg-cyan-500/10 rounded-lg">
                    <Briefcase className="w-4 h-4 text-cyan-400" />
                  </div>
                  <div>
                    <p className="text-sm text-white">Applied to Senior Go Engineer</p>
                    <p className="text-xs text-slate-500">2 hours ago</p>
                  </div>
                </div>
                <div className="flex items-start gap-3 p-3 bg-slate-800/50 rounded-lg">
                  <div className="p-2 bg-amber-500/10 rounded-lg">
                    <Flame className="w-4 h-4 text-amber-400" />
                  </div>
                  <div>
                    <p className="text-sm text-white">7-day streak achieved!</p>
                    <p className="text-xs text-slate-500">1 day ago</p>
                  </div>
                </div>
                <div className="flex items-start gap-3 p-3 bg-slate-800/50 rounded-lg">
                  <div className="p-2 bg-emerald-500/10 rounded-lg">
                    <BookOpen className="w-4 h-4 text-emerald-400" />
                  </div>
                  <div>
                    <p className="text-sm text-white">Enrolled in Data Science course</p>
                    <p className="text-xs text-slate-500">3 days ago</p>
                  </div>
                </div>
              </div>
            </div>
          </div>

          {/* Right Column - Achievements */}
          <div className="space-y-6">
            {/* Achievements */}
            <div className="bento-card-dark">
              <h3 className="text-lg font-semibold text-white mb-4 flex items-center gap-2">
                <Award className="w-5 h-5 text-amber-400" />
                Achievements
              </h3>
              
              <div className="space-y-4">
                {/* Completed */}
                <div>
                  <p className="text-xs text-slate-400 mb-2 font-semibold">EARNED ({completedAchievements.length})</p>
                  <div className="space-y-2">
                    {completedAchievements.map(achievement => (
                      <div 
                        key={achievement.id}
                        className="flex items-center gap-3 p-3 bg-emerald-500/5 border border-emerald-500/20 rounded-lg"
                      >
                        <div className="p-2 bg-emerald-500/10 rounded-lg text-emerald-400">
                          {achievement.icon}
                        </div>
                        <div className="flex-1">
                          <p className="text-sm font-medium text-white">{achievement.title}</p>
                          <p className="text-xs text-slate-400">{achievement.description}</p>
                        </div>
                        <div className="text-xs text-emerald-400 font-semibold">
                          +{achievement.xpReward} XP
                        </div>
                      </div>
                    ))}
                  </div>
                </div>

                {/* In Progress */}
                {inProgressAchievements.length > 0 && (
                  <div>
                    <p className="text-xs text-slate-400 mb-2 font-semibold">IN PROGRESS ({inProgressAchievements.length})</p>
                    <div className="space-y-2">
                      {inProgressAchievements.map(achievement => (
                        <div 
                          key={achievement.id}
                          className="flex items-center gap-3 p-3 bg-slate-800/50 rounded-lg"
                        >
                          <div className="p-2 bg-slate-700 rounded-lg text-slate-400">
                            {achievement.icon}
                          </div>
                          <div className="flex-1">
                            <p className="text-sm font-medium text-white">{achievement.title}</p>
                            <p className="text-xs text-slate-400">{achievement.description}</p>
                            <div className="mt-2 w-full h-1.5 bg-slate-700 rounded-full overflow-hidden">
                              <div 
                                className="h-full bg-cyan-500 transition-all duration-500"
                                style={{ width: `${(achievement.progress / achievement.target) * 100}%` }}
                              />
                            </div>
                          </div>
                          <div className="text-xs text-slate-500">
                            {achievement.progress}/{achievement.target}
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>
                )}
              </div>
            </div>

            {/* Quick Actions */}
            <div className="bento-card-dark">
              <h3 className="text-lg font-semibold text-white mb-4">Quick Actions</h3>
              <div className="space-y-2">
                <button className="w-full flex items-center gap-3 p-3 bg-slate-800/50 rounded-lg text-left hover:bg-slate-700 transition-colors">
                  <Settings className="w-4 h-4 text-slate-400" />
                  <span className="text-sm text-slate-300">Settings</span>
                </button>
                <button className="w-full flex items-center gap-3 p-3 bg-slate-800/50 rounded-lg text-left hover:bg-slate-700 transition-colors">
                  <Bell className="w-4 h-4 text-slate-400" />
                  <span className="text-sm text-slate-300">Notifications</span>
                </button>
                <button className="w-full flex items-center gap-3 p-3 bg-slate-800/50 rounded-lg text-left hover:bg-slate-700 transition-colors">
                  <LogOut className="w-4 h-4 text-slate-400" />
                  <span className="text-sm text-slate-300">Logout</span>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}