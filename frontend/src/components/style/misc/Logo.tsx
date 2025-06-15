import { GiMushroomsCluster } from "react-icons/gi";



export function Logo({ className, style }: { className?: string; style?: React.CSSProperties }) {
  return <GiMushroomsCluster className={className} style={style} />
}


export function HomepageLogo() {
  return <div className="flex flex-col items-center space-y-2" >
    <div className="flex flex-row items-center space-x-9">
      <Logo className="text-6xl lg:text-7xl xl:text-9xl text-base-content" />
      <h1 className="text-5xl lg:text-6xl xl:text-8xl font-bold font-serif tracking-tight">
        KESSLER
      </h1>
    </div>
    <p className="text-md xl:text-xl text-gray-600 text-center font-serif">
      Public Utility Commissions, Simplified.
    </p>
  </div>

}
