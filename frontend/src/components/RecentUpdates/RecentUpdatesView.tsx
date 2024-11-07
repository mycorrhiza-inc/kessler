"use client";
import axios from "axios";
import { useState, useEffect } from "react";
import { Filing, FilingTable } from "../Conversations/ConversationComponent";
import Navbar from "../Navbar";


function ConvertToFiling(data: any): Filing {
	const newFiling: Filing = {
		id: data.source_id,
	};

	return newFiling;
}

export default function RecentUpdatesView() {
	const [searchResults, setSearchResults] = useState([]);
	const [isSearching, setIsSearching] = useState(false);
	const [filings, setFilings] = useState<Filing[]>([]);
	const [page, setPage] = useState(0);
	const getRecentUpdates = async () => {
		setIsSearching(true);
		console.log("getting recent updates");
		try {
			const response = await axios.post("http://localhost/v2/recent_updates", {
				page: 0,
			});
			console.log(response.data);
			if (response.data.length > 0) {
				setFilings(response.data);
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
			const response = await axios.post("http://localhost/v2/recent_updates", {
				page: page + 1,
			});
			setPage(page + 1);
			console.log(response.data);
			if (response.data.length > 0) {
				setFilings([filings, ...response.data]);
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
					<FilingTable filings={filings} />
					<button onClick={() => getMore()}>Get More</button>
				</div>
			</div>
		</>
	);
}
