"use client";
import { ReactNode, useState } from "react";
import Sidebar from "./Sidebar";
import Header from "./Header";
import Link from "next/link";

interface LayoutProps {
  children: ReactNode;
}

export default function Layout({ children }: LayoutProps) {
  const [isSidebarVisible, setIsSidebarVisible] = useState(false);
  const [isSidebarPinned, setIsSidebarPinned] = useState(false);

  const [sidebarWidth, setSidebarWidth] = useState(200); // 16 * 16 = 256px (w-64)
  return (
    <>
      <Sidebar
        width={sidebarWidth}
        isVisible={isSidebarVisible}
        isPinned={isSidebarPinned}
        onPinChange={setIsSidebarPinned}
        onVisibilityChange={setIsSidebarVisible}
        onResize={setSidebarWidth}
      />
      <div
        className={`flex-1 p-6 transition-all bg-base-100 duration-200 ease-in-out w-full`}
        style={{
          marginLeft:
            isSidebarVisible || isSidebarPinned ? `${sidebarWidth}px` : "0",
          width:
            isSidebarVisible || isSidebarPinned
              ? `calc(100% - ${sidebarWidth}px)`
              : "100%",
        }}
      >
        <div className="flex flex-row bg-base-100 items-center h-15 pb-20" />
        {children}
      </div>
    </>
  );
}
