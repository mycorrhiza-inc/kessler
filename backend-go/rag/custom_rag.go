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
	instruct_string := `If it would be helpful to link to a Docket, Organization, or File, Include the following in your response
<LinkDocket docket_id="18-M-0084"/>
<LinkOrganization name="Public Service Comission"/>
<LinkOrganization uuid="b5009c5a-873c-44c5-8f7d-ab5b9fa8891b"/>
<LinkFile uuid="777b5c2d-d19e-4711-b2ed-2ba9bcfe449a" />
To create a link. Dont redirect to any other goverment resources except through this system.`
	return_list := make([]ChatMessage, 0)
	return_list = append(return_list, ChatMessage{
		Content: instruct_string,
		Role:    "system",
	})
	return_list = append(return_list, *chatHistory...)
	*chatHistory = return_list
	return *chatHistory
}
