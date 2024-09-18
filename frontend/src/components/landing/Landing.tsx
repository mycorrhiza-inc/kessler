import Feature from "./Features";
import Hero from "./Hero";
import Pricing from "./Pricing";

export default function Landing() {
  return (
    <main>
      <div
        className="items-center justify-center h-full"
        style={{
          width: "99vw",
          height: "60vh",
          padding: "20",
        }}
      >
        <div className="items-center justify-center h-full">
          <Hero></Hero>
        </div>

        <div
          className="flex items-center justify-center w-full"
          style={{ minWidth: "30vw" }}
        >
          <Feature></Feature>
          <Pricing></Pricing>
        </div>
      </div>
    </main>
  );
}
