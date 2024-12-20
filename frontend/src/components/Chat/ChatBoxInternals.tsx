import { useState } from "react";

import MarkdownRenderer from "../MarkdownRenderer";

import { QueryFilterFields } from "@/lib/filters";
import { getUpdatedChatHistory, Message } from "@/lib/chat";
import { publicAPIURL } from "@/lib/env_variables";
export const ChatMessages = ({
  messages,
  loading,
}: {
  messages: Message[];
  loading: boolean;
}) => {
  // const setMessageCitations = (index: number) => {
  //   setHighlighted(index);
  //   const message = messages[index];
  //   const isntUser = message.role != "user";
  //   const citationExists = message.citations && message.citations.length > 0;
  //   if (isntUser && citationExists) {
  //     setCitations(message.citations);
  //   }
  // };
  const MessageComponent = ({ message }: { message: Message }) => {
    const isUser = message.role === "user";
    const hasCitations = message.citations && message.citations.length > 0;
    if (hasCitations) {
      return (
        <div
          className={`flex w-full ${isUser ? "justify-end" : "justify-start"}`}
        >
          <div className="indicator">
            <span className="indicator-item badge badge-primary">
              This Message has Citations
            </span>
            <div
              className={`w-11/12 rounded-lg overflow-auto min-h-[100px] p-5 ${
                isUser ? "bg-success" : "bg-base-300"
              }`}
              // onClick={clickMessage}
            >
              <MarkdownRenderer
                color={isUser ? "success-content" : "base-content"}
              >
                {message.content}
              </MarkdownRenderer>
            </div>
          </div>
        </div>
      );
    }
    return (
      <div
        className={`flex w-full ${isUser ? "justify-end" : "justify-start"}`}
      >
        <div
          className={`w-11/12 rounded-lg overflow-auto min-h-[100px] p-5 ${
            isUser ? "bg-success" : "bg-base-300"
          }`}
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
            // clickMessage={() => setMessageCitations(index)}
            // highlighted={highlighted === index}
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

export interface ChatBoxInternalsState {
  highlighted: number;
  messages: Message[];
  loadingResponse: boolean;
  selectedModel: string;
  ragMode: boolean;
  draftText: string;
}

export const initialChatState: ChatBoxInternalsState = {
  highlighted: -1,
  messages: [],
  loadingResponse: false,
  selectedModel: "default",
  ragMode: true,
  draftText: "",
};

export const ChatBoxInternals = ({
  setCitations,
  ragFilters,
}: {
  setCitations: React.Dispatch<React.SetStateAction<any[]>>;
  ragFilters: QueryFilterFields;
}) => {
  const [chatState, setChatState] =
    useState<ChatBoxInternalsState>(initialChatState);
  return (
    <ChatBoxInternalsStateless
      setCitations={setCitations}
      ragFilters={ragFilters}
      chatState={chatState}
      setChatState={setChatState}
    />
  );
};

export const ChatBoxInternalsStateless = ({
  setCitations,
  ragFilters,
  chatState,
  setChatState,
}: {
  setCitations: React.Dispatch<React.SetStateAction<any[]>>;
  ragFilters: QueryFilterFields;
  chatState: ChatBoxInternalsState;
  setChatState: React.Dispatch<React.SetStateAction<ChatBoxInternalsState>>;
}) => {
  const { messages, loadingResponse, selectedModel, ragMode, draftText } =
    chatState;

  const setHighlighted = (value: number) => {
    setChatState((prev) => ({ ...prev, highlighted: value }));
  };
  const setMessages = (value: Message[]) =>
    setChatState((prev) => ({ ...prev, messages: value }));
  const setLoadingResponse = (value: boolean) =>
    setChatState((prev) => ({ ...prev, loadingResponse: value }));
  const setSelectedModel = (value: string) =>
    setChatState((prev) => ({ ...prev, selectedModel: value }));
  const setRagMode = (value: boolean) =>
    setChatState((prev) => ({ ...prev, ragMode: value }));
  const setDraftText = (value: string) =>
    setChatState((prev) => ({ ...prev, draftText: value }));
  const chatUrl = ragMode
    ? `${publicAPIURL}/v2/rag/chat`
    : `${publicAPIURL}/v2/rag/basic_chat`;

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

    const modelToSend = selectedModel === "default" ? undefined : selectedModel;
    setLoadingResponse(true);

    // Should this fetch get refactored out into lib as something that calls a chat endpoint?

    newMessages = await getUpdatedChatHistory(
      newMessages,
      ragFilters,
      chatUrl,
      modelToSend,
    );
    setLoadingResponse(false);
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
        <ChatMessages messages={messages} loading={loadingResponse} />
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
