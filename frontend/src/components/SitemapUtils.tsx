import {
  getStateDisplayName,
  rootApplicationSlug,
} from "@/lib/page_context";
import Link from "next/link";

export interface BreadcrumbValues {
  state?: string;
  breadcrumbs: {
    value: string;
    title: string;
  }[];
}

export const HeaderBreadcrumbs = ({
  breadcrumbs,
}: {
  breadcrumbs: BreadcrumbValues;
}) => {
  const statename = getStateDisplayName(breadcrumbs.state || "");
  const new_breadcrumbs = [
    { value: "", title: "" },
    { value: "/", title: statename },
    ...breadcrumbs.breadcrumbs,
  ];
  var breadcrumb_values = new_breadcrumbs.map((b, index) => {
    if (index === 0) {
      return { title: b.title, url: "/" };
    }
    if (index === breadcrumbs.breadcrumbs.length) {
      return { title: b.title };
    }
    const url = new_breadcrumbs
      .slice(0, index + 1)
      .map((item) => item.value)
      .join("/");
    return { title: b.title, url: url };
  });
  // Hack to get the breadcrumb to redirect to the root app page and not the landing page
  breadcrumb_values[0].url = rootApplicationSlug;
  return (
    <div className="breadcrumbs text-xl">
      <ul>
        {breadcrumb_values.map((b) => {
          if (b.url) {
            return (
              <li key={b.title}>
                <Link href={b.url}>{b.title}</Link>
              </li>
            );
          }
          return <li key={b.title}>{b.title}</li>;
        })}
      </ul>
    </div>
  );
};
