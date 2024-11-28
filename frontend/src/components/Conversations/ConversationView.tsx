// we want to load the most recent conversations in the database
import { use, useEffect, useState } from "react";

import axios from "axios";

import { Conversation } from "@/lib/conversations";
import ConversationComponent from "@/components/Conversations/ConversationComponent";
import { NYConversationDescription } from "@/components/Conversations/ConversationDescription";
import { GetConversationInformation } from "@/lib/requests/conversations";

import { FilterField, InheritedFilterValues } from "@/lib/filters";
import { PageContext } from "@/lib/page_context";

export const ConversationView = ({
  pageContext,
}: {
  pageContext: PageContext;
}) => {
  const conversation_id = pageContext.final_identifier || "bobloblawslawblog";
  const inheritedFilters: InheritedFilterValues = conversation_id
    ? [{ filter: FilterField.MatchDocketId, value: conversation_id }]
    : [];

  // useEffect(() => {
  //   if (conversation_id) {
  //     GetConversationInformation(conversation_id).then((data) => {
  //       setConversation(data);
  //     });
  //   }
  // }, []);

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
        <NYConversationDescription docket_id={conversation_id} />
        <ConversationComponent
          inheritedFilters={inheritedFilters}
          pageContext={pageContext}
        />
      </div>
    </>
  );
};
