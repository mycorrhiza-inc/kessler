import { User } from "@supabase/supabase-js";
import Navbar from "./Navbar";
import { BreadcrumbValues } from "../SitemapUtils";

const PageContainer = ({
  breadcrumbs,
  children,
}: {
  breadcrumbs: BreadcrumbValues;
  children: React.ReactNode;
}) => {
  return (
    <div className="w-full">
      <Navbar breadcrumbs={breadcrumbs} />
      <div className="w-full h-full p-20">{children}</div>
    </div>
  );
};
// Maybe use this?
// <div className="w-full h-full lg:pr-20 lg:pl-20">

export default PageContainer;
