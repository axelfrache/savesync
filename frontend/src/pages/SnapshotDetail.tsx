import { useParams, Link } from 'react-router-dom';
import { useSnapshot } from '@/hooks/useSnapshots';
import { useSnapshotFiles } from '@/hooks/useSnapshotFiles';
import LoadingState from '@/components/shared/LoadingState';
import ErrorState from '@/components/shared/ErrorState';
import StatusBadge from '@/components/shared/StatusBadge';
import SnapshotFileTree from '@/components/features/snapshots/SnapshotFileTree';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { ScrollArea } from '@/components/ui/scroll-area';
import { ArrowLeft, Download, Loader2 } from 'lucide-react';
import { formatDistanceToNow } from 'date-fns';

export default function SnapshotDetail() {
    const { id } = useParams<{ id: string }>();
    const { data: snapshot, isLoading, error } = useSnapshot(Number(id));
    const { data: fileTree, isLoading: filesLoading, error: filesError } = useSnapshotFiles(Number(id));

    if (isLoading) return <LoadingState />;
    if (error || !snapshot) return <ErrorState message="Snapshot not found" />;

    const handleDownloadManifest = async () => {
        if (!snapshot) return;

        try {
            // Create a temporary link to trigger the download
            // We use the API endpoint directly
            const url = `${import.meta.env.VITE_API_URL || 'http://localhost:8080'}/api/snapshots/${snapshot.id}/manifest`;
            const link = document.createElement('a');
            link.href = url;
            link.setAttribute('download', `manifest-${snapshot.id}.json`);
            document.body.appendChild(link);
            link.click();
            document.body.removeChild(link);
        } catch (err) {
            console.error('Failed to download manifest', err);
        }
    };

    return (
        <div className="space-y-6">
            <div className="flex items-center gap-4">
                <Link to="/snapshots">
                    <Button variant="ghost" size="icon">
                        <ArrowLeft className="h-5 w-5" />
                    </Button>
                </Link>
                <div className="flex-1">
                    <h1 className="text-3xl font-bold text-foreground">Snapshot #{snapshot.id}</h1>
                    <p className="text-muted-foreground mt-1">
                        Created {formatDistanceToNow(new Date(snapshot.created_at), { addSuffix: true })}
                    </p>
                </div>
                <StatusBadge status={snapshot.status} />
            </div>

            <div className="grid gap-6 md:grid-cols-3">
                <Card>
                    <CardHeader>
                        <CardTitle className="text-sm text-muted-foreground">Files</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <p className="text-2xl font-bold text-foreground">{snapshot.file_count}</p>
                    </CardContent>
                </Card>

                <Card>
                    <CardHeader>
                        <CardTitle className="text-sm text-muted-foreground">Total Size</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <p className="text-2xl font-bold text-foreground">
                            {(snapshot.total_bytes / 1024 / 1024).toFixed(2)} MB
                        </p>
                    </CardContent>
                </Card>

                <Card>
                    <CardHeader>
                        <CardTitle className="text-sm text-muted-foreground">Delta</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <p className="text-2xl font-bold text-foreground">
                            {(snapshot.delta_bytes / 1024 / 1024).toFixed(2)} MB
                        </p>
                    </CardContent>
                </Card>
            </div>

            <Card>
                <CardHeader>
                    <div className="flex items-center justify-between">
                        <CardTitle className="text-foreground">Files</CardTitle>
                        <Button size="sm" variant="outline" onClick={handleDownloadManifest}>
                            <Download className="mr-2 h-4 w-4" />
                            Download Manifest
                        </Button>
                    </div>
                </CardHeader>
                <CardContent>
                    {filesLoading && (
                        <div className="flex items-center justify-center py-8">
                            <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
                        </div>
                    )}
                    {filesError && (
                        <p className="text-sm text-destructive">
                            Failed to load file tree. Please try again.
                        </p>
                    )}
                    {fileTree && (
                        <ScrollArea className="h-[600px] w-full rounded-md border border-border p-4">
                            <SnapshotFileTree root={fileTree} />
                        </ScrollArea>
                    )}
                </CardContent>
            </Card>
        </div>
    );
}
