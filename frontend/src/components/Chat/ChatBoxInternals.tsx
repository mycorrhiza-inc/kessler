// import { useState } from "react";
//
// import MarkdownRenderer from "../MarkdownRenderer";
//
// import { getUpdatedChatHistory, Message } from "@/lib/chat";
// import clsx from "clsx";
// import { getClientRuntimeEnv } from "@/lib/env_variables/env_variables_hydration_script";
// import { GenericSearchInfo, GenericSearchType } from "@/lib/adapters/genericSearchCallback";
// export const ChatMessages = ({
//   messages,
//   loading,
// }: {
//   messages: Message[];
//   loading: boolean;
// }) => {
//   const MessageComponent = ({ message }: { message: Message }) => {
//     const isUser = message.role === "user";
//     const hasCitations = message.citations && message.citations.length > 0;
//     const clickCitation = () => {
//       console.log("Yay you clicked the citation component:", message.citations);
//     };
//     const showCitationsButton = !isUser && hasCitations;
//     return (
//       <div
//         className={clsx(
//           "flex w-full",
//           isUser ? "justify-end" : "justify-start",
//         )}
//       >
//         <div className="indicator w-4/5">
//           {showCitationsButton && (
//             <button
//               className="indicator-item badge badge-primary"
//               onClick={clickCitation}
//             >
//               <div
//                 className="tooltip"
//                 data-tip={message.citations.map((citation) => String(citation))}
//               >
//                 View Citations
//               </div>
//             </button>
//           )}
//           <div
//             className={`w-full rounded-lg overflow-auto min-h-[100px] p-5 ${isUser ? "bg-success" : "bg-base-300"
//               }`}
//           >
//             <MarkdownRenderer
//               color={isUser ? "success-content" : "base-content"}
//             >
//               {message.content}
//             </MarkdownRenderer>
//           </div>
//         </div>
//       </div>
//     );
//   };
//   return (
//     <div className="flex flex-col gap-8 h-[85%] overflow-y-auto">
//       {messages.length === 0 && (
//         <div className="p-5 text-center text-base-content">
//           <h2 className="text-lg font-bold">Welcome to the Chatbot!</h2>
//           <p>
//             Type your message in the input box below and press Enter to send.
//           </p>
//         </div>
//       )}
//       {messages.map((m: Message, index: number) => {
//         return <MessageComponent message={m} />;
//       })}
//       {loading && (
//         <div className="w-11/12 bg-base-300 rounded-lg min-h-[100px] p-5">
//           <div className="animate-pulse">
//             <div className="h-2 bg-accent my-4 rounded-sm"></div>
//             <div className="h-2 bg-accent my-4 rounded-sm"></div>
//             <div className="h-2 bg-accent my-4 rounded-sm"></div>
//           </div>
//         </div>
//       )}
//     </div>
//   );
// };
//
// export interface ChatBoxInternalsState {
//   highlighted: number;
//   messages: Message[];
//   loadingResponse: boolean;
//   selectedModel: string;
//   ragMode: boolean;
//   draftText: string;
// }
//
// export const initialChatState: ChatBoxInternalsState = {
//   highlighted: -1,
//   messages: [],
//   loadingResponse: false,
//   selectedModel: "default",
//   ragMode: true,
//   draftText: "",
// };
//
// export const ChatBoxInternals = ({
//   setCitations,
//   ragFilters,
// }: {
//   setCitations: React.Dispatch<React.SetStateAction<any[]>>;
//   ragFilters: QueryFileFilterFields;
// }) => {
//   const [chatState, setChatState] =
//     useState<ChatBoxInternalsState>(initialChatState);
//   const initial_info: GenericSearchInfo = {
//     search_type: GenericSearchType.Filling,
//     query: "",
//     filters: []
//   }
//   const [searchInfo, setSearchInfo] = useState<GenericSearchInfo>(initial_info)
//   return (
//     <ChatBoxInternalsStateless
//       setCitations={setCitations}
//       ragFilters={ragFilters}
//       chatState={chatState}
//       setChatState={setChatState}
//     />
//   );
// };
//
// export const ChatBoxInternalsStateless = ({
//   setCitations,
//   ragFilters,
//   chatState,
//   setChatState,
// }: {
//   setCitations: React.Dispatch<React.SetStateAction<any[]>>;
//   ragFilters: QueryFileFilterFields;
//   chatState: ChatBoxInternalsState;
//   setChatState: React.Dispatch<React.SetStateAction<ChatBoxInternalsState>>;
// }) => {
//   const { messages, loadingResponse, selectedModel, ragMode, draftText } =
//     chatState;
//
//   // const setHighlighted = (value: number) => {
//   //   setChatState((prev) => ({ ...prev, highlighted: value }));
//   // };
//   const setMessages = (value: Message[]) =>
//     setChatState((prev) => ({ ...prev, messages: value }));
//   const setLoadingResponse = (value: boolean) =>
//     setChatState((prev) => ({ ...prev, loadingResponse: value }));
//   const setSelectedModel = (value: string) =>
//     setChatState((prev) => ({ ...prev, selectedModel: value }));
//   const setRagMode = (value: boolean) =>
//     setChatState((prev) => ({ ...prev, ragMode: value }));
//   const setDraftText = (value: string) =>
//     setChatState((prev) => ({ ...prev, draftText: value }));
//   const runtimeConfig = getClientRuntimeEnv();
//   const chatUrl = ragMode
//     ? `${runtimeConfig.public_api_url}/v2/rag/chat`
//     : `${runtimeConfig.public_api_url}/v2/rag/basic_chat`;
//
//   const getResponse = async (responseText: string) => {
//     if (responseText == "") {
//       return;
//     }
//     const newMessage: Message = {
//       role: "user",
//       key: Symbol(),
//       content: responseText,
//       citations: [],
//     };
//     console.log(newMessage);
//     var newMessages = [...messages, newMessage];
//     console.log(newMessages);
//     setMessages(newMessages);
//
//     const modelToSend = selectedModel === "default" ? undefined : selectedModel;
//     setLoadingResponse(true);
//
//     // Should this fetch get refactored out into lib as something that calls a chat endpoint?
//
//     newMessages = await getUpdatedChatHistory(
//       newMessages,
//       ragFilters,
//       chatUrl,
//       modelToSend,
//     );
//     setLoadingResponse(false);
//     setMessages(newMessages);
//   };
//   const model_list = [
//     "default",
//   ];
//
//   // This isnt working fix problem with type
//   const handleSubmit = async () => {
//     // Check if e.target is a form element
//     const chatText = `${draftText}`;
//     setDraftText("");
//     await getResponse(chatText);
//   };
//
//   // More BS
//   const handleModelSelect = (model: any) => {
//     setSelectedModel(model);
//   };
//
//   const handleRagModeToggle = (e: any) => {
//     setRagMode(e.target.checked);
//   };
//
//   const handleKeyDown = (e: any) => {
//     if (e.key === "Enter" && !e.shiftKey) {
//       console.log("Hit enter, without shift key down");
//       e.preventDefault();
//       handleSubmit();
//     }
//   };
//
//   const createNewChat = () => {
//     setMessages([]);
//     setCitations([]);
//   };
//
//   return (
//     <div className="flex flex-col" style={{ height: "85vh" }}>
//       <div className="flex-none flex flex-row justify-center bg-base-100 text-base-content gap-11">
//         <div className="dropdown dropdown-hover">
//           <div tabIndex={0} role="button" className="btn m-1 bg-base-300">
//             {selectedModel !== "default" ? selectedModel : "Select Model"}
//           </div>
//           <ul
//             tabIndex={0}
//             className="dropdown-content menu bg-base-100 rounded-box z-1 w-52 p-2 shadow-sm"
//           >
//             {model_list.map((model) => (
//               <li key={model}>
//                 <a onClick={() => handleModelSelect(model)}>{model}</a>
//               </li>
//             ))}
//           </ul>
//         </div>
//         <div className="space-y-1">
//           <label className="label cursor-pointer flex flex-col">
//             <input
//               type="checkbox"
//               className="toggle toggle-accent"
//               checked={ragMode}
//               onChange={handleRagModeToggle}
//             />
//             <span className="label-text">Rag Mode</span>
//           </label>
//         </div>
//         <button className="btn btn-primary" onClick={createNewChat}>
//           New Chat
//         </button>
//       </div>
//       <ChatMessages messages={messages} loading={loadingResponse} />
//       <div className="flex-none h-[15%]">
//         <textarea
//           name="userMessage"
//           className="textarea textarea-accent w-full h-full bg-base-300"
//           placeholder={`Type Here to Chat\nEnter to Send, Shift+Enter for New Line`}
//           onKeyDown={handleKeyDown}
//           value={draftText} // ...force the input's value to match the state variable...
//           onChange={(e) => setDraftText(e.target.value)} // ... and update the state variable on any edits!
//           disabled={loadingResponse}
//         ></textarea>
//       </div>
//     </div>
//   );
// };
// export default ChatBoxInternals;
