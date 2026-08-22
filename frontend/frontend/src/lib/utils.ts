import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function formatDate(dateString: string): string {
  const date = new Date(dateString);
  return date.toLocaleDateString('en-IN', { 
    day: 'numeric', 
    month: 'short', 
    year: 'numeric' 
  });
}

export function getDaysLeft(dateString: string): number | null {
  if (!dateString) return null;
  const lastDate = new Date(dateString);
  const today = new Date();
  const diffTime = lastDate.getTime() - today.getTime();
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
  return diffDays;
}

export function formatSalary(salaryMin?: number, salaryMax?: number): string {
  if (!salaryMin && !salaryMax) return 'Not disclosed';
  if (salaryMin && salaryMax) {
    return `₹${(salaryMin / 100000).toFixed(1)}L - ₹${(salaryMax / 100000).toFixed(1)}L`;
  }
  if (salaryMin) return `₹${(salaryMin / 100000).toFixed(1)}L+`;
  return `Up to ₹${(salaryMax! / 100000).toFixed(1)}L`;
}

export function formatFees(fees?: string | number): string {
  if (!fees || fees === '0' || fees === 'free') return 'Free';
  if (typeof fees === 'number') return `₹${fees.toLocaleString()}`;
  return fees;
}

export function getMatchScoreColor(score: number): string {
  if (score >= 90) return 'text-emerald-500';
  if (score >= 75) return 'text-cyan-500';
  if (score >= 60) return 'text-amber-500';
  return 'text-slate-500';
}

export function getMatchScoreBg(score: number): string {
  if (score >= 90) return 'bg-emerald-500/10 border-emerald-500/20';
  if (score >= 75) return 'bg-cyan-500/10 border-cyan-500/20';
  if (score >= 60) return 'bg-amber-500/10 border-amber-500/20';
  return 'bg-slate-500/10 border-slate-500/20';
}