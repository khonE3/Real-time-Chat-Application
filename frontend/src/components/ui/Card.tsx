"use client";

import { HTMLAttributes, forwardRef } from "react";

interface CardProps extends HTMLAttributes<HTMLDivElement> {
  variant?: "default" | "elevated";
}

const Card = forwardRef<HTMLDivElement, CardProps>(
  ({ className = "", variant = "default", children, ...props }, ref) => {
    const baseStyles = "bg-white rounded-xl overflow-hidden";

    const variants = {
      default: "border border-[var(--color-earth-200)] shadow-sm",
      elevated: "shadow-lg shadow-[var(--color-earth-200)]/50",
    };

    return (
      <div
        ref={ref}
        className={`${baseStyles} ${variants[variant]} ${className}`}
        {...props}
      >
        {children}
      </div>
    );
  }
);

Card.displayName = "Card";

interface CardHeaderProps extends HTMLAttributes<HTMLDivElement> {
  gradient?: boolean;
}

const CardHeader = forwardRef<HTMLDivElement, CardHeaderProps>(
  ({ className = "", gradient = true, children, ...props }, ref) => {
    const baseStyles = "px-4 py-3 font-semibold";
    const gradientStyles = gradient
      ? "bg-gradient-to-r from-[var(--color-gold-500)] to-[var(--color-gold-600)] text-white"
      : "bg-[var(--color-earth-50)] text-[var(--color-earth-800)] border-b border-[var(--color-earth-200)]";

    return (
      <div
        ref={ref}
        className={`${baseStyles} ${gradientStyles} ${className}`}
        {...props}
      >
        {children}
      </div>
    );
  }
);

CardHeader.displayName = "CardHeader";

const CardContent = forwardRef<HTMLDivElement, HTMLAttributes<HTMLDivElement>>(
  ({ className = "", children, ...props }, ref) => {
    return (
      <div ref={ref} className={`p-4 ${className}`} {...props}>
        {children}
      </div>
    );
  }
);

CardContent.displayName = "CardContent";

const CardFooter = forwardRef<HTMLDivElement, HTMLAttributes<HTMLDivElement>>(
  ({ className = "", children, ...props }, ref) => {
    return (
      <div
        ref={ref}
        className={`px-4 py-3 border-t border-[var(--color-earth-200)] bg-[var(--color-earth-50)] ${className}`}
        {...props}
      >
        {children}
      </div>
    );
  }
);

CardFooter.displayName = "CardFooter";

export { Card, CardHeader, CardContent, CardFooter };
