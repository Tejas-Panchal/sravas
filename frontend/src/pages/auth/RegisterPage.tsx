import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { TextInput } from "../../components/common/TextInput.tsx";
import { FileInput } from "../../components/common/FileInput.tsx";
import { SubmitButton } from "../../components/common/SubmitButton.tsx";
import { Alert } from "../../components/common/Alert.tsx";
import { AuthLayout } from "../../components/layout/AuthLayout.tsx";
import { useAuth } from "../../hooks/useAuth.ts";
import { required, isEmail } from "../../utils/validators.ts";

interface FormErrors {
  fullName?: string | null;
  username?: string | null;
  email?: string | null;
  password?: string | null;
  avatar?: string | null;
}

export function RegisterPage() {
  const { register } = useAuth();
  const navigate = useNavigate();
  const [fullName, setFullName] = useState("");
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [avatar, setAvatar] = useState<File | null>(null);
  const [errors, setErrors] = useState<FormErrors>({});
  const [serverError, setServerError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  function validate(): FormErrors {
    return {
      fullName: required(fullName, "Full name"),
      username: required(username, "Username"),
      email: !required(email, "Email") && !isEmail(email) ? "Invalid email" : required(email, "Email"),
      password: required(password, "Password"),
      avatar: !avatar ? "Avatar is required" : null,
    };
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setServerError(null);

    const errs = validate();
    setErrors(errs);
    if (Object.values(errs).some(Boolean)) return;

    setLoading(true);
    try {
      const fd = new FormData();
      fd.append("fullName", fullName);
      fd.append("username", username);
      fd.append("email", email);
      fd.append("password", password);
      fd.append("avatar", avatar!);
      await register(fd);
      navigate("/login");
    } catch (err: unknown) {
      const msg =
        err && typeof err === "object" && "response" in err
          ? (err as { response?: { data?: { message?: string } } }).response?.data?.message
          : "Registration failed";
      setServerError(msg ?? "Registration failed");
    } finally {
      setLoading(false);
    }
  }

  return (
    <AuthLayout title="Create your account">
      <form onSubmit={handleSubmit} className="space-y-4">
        {serverError && <Alert message={serverError} />}
        <TextInput
          label="Full name"
          placeholder="John Doe"
          value={fullName}
          onChange={(e) => setFullName(e.target.value)}
          error={errors.fullName}
        />
        <TextInput
          label="Username"
          placeholder="johndoe"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          error={errors.username}
        />
        <TextInput
          label="Email"
          type="email"
          placeholder="you@example.com"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          error={errors.email}
        />
        <TextInput
          label="Password"
          type="password"
          placeholder="••••••••"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          error={errors.password}
        />
        <FileInput label="Avatar" file={avatar} onChange={setAvatar} />
        {errors.avatar && <p className="text-sm text-red-400">{errors.avatar}</p>}
        <SubmitButton loading={loading}>Create account</SubmitButton>
        <p className="text-center text-sm text-gray-400">
          Already have an account?{" "}
          <Link to="/login" className="text-blue-400 hover:underline">
            Sign in
          </Link>
        </p>
      </form>
    </AuthLayout>
  );
}
