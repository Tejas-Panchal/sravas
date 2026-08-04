import { forwardRef, useId } from "react";
import type { InputHTMLAttributes } from "react";

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  icon?: React.ElementType;
  className?: string;
}

const Input = forwardRef<HTMLInputElement, InputProps>(function Input(
  { label, icon: Icon, type = "text", className = "", ...props },
  ref,
) {
  const id = useId();
  return (
    <div className="w-full">
      {label && (
        <label className="inline-block mb-1 pl-1 text-sm text-gray-300" htmlFor={id}>
          {label}
        </label>
      )}
      <div className="relative">
        {Icon && (
          <span className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500">
            <Icon size={18} />
          </span>
        )}
        <input
          id={id}
          ref={ref}
          type={type}
          className={
            Icon
              ? `w-full rounded-lg border border-gray-700 bg-gray-800 px-4 py-2.5 pl-10 text-gray-100 placeholder-gray-500 outline-none transition-colors focus:border-blue-500 focus:ring-1 focus:ring-blue-500 ${className}`
              : `w-full rounded-lg border border-gray-700 bg-gray-800 px-4 py-2.5 text-gray-100 placeholder-gray-500 outline-none transition-colors focus:border-blue-500 focus:ring-1 focus:ring-blue-500 ${className}`
          }
          {...props}
        />
      </div>
    </div>
  );
});

export default Input;
