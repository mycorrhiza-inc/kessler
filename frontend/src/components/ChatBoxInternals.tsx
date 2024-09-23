import { CloseIcon, HamburgerIcon } from "@/components/Icons";
import { ChatMessages, exampleChatHistory } from "./ChatHistory";
import { useState } from "react";

interface Message {
  role: string;
  content: string;
  key: symbol;
}

interface ChatBoxInternalsProps {
  setCitations: (citations: any[]) => void;
}

const ChatBoxInternals = ({ setCitations }: ChatBoxInternalsProps) => {
  const [messages, setMessages] = useState<Message[]>([]);
  const [needsResponse, setResponse] = useState(false);
  const [loadingResponse, setLoadingResponse] = useState(false);
  const [selectedModel, setSelectedModel] = useState("default");
  const [ragMode, setRagMode] = useState(true);
  const [draftText, setDraftText] = useState("");
  const chatUrl = ragMode ? "/api/v1/rag/rag_chat" : "/api/v1/rag/basic_chat";

  const getResponse = async (responseText: string) => {
    if (responseText == "") {
      return;
    }
    const newMessage: Message = {
      role: "user",
      key: Symbol(),
      content: responseText,
    };
    var newMessages = [...messages, newMessage];
    setMessages(newMessages);

    let chat_hist = newMessages.map((m) => {
      let { key, ...rest } = m;
      return rest;
    });

    const modelToSend = selectedModel === "default" ? undefined : selectedModel;
    setLoadingResponse(true);

    // Should this fetch get refactored out into lib as something that calls a chat endpoint?
    let result = await fetch(chatUrl, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        accept: "application/json",
      },
      body: JSON.stringify({
        model: modelToSend,
        chat_history: chat_hist,
      }),
    })
      .then((resp) => {
        setLoadingResponse(false);
        if (resp.status < 200 || resp.status > 299) {
          return "failed request";
        }
        return resp.json();
      })
      .then((data) => {
        setCitations(data.citations || []);
        return data;
      })
      .catch((e) => {
        console.log("error making request");
        console.log(JSON.stringify(e));
      });
    const chat_response: Message = {
      role: "assistant",
      key: Symbol(),
      content: result == "failed request" ? result : result.message.content,
    };
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
            Select Model
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
        ></ChatMessages>
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
// Added complexity, removing for a bit.
// <div
//   className="chatbox-banner"
//   style={{
//     position: "sticky",
//     top: "0",
//     padding: "20px",
//     textAlign: "center",
//     zIndex: "1000",
//     borderBottom: "1px solid ",
//     height: "auto",
//     pointerEvents: "auto",
//   }}
// >
//   <Stack direction="row" justifyContent="space-between">
//     <button
//       onClick={() => setChatSidebarVisible((prevState) => !prevState)}
//     ></button>
//   </Stack>
// </div>
// <div className="chatbox-banner sticky top-0 p-5 text-center z-50 border-b border-accent h-auto">
//   <div className="flex flex-row justify-between">
//     <button>
//       <HamburgerIcon />
//     </button>
//     <button
//       onClick={() => {
//         setChatVisible((prev) => !prev);
//       }}
//     >
//       <CloseIcon />
//     </button>
//   </div>
// </div>
//
// <div
//   className="chatContainer"
//   style={{
//     display: "flex",
//     flexDirection: "row",
//     height: "90%",
//     width: "100%",
//     padding: "2px",
//   }}
// >
//   <div
//     className="chatSidebar"
//     style={{
//       width: chatSidebarVisible ? "20%" : "0%",
//       backgroundColor: "red",
//       overflow: "scroll",
//     }}
//   >
//     sidebar contents
//   </div>
//   <div> chat contents</div>
// </div>
