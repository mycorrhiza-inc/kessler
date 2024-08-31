
const TagList = ({ tags }: { tags: Array<any> }) => {
	  return (
	<div>
	  {tags.map(tag => (
		<div key={tag.id} >

		</div>
	  ))}
	</div>
  );
}

export default TagList;