import { User } from "@supabase/supabase-js";
import Feature from "./Features";
import Hero from "./Hero";
import Pricing from "./Pricing";
import { SupportedStates } from "@/components/landing/DataSetList";
import { GiMushroomsCluster } from "react-icons/gi";
import "./Landing.css";

function LandingFooter({}) {
  return (
    <>
      <footer className="footer bg-base-200 text-base-content p-10">
        <aside>
          <GiMushroomsCluster style={{ fontSize: "4em" }} />
          <p>
            Mycorrhiza Inc
            <br />
            Providing reliable tech since 5 hours ago.
          </p>
        </aside>
        <nav>
          <h6 className="footer-title">Contact Founders</h6>
          <a href="mailto:n@mycor.io" className="link link-hover">
            Nicole Venner
          </a>
          <a href="mailto:m@mycor.io" className="link link-hover">
            Mirri Bright
          </a>
        </nav>
        <nav>
          <h6 className="footer-title">Company (Links dont work yet)</h6>
          <a className="link link-hover">About us</a>
          <a className="link link-hover">Contact</a>
          <a className="link link-hover">Jobs</a>
          <a className="link link-hover">Blogs </a>
        </nav>
        {/* <nav> */}
        {/*   <h6 className="footer-title">Legal</h6> */}
        {/*   <a className="link link-hover">Terms of use</a> */}
        {/*   <a className="link link-hover">Privacy policy</a> */}
        {/*   <a className="link link-hover">Cookie policy</a> */}
        {/* </nav> */}
      </footer>
    </>
  );
}

export default function Landing({ user }: { user: User | null }) {
  const isLoggedIn = user ? true : false;
  return (
    <>
      <div data-theme="light">
        <Hero isLoggedIn={isLoggedIn} data-theme="light"></Hero>
        <div
          className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-7 "
          data-theme="light"
        >
          <Feature className="col-span-1 md:col-start-1 md:col-span-2 lg:col-start-2 lg:col-span-5 items-center justify-center" />
          {/* <SupportedStates className="col-span-3 md:col-start-2 md:col-span-3  sm:col-span-1" /> */}
        </div>
        <SupportedStates className="col-span-3 md:col-start-2 md:col-span-3  sm:col-span-1" />
        <Pricing className="col-span-1 md:col-start-1 md:col-span-2 lg:col-start-2 lg:col-span-5 items-center justify-center" />
        <LandingFooter />
      </div>
    </>
  );
}
