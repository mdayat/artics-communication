import { BrowserRouter, Navigate, Route, Routes } from "react-router";
import { lazy, Suspense } from "react";
import { Toaster } from "./components/ui/Sonner";

const Registration = lazy(() =>
  import("./pages/registration").then(({ Registration }) => ({
    default: Registration,
  }))
);

function App() {
  return (
    <BrowserRouter>
      <Toaster />
      <Routes>
        <Route
        // element={
        //   <Suspense fallback={<></>}>
        //     <ProtectedRouteView />
        //   </Suspense>
        // }
        >
          <Route
            path="/login"
            element={
              <Suspense fallback={<></>}>
                <Registration />
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
    </BrowserRouter>
  );
}

export default App;
