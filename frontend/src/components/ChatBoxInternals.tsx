import { CloseIcon, HamburgerIcon } from "@/components/Icons";
import { ChatMessages, exampleChatHistory } from "./ChatHistory";
import { useState } from "react";
import { Stack } from "@mui/joy";

interface Message {
  role: string;
  content: string;
  key: symbol;
}
const ChatBoxInternals = ({
  chatSidebarVisible,
  setChatSidebarVisible,
}: {
  chatSidebarVisible: boolean;
  setChatSidebarVisible: React.Dispatch<React.SetStateAction<boolean>>;
}) => {
  const [messages, setMessages] = useState<Message[]>([]);
  const [needsResponse, setResponse] = useState(false);
  const [loadingResponse, setLoadingResponse] = useState(false);
  const [selectedModel, setSelectedModel] = useState("default");
  const [citations, setCitations] = useState<any[]>([]);
  const [isRag, setIsRag] = useState(true);
  const chatUrl = isRag ? "/api/v1/rag/rag_chat" : "/api/v1/rag/basic_chat";

  const getResponse = async () => {
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

    setMessages([...messages, result.message]);
  };
  const model_list = ["llama-405b", "gpt-4o", "llama-70b", "llama-8b"];
  return (
    <>
      <div className="w-full flex flex-row justify-center">
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
                <a>{model}</a>
              </li>
            ))}
          </ul>
        </div>
        <div className="form-control w-52">
          <label className="label cursor-pointer flex flex-col">
            <input
              type="checkbox"
              className="toggle toggle-accent"
              defaultChecked
            />
            <span className="label-text">Rag Mode</span>
          </label>
        </div>
      </div>
      <div>
        <ChatMessages
          messages={exampleChatHistory}
          loading={false}
        ></ChatMessages>
      </div>
      <textarea
        className="textarea textarea-accent w-full"
        placeholder={`Type Here to Chat\nEnter to Send, Shift+Enter for New Line`}
      ></textarea>
    </>
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
