"use client";

import { InputHTMLAttributes, forwardRef } from "react";

export const Input = forwardRef<HTMLInputElement, InputHTMLAttributes<HTMLInputElement>>(
  function Input({ className = "", ...rest }, ref) {
    return (
      <input
        ref={ref}
        className={`h-10 w-full rounded-lg border border-[var(--border)] bg-transparent px-3 text-sm outline-none transition placeholder:text-muted focus:ring-2 focus:ring-[var(--ring)] ${className}`}
        {...rest}
      />
    );
  }
);