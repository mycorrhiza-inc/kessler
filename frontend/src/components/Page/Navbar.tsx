"use client";
import { AngleDownIcon, HamburgerIcon, UserIcon } from "@/components/Icons";
import { useState } from "react";
import Modal from "@/components/styled-components/Modal";
import SettingsContent from "@/components/SettingsContent";
import { useKesslerStore } from "@/lib/store";
import Link from "next/link";
import { BreadcrumbValues, HeaderBreadcrumbs } from "@/components/SitemapUtils";

export function HeaderAuth({ children }: { children?: React.ReactNode }) {
  const [settingsOpen, setSettingsOpen] = useState(false);
  const globalStore = useKesslerStore();

  // useEffect(() => {
  //   console.log("Is logged in:", globalStore.isLoggedIn);
  // }, []);
  return (
    <div className="flex gap-2">
      <div>
        {children}
      </div>
      <div
        tabIndex={0}
        role="button"
        className="btn btn-primary rounded-btn"
        onClick={() => setSettingsOpen((prev) => !prev)}
      >
        <UserIcon />
      </div>
      <Modal open={settingsOpen} setOpen={setSettingsOpen}>
        <SettingsContent />
      </Modal>
    </div>
  );
}

const HeaderMenus = () => {
  return (
    <div className="dropdown z-100">
      <div tabIndex={0} role="button" className="btn m-1 mr-5">
        <HamburgerIcon />
      </div>
      <ul
        tabIndex={0}
        className="dropdown-content menu bg-base-100 rounded-box z-50 w-52 p-2 shadow"
      >
        <li>
          <Link href="/home">Home</Link>
        </li>
        <li>
          <Link href="/dockets">Dockets</Link>
        </li>
        <li>
          <Link href="/orgs">Organizations</Link>
        </li>
        <li>
          <Link href="/files">All Files</Link>
        </li>
        <li>
          <Link href="/">Landing Page</Link>
        </li>
      </ul>
    </div>
  );
};
const Navbar = ({ breadcrumbs }: { breadcrumbs: BreadcrumbValues }) => {
  return (
    <div className="navbar bg-base-200 w-max-50">
      <div className="flex-1 font-semibold">
        <HeaderMenus />
        <span />
      </div>
      <div className="flex-none">
        <HeaderAuth />
      </div>
    </div>
  );
};
export default HeaderMenus;
