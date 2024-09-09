import { Stack } from "@mui/joy";
import { motion } from "framer-motion";
import { Dispatch, SetStateAction, useState } from "react";
import { CloseIcon, HamburgerIcon } from "@/components/Icons";
import { ChatMessages, exampleChatHistory } from "./ChatHistory";

interface ChatBoxProps {
  chatVisible: boolean;
  setChatVisible: Dispatch<SetStateAction<boolean>>;
}

const ChatBox = ({ chatVisible, setChatVisible }: ChatBoxProps) => {
  return (
    <motion.div
      animate={{
        height: "60vh",
        minHeight: "40vh",
        position: "fixed",
        y: chatVisible ? "20vh" : "150vh", // if mobile, animate to the middle of the screen
        x: "0", // if mobile, animate to the middle of the screen
        display: "block",
      }}
      transition={{ type: "spring", stiffness: 100, damping: 20 }} // Customize the animation
      data-isOpen={!chatVisible}
      style={{
        zIndex: "9999",
        width: "900px",
        maxWidth: "90vw",
        minHeight: "40vh",
        position: "fixed",
        bottom: "30px",
        backgroundColor: "white",
        borderRadius: "10px",
        border: "2px solid grey",
        padding: "10px",
        color: "black",
        display: "none",
      }}
    >
      <div
        className="chatbox-banner"
        style={{
          position: "sticky",
          top: "0",
          backgroundColor: "#f1f1f1",
          padding: "20px",
          textAlign: "center",
          zIndex: "1000" /* Ensures it's on top of other content */,
          borderBottom: "1px solid #ccc",
          height: "auto",
        }}
      >
        <Stack direction="row" justifyContent="space-between">
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
        </Stack>
      </div>
      <div>
        <ChatMessages
          messages={exampleChatHistory}
          loading={false}
        ></ChatMessages>
      </div>
    </motion.div>
  );
};

export default ChatBox;
