import { useAppSelector } from "../store/hooks.ts";

function Home() {
  const userData = useAppSelector((state) => state.auth.userData);

  return (
    <main className="mx-auto max-w-6xl px-4 py-8">
      <div className="rounded-2xl bg-gray-900 p-8">
        <h2 className="text-2xl font-bold text-white">Welcome back, {userData?.fullName}</h2>
        <p className="mt-2 text-gray-400">@{userData?.username}</p>
        <div className="mt-6 flex gap-6">
          {userData?.avatar && (
            <img
              src={userData.avatar}
              alt={`${userData.username}'s avatar`}
              className="h-20 w-20 rounded-full object-cover ring-2 ring-gray-700"
            />
          )}
          {userData?.coverImage && (
            <img
              src={userData.coverImage}
              alt="Cover"
              className="h-20 rounded-xl object-cover ring-2 ring-gray-700"
            />
          )}
        </div>
      </div>
    </main>
  );
}

export default Home;
