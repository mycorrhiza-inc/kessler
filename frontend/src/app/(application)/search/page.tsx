"use client";
import HomeSearchBar from "@/components/NewSearch/HomeSearch";
import React, { useState } from "react";

import { motion, AnimatePresence } from "framer-motion";
export default async function Page() {
  const [isSearching, setIsSearching] = useState(false);
  return (
    <motion.div
      initial={{ height: "70vh" }}
      animate={{ height: isSearching ? "30vh" : "70vh" }}
      transition={{ duration: 0.5, ease: "easeInOut" }}
      className="flex flex-col items-center justify-center bg-base-100 p-4"
      style={{ overflow: "hidden" }}
    >
      <HomeSearchBar />
    </motion.div>
  );
}
