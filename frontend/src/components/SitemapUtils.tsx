import { PageContext } from "@/lib/page_context";
import { Link } from "lucide-react";

interface Breadcrumbs {
  state?: string;
  breadcrumbs: {
    value: string;
    title: string;
  }[];
}
interface BreadcrumbValue {
  title: string;
  url?: string;
}

const HeaderBreadcrumbs = ({ breadcrumbs }: { breadcrumbs: Breadcrumbs }) => {
  const breadcrumb_values = breadcrumbs.breadcrumbs.map(
    (b, index): BreadcrumbValue => {
      if (index === breadcrumbs.breadcrumbs.length - 1) {
        return { title: b.title };
      }
      const url = breadcrumbs.breadcrumbs
        .slice(0, index + 1)
        .map((item) => item.value)
        .join("/");
      return { title: b.title, url: url };
    },
  );
  return (
    <div className="breadcrumbs text-xl">
      <ul>
        {breadcrumb_values.map((b) => {
          if (b.url) {
            return (
              <li>
                <Link href={b.url}>{b.title}</Link>
              </li>
            );
          }
          return <li>{b.title}</li>;
        })}
      </ul>
    </div>
  );
};
