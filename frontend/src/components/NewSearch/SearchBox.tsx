"use client";

import { FaSearch } from "react-icons/fa";
import React from "react";

export default function SearchBox({
  value,
  onChange,
  placeholder,
}: {
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
}) {
  return (
    <div className="relative">
      <input
        type="text"
        placeholder={placeholder}
        className="w-full px-4 py-3 pr-10 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent"
        value={value}
        onChange={(e) => onChange(e.target.value)}
      />
      <FaSearch className="h-6 w-6 text-gray-400 absolute right-3 top-1/2 transform -translate-y-1/2" />
    </div>
  );
}
