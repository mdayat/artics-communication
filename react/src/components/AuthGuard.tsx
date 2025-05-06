import { Outlet, useLocation, useNavigate } from "react-router";
import { useEffect } from "react";
import { useAuthContext } from "@/contexts/AuthProvider";
import { Layout } from "./Layout";

function AuthGuard() {
  const location = useLocation();
  const navigate = useNavigate();
  const { user, isLoading } = useAuthContext();

  useEffect(() => {
    if (isLoading) {
      return;
    }

    if (
      !user &&
      location.pathname !== "/login" &&
      location.pathname !== "/registration"
    ) {
      navigate("/login", { replace: true });
    } else if (
      user &&
      (location.pathname === "/login" || location.pathname === "/registration")
    ) {
      navigate("/", { replace: true });
    } else if (
      user &&
      user.role === "admin" &&
      location.pathname === "/history"
    ) {
      navigate("/", { replace: true });
    }
  }, [isLoading, user, location, navigate]);

  if (isLoading) {
    return <></>;
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
