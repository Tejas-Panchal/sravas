import { useAuth } from "../hooks/useAuth.ts";

export function DashboardPage() {
  const { user, logout } = useAuth();

  return (
    <div className="min-h-screen bg-gray-950 text-gray-100">
      <header className="border-b border-gray-800 bg-gray-900">
        <div className="mx-auto flex max-w-5xl items-center justify-between px-4 py-4">
          <h1 className="text-xl font-bold text-white">Sravas</h1>
          <div className="flex items-center gap-4">
            <span className="text-sm text-gray-300">{user?.fullName}</span>
            <button
              onClick={() => logout()}
              className="rounded-lg bg-gray-800 px-3 py-1.5 text-sm text-gray-300 transition-colors hover:bg-gray-700"
            >
              Logout
            </button>
          </div>
        </div>
      </header>
      <main className="mx-auto max-w-5xl px-4 py-8">
        <div className="rounded-2xl bg-gray-900 p-8">
          <h2 className="text-2xl font-bold text-white">Welcome back, {user?.fullName}</h2>
          <p className="mt-2 text-gray-400">@{user?.username}</p>
          <div className="mt-6 flex gap-6">
            <img
              src={user?.avatar}
              alt={`${user?.username}'s avatar`}
              className="h-20 w-20 rounded-full object-cover ring-2 ring-gray-700"
            />
            {user?.coverImage && (
              <img
                src={user.coverImage}
                alt="Cover"
                className="h-20 rounded-xl object-cover ring-2 ring-gray-700"
              />
            )}
          </div>
        </div>
      </main>
    </div>
  );
}
