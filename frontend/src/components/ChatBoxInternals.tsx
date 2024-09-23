import { CloseIcon, HamburgerIcon } from "@/components/Icons";
import { ChatMessages, exampleChatHistory } from "./ChatHistory";
import { useState } from "react";
import { Stack } from "@mui/joy";

interface Message {
  role: string;
  content: string;
  key: symbol;
}
const ChatBoxInternals = ({}: {}) => {
  const [messages, setMessages] = useState<Message[]>([]);
  const [needsResponse, setResponse] = useState(false);
  const [loadingResponse, setLoadingResponse] = useState(false);
  const [selectedModel, setSelectedModel] = useState("default");
  const [citations, setCitations] = useState<any[]>([]);
  const [ragMode, setRagMode] = useState(true);
  const [draftText, setDraftText] = useState("");
  const chatUrl = ragMode ? "/api/v1/rag/rag_chat" : "/api/v1/rag/basic_chat";

  const getResponse = async (responseText: string) => {
    const newMessage: Message = {
      role: "user",
      content: responseText,
      key: Symbol(),
    };
    console.log(newMessage);
    setMessages([...messages, newMessage]);
    let chat_hist = messages.map((m) => {
      let { key, ...rest } = m;
      return rest;
    });

    const modelToSend = selectedModel === "default" ? undefined : selectedModel;
    setLoadingResponse(true);

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
    console.log(chat_response);
    setMessages([...messages, chat_response]);
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
      console.log("Hit enter, with shift key down");
      e.preventDefault();
      handleSubmit();
    }
  };

  return (
    <form
      className="flex flex-col"
      onSubmit={handleSubmit}
      style={{ height: "85vh" }}
    >
      <div className="flex-none h-[5%] flex flex-row justify-center">
        <div className="dropdown dropdown-hover">
          <div tabIndex={0} role="button" className="btn m-1">
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
        <div className="form-control w-52">
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
      </div>
      <div className="flex-1 h-[70%] overflow-y-auto">
        <ChatMessages
          messages={messages}
          loading={loadingResponse}
        ></ChatMessages>
      </div>
      <div className="flex-none h-[25%]">
        <textarea
          name="userMessage"
          className="textarea textarea-accent w-full h-full"
          placeholder={`Type Here to Chat\nEnter to Send, Shift+Enter for New Line`}
          onKeyDown={handleKeyDown}
          value={draftText} // ...force the input's value to match the state variable...
          onChange={(e) => setDraftText(e.target.value)} // ... and update the state variable on any edits!
        ></textarea>
      </div>
    </form>
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
