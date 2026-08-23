import React from "react";
import { cn } from "@/lib/utils";

export interface BadgeProps extends React.HTMLAttributes<HTMLDivElement> {
  variant?: "emerald" | "amber" | "cyan" | "slate";
}

export const Badge = React.forwardRef<HTMLDivElement, BadgeProps>(
  ({ className, variant = "slate", ...props }, ref) => {
    const variants = {
      emerald: "bg-emerald-500/10 text-emerald-500 border-emerald-500/20",
      amber: "bg-amber-500/10 text-amber-500 border-amber-500/20",
      cyan: "bg-cyan-500/10 text-cyan-500 border-cyan-500/20",
      slate: "bg-slate-500/10 text-slate-500 border-slate-500/20"
    };

    return (
      <div
        ref={ref}
        className={cn(
          "px-3 py-1 rounded-full text-xs font-medium border",
          variants[variant],
          className
        )}
        {...props}
      />
    );
  }
);

Badge.displayName = "Badge";