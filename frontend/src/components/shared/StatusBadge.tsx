import { cn } from '@/lib/utils';
import { CheckCircle2, XCircle, Clock, Loader2 } from 'lucide-react';

interface StatusBadgeProps {
    status: 'pending' | 'running' | 'success' | 'failed';
    className?: string;
}

export default function StatusBadge({ status, className }: StatusBadgeProps) {
    const config = {
        pending: {
            icon: Clock,
            label: 'Pending',
            className: 'bg-yellow-500/10 text-yellow-500 border-yellow-500/20',
        },
        running: {
            icon: Loader2,
            label: 'Running',
            className: 'bg-blue-500/10 text-blue-500 border-blue-500/20',
        },
        success: {
            icon: CheckCircle2,
            label: 'Success',
            className: 'bg-green-500/10 text-green-500 border-green-500/20',
        },
        failed: {
            icon: XCircle,
            label: 'Failed',
            className: 'bg-red-500/10 text-red-500 border-red-500/20',
        },
    };

    const { icon: Icon, label, className: statusClassName } = config[status];

    return (
        <span
            className={cn(
                'inline-flex items-center gap-1.5 rounded-full border px-2.5 py-0.5 text-xs font-medium',
                statusClassName,
                className
            )}
        >
            <Icon className={cn('h-3 w-3', status === 'running' && 'animate-spin')} />
            {label}
        </span>
    );
}
