
import ThemeSelector from "@/components/ThemeSelector";


export default function Page() {
  return (
    <div className=" p-5 m-5 justify-center border-2 border-['accent'] rounded-box w-full">

      <h1 className="text-5xl font-extrabold">Settings</h1>
      <ThemeSelector />
    </div>
  );
}

