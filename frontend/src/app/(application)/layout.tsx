// import "@/app/globals.css";
import AuthGuard from "@/components/AuthGuard";
import Header from "@/components/Layout/Header";
import Layout from "@/components/Layout/Layout";
import { BreadcrumbValues } from "@/components/SitemapUtils";
import { getUniversalEnvConfig } from "@/lib/env_variables/env_variables";

import { createClient } from "@/utils/supabase/server";
import { SignedIn, SignedOut, SignInButton, SignUp, SignUpButton } from "@clerk/nextjs";
import Link from "next/link";

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
    if (getUniversalEnvConfig().public_api_url == "http://localhost") {
      return true;
    }
    return true;
  };

  const isLoggedIn = await checkLoggedIn();
  return (
    <AuthGuard isLoggedIn={isLoggedIn}>
      <Header
        breadcrumbs={{ state: "ny", breadcrumbs: [] } as BreadcrumbValues}
      >
        {/* <div className="relative flex gap-3"> */}
        <SignedIn>
          <Link
            href="/dashboard"
            className="px-4 py-2 rounded-full bg-[#131316] text-white text-sm font-semibold"
          >
            Dashboard
          </Link>
        </SignedIn>
        <SignedOut>
          <SignUpButton />
          {/* <button className="px-4 py-2 rounded-full bg-[#131316] text-white text-sm font-semibold"> */}
          {/*   Sign up */}
          {/* </button> */}
          <SignInButton />
          {/*   <button className="px-4 py-2 rounded-full bg-[#131316] text-white text-sm font-semibold"> */}
          {/*     Sign in */}
          {/*   </button> */}
          {/* </SignInButton> */}
        </SignedOut>
        {/* </div> */}
      </Header>
      <Layout>{children}</Layout>
    </AuthGuard>
  );
}
