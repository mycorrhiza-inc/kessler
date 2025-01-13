// import "@/app/globals.css";
import AuthGuard from "@/components/AuthGuard";
import { runtimeConfig } from "@/lib/env_variables";

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
    // Make the apps and pages load when you are running it locally even if offline, otherwise it errors out
    // when trying to connect to supabase.
    if (runtimeConfig.public_api_url == "http://localhost") {
      return true;
    }
    try {
      const supabase = createClient();
      const {
        data: { user },
      } = await supabase.auth.getUser();
      const userPresent = Boolean(user);
      return userPresent;
    } catch {
      return false;
    }
  };

  const isLoggedIn = await checkLoggedIn();
  return <AuthGuard isLoggedIn={isLoggedIn}>{children}</AuthGuard>;
}
