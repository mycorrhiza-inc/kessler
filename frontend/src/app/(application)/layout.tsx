// import "@/app/globals.css";
import AuthGuard from "@/components/AuthGuard";

import { createClient } from "@/utils/supabase/server";

const defaultUrl = "https://kessler.xyz";
export const metadata = {
  metadataBase: new URL(defaultUrl),
  title: "Kessler",
  description: "Inteligence and Research tools for Lobbying",
};

export default async function ApplicationLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const checkLoggedIn = async () => {
    const supabase = createClient();
    const {
      data: { user },
    } = await supabase.auth.getUser();
    const userPresent = Boolean(user);
    return userPresent;
  };

  const isLoggedIn = await checkLoggedIn();
  return <AuthGuard isLoggedIn={isLoggedIn}>{children}</AuthGuard>;
}
