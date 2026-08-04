import { useEffect, useMemo, useRef } from "react";

interface FileInputProps {
  file: File | null;
  onChange: (file: File | null) => void;
  className?: string;
}

function FileInput({ file, onChange, className = "" }: FileInputProps) {
  const inputRef = useRef<HTMLInputElement>(null);
  const preview = useMemo(
    () => (file ? URL.createObjectURL(file) : null),
    [file],
  );

  useEffect(() => {
    return () => {
      if (preview) URL.revokeObjectURL(preview);
    };
  }, [preview]);

  return (
    <div className={`flex flex-col items-center pt-14 ${className}`}>
      <button
        type="button"
        onClick={() => inputRef.current?.click()}
        className="group relative h-30 w-30 overflow-hidden rounded-full border-2 border-dashed border-gray-600 bg-gray-800 transition-colors hover:border-blue-500"
      >
        {preview ? (
          <img
            src={preview}
            alt="Avatar preview"
            className="h-full w-full object-cover"
          />
        ) : (
          <div className="flex h-full w-full items-center justify-center text-gray-500 group-hover:text-gray-300">
            <svg
              width="40"
              height="40"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="1.5"
              strokeLinecap="round"
              strokeLinejoin="round"
            >
              <path d="M20 21v-2a4 4 0 00-4-4H8a4 4 0 00-4 4v2" />
              <circle cx="12" cy="7" r="4" />
            </svg>
          </div>
        )}
      </button>
      <input
        ref={inputRef}
        type="file"
        accept="image/*"
        className="hidden"
        onChange={(e) => onChange(e.target.files?.[0] ?? null)}
      />
    </div>
  );
}

export default FileInput;
