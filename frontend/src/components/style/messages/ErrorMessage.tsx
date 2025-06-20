const ErrorMessage = ({
  message,
  error,
  reload,
}: {
  message?: string;
  error?: any;
  reload?: boolean;
}) => {
  if (!message) {
    message = `We encountered some kind of error. We are really sorry: ${JSON.stringify(error)}`;
  }
  return (
    <div className="flex flex-col items-center justify-center p-8 m-4 rounded-lg bg-error/10 text-error">
      <div className="text-5xl mb-4">⚠️</div>
      <h3 className="text-xl font-bold mb-2">Oops! Something went wrong</h3>
      <p className="text-center mb-4">{message}</p>
      {error && (
        <p className="text-center text-xs mb-4">
          If this keeps happening,{" "}
          <a
            href="https://github.com/mycorrhiza-inc/kessler/issues/new"
            target="_blank"
            className="link underline"
          >
            please open an issue on our GitHub repository{" "}
          </a>{" "}
          so we can fix it as soon as possible. Error details:{" "}
          <i>{error.slice(0, 1000)}</i>
        </p>
      )}
      {reload && (
        <button
          onClick={() => window.location.reload()}
          className="btn btn-error btn-outline"
        >
          Reload Page
        </button>
      )}
    </div>
  );
};

export default ErrorMessage;
