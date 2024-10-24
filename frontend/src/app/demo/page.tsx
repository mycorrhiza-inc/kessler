import { createClient } from "@/utils/supabase/server";
import { redirect } from "next/navigation";

import SearchApp from "@/components/SearchApp";
export default async function ProtectedPage() {
  const supabase = createClient();

  const {
    data: { user },
  } = await supabase.auth.getUser();

  if (user) {
    return redirect("/app");
  }

  return (
    <div className="w-full">
      <SearchApp user={user}></SearchApp>
    </div>
  );
}
