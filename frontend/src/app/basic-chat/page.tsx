import ChatUI from "../../lib/components/ChatUI";
import DefaultShell from "../../lib/components/DefaultShell";

export default function Page() {
  return (
    <DefaultShell>
      <ChatUI chatUrl="http://127.0.0.1/api/rag/basic_chat" />
    </DefaultShell>
  );
}
