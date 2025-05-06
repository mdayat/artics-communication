import type { UserResponse } from "@/dtos/user";
import { axiosInstance } from "@/lib/axios";
import {
  createContext,
  type Dispatch,
  type SetStateAction,
  useContext,
  useEffect,
  useMemo,
  useState,
  type PropsWithChildren,
} from "react";
import { toast } from "sonner";

interface AuthContextType {
  user: UserResponse | null;
  setUser: Dispatch<SetStateAction<UserResponse | null>>;
  isLoading: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

function AuthProvider({ children }: PropsWithChildren) {
  const [user, setUser] = useState<UserResponse | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    (async () => {
      try {
        const res = await axiosInstance.get<UserResponse>("/users/me");
        if (res.status === 200) {
          setUser(res.data);
        } else if (res.status === 401) {
          // not doing anything
        } else if (res.status === 404) {
          toast.error("User not found", { richColors: true });
        } else {
          throw new Error(`unknown status code: ${res.status}`);
        }
      } catch (error) {
        console.error(new Error("failed to get user data", { cause: error }));
        toast.error("Something is wrong, please refresh your browser");
      } finally {
        setIsLoading(false);
      }
    })();
  }, []);

  const value = useMemo((): AuthContextType => {
    return { user, setUser, isLoading };
  }, [user, setUser, isLoading]);

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

function useAuthContext() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuthContext must be used within a AuthProvider");
  }
  return context;
}

export { AuthProvider, useAuthContext };
