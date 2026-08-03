import { useFilePreview } from "../../hooks/useFilePreview.ts";

interface FileInputProps {
  label: string;
  file: File | null;
  onChange: (file: File | null) => void;
  accept?: string;
}

export function FileInput({ label, file, onChange, accept = "image/*" }: FileInputProps) {
  const preview = useFilePreview(file);

  return (
    <div className="space-y-1">
      <label className="block text-sm font-medium text-gray-300">{label}</label>
      <div className="flex items-center gap-4">
        {preview && (
          <img
            src={preview}
            alt={`${label} preview`}
            className="h-16 w-16 rounded-lg object-cover"
          />
        )}
        <input
          type="file"
          accept={accept}
          onChange={(e) => onChange(e.target.files?.[0] ?? null)}
          className="text-sm text-gray-400 file:mr-4 file:rounded-lg file:border-0 file:bg-blue-600 file:px-4 file:py-2 file:text-sm file:font-medium file:text-white hover:file:bg-blue-500"
        />
      </div>
    </div>
  );
}
