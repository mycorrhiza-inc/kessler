import { GiMushroomsCluster } from "react-icons/gi";
import Navbar from "@/components/Page/Navbar";
import { HeaderBreadcrumbs } from "../SitemapUtils";

export default function Header() {
	return (
		<div className="fixed top-0 left-0 flex flex-row justify-between items-center h-15 pb-10 w-full">
			<div
				className="  z-50"
				style={{ width: `200px` }}
			>
				<div className="flex flex-row items-center p-4 m-4">
					<GiMushroomsCluster style={{ fontSize: "2em" }} />
					<span className='w-10' />
					<span className="font-bold text-lg">KESSLER</span>
				</div>
			</div>
			<div className="pr-10 mr-10">
				<HeaderBreadcrumbs breadcrumbs={{
					state: "ny",
					breadcrumbs: [{ title: "Files", value: "/files" }],
				}} />
			</div>
		</div>
	)
}