"use client";
import axios from "axios";
import { useState, useEffect } from "react";
import { Filing, FilingTable } from "../Conversations/ConversationComponent";
import Navbar from "../Navbar";


function ConvertToFiling(data: any): Filing {
	const newFiling: Filing = {
		id: data.sourceID,
	};

	return newFiling;
}

export default function RecentUpdatesView() {
	const [searchResults, setSearchResults] = useState([]);
	const [isSearching, setIsSearching] = useState(false);
	const [filing_ids, setFilingIds] = useState<string[]>([]);
	const [page, setPage] = useState(0);
	const getRecentUpdates = async () => {
		setIsSearching(true);
		console.log("getting recent updates");
		try {
			const response = await axios.post("http://api.kessler.xyz/v2/recent_updates", {
				page: 0,
			});
			console.log(response.data);
			if (response.data.length > 0) {
				setFilingIds(response.data.map((item: any) => item.sourceID));
			}
		} catch (error) {
			console.log(error);
		} finally {
			setIsSearching(false);
		}
	};

	const getMore = async () => {
		setIsSearching(true);
		try {
			console.log("getting page ", page + 1);
			const response = await axios.post("http://api.kessler.xyz/v2/recent_updates", {
				page: page + 1,
			});
			setPage(page + 1);
			console.log(response.data);
			if (response.data.length > 0) {
				setFilingIds([...filing_ids, ...response.data.map((item: any) => item.sourceID)]);
			}
		} catch (error) {
			console.log(error);
		} finally {
			setIsSearching(false);
		}
	};

	useEffect(() => {
		getRecentUpdates();
	}, []);

	return (
		<>
			<Navbar user={null} />

			<div className="w-full h-full p-20">
				<div className="w-full h-full p-10 card grid grid-flow-rows box-border border-2 border-black ">
					<h1 className=" text-2xl font-bold">Recent Updates</h1>
					<FilingTable filing_ids={filing_ids} />
					<button onClick={() => getMore()}>Get More</button>
				</div>
			</div>
		</>
	);
}
