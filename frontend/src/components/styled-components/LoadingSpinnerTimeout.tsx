import { useEffect, useState } from "react";
import LoadingSpinner from "./LoadingSpinner";
const LoadingSpinnerTimeout = ({
  loadingText,
  timeoutSeconds,
  replacement,
}: {
  loadingText?: string;
  timeoutSeconds: number;
  replacement?: React.ReactElement;
}) => {
  const [showSpinner, setShowSpinner] = useState(true);

  useEffect(() => {
    const timer = setTimeout(() => {
      setShowSpinner(false);
    }, timeoutSeconds * 1000);

    return () => clearTimeout(timer);
  }, [timeoutSeconds]);

  if (!showSpinner) {
    return replacement || null;
  }

  return <LoadingSpinner loadingText={loadingText} />;
};

export default LoadingSpinnerTimeout;
