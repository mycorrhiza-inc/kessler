package quickwit_test

func test_quickwit_ingest() {
	// Go ahead and make a new Quickwit Test Index. And upload a bunch of example data to it. (Maybe 1000 documents, each with a 1000 word body, and a 1000 word title, and some numbers as an example.)
	// Nevermind, Its good to actually hit postgres and ingest maybe 100 regular files, or 100 organizations.
	// Go ahead and upload that to quickwit. Using the generic ingest pipeline:
	// Then after waiting a minute or so, go ahead and make a search request to it. Then see if the data you get back is the same as the data you ingested.
}
