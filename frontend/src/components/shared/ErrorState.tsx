import { AlertCircle } from 'lucide-react';

interface ErrorStateProps {
    message?: string;
}

export default function ErrorState({ message = 'An error occurred' }: ErrorStateProps) {
    return (
        <div className="rounded-lg border border-red-500/20 bg-red-500/10 p-4">
            <div className="flex items-start gap-3">
                <AlertCircle className="h-5 w-5 text-red-500 flex-shrink-0 mt-0.5" />
                <div>
                    <h3 className="font-semibold text-red-500">Error</h3>
                    <p className="text-sm text-red-400 mt-1">{message}</p>
                </div>
            </div>
        </div>
    );
}
