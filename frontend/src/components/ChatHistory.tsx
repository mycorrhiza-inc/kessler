interface Message {
  role: string;
  content: string;
  key: symbol;
}
const chatHistory: Message[] = [
  {
    role: "user",
    content: "What is a black hole?",
    key: Symbol(),
  },
  {
    role: "assistant",
    content:
      "A black hole is a region of space where the gravitational pull is so strong that not even light can escape from it.",
    key: Symbol(),
  },
  {
    role: "user",
    content: "How are black holes formed?",
    key: Symbol(),
  },
  {
    role: "assistant",
    content:
      "Black holes are formed when massive stars collapse under their own gravity at the end of their life cycle.",
    key: Symbol(),
  },
  {
    role: "user",
    content: "Can black holes be seen?",
    key: Symbol(),
  },
  {
    role: "assistant",
    content:
      "No, black holes cannot be seen directly because their gravitational pull prevents light from escaping, but their presence can be inferred by observing the effects on nearby objects.",
    key: Symbol(),
  },
  {
    role: "user",
    content: "What would happen if you fell into a black hole?",
    key: Symbol(),
  },
  {
    role: "assistant",
    content:
      "If you fell into a black hole, you would experience extreme gravitational forces and time dilation. Ultimately, you would be stretched and compressed in a process known as spaghettification.",
    key: Symbol(),
  },
];

function MessageComponent({ message }: { message: Message }) {
  const isUser = message.role === "user";

  return (
    <div className={`flex w-full ${isUser ? "justify-end" : "justify-start"}`}>
      <div
        className={`w-11/12 rounded-lg overflow-auto min-h-[100px] p-5 ${
          isUser
            ? "bg-teal-100 dark:bg-teal-700"
            : "bg-gray-200 dark:bg-gray-700"
        }`}
      >
        <MarkdownRenderer>{message.content}</MarkdownRenderer>
      </div>
    </div>
  );
}

function AwaitingMessageSkeleton() {
  return (
    <div className="w-11/12 bg-gray-200 dark:bg-gray-700 rounded-lg min-h-[100px] p-5">
      <div className="animate-pulse">
        <div className="h-2 bg-pink-500 my-4 rounded"></div>
        <div className="h-2 bg-pink-500 my-4 rounded"></div>
        <div className="h-2 bg-pink-500 my-4 rounded"></div>
      </div>
    </div>
  );
}

export function ChatMessages({
  messages,
  loading,
}: {
  messages: Message[];
  loading: boolean;
}) {
  return (
    <div className="flex flex-col h-screen p-4 space-y-4 overflow-y-auto bg-white border divide-y divide-gray-200 rounded-md">
      {messages.length === 0 && (
        <div className="p-5 text-center text-gray-500">
          <h2 className="text-lg font-bold">Welcome to the Chatbot!</h2>
          <p>
            Type your message in the input box below and press Enter to send.
          </p>
        </div>
      )}
      {messages.map((m: Message) => (
        <MessageComponent message={m} />
      ))}
      {loading && <AwaitingMessageSkeleton />}
      <div className="w-full min-h-[300px] text-red-500"></div>
    </div>
  );
}
