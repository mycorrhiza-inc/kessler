export type PageContext = {
  state?: string;
  slug: string[];
  final_identifier?: string;
};

const getStateDisplayName = (state?: string) => {
  if (state === "ny") {
    return "New York State";
  }
  return "Unknown";
};

const ConversationHeader = ({ context }: { context: PageContext }) => {
  const displayState = getStateDisplayName(context.state);
  const proceeding_id = context.final_identifier;
  return (
    <div className="breadcrumbs text-xl">
      <ul>
        <li>
          <a>{displayState}</a>
        </li>
        <li>
          <a>Proceedings</a>
        </li>
        {proceeding_id && <li>{proceeding_id}</li>}
      </ul>
    </div>
  );
};
