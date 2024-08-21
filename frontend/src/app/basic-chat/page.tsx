import ChatUI from "../../lib/components/ChatUI";
import DefaultShell from "../../lib/components/DefaultShell";

export default function Page() {
  return (
    <DefaultShell>
      <ChatUI
        chatUrl="/api/rag/basic_chat"
        // LLama 405b isnt availible yet on groq, leaving it here for safety and so I can remember to enable it
        modelOptions={["llama-70b", "llama-405b", "gpt-4o", "llama-8b"]}
        useCitations={false}
      />
    </DefaultShell>
  );
}
