"use client";
import { AuroraBackground } from "../aceternity/aurora-background";

import { motion } from "framer-motion";
import { Highlight } from "../aceternity/hero-highlight";

export function HeroHighlightText() {
  return (
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
      className="text-2xl px-4 md:text-4xl lg:text-5xl font-bold text-neutral-700 dark:text-white max-w-4xl leading-relaxed lg:leading-snug text-center mx-auto "
    >
      New York PUC Proceedings, now with a{" "}
      <Highlight className="text-black dark:text-white">
        Fast, Modern Interface
      </Highlight>
    </motion.h1>
  );
}
export default function Hero({ isLoggedIn }: { isLoggedIn: boolean }) {
  // Fix the broken min-h-screen stuff and make it actually work
  return (
    <AuroraBackground showRadialGradient={false}>
      <HeroHighlightText />
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
    </AuroraBackground>
  );
}
