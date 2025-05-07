import { useState, type Dispatch, type SetStateAction } from "react";
import { Home, Menu, LogOut, History, Loader2 } from "lucide-react";
import { Button } from "@/components/ui/Button";
import { Sheet, SheetContent } from "@/components/ui/Sheet";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/Avatar";
import { Link, Outlet } from "react-router";
import { useAuthContext } from "@/contexts/AuthProvider";
import { toast } from "sonner";
import { axiosInstance } from "@/lib/axios";

function Layout() {
  const [isLoading, setIsLoading] = useState(false);
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);
  const { user, setUser } = useAuthContext();

  const handleLogout = async () => {
    setIsLoading(true);
    try {
      const res = await axiosInstance.post("/auth/logout");
      if (res.status === 204) {
        toast.success("Logout successful", { richColors: true });
        setUser(null);
      } else {
        throw new Error(`unknown status code: ${res.status}`);
      }
    } catch (error) {
      console.error(new Error("failed to logout", { cause: error }));
      toast.error("Logout failed, please try again", { richColors: true });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="flex min-h-screen bg-background">
      <Sheet open={isSidebarOpen} onOpenChange={setIsSidebarOpen}>
        <SheetContent
          side="left"
          className="p-0 bg-slate-800 text-white border-r-0 w-fit"
          onCloseAutoFocus={(event) => event.preventDefault()}
        >
          <Sidebar
            isLoading={isLoading}
            handleLogout={handleLogout}
            setIsSidebarOpen={setIsSidebarOpen}
          />
        </SheetContent>
      </Sheet>

      <div className="hidden md:block">
        <Sidebar isLoading={isLoading} handleLogout={handleLogout} />
      </div>

      <div className="flex flex-col flex-1">
        <header className="sticky top-0 z-10 flex h-16 items-center gap-4 border-b-2 bg-background px-6">
          <Button
            variant="outline"
            size="icon"
            className="cursor-pointer md:hidden"
            onClick={() => setIsSidebarOpen(true)}
          >
            <Menu className="h-5 w-5" />
            <span className="sr-only">Toggle Menu</span>
          </Button>

          <div className="ml-auto flex items-center gap-4">
            <div className="flex items-center gap-2">
              <span className="text-sm font-medium">
                {user ? user.name : "John Doe"}
              </span>

              <Avatar>
                <AvatarImage
                  src="/placeholder-user.jpg"
                  alt={user ? user.name : "John Doe"}
                />

                <AvatarFallback>
                  {user ? user.name.split(" ").map((word) => word[0]) : "JD"}
                </AvatarFallback>
              </Avatar>
            </div>
          </div>
        </header>

        <main className="flex-1 p-6">
          <Outlet />
        </main>
      </div>
    </div>
  );
}

interface SidebarProps {
  isLoading: boolean;
  handleLogout: () => void;
  setIsSidebarOpen?: Dispatch<SetStateAction<boolean>>;
}

function Sidebar({ isLoading, handleLogout, setIsSidebarOpen }: SidebarProps) {
  const { user } = useAuthContext();

  return (
    <div className="flex h-full w-64 flex-col bg-slate-800 text-white">
      <div className="flex h-16 items-center px-6">
        <span className="font-semibold text-xl">Reservation App</span>
      </div>

      <nav className="flex-1 overflow-auto py-6 px-3">
        <div className="space-y-1">
          <Link
            onClick={() => setIsSidebarOpen && setIsSidebarOpen(false)}
            to="/"
            className="flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium hover:bg-accent hover:text-accent-foreground"
          >
            <Home className="h-4 w-4" />
            <span>Home</span>
          </Link>

          {user && user.role === "user" ? (
            <Link
              onClick={() => setIsSidebarOpen && setIsSidebarOpen(false)}
              to="/history"
              className="flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium hover:bg-accent hover:text-accent-foreground"
            >
              <History className="h-4 w-4" />
              <span>History</span>
            </Link>
          ) : (
            <></>
          )}
        </div>
      </nav>

      <div className="p-3">
        <Button
          disabled={isLoading}
          onClick={handleLogout}
          variant="ghost"
          className="cursor-pointer w-full justify-start text-sm font-medium"
        >
          {isLoading ? (
            <>
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              Ending Session
            </>
          ) : (
            <>
              <LogOut className="mr-2 h-4 w-4" />
              Logout
            </>
          )}
        </Button>
      </div>
    </div>
  );
}

export { Layout };
