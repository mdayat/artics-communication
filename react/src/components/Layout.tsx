import type React from "react";
import { useState, type Dispatch, type SetStateAction } from "react";
import { Home, Menu, LogOut, History } from "lucide-react";
import { Button } from "@/components/ui/Button";
import { Sheet, SheetContent } from "@/components/ui/Sheet";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/Avatar";
import { Link } from "react-router";
import { useAuthContext } from "@/contexts/AuthProvider";

function Layout({ children }: { children: React.ReactNode }) {
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);
  const { user } = useAuthContext();

  return (
    <div className="flex min-h-screen bg-background">
      <Sheet open={isSidebarOpen} onOpenChange={setIsSidebarOpen}>
        <SheetContent
          side="left"
          className="p-0 bg-slate-800 text-white border-r-0 w-fit"
          onCloseAutoFocus={(event) => event.preventDefault()}
        >
          <Sidebar setIsSidebarOpen={setIsSidebarOpen} />
        </SheetContent>
      </Sheet>

      <div className="hidden md:block">
        <Sidebar />
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
                {user?.name ?? "John Doe"}
              </span>

              <Avatar>
                <AvatarImage
                  src="/placeholder-user.jpg"
                  alt={user?.name ?? "John Doe"}
                />

                <AvatarFallback>
                  {user?.name.split(" ").map((word) => word[0]) ?? "JD"}
                </AvatarFallback>
              </Avatar>
            </div>
          </div>
        </header>

        <main className="flex-1 p-6">{children}</main>
      </div>
    </div>
  );
}

interface SidebarProps {
  setIsSidebarOpen?: Dispatch<SetStateAction<boolean>>;
}

function Sidebar({ setIsSidebarOpen }: SidebarProps) {
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

          {user?.role === "user" ? (
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
          variant="ghost"
          className="cursor-pointer w-full justify-start text-sm font-medium"
        >
          <LogOut className="mr-2 h-4 w-4" />
          Logout
        </Button>
      </div>
    </div>
  );
}

export { Layout };
