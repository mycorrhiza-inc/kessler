export const LoadingSpinner = ({ text }: { text?: string }) => {
  text = text || "Loading...";
  return (
    <div className="">
      <span className="loading loading-spinner text-primary"></span>
      <span className="loading loading-spinner text-secondary"></span>
      <span className="loading loading-spinner text-accent"></span>
      <span className="loading loading-spinner text-neutral"></span>
      <br />
      <p>{text}</p>
    </div>
  );
};
