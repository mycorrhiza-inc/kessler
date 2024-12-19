package rag


func rerankSearchResults(searchResults []SearchData, query string) ([]SearchData, error) {
	var documents []string
	for _, result := range searchResults {
		documents = append(documents, result.Snippet)
	}
	permutation, err := rerankStringsAndQueryPermutation(context.Background(), query, documents)
	if err != nil {
		return nil, err
	}
	rerankedResults := make([]SearchData, len(searchResults))
	for i, permutation := range permutation {
		rerankedResults[i] = searchResults[permutation]
	}
	return rerankedResults, nil
}


func AppendInstructionHeaderToChathistory(chatHistory *ChatMessage[]) ChatMessage[] {
	instruct_string := `If it would be helpful to link to a Docket, Organization, or File, Include the following in your response
LinkDocket(18-M-0084)
LinkOrganization(Public Service Comission)
LinkFile(777b5c2d-d19e-4711-b2ed-2ba9bcfe449a)
To create a link. Dont redirect to any other goverment resources except through this system.`
	instruct_message := ChatMessage{
		Content: instruct_string, 
		Role: "system"
	}
  chatHistory = append([instruct_message],chatHistory)
	return chatHistory
}

