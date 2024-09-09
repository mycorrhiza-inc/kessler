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
      className="z-[9999] w-[900px] max-w-[90vw] min-h-[40vh] fixed bottom-[30px] bg-white dark:bg-gray-700 rounded-[10px] border-2 border-grey-500 p-[10px] text-black dark:text-white hidden"
    >
      <div className="chatbox-banner sticky top-0 bg-[#f5f5f5] dark:bg-gray-400 p-5 text-center z-50 border-b border-gray-300 h-auto">
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
    </motion.div>
  );
};

export default ChatBox;
