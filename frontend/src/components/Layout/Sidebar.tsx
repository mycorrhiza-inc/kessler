// Sidebar.tsx
"use client"
import React, { useState } from 'react';
import { Code, Database, FileText, Menu, Settings } from 'lucide-react';
import { GiMushroomsCluster } from "react-icons/gi";
import { BsArrowBarLeft, BsArrowBarRight } from 'react-icons/bs';
import Link from 'next/link';

interface SidebarButtonProps {
  icon?: React.ElementType;
  label: string;
  href: string;
}

const SidebarLink: React.FC<SidebarButtonProps> = ({ icon: Icon, label, href }) => (
  <Link
    className="w-full flex items-center space-x-2 px-3 py-2 text-sm text-base-700 dark:text-base-300 hover:bg-base-300 dark:hover:bg-base-800 rounded"
    href={href}
  >
    {Icon && <Icon size={16} />}
    <span>{label}</span>
  </Link>
);

interface SidebarProps {
  width: number;
  isVisible: boolean;
  isPinned: boolean;
  onPinChange: (pinned: boolean) => void;
  onVisibilityChange: (visible: boolean) => void;
  onResize: (width: number) => void;
}

const Sidebar: React.FC<SidebarProps> = ({
  width,
  isVisible,
  isPinned,
  onPinChange,
  onVisibilityChange,
  onResize
}) => {
  const [isResizing, setIsResizing] = useState(false);

  const handleMouseDown = (e: React.MouseEvent) => {
    e.preventDefault();
    setIsResizing(true);
    document.addEventListener('mousemove', handleMouseMove);
    document.addEventListener('mouseup', handleMouseUp);
  };

  const handleMouseMove = (e: MouseEvent) => {
    if (isResizing) {
      const newWidth = e.clientX;
      if (newWidth >= 200 && newWidth <= 600) {
        onResize(newWidth);
      }
    }
  };

  const handleMouseUp = () => {
    setIsResizing(false);
    document.removeEventListener('mousemove', handleMouseMove);
    document.removeEventListener('mouseup', handleMouseUp);
  };

  return (
    <div className="z-20">
      {/* Sidebar Toggle Button */}
      <div
        className="fixed bottom-4 left-4 z-30"
        onMouseEnter={() => onVisibilityChange(true)}
      >
        {isVisible || isPinned ? (
          <button
            onClick={() => {
              onPinChange(!isPinned);
              onVisibilityChange(true);
            }}
            className="p-2 bg-base-100 dark:bg-base-800 rounded  hover:bg-base-200 dark:hover:bg-base-700 transition-colors border-2 border-secondary"
          >
            {isPinned ? (
              <BsArrowBarLeft size={20} className="text-base-600 dark:text-base-400 left-4" />
            ) : (
              <BsArrowBarRight size={20} className="text-base-600 dark:text-base-400" />
            )}
          </button>
        ) : (
          <button className="p-2 bg-base-100 dark:bg-base-800 rounded-full  hover:bg-base-200 dark:hover:bg-base-700 transition-colors border-2 border-secondary">
            <Menu size={24} />
          </button>
        )}
      </div>

      {/* Sidebar Content */}
      <div
        className={`fixed left-0 top-0 h-full bg-primary-100 transition-transform duration-100 ease-in-out transform ${isVisible || isPinned ? 'translate-x-0' : '-translate-x-full z-100'
          }`}
        style={{ width: `${width}px` }}
        onMouseLeave={() => !isPinned && onVisibilityChange(false)}
      >
        {/* Sidebar Header */}
        <div className="h-full bg-base-200 dark:bg-base-900 border-r border-base-200 dark:border-base-800 shadow-lag opacity-100">
          <div className="flex flex-col items-center  p-4 border-b border-base-200 dark:border-base-800">
            <div className="h-5 p-4 m-4">
            </div>
          </div>

          <div className="p-4">
            <nav className="space-y-6">
              <div>
                <SidebarLink icon={Code} label="Home" href="/home" />
                <SidebarLink href="/dockets" label="Dockets" />
                <SidebarLink href="/orgs" label="Organizations" />
                <SidebarLink href="/files" label="All Files" />
                <SidebarLink href="/" lable="Landing Page" />
                {/* <SidebarLink icon={Code} label="Query editor" />
                  <SidebarLink icon={Code} label="Query editor" />
                  <SidebarLink icon={Code} label="Query editor" />
                  <SidebarLink icon={Code} label="Query editor" /> */}

              </div>
              {/* <div>
                <h3 className="text-sm font-medium text-base-500 dark:text-base-400 mb-2">Discover</h3>
                <div className="space-y-1">
                  <SidebarLink icon={Code} label="Query editor" />
                </div>
              </div>

              <div>
                <h3 className="text-sm font-medium text-base-500 dark:text-base-400 mb-2">Admin</h3>
                <div className="space-y-1">
                  <SidebarButton icon={Database} label="Indexes" />
                  <SidebarButton icon={Settings} label="Cluster" />
                  <SidebarButton icon={FileText} label="Node info" />
                  <SidebarButton icon={Code} label="API" />
                </div>
              </div> */}
            </nav>
          </div>
        </div>

        {/* Resize Handle */}
        <div
          className={`w-2 h-full cursor-col-resize hover:bg-base-400/50 active:bg-base-400 relative group ${isResizing ? 'bg-base-400' : 'bg-transparent'
            }`}
          onMouseDown={handleMouseDown}
        >
          <div className="absolute inset-y-0 left-1/2 w-0.5 bg-base-300 dark:bg-base-600 group-hover:bg-base-400 dark:group-hover:bg-base-500" />
        </div>
      </div>
    </div>
  );
};

export default Sidebar;