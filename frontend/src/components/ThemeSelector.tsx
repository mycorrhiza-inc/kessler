"use client";
import { useTheme } from "next-themes";
import { useState, useEffect } from "react";

const ThemeSelector = () => {
  const [mounted, setMounted] = useState(false);
  const { theme, setTheme } = useTheme();
  // useEffect only runs on the client, so now we can safely show the UI
  useEffect(() => {
    setMounted(true);
  }, []);

  // Prevent theme selection if js not loaded
  if (!mounted) {
    return (
      <>
        <div className="dropdown dropdown-hover">
          <div tabIndex={0} role="button" className="btn m-1">
            Theme Big
          </div>
        </div>
        <div className="dropdown dropdown-hover">
          <div tabIndex={0} role="button" className="btn m-1">
            Theme Small
          </div>
        </div>
      </>
    );
  }
  const theme_list: string[] = [
    "light",
    "dark",
    "forest",
    "lemonade",
    "retro",
    // "colorblind",
    "cyberpunk",
    "valentine",
    "aqua",
  ];

  return (
    <>
      <div className="dropdown dropdown-hover">
        <div tabIndex={0} role="button" className="btn m-1">
          Theme Big
        </div>
        <ul className="menu dropdown-content bg-base-100 rounded-box z-[1] w-52 p-2 shadow">
          {theme_list.map((themeValue) => (
            <li
              key={themeValue}
              data-theme={themeValue}
              onClick={() => setTheme(themeValue)}
            >
              <div
                className="bg-base-100 text-base-content w-full cursor-pointer font-sans"
                data-theme={themeValue}
              >
                <div className="grid grid-cols-5 grid-rows-3">
                  <div className="bg-base-200 col-start-1 row-span-2 row-start-1"></div>{" "}
                  <div className="bg-base-300 col-start-1 row-start-3"></div>{" "}
                  <div className="bg-base-100 col-span-4 col-start-2 row-span-3 row-start-1 flex flex-col gap-1 p-2">
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      width="16"
                      height="16"
                      viewBox="0 0 24 24"
                      fill="currentColor"
                      className={
                        (themeValue === theme ? "" : "invisible ") +
                        "h-3 w-3 shrink-0"
                      }
                    >
                      <path d="M20.285 2l-11.285 11.567-5.286-5.011-3.714 3.716 9 8.728 15-15.285z"></path>
                    </svg>
                    <div className="font-bold">{themeValue}</div>
                    <div className="flex flex-wrap gap-1">
                      <div className="bg-primary flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                        <div className="text-primary-content text-sm font-bold">
                          A
                        </div>
                      </div>
                      <div className="bg-secondary flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                        <div className="text-secondary-content text-sm font-bold">
                          A
                        </div>
                      </div>
                      <div className="bg-accent flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                        <div className="text-accent-content text-sm font-bold">
                          A
                        </div>
                      </div>
                      <div className="bg-neutral flex aspect-square w-5 items-center justify-center rounded lg:w-6">
                        <div className="text-neutral-content text-sm font-bold">
                          A
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </li>
          ))}
        </ul>
      </div>
      <div className="dropdown dropdown-hover">
        <div tabIndex={0} role="button" className="btn m-1">
          Theme Small
        </div>
        <div className="dropdown-content bg-base-200 text-base-content rounded-box top-px  w-56 overflow-y-auto border border-white/5 shadow-2xl outline outline-1 outline-black/5 mt-16">
          {theme_list.map((themeValue) => (
            <span
              className="bg-base-100 rounded-btn text-base-content block w-full cursor-pointer font-sans"
              data-theme={themeValue}
              onClick={() => setTheme(themeValue)}
            >
              <span className="grid grid-cols-5 grid-rows-3">
                <span className="col-span-5 row-span-3 row-start-1 flex items-center gap-2 px-4 py-3">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="16"
                    height="16"
                    viewBox="0 0 24 24"
                    fill="currentColor"
                    className={
                      (themeValue === theme ? "" : "invisible") +
                      " h-3 w-3 shrink-0"
                    }
                  >
                    <path d="M20.285 2l-11.285 11.567-5.286-5.011-3.714 3.716 9 8.728 15-15.285z"></path>
                  </svg>{" "}
                  <span className="flex-grow text-sm">{themeValue}</span>{" "}
                  <span className="flex h-full shrink-0 flex-wrap gap-1">
                    <span className="bg-primary rounded-badge w-2"></span>{" "}
                    <span className="bg-secondary rounded-badge w-2"></span>{" "}
                    <span className="bg-accent rounded-badge w-2"></span>{" "}
                    <span className="bg-neutral rounded-badge w-2"></span>
                  </span>
                </span>
              </span>
            </span>
          ))}
        </div>
      </div>
    </>
  );
};
//<div className="dropdown mb-72" style={{ zIndex: 3001 }}>
//  <div tabIndex={0} role="button" className="btn m-1">
//    Theme
//  </div>
//  <ul
//    tabIndex={0}
//    className="dropdown-content bg-base-300 rounded-box z-[1] w-52 p-2 shadow-2xl"
//  >
//    {themes.map((theme) => (
//      <li key={theme}>
//        <input
//          type="radio"
//          name="theme-dropdown"
//          className="theme-controller btn btn-sm btn-block btn-ghost justify-start"
//          aria-label={theme.charAt(0).toUpperCase() + theme.slice(1)}
//          value={theme.toLowerCase()}
//        />
//      </li>
//    ))}
//    );
//  </ul>
//</div>
export default ThemeSelector;
