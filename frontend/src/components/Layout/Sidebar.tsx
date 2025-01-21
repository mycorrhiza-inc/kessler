// Sidebar.tsx
"use client"
import React, { useState } from 'react';
import { Code, Database, FileText, Menu, Settings } from 'lucide-react';
import { GiMushroomsCluster } from "react-icons/gi";
import { BsArrowBarLeft, BsArrowBarRight } from 'react-icons/bs';
import Link from 'next/link';

interface SidebarButtonProps {
  icon: React.ElementType;
  label: string;
  href: string;
}

const SidebarLink: React.FC<SidebarButtonProps> = ({ icon: Icon, label, href }) => (
  <Link
    className="w-full flex items-center space-x-2 px-3 py-2 text-sm text-neutral-700 dark:text-neutral-300 hover:bg-neutral-200 dark:hover:bg-neutral-800 rounded"
    href={href}
  >
    <Icon size={16} />
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
    <>
      {/* Sidebar Toggle Button */}
      <div
        className="fixed bottom-4 left-4 z-50"
        onMouseEnter={() => onVisibilityChange(true)}
      >
        {isVisible || isPinned ? (
          <button
            onClick={() => {
              onPinChange(!isPinned);
              onVisibilityChange(true);
            }}
            className="p-1 hover:bg-neutral-200 dark:hover:bg-neutral-700 rounded transition-colors"
          >
            {isPinned ? (
              <BsArrowBarLeft size={20} className="text-neutral-600 dark:text-neutral-400 left-4" />
            ) : (
              <BsArrowBarRight size={20} className="text-neutral-600 dark:text-neutral-400" />
            )}
          </button>
        ) : (
          <button className="p-2 bg-neutral-100 dark:bg-neutral-800 rounded-full shadow-lg hover:bg-neutral-200 dark:hover:bg-neutral-700 transition-colors">
            <Menu size={24} />
          </button>
        )}
      </div>

      {/* Sidebar Content */}
      <div
        className={`fixed left-0 top-0 h-full transition-transform duration-100 ease-in-out transform ${isVisible || isPinned ? 'translate-x-0' : '-translate-x-full'
          }`}
        style={{ width: `${width}px` }}
        onMouseLeave={() => !isPinned && onVisibilityChange(false)}
      >
        {/* Sidebar Header */}
        <div className="h-full bg-neutral-50 dark:bg-neutral-900 border-r border-neutral-200 dark:border-neutral-800 shadow-lg">
          <div className="flex flex-col items-center justify-between p-4 border-b border-neutral-200 dark:border-neutral-800">
            <div className="h-5 p-4 m-4">
            </div>
          </div>

          <div className="p-4">
            <nav className="space-y-6">
              <div>
                  <SidebarLink icon={Code} label="Home" href="/home" />
                  {/* <SidebarLink icon={Code} label="Query editor" />
                  <SidebarLink icon={Code} label="Query editor" />
                  <SidebarLink icon={Code} label="Query editor" />
                  <SidebarLink icon={Code} label="Query editor" /> */}
                 
              </div>
              {/* <div>
                <h3 className="text-sm font-medium text-neutral-500 dark:text-neutral-400 mb-2">Discover</h3>
                <div className="space-y-1">
                  <SidebarLink icon={Code} label="Query editor" />
                </div>
              </div>

              <div>
                <h3 className="text-sm font-medium text-neutral-500 dark:text-neutral-400 mb-2">Admin</h3>
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
          className={`w-2 h-full cursor-col-resize hover:bg-neutral-400/50 active:bg-neutral-400 relative group ${isResizing ? 'bg-neutral-400' : 'bg-transparent'
            }`}
          onMouseDown={handleMouseDown}
        >
          <div className="absolute inset-y-0 left-1/2 w-0.5 bg-neutral-300 dark:bg-neutral-600 group-hover:bg-neutral-400 dark:group-hover:bg-neutral-500" />
        </div>
      </div>
    </>
  );
};

export default Sidebar;