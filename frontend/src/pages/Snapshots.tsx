import { Link } from 'react-router-dom';
import { useSnapshots } from '@/hooks/useSnapshots';
import LoadingState from '@/components/shared/LoadingState';
import ErrorState from '@/components/shared/ErrorState';
import StatusBadge from '@/components/shared/StatusBadge';
import { Card, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Camera, Download } from 'lucide-react';
import { formatDistanceToNow } from 'date-fns';

export default function Snapshots() {
    const { data: snapshots, isLoading, error } = useSnapshots();

    if (isLoading) return <LoadingState />;
    if (error) return <ErrorState message="Failed to load snapshots" />;

    const handleDownloadManifest = async (e: React.MouseEvent, snapshotId: number) => {
        e.preventDefault(); // Prevent navigation
        e.stopPropagation();

        try {
            const url = `${import.meta.env.VITE_API_URL || 'http://localhost:8080'}/api/snapshots/${snapshotId}/manifest`;
            const link = document.createElement('a');
            link.href = url;
            link.setAttribute('download', `manifest-${snapshotId}.json`);
            document.body.appendChild(link);
            link.click();
            document.body.removeChild(link);
        } catch (err) {
            console.error('Failed to download manifest', err);
        }
    };

    return (
        <div className="space-y-6">
            <div>
                <h1 className="text-3xl font-bold text-foreground">Snapshots</h1>
                <p className="text-muted-foreground mt-1">Browse all backup snapshots</p>
            </div>

            {snapshots?.length === 0 ? (
                <Card>
                    <CardContent className="flex flex-col items-center justify-center py-12">
                        <Camera className="h-12 w-12 text-muted-foreground mb-4" />
                        <p className="text-muted-foreground text-center">
                            No snapshots yet.
                            <br />
                            Run a backup to create your first snapshot.
                        </p>
                    </CardContent>
                </Card>
            ) : (
                <div className="space-y-3">
                    {snapshots?.map((snapshot) => (
                        <Link key={snapshot.id} to={`/snapshots/${snapshot.id}`}>
                            <Card className="hover:bg-accent/50 transition-colors">
                                <CardContent className="flex items-center justify-between p-4">
                                    <div className="flex items-center gap-4">
                                        <Camera className="h-5 w-5 text-muted-foreground" />
                                        <div>
                                            <div className="flex items-center gap-2">
                                                <p className="text-sm font-medium text-foreground">Snapshot #{snapshot.id}</p>
                                                <StatusBadge status={snapshot.status} />
                                            </div>
                                            <p className="text-xs text-muted-foreground">
                                                {formatDistanceToNow(new Date(snapshot.created_at), { addSuffix: true })}
                                            </p>
                                        </div>
                                    </div>
                                    <div className="flex items-center gap-4">
                                        <div className="text-right text-sm">
                                            <p className="text-foreground">{snapshot.file_count} files</p>
                                            <p className="text-muted-foreground">
                                                {(snapshot.total_bytes / 1024 / 1024).toFixed(2)} MB
                                            </p>
                                        </div>
                                        <Button
                                            variant="ghost"
                                            size="icon"
                                            onClick={(e) => handleDownloadManifest(e, snapshot.id)}
                                            title="Download Manifest"
                                        >
                                            <Download className="h-4 w-4" />
                                        </Button>
                                    </div>
                                </CardContent>
                            </Card>
                        </Link>
                    ))}
                </div>
            )}
        </div>
    );
}
