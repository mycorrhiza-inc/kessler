import { PageContext, getStateDisplayName } from "@/lib/page_context";

export const ConversationHeader = ({ context }: { context: PageContext }) => {
  const displayState = getStateDisplayName(context.state);
  const proceeding_id = context.final_identifier;
  return (
    <div className="breadcrumbs text-xl">
      <ul>
        <li>
          <Link{displayState}
        </li>
        <li>
          <a>Proceedings</a>
        </li>
        {proceeding_id && <li>{proceeding_id}</li>}
      </ul>
    </div>
  );
};
