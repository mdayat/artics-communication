import { useRouteError } from "react-router";

function ErrorBoundary() {
  const error = useRouteError();
  let errMsg = "";
  if (error instanceof Error) {
    errMsg = error.message;
  } else if (
    (typeof error !== "undefined" ||
      typeof error !== "symbol" ||
      typeof error !== "function") &&
    typeof error === "object"
  ) {
    errMsg = JSON.stringify(error);
  } else {
    errMsg = "Unknown Error Type";
    console.error(error);
  }

  return (
    <div className="absolute top-1/2 -translate-y-1/2 left-1/2 -translate-x-1/2 flex flex-col gap-y-4">
      <h1 className="text-2xl font-bold text-center">Oops!</h1>
      <p className="font-medium text-center">
        Sorry, an unexpected error has occurred.
      </p>
      <p className="italic text-center">{errMsg}</p>
    </div>
  );
}

export { ErrorBoundary };
