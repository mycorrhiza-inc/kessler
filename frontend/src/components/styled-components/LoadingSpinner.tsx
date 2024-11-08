const LoadingSpinner = ({ loadingText }: { loadingText?: string }) => {
  const text = loadingText || "Loading...";
  return (
    <div className="flex flex-col justify-center items-center">
      <span className="loading loading-infinity loading-lg text-primary"></span>
      <span className="loading loading-infinity loading-lg text-secondary"></span>
      <span className="loading loading-infinity loading-lg text-accent"></span>
      <br />
      <p>{text}</p>
    </div>
  );
};
export default LoadingSpinner;
