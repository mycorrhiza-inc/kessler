"use client";
// Is this even a good idea/acceptable?
import { createClient } from "@/utils/supabase/server";
import { UserIcon } from "@/components/Icons";
import { useState } from "react";
import Modal from "./styled-components/Modal";
import SettingsContent from "./SettingsContent";
import { User } from "@supabase/supabase-js";
import Link from "next/link";

function HeaderAuth({ user }: { user: User | null }) {
  const [settingsOpen, setSettingsOpen] = useState(false);

  return user ? (
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
        <SettingsContent user={user} />
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
const Header = ({ user }: { user: User | null }) => {
  return (
    <nav
      className="w-full flex justify-center border-b border-b-foreground/10 h-16 bg-base-200 text-base-content"
      style={{ zIndex: 3000 }}
    >
      <div
        className="w-full max-w-5xl flex justify-between items-center bg-base-200 p-3 px-5 text-sm"
        style={{ zIndex: 3000 }}
      >
        <div className="flex gap-5 items-center font-semibold">
          <Link href="/">Kessler</Link>
        </div>
        <HeaderAuth user={user} />
      </div>
    </nav>
  );
};
export default Header;
