"use client";
import clsx from "clsx";
import { useTheme } from "next-themes";
import { useState, useEffect } from "react";

interface ThemeData {
  name: string;
  value: string;
  lightdark: string;
}

export const themeDataDictionary: { [key: string]: ThemeData } = {
  light: { name: "kessler", value: "kessler", lightdark: "light" },
};
export const themeDataList: ThemeData[] = Object.values(themeDataDictionary);
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
      <div className=" p-5 m-5 justify-center border-2 border-['accent'] rounded-box">
        <h1 className="text-3xl font-bold">Loading Themes</h1>
      </div>
    );
  }

  // These extra themes are included in TW, but not selectable, is this intentional?
  // "forest",
  // "corporate",
  // "sunset",
  // "acid",
  //

  return (
    <>
      <h1 className="text-3xl font-bold">Themes</h1>
      <div className=" p-5 m-5 justify-center border-2 border-['accent'] rounded-box">
        <div className="flex flex-row flex-wrap gap-5">
          {themeDataList.map((themeData) => (
            <ThemeBigButton
              setTheme={setTheme}
              theme={theme as string}
              themeData={themeData}
            />
          ))}
        </div>
      </div>
    </>
  );
};

const ThemeBigButton = ({
  setTheme,
  theme,
  themeData,
}: {
  setTheme: any;
  theme: string;
  themeData: ThemeData;
}) => {
  return (
    <div
      key={themeData.value}
      onClick={() => setTheme(themeData.value)}
      className="rounded-box"
      data-theme={themeData.value}
      data-act-class="ACTIVECLASS"
    >
      <div className="bg-base-100 text-base-content w-full cursor-pointer font-sans rounded-box shadow-lg p-2">
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
              className={clsx(
                themeData.value === theme ? "" : "invisible",
                "h-3 w-3 shrink-0",
              )}
            >
              <path d="M20.285 2l-11.285 11.567-5.286-5.011-3.714 3.716 9 8.728 15-15.285z"></path>
            </svg>
            <div className="font-bold">{themeData.name}</div>
            <div className="flex flex-wrap gap-1">
              <div className="bg-primary flex aspect-square w-5 items-center justify-center rounded-sm lg:w-6">
                <div className="text-primary-content text-sm font-bold">A</div>
              </div>
              <div className="bg-secondary flex aspect-square w-5 items-center justify-center rounded-sm lg:w-6">
                <div className="text-secondary-content text-sm font-bold">
                  A
                </div>
              </div>
              <div className="bg-accent flex aspect-square w-5 items-center justify-center rounded-sm lg:w-6">
                <div className="text-accent-content text-sm font-bold">A</div>
              </div>
              <div className="bg-neutral flex aspect-square w-5 items-center justify-center rounded-sm lg:w-6">
                <div className="text-neutral-content text-sm font-bold">A</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

const ThemeSmallButton = ({
  setTheme,
  theme,
  themeData,
}: {
  setTheme: any;
  theme: string;
  themeData: ThemeData;
}) => {
  return (
    <span
      className="grid grid-cols-5 grid-rows-3"
      onClick={() => setTheme(themeData.value)}
      data-theme={themeData.value}
    >
      <span className="col-span-5 row-span-3 row-start-1 flex items-center gap-2 px-4 py-3">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="16"
          height="16"
          viewBox="0 0 24 24"
          fill="currentColor"
          className={clsx(
            "h-3 w-3 shrink-0",
            themeData.value === theme ? "" : "invisible",
          )}
        >
          <path d="M20.285 2l-11.285 11.567-5.286-5.011-3.714 3.716 9 8.728 15-15.285z"></path>
        </svg>
        <span className="grow text-sm">{themeData.name}</span>{" "}
        <span className="flex h-full shrink-0 flex-wrap gap-1">
          <span className="bg-primary rounded-badge w-2"></span>{" "}
          <span className="bg-secondary rounded-badge w-2"></span>
          <span className="bg-accent rounded-badge w-2"></span>{" "}
          <span className="bg-neutral rounded-badge w-2"></span>
        </span>
      </span>
    </span>
  );
};
export default ThemeSelector;
