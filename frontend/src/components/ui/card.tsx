import React from "react";
import { cn } from "@/lib/utils";

export interface CardProps extends React.HTMLAttributes<HTMLDivElement> {
  variant?: "default" | "dark" | "gradient";
}

export const Card = React.forwardRef<HTMLDivElement, CardProps>(
  ({ className, variant = "default", ...props }, ref) => {
    const variants = {
      default: "bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800",
      dark: "bg-slate-900/90 border border-slate-800 hover:border-slate-700",
      gradient: "bg-gradient-to-br from-slate-900 to-slate-800 border border-slate-700"
    };

    return (
      <div
        ref={ref}
        className={cn(
          "rounded-bento-lg p-6 shadow-bento transition-all duration-300 hover:shadow-bento-hover hover:scale-[1.02]",
          variants[variant],
          className
        )}
        {...props}
      />
    );
  }
);

Card.displayName = "Card";