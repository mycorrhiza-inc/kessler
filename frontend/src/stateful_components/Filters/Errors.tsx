import React from 'react'
import { useKesslerStore } from '@/lib/store'

const DefaultFilterErrorFallback: React.FC<{ error: Error }> = ({ error }) => {
	return (
		<div className="card bg-error text-error-content">
			<div className="card-body">
				<h2 className="card-title">🚨 Filter System Error</h2>
				<p className="mb-4">Something went wrong with the filter system:</p>
				<div className="bg-error-content bg-opacity-10 p-3 rounded mb-4">
					<code className="text-sm break-all">{error.message}</code>
				</div>
				<div className="card-actions justify-end">
					<button
						className="btn btn-primary"
						onClick={() => {
							try {
								const store = useKesslerStore.getState();
								store.cleanup();
							} catch (e) {
								console.warn('Failed to cleanup store:', e);
							}
							window.location.reload();
						}}
					>
						🔄 Reload Page
					</button>
					<button
						className="btn btn-secondary"
						onClick={() => {
							try {
								const store = useKesslerStore.getState();
								store.cleanup();
								store.resetFilters();
								window.location.href = window.location.pathname;
							} catch (e) {
								console.warn('Failed to reset filters:', e);
								window.location.reload();
							}
						}}
					>
						🗑️ Reset Filters
					</button>
				</div>
			</div>
		</div>
	)
};

const DetailedFilterErrorFallback: React.FC<{ error: Error }> = ({ error }) => {
	const [showDetails, setShowDetails] = React.useState(false);

	return (
		<div className="card bg-error text-error-content max-w-2xl mx-auto">
			<div className="card-body">
				<h2 className="card-title">🚨 Filter System Critical Error</h2>

				<div className="space-y-4">
					<p>The filter system encountered a critical error and needs to be reset.</p>

					<div className="bg-error-content bg-opacity-10 p-3 rounded">
						<div className="flex justify-between items-center mb-2">
							<strong>Error Details:</strong>
							<button
								className="btn btn-xs btn-ghost"
								onClick={() => setShowDetails(!showDetails)}
							>
								{showDetails ? '▲ Hide' : '▼ Show'}
							</button>
						</div>

						<code className="text-sm break-all">
							{showDetails ? error.stack || error.message : error.message}
						</code>
					</div>

					<div className="alert alert-warning">
						<svg xmlns="http://www.w3.org/2000/svg" className="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
							<path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.96-.833-2.732 0L3.732 16c-.77.833.192 2.5 1.732 2.5z" />
						</svg>
						<div>
							<h3 className="font-bold">Recovery Options:</h3>
							<div className="text-sm mt-1">
								• <strong>Soft Reset:</strong> Clears filters but keeps page state<br />
								• <strong>Hard Reset:</strong> Reloads entire page<br />
								• <strong>Report:</strong> Copies error details for debugging
							</div>
						</div>
					</div>
				</div>

				<div className="card-actions justify-between">
					<button
						className="btn btn-ghost btn-sm"
						onClick={() => {
							navigator.clipboard.writeText(`Filter System Error:\n${error.stack || error.message}\n\nTimestamp: ${new Date().toISOString()}`);
							alert('Error details copied to clipboard');
						}}
					>
						📋 Copy Error
					</button>

					<div className="space-x-2">
						<button
							className="btn btn-secondary"
							onClick={() => {
								try {
									const store = useKesslerStore.getState();
									store.cleanup();
									store.resetFilters();
									window.location.href = window.location.pathname;
								} catch (e) {
									console.warn('Soft reset failed, doing hard reset:', e);
									window.location.reload();
								}
							}}
						>
							🔧 Soft Reset
						</button>
						<button
							className="btn btn-primary"
							onClick={() => window.location.reload()}
						>
							🔄 Hard Reset
						</button>
					</div>
				</div>
			</div>
		</div>
	);
};

export class FilterErrorBoundary extends React.Component<
	{
		children: React.ReactNode;
		fallback?: React.ComponentType<{ error: Error }>;
		detailed?: boolean;
	},
	{ hasError: boolean; error: Error | null }
> {
	constructor(props: any) {
		super(props);
		this.state = { hasError: false, error: null };
	}

	static getDerivedStateFromError(error: Error) {
		return { hasError: true, error };
	}

	componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
		console.error('Filter system error boundary caught:', error, errorInfo);

		// Try to report error to store if possible
		try {
			const store = useKesslerStore.getState();
			store.setFilterError(`Critical Error: ${error.message}`);

			// Auto-reset filters for certain error types to prevent cascading failures
			if (
				error.message.includes('configuration') ||
				error.message.includes('validation') ||
				error.message.includes('Cannot read properties')
			) {
				console.warn('Auto-resetting filters due to critical error');
				store.resetFilters().catch(console.error);
			}
		} catch (storeError) {
			console.error('Failed to update store from error boundary:', storeError);
		}
	}

	render() {
		if (this.state.hasError) {
			// Use custom fallback, detailed fallback, or default fallback
			const FallbackComponent =
				this.props.fallback ||
				(this.props.detailed ? DetailedFilterErrorFallback : DefaultFilterErrorFallback);

			return <FallbackComponent error={this.state.error!} />;
		}

		return this.props.children;
	}
}