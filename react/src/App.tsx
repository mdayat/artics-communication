import { BrowserRouter, Navigate, Route, Routes } from "react-router";
import { lazy, Suspense } from "react";
import { Toaster } from "./components/ui/Sonner";
import { AuthProvider } from "./contexts/AuthProvider";
import { AuthGuard } from "./components/AuthGuard";

const Home = lazy(() =>
  import("./pages/home").then(({ Home }) => ({
    default: Home,
  }))
);

const History = lazy(() =>
  import("./pages/history").then(({ History }) => ({
    default: History,
  }))
);

const Registration = lazy(() =>
  import("./pages/registration").then(({ Registration }) => ({
    default: Registration,
  }))
);

const Login = lazy(() =>
  import("./pages/login").then(({ Login }) => ({
    default: Login,
  }))
);

function App() {
  return (
    <BrowserRouter>
      <Toaster />
      <AuthProvider>
        <Routes>
          <Route
            element={
              <Suspense fallback={<></>}>
                <AuthGuard />
              </Suspense>
            }
          >
            <Route
              path="/"
              element={
                <Suspense fallback={<></>}>
                  <Home />
                </Suspense>
              }
            />

            <Route
              path="/history"
              element={
                <Suspense fallback={<></>}>
                  <History />
                </Suspense>
              }
            />

            <Route
              path="/login"
              element={
                <Suspense fallback={<></>}>
                  <Login />
                </Suspense>
              }
            />

            <Route
              path="/registration"
              element={
                <Suspense fallback={<></>}>
                  <Registration />
                </Suspense>
              }
            />
          </Route>

          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </AuthProvider>
    </BrowserRouter>
  );
}

export default App;
