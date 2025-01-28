"use client";
import { ReactNode, useState } from "react";
import Sidebar from "./Sidebar";
import Header from "./Header";
import { BreadcrumbValues } from "../SitemapUtils";

interface LayoutProps {
  children: ReactNode;
  breadcrumbs: BreadcrumbValues;
}

export default function Layout({ children, breadcrumbs }: LayoutProps) {
  const [isSidebarVisible, setIsSidebarVisible] = useState(true);
  const [isSidebarPinned, setIsSidebarPinned] = useState(true);
  const [sidebarWidth, setSidebarWidth] = useState(200); // 16 * 16 = 256px (w-64)
  return (
    <>
      <Header breadcrumbs={breadcrumbs} />
      <Sidebar
        width={sidebarWidth}
        isVisible={isSidebarVisible}
        isPinned={isSidebarPinned}
        onPinChange={setIsSidebarPinned}
        onVisibilityChange={setIsSidebarVisible}
        onResize={setSidebarWidth}
      />
      <div
        className={`flex-1 p-6 transition-all duration-200 ease-in-out w-full`}
        style={{
          marginLeft:
            isSidebarVisible || isSidebarPinned ? `${sidebarWidth}px` : "0",
          width:
            isSidebarVisible || isSidebarPinned
              ? `calc(100% - ${sidebarWidth}px)`
              : "100%",
        }}
      >
        <div className="flex flex-row items-center h-15 pb-20" />
        {children}
      </div>
    </>
  );
}

