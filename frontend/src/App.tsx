import { useEffect, useState } from "react";
import { Outlet } from "react-router-dom";
import { Header, Footer } from "./components/index.ts";
import { useAppDispatch } from "./store/hooks.ts";
import { login, logout } from "./store/authSlice.ts";
import authService from "./api/auth.ts";

function App() {
  const [loading, setLoading] = useState(true);
  const dispatch = useAppDispatch();

  useEffect(() => {
    authService
      .getCurrentUser()
      .then((userData) => {
        if (userData) {
          dispatch(login(userData));
        } else {
          dispatch(logout());
        }
      })
      .finally(() => setLoading(false));
  }, [dispatch]);

  return !loading ? (
    <div className="min-h-screen flex flex-col bg-gray-950 text-gray-100">
      <Header />
      <main className="flex-1">
        <Outlet />
      </main>
      <Footer />
    </div>
  ) : (
    <div className="flex min-h-screen items-center justify-center bg-gray-950 text-gray-100">
      Loading...
    </div>
  );
}

export default App;
