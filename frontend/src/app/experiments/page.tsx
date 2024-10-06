import TestFlowVisuals from "@/components/SeriousGaming/SeriousGameAdmin";

import { createClient } from "@/utils/supabase/server";
export default async function Page() {
  const supabase = createClient();
  const {
    data: { user },
  } = await supabase.auth.getUser();
  return <TestFlowVisuals user={user}></TestFlowVisuals>;
}
