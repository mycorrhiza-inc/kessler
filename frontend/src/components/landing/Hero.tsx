import { AuroraBackground } from "../aceternity/aurora-background";

export default function Hero({ isLoggedIn }: { isLoggedIn: boolean }) {
  // Fix the broken min-h-screen stuff and make it actually work
  return (
    <AuroraBackground>
      <h1 className="mb-5 text-5xl font-bold">
        Welcome to <br /> Kessler
      </h1>
      <p className="mb-5">Please use our application 🙏 Namaste</p>
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
