// we want to load the most recent conversations in the database

import ConversationComponent from "@/components/Conversations/ConversationComponent";
import { User } from "@supabase/supabase-js";
import { FilterField, InheritedFilterValues } from "@/lib/filters";
import Navbar from "./Navbar";

export const ConversationView = ({
  conversation_id,
  user,
}: {
  conversation_id?: string;
  user: User | null;
}) => {
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
          overflow: "scroll",
        }}
      >
        <ConversationComponent inheritedFilters={inheritedFilters} />
      </div>
    </>
  );
};
