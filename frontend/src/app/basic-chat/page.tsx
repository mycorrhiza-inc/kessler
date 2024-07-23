import ChatUI from "../../lib/components/ChatUI";
import DefaultShell from "../../lib/components/DefaultShell";

export default function Page() {
  return (
    <DefaultShell>
      <ChatUI
        chatUrl="/api/rag/basic_chat"
        modelOptions={["llama-70b", "llama-405b", "gpt-4o", "llama-8b"]}
      />
    </DefaultShell>
  );
}
