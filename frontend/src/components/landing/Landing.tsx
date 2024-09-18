import Feature from "./Features";
import Hero from "./Hero";

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

        <div className=" items-center justify-center">
          <Feature></Feature>
        </div>
      </div>
    </main>
  );
}
