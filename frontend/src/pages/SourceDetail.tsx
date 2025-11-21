import { useParams, Link } from 'react-router-dom';
import { useSource } from '@/hooks/useSources';
import { useSnapshots } from '@/hooks/useSnapshots';
import LoadingState from '@/components/shared/LoadingState';
import ErrorState from '@/components/shared/ErrorState';
import StatusBadge from '@/components/shared/StatusBadge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { ArrowLeft, Play } from 'lucide-react';
import { formatDistanceToNow } from 'date-fns';

export default function SourceDetail() {
    const { id } = useParams<{ id: string }>();
    const { data: source, isLoading, error } = useSource(Number(id));
    const { data: snapshots } = useSnapshots();

    if (isLoading) return <LoadingState />;
    if (error || !source) return <ErrorState message="Source not found" />;

    const sourceSnapshots = snapshots?.filter((s) => s.source_id === source.id) || [];

    return (
        <div className="space-y-6">
            <div className="flex items-center gap-4">
                <Link to="/sources">
                    <Button variant="ghost" size="icon">
                        <ArrowLeft className="h-5 w-5" />
                    </Button>
                </Link>
                <div className="flex-1">
                    <h1 className="text-3xl font-bold text-foreground">{source.name}</h1>
                    <p className="text-muted-foreground mt-1">{source.path}</p>
                </div>
                <Button>
                    <Play className="mr-2 h-4 w-4" />
                    Run Backup
                </Button>
            </div>

            <div className="grid gap-6 md:grid-cols-2">
                <Card>
                    <CardHeader>
                        <CardTitle>Configuration</CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-2 text-sm">
                        <div>
                            <span className="text-muted-foreground">Path:</span>
                            <span className="ml-2 text-foreground">{source.path}</span>
                        </div>
                        <div>
                            <span className="text-muted-foreground">Exclusions:</span>
                            <span className="ml-2 text-foreground">{source.exclusions.length || 'None'}</span>
                        </div>
                        {source.exclusions.length > 0 && (
                            <div className="flex flex-wrap gap-1 mt-2">
                                {source.exclusions.map((pattern, i) => (
                                    <span key={i} className="px-2 py-1 bg-muted rounded text-xs">
                                        {pattern}
                                    </span>
                                ))}
                            </div>
                        )}
                    </CardContent>
                </Card>

                <Card>
                    <CardHeader>
                        <CardTitle>Statistics</CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-2 text-sm">
                        <div>
                            <span className="text-muted-foreground">Total Snapshots:</span>
                            <span className="ml-2 text-foreground">{sourceSnapshots.length}</span>
                        </div>
                        <div>
                            <span className="text-muted-foreground">Last Backup:</span>
                            <span className="ml-2 text-foreground">
                                {sourceSnapshots[0]
                                    ? formatDistanceToNow(new Date(sourceSnapshots[0].created_at), { addSuffix: true })
                                    : 'Never'}
                            </span>
                        </div>
                    </CardContent>
                </Card>
            </div>

            <Card>
                <CardHeader>
                    <CardTitle>Snapshot History</CardTitle>
                </CardHeader>
                <CardContent>
                    {sourceSnapshots.length === 0 ? (
                        <p className="text-sm text-muted-foreground">No snapshots yet</p>
                    ) : (
                        <div className="space-y-2">
                            {sourceSnapshots.map((snapshot) => (
                                <Link
                                    key={snapshot.id}
                                    to={`/snapshots/${snapshot.id}`}
                                    className="flex items-center justify-between rounded-lg border border-border p-3 hover:bg-accent/50 transition-colors"
                                >
                                    <div className="flex items-center gap-3">
                                        <StatusBadge status={snapshot.status} />
                                        <div>
                                            <p className="text-sm font-medium text-foreground">Snapshot #{snapshot.id}</p>
                                            <p className="text-xs text-muted-foreground">
                                                {formatDistanceToNow(new Date(snapshot.created_at), { addSuffix: true })}
                                            </p>
                                        </div>
                                    </div>
                                    <div className="text-right text-sm">
                                        <p className="text-foreground">{snapshot.file_count} files</p>
                                        <p className="text-muted-foreground">{(snapshot.total_bytes / 1024 / 1024).toFixed(2)} MB</p>
                                    </div>
                                </Link>
                            ))}
                        </div>
                    )}
                </CardContent>
            </Card>
        </div>
    );
}
