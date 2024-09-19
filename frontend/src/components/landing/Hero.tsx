export default function Hero() {
  // Fix the broken min-h-screen stuff and make it actually work
  return (
    <div
      className="hero min-h-full min-w-screen"
      style={{
        height: "60vh",
        backgroundImage: "url(/landing-background.webp)",
      }}
    >
      <div className="hero-overlay bg-opacity-40"></div>
      <div className="hero-content text-neutral-content text-center flex flex-col items-center w-max">
        <h1 className="mb-5 text-5xl font-bold">
          Welcome to <br /> Kessler
        </h1>
        <p className="mb-5">Please use our application 🙏 Namaste</p>
        <div className="flex justify-center space-x-4">
          <a href="/demo" className="btn glass shadow-xl btn-lg">
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
      </div>
    </div>
  );
}
