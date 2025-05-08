There are two documents currently describing the search architecture.

One describing my vision for the search libraries in the project
/docs/frontend_search_design.md


And one describing the current state of all the search components.
/docs/frontend_search_current_design.md

For this part of the project I want you to brainstorm out how to add server side rendering to this by utilizing server components. While still retaining functionality with infinite scroll. 

If you need to look up documentation on Server Components or Suspense you can always use context7. For this you might find it very useful to know that if you have a parent server component in react, and create a sub component that is also a server component, you can pass it as a child or a  prop to a client, and it will still render properly

```tsx 

const SomeServerComponent = () => {
  const server_subcomponent = (<SomeServerSubComponent/>)
  // both of these work
  return <>
  <ClientComponent prop={server_subcomponent}/>
  <ClientComponent>{server_subcomponent}<ClientComponent/>
  </>
}
```

Specifically try to think about the overall architecture and how to handle and abstract the flow of all these components so it doesnt increase or decreases complexity from the existing implementation. 

Keep a running thought of your architecture ideas in docs/llm_ideas.md , and once you have finalized everything write up your proposal in docs/srr_search_design.md
