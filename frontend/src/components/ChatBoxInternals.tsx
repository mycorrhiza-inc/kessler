import { useState } from "react";

import MarkdownRenderer from "./MarkdownRenderer";

import { extraProperties } from "@/utils/interfaces";
import { exampleChatHistory, Message } from "@/lib/chat";
export const ChatMessages = ({
  messages,
  loading,
  setCitations,
  highlighted,
  setHighlighted,
}: {
  messages: Message[];
  loading: boolean;
  setCitations: (citations: any[]) => void;
  highlighted: number;
  setHighlighted: (index: number) => void;
}) => {
  const setMessageCitations = (index: number) => {
    setHighlighted(index);
    const message = messages[index];
    const isntUser = message.role != "user";
    const citationExists = message.citations && message.citations.length > 0;
    if (isntUser && citationExists) {
      setCitations(message.citations);
    }
  };
  const MessageComponent = ({
    message,
    clickMessage,
    highlighted,
  }: {
    message: Message;
    clickMessage: any; // This makes me sad
    highlighted: boolean;
  }) => {
    const isUser = message.role === "user";
    return (
      <div
        className={`flex w-full ${isUser ? "justify-end" : "justify-start"}`}
      >
        <div
          className={`w-11/12 rounded-lg overflow-auto min-h-[100px] p-5 ${
            isUser ? "bg-success" : "bg-base-300"
          } ${highlighted ? "highlighted" : "not-highlighted"}`}
          onClick={clickMessage}
        >
          <MarkdownRenderer color={isUser ? "success-content" : "base-content"}>
            {message.content}
          </MarkdownRenderer>
        </div>
      </div>
    );
  };
  return (
    <>
      {messages.length === 0 && (
        <div className="p-5 text-center text-base-content">
          <h2 className="text-lg font-bold">Welcome to the Chatbot!</h2>
          <p>
            Type your message in the input box below and press Enter to send.
          </p>
        </div>
      )}
      {messages.map((m: Message, index: number) => {
        return (
          <MessageComponent
            message={m}
            clickMessage={() => setMessageCitations(index)}
            highlighted={highlighted === index}
          />
        );
      })}
      {loading && (
        <div className="w-11/12 bg-base-300 rounded-lg min-h-[100px] p-5">
          <div className="animate-pulse">
            <div className="h-2 bg-accent my-4 rounded"></div>
            <div className="h-2 bg-accent my-4 rounded"></div>
            <div className="h-2 bg-accent my-4 rounded"></div>
          </div>
        </div>
      )}
    </>
  );
};
interface ChatBoxInternalsProps {
  setCitations: (citations: any[]) => void;
  ragFilters: extraProperties;
}

const ChatBoxInternals = ({
  setCitations,
  ragFilters,
}: ChatBoxInternalsProps) => {
  const [highlighted, setHighlighted] = useState<number>(-1); // -1 Means no message is highlighted
  const [messages, setMessages] = useState<Message[]>([]);
  const [loadingResponse, setLoadingResponse] = useState(false);
  const [selectedModel, setSelectedModel] = useState("default");
  const [ragMode, setRagMode] = useState(true);
  const [draftText, setDraftText] = useState("");
  const chatUrl = ragMode ? "/api/v2/rag/chat" : "/api/v2/rag/basic_chat";

  const getResponse = async (responseText: string) => {
    if (responseText == "") {
      return;
    }
    const newMessage: Message = {
      role: "user",
      key: Symbol(),
      content: responseText,
      citations: [],
    };
    console.log(newMessage);
    var newMessages = [...messages, newMessage];
    console.log(newMessages);
    setMessages(newMessages);

    let chat_hist = newMessages.map((m) => {
      let { key, ...rest } = m;
      return rest;
    });

    const modelToSend = selectedModel === "default" ? undefined : selectedModel;
    setLoadingResponse(true);

    // Should this fetch get refactored out into lib as something that calls a chat endpoint?
    let result_message = await fetch(chatUrl, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        accept: "application/json",
      },
      body: JSON.stringify({
        model: modelToSend,
        chat_history: chat_hist,
        filters: ragFilters,
      }),
    })
      .then((resp) => {
        setLoadingResponse(false);
        if (resp.status < 200 || resp.status > 299) {
          console.log("failed request with status " + resp.status);
          console.log(resp);
          return "failed request";
        }
        return resp.json();
      })
      .then((data) => {
        if (!data.message) {
          console.log("no message in data");
          console.log(data);
          console.log(data["message"]);
          console.log(data.message);
          return "failed request";
        }
        console.log("got data");
        console.log(data);
        if (data.message.citations) {
          console.log("got citations");
          setHighlighted(newMessages.length); // You arent subtracting one here, since you want it to highlight the last message added to the list.

          console.log("set highlighted message");
          setCitations(data.message.citations);
          console.log("set citations");
        }
        console.log("Returning Message:");
        console.log(data.message);
        return data.message;
      })
      .catch((e) => {
        console.log("error making request");
        console.log(JSON.stringify(e));
        return "encountered exception while fetching data";
      });
    let chat_response: Message;

    if (typeof result_message === "string") {
      chat_response = {
        role: "assistant",
        key: Symbol(),
        content: result_message,
        citations: [],
      };
    } else {
      chat_response = {
        role: "assistant",
        key: Symbol(),
        content: result_message.content,
        citations: result_message.citations,
      };
    }

    newMessages = [...newMessages, chat_response];
    setMessages(newMessages);
  };
  const model_list = [
    "default",
    "llama-405b",
    "gpt-4o",
    "llama-70b",
    "llama-8b",
  ];

  // This isnt working fix problem with type
  const handleSubmit = async () => {
    // Check if e.target is a form element
    const chatText = `${draftText}`;
    setDraftText("");
    await getResponse(chatText);
  };

  // More BS
  const handleModelSelect = (model: any) => {
    setSelectedModel(model);
  };

  const handleRagModeToggle = (e: any) => {
    setRagMode(e.target.checked);
  };

  const handleKeyDown = (e: any) => {
    if (e.key === "Enter" && !e.shiftKey) {
      console.log("Hit enter, without shift key down");
      e.preventDefault();
      handleSubmit();
    }
  };

  const createNewChat = () => {
    setMessages([]);
    setCitations([]);
  };

  return (
    <div className="flex flex-col" style={{ height: "85vh" }}>
      <div className="flex-none flex flex-row justify-center bg-base-100 text-base-content gap-11">
        <div className="dropdown dropdown-hover">
          <div tabIndex={0} role="button" className="btn m-1 bg-base-300">
            {selectedModel !== "default" ? selectedModel : "Select Model"}
          </div>
          <ul
            tabIndex={0}
            className="dropdown-content menu bg-base-100 rounded-box z-[1] w-52 p-2 shadow"
          >
            {model_list.map((model) => (
              <li key={model}>
                <a onClick={() => handleModelSelect(model)}>{model}</a>
              </li>
            ))}
          </ul>
        </div>
        <div className="space-y-1">
          <label className="label cursor-pointer flex flex-col">
            <input
              type="checkbox"
              className="toggle toggle-accent"
              checked={ragMode}
              onChange={handleRagModeToggle}
            />
            <span className="label-text">Rag Mode</span>
          </label>
        </div>
        <button className="btn btn-primary" onClick={createNewChat}>
          New Chat
        </button>
      </div>
      <div className="flex-1 h-[85%] overflow-y-auto">
        <ChatMessages
          messages={messages}
          loading={loadingResponse}
          setCitations={setCitations}
          highlighted={highlighted}
          setHighlighted={setHighlighted}
        />
      </div>
      <div className="flex-none h-[15%]">
        <textarea
          name="userMessage"
          className="textarea textarea-accent w-full h-full bg-base-300"
          placeholder={`Type Here to Chat\nEnter to Send, Shift+Enter for New Line`}
          onKeyDown={handleKeyDown}
          value={draftText} // ...force the input's value to match the state variable...
          onChange={(e) => setDraftText(e.target.value)} // ... and update the state variable on any edits!
          disabled={loadingResponse}
        ></textarea>
      </div>
    </div>
  );
};
export default ChatBoxInternals;
