import Landing from "@/components/landing/Landing";

import { createClient } from "@/utils/supabase/server";
export default async function Page() {
  const supabase = createClient();
  const {
    data: { user },
  } = await supabase.auth.getUser();
  return (
    <div className="w-full">
      <Landing></Landing>
    </div>
  );
}
