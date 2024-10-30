"use client";
// Is this even a good idea/acceptable?
import { createClient } from "@/utils/supabase/server";
import { UserIcon } from "@/components/Icons";
import { useEffect, useState } from "react";
import Modal from "./styled-components/Modal";
import SettingsContent from "./SettingsContent";
import { User } from "@supabase/supabase-js";
import { useKesslerStore } from "@/lib/store";
import Link from "next/link";

function HeaderAuth({ user }: { user: User | null }) {
  const [settingsOpen, setSettingsOpen] = useState(false);
  const globalStore = useKesslerStore();

  useEffect(() => {
    console.log('Is logged in:', globalStore.isLoggedIn);
  }, [user]);
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
const Navbar = ({ user }: { user: User | null }) => {
  return (
    <div
      className="navbar bg-base-200 w-max-50"
    >
      <div className="flex-1 font-semibold">
        <a href="/">Kessler</a>
      </div>
      <div className="flex-none">

        <HeaderAuth user={user} />
      </div>
    </div>
  );
};
export default Navbar;
