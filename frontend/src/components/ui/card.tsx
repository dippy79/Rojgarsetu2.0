import React from "react";
import { cn } from "@/lib/utils";

export interface CardProps extends React.HTMLAttributes<HTMLDivElement> {
  variant?: "default" | "dark" | "gradient";
}

export const Card = React.forwardRef<HTMLDivElement, CardProps>(
  ({ className, variant = "default", ...props }, ref) => {
    const variants = {
      default: "bg-white dark:bg-stone-900 border border-stone-200 dark:border-stone-800",
      dark: "bg-stone-900/90 border border-stone-800 hover:border-stone-700",
      gradient: "bg-gradient-to-br from-stone-900 to-stone-800 border border-stone-700"
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