import { PageContext, getStateDisplayName } from "@/lib/page_context";
import Link from "next/link";

export const ConversationHeader = ({ context }: { context: PageContext }) => {
  const displayState = getStateDisplayName(context.state);
  const proceeding_id = context.final_identifier;
  return (
    <div className="breadcrumbs text-xl">
      <ul>
        <li>
          <Link href="/">{displayState}</Link>
        </li>
        <li>
          <Link href="/proceedings">Proceedings</Link>
        </li>
        {proceeding_id && <li>{proceeding_id}</li>}
      </ul>
    </div>
  );
};
