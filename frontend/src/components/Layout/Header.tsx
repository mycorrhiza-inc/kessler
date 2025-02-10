import { GiMushroomsCluster } from "react-icons/gi";
import Navbar, { HeaderAuth } from "@/components/Page/Navbar";
import { BreadcrumbValues, HeaderBreadcrumbs } from "../SitemapUtils";

export default function Header({
  children,
  breadcrumbs,
}: {
  children?: React.ReactNode;
  breadcrumbs: BreadcrumbValues;
}) {
  return (
    <div className="fixed top-0 left-0 flex flex-row  justify-start h-15 pt-5 w-full bg-base-100 z-50 pr-5">
      <div className="z-50" style={{ width: `200px` }}>
        <div className="flex flex-row items-center z-80 w-auto pl-5 ">
          <GiMushroomsCluster style={{ fontSize: "2.75em" }} />
          <span className="w-10" />
          <span className="font-bold text-lg">KESSLER</span>
        </div>
      </div>
      <div className="h-15 flex-1">
        <HeaderBreadcrumbs breadcrumbs={breadcrumbs} />
      </div>
      <div className="flex-none">
        <HeaderAuth>{children}</HeaderAuth>
      </div>
    </div>
  );
}

