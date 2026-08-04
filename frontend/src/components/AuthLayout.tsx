import { Navigate } from "react-router-dom";
import type { ReactNode } from "react";
import { useAppSelector } from "../store/hooks.ts";

function AuthLayout({ children, authentication = true }: { children: ReactNode; authentication?: boolean }) {
  const authStatus = useAppSelector((state) => state.auth.status);

  if (authentication && authStatus !== authentication) {
    return <Navigate to="/login" replace />;
  }
  if (!authentication && authStatus !== authentication) {
    return <Navigate to="/" replace />;
  }

  return <>{children}</>;
}

export default AuthLayout;
