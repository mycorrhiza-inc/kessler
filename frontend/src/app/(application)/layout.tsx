// import "@/app/globals.css";
import { CLIENT_API_URL } from "@/lib/env_variables";
import AuthGuard from "@/stateful_components/AuthGuard";
import Header from "@/style_components/Layout/Header";
import Layout from "@/style_components/Layout/Layout";

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
    if (CLIENT_API_URL == "http://localhost") {
      return true;
    }
    return true;
  };

  const isLoggedIn = await checkLoggedIn();
  return (
    <AuthGuard isLoggedIn={isLoggedIn}>
      <Header>
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
