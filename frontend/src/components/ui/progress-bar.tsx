import React from "react";
import { cn } from "@/lib/utils";

export interface ProgressBarProps extends React.HTMLAttributes<HTMLDivElement> {
  value: number;
  max?: number;
  variant?: "cyan" | "emerald" | "amber" | "gradient";
  size?: "sm" | "md" | "lg";
}

export const ProgressBar = React.forwardRef<HTMLDivElement, ProgressBarProps>(
  ({ className, value, max = 100, variant = "cyan", size = "md", ...props }, ref) => {
    const percentage = Math.min(Math.max((value / max) * 100, 0), 100);
    
    const variants = {
      cyan: "bg-cyan-500",
      emerald: "bg-emerald-500",
      amber: "bg-amber-500",
      gradient: "bg-gradient-to-r from-cyan-500 to-emerald-400"
    };
    
    const sizes = {
      sm: "h-1.5",
      md: "h-2",
      lg: "h-3"
    };

    return (
      <div
        ref={ref}
        className={cn(
          "w-full bg-slate-700 rounded-full overflow-hidden",
          sizes[size],
          className
        )}
        {...props}
      >
        <div
          className={cn(
            "h-full transition-all duration-500 rounded-full",
            variants[variant]
          )}
          style={{ width: `${percentage}%` }}
        />
      </div>
    );
  }
);

ProgressBar.displayName = "ProgressBar";