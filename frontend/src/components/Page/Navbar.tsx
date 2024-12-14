"use client";
import { UserIcon } from "@/components/Icons";
import { useState } from "react";
import Modal from "@/components/styled-components/Modal";
import SettingsContent from "@/components/SettingsContent";
import { useKesslerStore } from "@/lib/store";
import Link from "next/link";
import { BreadcrumbValues, HeaderBreadcrumbs } from "@/components/SitemapUtils";

function HeaderAuth() {
  const [settingsOpen, setSettingsOpen] = useState(false);
  const globalStore = useKesslerStore();

  // useEffect(() => {
  //   console.log("Is logged in:", globalStore.isLoggedIn);
  // }, []);
  return globalStore.isLoggedIn ? (
    <>
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
    </>
  ) : (
    <div className="flex gap-2">
      <Link href="/sign-in" className="btn btn-outline btn-secondary">
        Sign in
      </Link>
      <Link href="/sign-up" className="btn btn-outline btn-secondary">
        Sign up
      </Link>
    </div>
  );
}
const Navbar = ({ breadcrumbs }: { breadcrumbs: BreadcrumbValues }) => {
  return (
    <div className="navbar bg-base-200 w-max-50">
      <div className="flex-1 font-semibold">
        <HeaderBreadcrumbs breadcrumbs={breadcrumbs} />
      </div>
      <div className="flex-none">
        <HeaderAuth />
      </div>
    </div>
  );
};
export default Navbar;
