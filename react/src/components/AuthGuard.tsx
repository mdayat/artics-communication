import { Outlet, useLocation, useNavigate } from "react-router";
import { useEffect } from "react";
import { useAuthContext } from "@/contexts/AuthProvider";

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
    }
  }, [isLoading, user, location, navigate]);

  if (isLoading) {
    return <></>;
  }

  return <Outlet />;
}

export { AuthGuard };
