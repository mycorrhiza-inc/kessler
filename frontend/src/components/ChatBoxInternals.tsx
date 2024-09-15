import { CloseIcon, HamburgerIcon } from "@/components/Icons";
import { ChatMessages, exampleChatHistory } from "./ChatHistory";

import { Stack } from "@mui/joy";

const ChatBoxInternals = ({
  chatSidebarVisible,
  setChatSidebarVisible,
}: {
  chatSidebarVisible: boolean;
  setChatSidebarVisible: React.Dispatch<React.SetStateAction<boolean>>;
}) => {
  return (
    <>
      <div className="overflow-y-auto">
        <ChatMessages
          messages={exampleChatHistory}
          loading={false}
        ></ChatMessages>
      </div>
      <textarea
        className="textarea textarea-bordered"
        placeholder="Bio"
      ></textarea>
      <div
        className="chatContainer"
        style={{
          display: "flex",
          flexDirection: "row",
          height: "90%",
          width: "100%",
          padding: "2px",
        }}
      >
        <div
          className="chatSidebar"
          style={{
            width: chatSidebarVisible ? "20%" : "0%",
            backgroundColor: "red",
            overflow: "scroll",
          }}
        >
          sidebar contents
        </div>
        <div> chat contents</div>
      </div>
    </>
  );
};
export default ChatBoxInternals;
// <div
//   className="chatbox-banner"
//   style={{
//     position: "sticky",
//     top: "0",
//     backgroundColor: "#f1f1f1",
//     padding: "20px",
//     textAlign: "center",
//     zIndex: "1000",
//     borderBottom: "1px solid #ccc",
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
// <div className="chatbox-banner sticky top-0 bg-[#f5f5f5] dark:bg-gray-700 p-5 text-center z-50 border-b border-gray-300 h-auto">
//   <div className="flex flex-row justify-between">
//     <button>
//       <HamburgerIcon />
//     </button>
//     <button
//       onClick={() => {
//         // This is fixable, dont want to worry about it now, if encounter bugs remove and actually solve.
//         // @ts-ignore
//         setChatVisible((prev) => !prev);
//       }}
//     >
//       <CloseIcon />
//     </button>
//   </div>
// </div>
