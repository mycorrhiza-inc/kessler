import { ReactNode } from "react";

export default function DefaultContainer({ children }: { children: ReactNode }) {
  return <div className="flex justify-center items-center flex-col space-y-12">
    {children}
  </div >
}
