import { cn } from "../../utils/cn.ts";

interface TextInputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  label: string;
  error?: string | null;
}

export function TextInput({ label, error, className, ...props }: TextInputProps) {
  return (
    <div className="space-y-1">
      <label className="block text-sm font-medium text-gray-300">{label}</label>
      <input
        className={cn(
          "w-full rounded-lg border bg-gray-800 px-4 py-2.5 text-gray-100 placeholder-gray-500",
          "outline-none transition-colors focus:border-blue-500 focus:ring-1 focus:ring-blue-500",
          error ? "border-red-500" : "border-gray-700",
          className
        )}
        {...props}
      />
      {error && <p className="text-sm text-red-400">{error}</p>}
    </div>
  );
}
