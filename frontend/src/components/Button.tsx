import type { ButtonHTMLAttributes } from "react";

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  bgColor?: string;
  textColor?: string;
  className?: string;
}

function Button({
  children,
  type = "button",
  bgColor = "bg-blue-600",
  textColor = "text-white",
  className = "",
  ...props
}: ButtonProps) {
  return (
    <button
      className={`px-4 py-2.5 rounded-lg font-semibold ${bgColor} ${textColor} hover:opacity-90 disabled:cursor-not-allowed disabled:opacity-50 ${className}`}
      type={type}
      {...props}
    >
      {children}
    </button>
  );
}

export default Button;
