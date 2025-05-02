export type Filter = any;

export type FilterKey = string;
// I have a couple of requirements for this type. Namely it should be serializable to a list of all filters can be sent over the internet along with a search query. However it should also support efficent operations such as setting the filter value with a FilterKey. Or deleting a filter corresponding with a resulting FilterKey. I might want to change the data type in the future, so for now just go ahead and implement functions that will do all these primative operations, in such a way that I could change the datastructure in this file to something like a self balancing binary tree and I wouldnt have to adjust anything in the rest of the project.
export type Filters = Filter[];
