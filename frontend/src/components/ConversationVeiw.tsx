

// we want to load the most recent conversations in the database

import ConversationComponent from "@/components/Conversations/ConversationComponent";

export const ConversationView = () => {
	return <>
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
      <ConversationComponent battle={undefined} childBattles={[]} />
	  </div>
	</>
}