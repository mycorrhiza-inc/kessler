
'use client'
import { ReactNode, useState } from 'react';
import Sidebar from './sidebar';
import Header from './header';

interface LayoutProps {
	children: ReactNode;
}

export default function Layout({ children }: LayoutProps) {
	const [isSidebarVisible, setIsSidebarVisible] = useState(true);
	const [isSidebarPinned, setIsSidebarPinned] = useState(true);
	const [sidebarWidth, setSidebarWidth] = useState(200); // 16 * 16 = 256px (w-64)
	return (
		<>
			<Header />
			<Sidebar
				width={sidebarWidth}
				isVisible={isSidebarVisible}
				isPinned={isSidebarPinned}
				onPinChange={setIsSidebarPinned}
				onVisibilityChange={setIsSidebarVisible}
				onResize={setSidebarWidth}
			/>
			<div
				className={`flex-1 p-6 transition-all duration-200 ease-in-out`}
				style={{
					marginLeft: (isSidebarVisible || isSidebarPinned) ? `${sidebarWidth}px` : '-1px',
				}}
			>
				<div className="flex flex-row items-center h-15 pb-20 w-full"/>
				{children}
			</div>
		</>
	)
}