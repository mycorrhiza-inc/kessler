"use client";
import { AuroraBackground } from "../aceternity/aurora-background";

import { motion } from "framer-motion";
import { Highlight } from "../aceternity/hero-highlight";
import { Compare } from "../aceternity/compare";
import { useKesslerStore } from "@/lib/store";
import Link from "next/link";
import { rootApplicationSlug } from "@/lib/page_context";

export default function Hero() {
  // Fix the broken min-h-screen stuff and make it actually work
  const globalStore = useKesslerStore();
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
        className="gap-96"
      >
        <div className="text-3xl px-4 md:text-5xl lg:text-6xl font-bold text-neutral-700 dark:text-white max-w-4xl leading-relaxed lg:leading-snug text-center mx-auto ">
          New York PUC Dockets, now with a <br />
          <Highlight className="text-black dark:text-white">
            Fast, Modern Interface
          </Highlight>
        </div>
        <div className="flex justify-center space-x-4">
          <Compare
            firstImage="/ny-puc-ui.png"
            secondImage="/kessler-light-rag-search.png"
            firstImageClassName="object-cover object-top-left"
            secondImageClassname="object-cover object-top-left"
            className="h-[280px] w-[500px] md:h-[400px] md:w-[700px] lg:h-[650px] lg:w-[1000px]"
            slideMode="hover"
          />
        </div>
        {globalStore.isLoggedIn ? (
          <div className="flex justify-center space-x-4">
            <Link
              href={rootApplicationSlug}
              className="btn glass shadow-xl btn-lg btn-outline btn-neutral"
            >
              Go To App
            </Link>
          </div>
        ) : (
          <>
            <div className="flex justify-center space-x-4">
              <Link
                href={rootApplicationSlug}
                className="btn glass shadow-xl btn-lg btn-outline btn-neutral"
              >
                Try Now!
              </Link>
            </div>
            <div className="flex justify-center space-x-4">
              <Link href="/sign-in" className="btn glass shadow-xl">
                Sign In
              </Link>
              <Link href="/sign-up" className="btn glass shadow-xl">
                Sign Up
              </Link>
            </div>
          </>
        )}
      </motion.h1>
    </AuroraBackground>
  );
}
