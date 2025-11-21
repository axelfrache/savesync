import { useSources } from '@/hooks/useSources';
import { useSnapshots } from '@/hooks/useSnapshots';
import { useJobs } from '@/hooks/useJobs';
import StatsCard from '@/components/shared/StatsCard';
import LoadingState from '@/components/shared/LoadingState';
import ErrorState from '@/components/shared/ErrorState';
import StatusBadge from '@/components/shared/StatusBadge';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { FolderOpen, Camera, Clock, AlertCircle } from 'lucide-react';
import { formatDistanceToNow } from 'date-fns';

export default function Dashboard() {
    const { data: sources, isLoading: sourcesLoading, error: sourcesError } = useSources();
    const { data: snapshots, isLoading: snapshotsLoading } = useSnapshots();
    const { data: jobs, isLoading: jobsLoading } = useJobs();

    if (sourcesLoading || snapshotsLoading || jobsLoading) {
        return <LoadingState />;
    }

    if (sourcesError) {
        return <ErrorState message="Failed to load dashboard data" />;
    }

    const recentJobs = jobs?.slice(0, 5) || [];
    const failedJobs = jobs?.filter((j) => j.status === 'failed').length || 0;
    const lastBackup = jobs?.[0];

    return (
        <div className="space-y-6">
            <div>
                <h1 className="text-3xl font-bold text-foreground">Dashboard</h1>
                <p className="text-muted-foreground mt-1">Overview of your backup system</p>
            </div>

            {/* Stats Grid */}
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
                <StatsCard
                    title="Sources"
                    value={sources?.length || 0}
                    icon={FolderOpen}
                    description="Configured backup sources"
                />
                <StatsCard
                    title="Snapshots"
                    value={snapshots?.length || 0}
                    icon={Camera}
                    description="Total snapshots created"
                />
                <StatsCard
                    title="Last Backup"
                    value={lastBackup ? formatDistanceToNow(new Date(lastBackup.started_at), { addSuffix: true }) : 'Never'}
                    icon={Clock}
                />
                <StatsCard
                    title="Failed Jobs"
                    value={failedJobs}
                    icon={AlertCircle}
                    description="Jobs that failed"
                />
            </div>

            {/* Recent Backups */}
            <Card>
                <CardHeader>
                    <CardTitle>Recent Backups</CardTitle>
                </CardHeader>
                <CardContent>
                    {recentJobs.length === 0 ? (
                        <p className="text-sm text-muted-foreground">No backups yet</p>
                    ) : (
                        <div className="space-y-3">
                            {recentJobs.map((job) => (
                                <div
                                    key={job.id}
                                    className="flex items-center justify-between rounded-lg border border-border p-3"
                                >
                                    <div className="flex items-center gap-3">
                                        <StatusBadge status={job.status} />
                                        <div>
                                            <p className="text-sm font-medium text-foreground">
                                                {job.type === 'backup' ? 'Backup' : 'Restore'} #{job.id}
                                            </p>
                                            <p className="text-xs text-muted-foreground">
                                                {formatDistanceToNow(new Date(job.started_at), { addSuffix: true })}
                                            </p>
                                        </div>
                                    </div>
                                    {job.error && (
                                        <p className="text-xs text-destructive max-w-xs truncate">{job.error}</p>
                                    )}
                                </div>
                            ))}
                        </div>
                    )}
                </CardContent>
            </Card>
        </div>
    );
}
