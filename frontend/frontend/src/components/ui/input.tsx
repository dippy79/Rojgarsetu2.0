import React from "react";
import { cn } from "@/lib/utils";

export interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  variant?: "default" | "dark";
}

export const Input = React.forwardRef<HTMLInputElement, InputProps>(
  ({ className, variant = "default", ...props }, ref) => {
    const variants = {
      default: "bg-white dark:bg-slate-950 border-slate-300 dark:border-slate-800 text-slate-900 dark:text-slate-100",
      dark: "bg-slate-950 border-slate-800 text-slate-100"
    };

    return (
      <input
        ref={ref}
        className={cn(
          "w-full rounded-bento-sm px-4 py-3.5 text-sm focus:outline-none focus:ring-2 focus:ring-cyan-500 focus:border-transparent transition-colors border",
          variants[variant],
          className
        )}
        {...props}
      />
    );
  }
);

Input.displayName = "Input";