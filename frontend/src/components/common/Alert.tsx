import { cn } from "../../utils/cn.ts";

interface AlertProps {
  type?: "error" | "success";
  message: string;
  className?: string;
}

export function Alert({ type = "error", message, className }: AlertProps) {
  return (
    <div
      className={cn(
        "rounded-lg px-4 py-3 text-sm",
        type === "error" ? "bg-red-900/50 text-red-300" : "bg-green-900/50 text-green-300",
        className
      )}
    >
      {message}
    </div>
  );
}
