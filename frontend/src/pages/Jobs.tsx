import { useJobs } from '@/hooks/useJobs';
import LoadingState from '@/components/shared/LoadingState';
import ErrorState from '@/components/shared/ErrorState';
import StatusBadge from '@/components/shared/StatusBadge';
import { Card, CardContent } from '@/components/ui/card';
import { ListChecks } from 'lucide-react';
import { formatDistanceToNow } from 'date-fns';

export default function Jobs() {
    const { data: jobs, isLoading, error } = useJobs();

    if (isLoading) return <LoadingState />;
    if (error) return <ErrorState message="Failed to load jobs" />;

    return (
        <div className="space-y-6">
            <div>
                <h1 className="text-3xl font-bold text-foreground">Jobs</h1>
                <p className="text-muted-foreground mt-1">View backup and restore job history</p>
            </div>

            {jobs?.length === 0 ? (
                <Card>
                    <CardContent className="flex flex-col items-center justify-center py-12">
                        <ListChecks className="h-12 w-12 text-muted-foreground mb-4" />
                        <p className="text-muted-foreground text-center">
                            No jobs yet.
                            <br />
                            Run a backup to see job history.
                        </p>
                    </CardContent>
                </Card>
            ) : (
                <div className="space-y-3">
                    {jobs?.map((job) => (
                        <Card key={job.id}>
                            <CardContent className="p-4">
                                <div className="flex items-center justify-between">
                                    <div className="flex items-center gap-4">
                                        <div>
                                            <div className="flex items-center gap-2">
                                                <p className="text-sm font-medium text-foreground">
                                                    {job.type === 'backup' ? 'Backup' : 'Restore'} #{job.id}
                                                </p>
                                                <StatusBadge status={job.status} />
                                            </div>
                                            <p className="text-xs text-muted-foreground mt-1">
                                                Started {formatDistanceToNow(new Date(job.started_at), { addSuffix: true })}
                                            </p>
                                            {job.ended_at && (
                                                <p className="text-xs text-muted-foreground">
                                                    Ended {formatDistanceToNow(new Date(job.ended_at), { addSuffix: true })}
                                                </p>
                                            )}
                                        </div>
                                    </div>
                                    <div className="text-right">
                                        {job.source_id && (
                                            <p className="text-xs text-muted-foreground">Source #{job.source_id}</p>
                                        )}
                                        {job.snapshot_id && (
                                            <p className="text-xs text-muted-foreground">Snapshot #{job.snapshot_id}</p>
                                        )}
                                    </div>
                                </div>
                                {job.error && (
                                    <div className="mt-3 p-2 bg-destructive/10 border border-destructive/20 rounded">
                                        <p className="text-xs text-destructive">{job.error}</p>
                                    </div>
                                )}
                            </CardContent>
                        </Card>
                    ))}
                </div>
            )}
        </div>
    );
}
