import { PageContext, getStateDisplayName } from "@/lib/page_context";
import Link from "next/link";

export const ConversationHeader = ({ context }: { context: PageContext }) => {
  const displayState = getStateDisplayName(context.state);
  const docket_id = context.final_identifier;
  return (
    <div className="breadcrumbs text-xl">
      <ul>
        <li>
          <Link href="/">{displayState}</Link>
        </li>
        <li>
          <Link href="/dockets">Dockets</Link>
        </li>
        {docket_id && <li>{docket_id}</li>}
      </ul>
    </div>
  );
};
