// we want to load the most recent conversations in the database

import ConversationComponent from "@/components/Conversations/ConversationComponent";
import { FilterField, InheritedFilterValues } from "@/lib/filters";
import { PageContext } from "@/lib/page_context";

export const ConversationView = ({
  pageContext,
}: {
  pageContext: PageContext;
}) => {
  const conversation_id = pageContext.final_identifier;
  const inheritedFilters: InheritedFilterValues = conversation_id
    ? [{ filter: FilterField.MatchDocketId, value: conversation_id }]
    : [];
  return (
    <>
      <div
        className="conversationContainer contents-center"
        style={{
          position: "relative",
          width: "99vw",
          height: "90vh",
          padding: "20px",
          // overflow: "scroll",
        }}
      >
        <ConversationComponent
          inheritedFilters={inheritedFilters}
          pageContext={pageContext}
        />
      </div>
    </>
  );
};
