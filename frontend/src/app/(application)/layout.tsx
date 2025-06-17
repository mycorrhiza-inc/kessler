// import "@/app/globals.css";
import { CLIENT_API_URL } from "@/lib/env_variables";
import AuthGuard from "@/components/stateful/tracking/AuthGuard";
import Header from "@/components/style/Layout/Header";
import Layout from "@/components/style/Layout/Layout";


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
        <div className="relative flex gap-3">
          {/* <SignedIn> */}
          {/*   <Link */}
          {/*     href="/dashboard" */}
          {/*     className="px-4 py-2 rounded-full bg-[#131316] text-white text-sm font-semibold" */}
          {/*   > */}
          {/*     Dashboard */}
          {/*   </Link> */}
          {/* </SignedIn> */}
          {/* <SignedOut> */}
          {/*   <SignUpButton /> */}
          {/*   <SignInButton /> */}
          {/* </SignedOut> */}
        </div>
      </Header>
      <Layout>{children}</Layout>
    </AuthGuard>
  );
}
