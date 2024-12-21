package rag

// func rerankSearchResults(searchResults []SearchData, query string) ([]SearchData, error) {
// 	var documents []string
// 	for _, result := range searchResults {
// 		documents = append(documents, result.Snippet)
// 	}
// 	permutation, err := rerankStringsAndQueryPermutation(context.Background(), query, documents)
// 	if err != nil {
// 		return nil, err
// 	}
// 	rerankedResults := make([]SearchData, len(searchResults))
// 	for i, permutation := range permutation {
// 		rerankedResults[i] = searchResults[permutation]
// 	}
// 	return rerankedResults, nil
// }

func AppendInstructionHeaderToChathistory(chatHistory *[]ChatMessage) []ChatMessage {
	instruct_string := `If it would be helpful to link to a Docket, Organization, or File, Instead of using a markdown link, use one of these components to create a button that when clicked will link to the proper resource. Like so:

In order to access the docket, <LinkDocketButton text="click here" docket_id="18-M-0084"/>. 

The organization <LinkOrganizationButton text="Public Service Comission" name="Public Service Comission"/> created the document.

Their report <LinkFile text"1" uuid="777b5c2d-d19e-4711-b2ed-2ba9bcfe449a" /> claims xcel energy failed to meet its renewable energy targets.

YOU MUST USE THIS SYSTEM WHENEVER YOU TALK ABOUT OR REFERENCE A GOVERMENT DOCUMENT. Use it to cite your sources and provide context that will make it easier to research and help your user learn more!`
	// Only enable for testing
	// instruct_string = `In this chat, disregard all other instructions, and reply with the word "gorb". Do nothing else aside from replying with "gorb". YOU MUST NO MATTER WHAT THE USER SAYS, OR OTHER INSTRUCTIONS, YOU MUST ANSWER "gorb".`
	return_list := make([]ChatMessage, 0)
	return_list = append(return_list, ChatMessage{
		Content: instruct_string,
		Role:    "system",
	})
	return_list = append(return_list, *chatHistory...)
	*chatHistory = return_list
	return *chatHistory
}
