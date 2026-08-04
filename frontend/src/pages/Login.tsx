import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useForm } from "react-hook-form";
import { Button, Input } from "../components/index.ts";
import AuthIcon from "../components/AuthIcon.tsx";
import { useAppDispatch } from "../store/hooks.ts";
import { login as authLogin } from "../store/authSlice.ts";
import authService from "../api/auth.ts";

interface LoginForm {
  identifier: string;
  password: string;
}

function UserIcon({ size = 18 }: { size?: number }) {
  return (
    <svg
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    >
      <path d="M20 21v-2a4 4 0 00-4-4H8a4 4 0 00-4 4v2" />
      <circle cx="12" cy="7" r="4" />
    </svg>
  );
}

function LockIcon({ size = 18 }: { size?: number }) {
  return (
    <svg
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    >
      <rect x="3" y="11" width="18" height="11" rx="2" ry="2" />
      <path d="M7 11V7a5 5 0 0110 0v4" />
    </svg>
  );
}

function Login() {
  const navigate = useNavigate();
  const dispatch = useAppDispatch();
  const { register, handleSubmit } = useForm<LoginForm>();
  const [error, setError] = useState("");

  const submit = handleSubmit(async (data) => {
    setError("");
    try {
      const user = await authService.login({
        identifier: data.identifier,
        password: data.password,
      });
      dispatch(authLogin(user));
      navigate("/");
    } catch (err: unknown) {
      const msg =
        err && typeof err === "object" && "response" in err
          ? (err as { response?: { data?: { message?: string } } }).response
              ?.data?.message
          : "Login failed";
      setError(msg ?? "Login failed");
    }
  });

  return (
    <div className="flex min-h-screen items-center justify-center px-4 bg-gray-950">
      <div className="w-full max-w-sm">
        <div className="flex gap-6 items-center justify-center">
          <AuthIcon />
          <div className="mb-6 text-center">
            <Link to="/login" className="text-lg font-semibold text-white">
              Login
            </Link>
          </div>
        </div>
        {error && (
          <p className="mb-4 text-center text-sm text-red-400">{error}</p>
        )}
        <form onSubmit={submit}>
          <div className="space-y-4 text-center">
            <Input
              placeholder="USERNAME"
              icon={UserIcon}
              {...register("identifier", { required: true })}
            />
            <Input
              type="password"
              placeholder="PASSWORD"
              icon={LockIcon}
              {...register("password", { required: true })}
            />
            <Button type="submit" className="w-full">
              LOGIN
            </Button>
            <Link to="/register" className="text-lg text-gray-400">
              Don't have an account?{" "}
              <span className="hover:text-white hover:underline">Sign Up</span>
            </Link>
          </div>
        </form>
      </div>
    </div>
  );
}

export default Login;
