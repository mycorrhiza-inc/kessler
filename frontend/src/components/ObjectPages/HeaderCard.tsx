
type HeaderCardProps = {
	title: String
	children: React.ReactNode
}
const HeaderCard = ({ title, children }: HeaderCardProps) => {
	return (
		<div className="conversation-description  card  rounded-lg border-2">
			<div className="card-body">
				<h1 className="text-2xl font-bold  card-title">
					{title} <br />
				</h1>
				<div className="divider" />
				{children}
			</div>
		</div>
	)
}

export default HeaderCard