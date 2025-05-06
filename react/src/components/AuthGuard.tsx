import { Navigate, Outlet, useLocation } from "react-router";
import { useAuthContext } from "@/contexts/AuthProvider";
import { Layout } from "./Layout";

function AuthGuard() {
  const location = useLocation();
  const { user, isLoading } = useAuthContext();

  if (isLoading) {
    return <></>;
  }

  if (
    !user &&
    location.pathname !== "/login" &&
    location.pathname !== "/registration"
  ) {
    return <Navigate to="/login" replace />;
  }

  if (
    user &&
    (location.pathname === "/login" || location.pathname === "/registration")
  ) {
    return <Navigate to="/" replace />;
  }

  if (user && user.role === "admin" && location.pathname === "/history") {
    return <Navigate to="/" replace />;
  }

  if (location.pathname === "/login" || location.pathname === "/registration") {
    return <Outlet />;
  }

  return (
    <Layout>
      <Outlet />
    </Layout>
  );
}

export { AuthGuard };
