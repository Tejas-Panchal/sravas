import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useForm } from "react-hook-form";
import { Button, Input, FileInput } from "../components/index.ts";
import AuthIcon from "../components/AuthIcon.tsx";
import authService from "../api/auth.ts";

interface SignupForm {
  fullName: string;
  username: string;
  email: string;
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

function MailIcon({ size = 18 }: { size?: number }) {
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
      <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z" />
      <polyline points="22,6 12,13 2,6" />
    </svg>
  );
}

function Signup() {
  const navigate = useNavigate();
  const { register, handleSubmit } = useForm<SignupForm>();
  const [avatar, setAvatar] = useState<File | null>(null);
  const [error, setError] = useState("");

  const create = handleSubmit(async (data) => {
    setError("");
    if (!avatar) {
      setError("Avatar is required");
      return;
    }
    const formData = new FormData();
    formData.append("fullName", data.fullName);
    formData.append("username", data.username);
    formData.append("email", data.email);
    formData.append("password", data.password);
    formData.append("avatar", avatar);
    try {
      await authService.register(formData);
      navigate("/login");
    } catch (err: unknown) {
      const msg =
        err && typeof err === "object" && "response" in err
          ? (err as { response?: { data?: { message?: string } } }).response
              ?.data?.message
          : "Registration failed";
      setError(msg ?? "Registration failed");
    }
  });

  return (
    <div className="flex min-h-screen items-center justify-center px-4 bg-gray-950">
      <div className="w-full max-w-lg">
        <div className="flex gap-6 items-center justify-center">
          <AuthIcon />
          <div className="mb-6 text-center">
            <Link to="/register" className="text-lg font-semibold text-white">
              Sign Up
            </Link>
          </div>
          {error && (
            <p className="mb-4 text-center text-sm text-red-400">{error}</p>
          )}
        </div>
        <form onSubmit={create}>
          <div className="flex gap-6">
            <FileInput file={avatar} onChange={setAvatar} />
            <div className="flex-1 space-y-4 text-center">
              <Input
                placeholder="FULL NAME"
                icon={UserIcon}
                {...register("fullName", { required: true })}
              />
              <Input
                placeholder="USERNAME"
                icon={UserIcon}
                {...register("username", { required: true })}
              />
              <Input
                placeholder="EMAIL"
                type="email"
                icon={MailIcon}
                {...register("email", { required: true })}
              />
              <Input
                placeholder="PASSWORD"
                type="password"
                icon={LockIcon}
                {...register("password", { required: true })}
              />
              <Button type="submit" className="w-full">
                SIGN UP
              </Button>
              <Link to="/login" className="text-lg text-gray-400">
                Already have an account?{" "}
                <span className="hover:text-white hover:underline">
                  Sign In
                </span>
              </Link>
            </div>
          </div>
        </form>
      </div>
    </div>
  );
}

export default Signup;
