import { CloseIcon, HamburgerIcon } from "@/components/Icons";
import { ChatMessages, exampleChatHistory } from "./ChatHistory";

import { Stack } from "@mui/joy";

const ChatBoxInternals = ({
  chatVisible,
  setChatVisible,
  chatSidebarVisible,
  setChatSidebarVisible,
}: {
  chatSidebarVisible: boolean;
  setChatSidebarVisible: React.Dispatch<React.SetStateAction<boolean>>;
}) => {
  return (
    <>
      <div
        className="chatbox-banner"
        style={{
          position: "sticky",
          top: "0",
          padding: "20px",
          textAlign: "center",
          zIndex: "1000",
          borderBottom: "1px solid ",
          height: "auto",
          pointerEvents: "auto",
        }}
      >
        <Stack direction="row" justifyContent="space-between">
          <button
            onClick={() => setChatSidebarVisible((prevState) => !prevState)}
          ></button>
        </Stack>
      </div>
      <div className="chatbox-banner sticky top-0 p-5 text-center z-50 border-b border-accent h-auto">
        <div className="flex flex-row justify-between">
          <button>
            <HamburgerIcon />
          </button>
          <button
            onClick={() => {
              setChatVisible((prev) => !prev);
            }}
          >
            <CloseIcon />
          </button>
        </div>
      </div>
      <div>
        <ChatMessages
          messages={exampleChatHistory}
          loading={false}
        ></ChatMessages>
      </div>

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
