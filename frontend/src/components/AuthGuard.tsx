"use client";
import { useKesslerStore } from "@/lib/store";
import { useState } from "react";
import { useEffect } from "react";

export default function AuthGuard({
  isLoggedIn,
  children,
}: {
  isLoggedIn: boolean;
  children: React.ReactNode;
}) {
  const [loaded, setLoaded] = useState(false);
  const globalStore = useKesslerStore();

  useEffect(() => {
    globalStore.setIsLoggedIn(isLoggedIn);
    setLoaded(true);
  }, []);

  return <>{children}</>;
}
