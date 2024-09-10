export default function Hero() {
  // TODO: WEBP the background image to not ddos a very good website unsplash
  // TODO: Replace with our own photo.
  return (
    <div
      className="hero min-h-screen w-max"
      style={{
        backgroundImage:
          "url(https://images.unsplash.com/photo-1506475064951-f1c5640c1e59)",
      }}
    >
      <div className="hero-overlay bg-opacity-60"></div>
      <div className="hero-content text-neutral-content text-center flex flex-col items-center">
        <div className="max-w-md">
          <h1 className="mb-5 text-5xl font-bold">Welcome to Kessler</h1>
          <p className="mb-5">Please use our application 🙏 Namaste</p>
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
    </div>
  );
}
