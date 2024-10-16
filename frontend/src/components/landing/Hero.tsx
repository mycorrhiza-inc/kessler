"use client";
import { AuroraBackground } from "../aceternity/aurora-background";

import { motion } from "framer-motion";
import { Highlight } from "../aceternity/hero-highlight";
import { Compare } from "../aceternity/compare";

export default function Hero({ isLoggedIn }: { isLoggedIn: boolean }) {
  // Fix the broken min-h-screen stuff and make it actually work
  return (
    <AuroraBackground showRadialGradient={false}>
      <motion.h1
        initial={{
          opacity: 0,
          y: 20,
        }}
        animate={{
          opacity: 1,
          y: [20, -5, 0],
        }}
        transition={{
          duration: 0.5,
          ease: [0.4, 0.0, 0.2, 1],
        }}
      >
        <div className="text-3xl px-4 md:text-5xl lg:text-6xl font-bold text-neutral-700 dark:text-white max-w-4xl leading-relaxed lg:leading-snug text-center mx-auto ">
          New York PUC Proceedings, now with a <br />
          <Highlight className="text-black dark:text-white">
            Fast, Modern Interface
          </Highlight>
        </div>
        {/* <div className="p-4 border rounded-3xl dark:bg-neutral-900 bg-neutral-100  border-neutral-200 dark:border-neutral-800 px-4"> */}
        <div className="flex justify-center space-x-4">
          <Compare
            firstImage="/ny-puc-ui.png"
            secondImage="/kessler-light-rag-search.png"
            firstImageClassName="object-cover object-left-top"
            secondImageClassname="object-cover object-left-top"
            className="h-[280px] w-[500px] md:h-[400px] md:w-[700px] lg:h-[650px] lg:w-[1000px]"
            slideMode="hover"
          />
        </div>
        {/* </div> */}
        {isLoggedIn ? (
          <div className="flex justify-center space-x-4">
            <a
              href="/app"
              className="btn glass shadow-xl btn-lg btn-outline btn-neutral"
            >
              Go To App
            </a>
          </div>
        ) : (
          <>
            <div className="flex justify-center space-x-4">
              <a
                href="/demo"
                className="btn glass shadow-xl btn-lg btn-outline btn-neutral"
              >
                Try Now!
              </a>
            </div>
            <div className="flex justify-center space-x-4">
              <a href="/sign-in" className="btn glass shadow-xl">
                Sign In
              </a>
              <a href="/sign-up" className="btn glass shadow-xl">
                Sign Up
              </a>
            </div>
          </>
        )}
      </motion.h1>
    </AuroraBackground>
  );
}
