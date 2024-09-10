import Hero from "./Hero";

export default function Landing() {
  return (
    <main>
      <div
        className="searchContainer"
        style={{
          position: "relative",
          width: "99vw",
          height: "80vh",
          padding: "20",
        }}
      >
        <div className="flex items-center justify-center h-full">
          <Hero></Hero>
        </div>
      </div>
    </main>
  );
}
