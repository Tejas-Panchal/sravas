import type { ReactNode } from "react";

interface AuthLayoutProps {
  title: string;
  children: ReactNode;
}

export function AuthLayout({ title, children }: AuthLayoutProps) {
  return (
    <div className="flex min-h-screen items-center justify-center px-4">
      <div className="w-full max-w-md rounded-2xl bg-gray-900 p-8 shadow-xl">
        <h1 className="mb-6 text-center text-2xl font-bold text-white">{title}</h1>
        {children}
      </div>
    </div>
  );
}
