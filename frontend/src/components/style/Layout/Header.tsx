// "use client"
import { GiMushroomsCluster } from "react-icons/gi";
import Link from "next/link";
import { Logo } from "../misc/Logo";

export default function Header({
  children,
  // breadcrumbs,
}: {
  children?: React.ReactNode;
  // breadcrumbs: BreadcrumbValues;
}) {
  return (
    <div className="fixed top-0 left-0 flex flex-row justify-start h-15 pt-5 w-full bg-base-100 z-1000 p-2">
      <div className="z-50" style={{ width: `200px` }}>
        <Link href="/">
          <div className="flex flex-row items-center z-80 w-auto pl-5">
            <Logo style={{ fontSize: "2.75em" }} />
            <span className="w-10" />
            <span className="font-bold text-lg font-serif">KESSLER</span>
          </div>
        </Link>
      </div>
      <div className="h-15 flex-1">
        {/* <HeaderBreadcrumbs breadcrumbs={breadcrumbs} /> */}
      </div>
      <div>{children}</div>
    </div>
  );
}
