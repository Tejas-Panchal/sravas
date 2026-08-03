import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { TextInput } from "../../components/common/TextInput.tsx";
import { SubmitButton } from "../../components/common/SubmitButton.tsx";
import { Alert } from "../../components/common/Alert.tsx";
import { AuthLayout } from "../../components/layout/AuthLayout.tsx";
import { useAuth } from "../../hooks/useAuth.ts";
import { required } from "../../utils/validators.ts";

export function LoginPage() {
  const { login } = useAuth();
  const navigate = useNavigate();
  const [identifier, setIdentifier] = useState("");
  const [password, setPassword] = useState("");
  const [errors, setErrors] = useState<{ identifier?: string | null; password?: string | null }>({});
  const [serverError, setServerError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setServerError(null);

    const idErr = required(identifier, "Username or email");
    const pwErr = required(password, "Password");
    setErrors({ identifier: idErr, password: pwErr });
    if (idErr || pwErr) return;

    setLoading(true);
    try {
      await login(identifier, password);
      navigate("/dashboard", { replace: true });
    } catch (err: unknown) {
      const msg =
        err && typeof err === "object" && "response" in err
          ? (err as { response?: { data?: { message?: string } } }).response?.data?.message
          : "Login failed";
      setServerError(msg ?? "Login failed");
    } finally {
      setLoading(false);
    }
  }

  return (
    <AuthLayout title="Sign in to Sravas">
      <form onSubmit={handleSubmit} className="space-y-4">
        {serverError && <Alert message={serverError} />}
        <TextInput
          label="Username or email"
          placeholder="you@example.com"
          value={identifier}
          onChange={(e) => setIdentifier(e.target.value)}
          error={errors.identifier}
        />
        <TextInput
          label="Password"
          type="password"
          placeholder="••••••••"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          error={errors.password}
        />
        <SubmitButton loading={loading}>Sign in</SubmitButton>
        <p className="text-center text-sm text-gray-400">
          Don&apos;t have an account?{" "}
          <Link to="/register" className="text-blue-400 hover:underline">
            Create one
          </Link>
        </p>
      </form>
    </AuthLayout>
  );
}
