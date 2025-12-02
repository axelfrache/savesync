import { useState } from 'react';
import { Link } from 'react-router-dom';
import { useSources, useDeleteSource, useRunBackup } from '@/hooks/useSources';
import { useTargets } from '@/hooks/useTargets';
import LoadingState from '@/components/shared/LoadingState';
import ErrorState from '@/components/shared/ErrorState';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle, CardFooter } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Plus, Play, Pencil, Trash2, FolderOpen, MoreHorizontal } from 'lucide-react';
import SourceDialog from '@/components/features/sources/SourceDialog';
import type { Source } from '@/types/api';
import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
} from '@/components/ui/alert-dialog';
import { useToast } from '@/hooks/use-toast';

export default function Sources() {
    const { data: sources, isLoading, error } = useSources();
    const { data: targets } = useTargets();
    const deleteSource = useDeleteSource();
    const runBackup = useRunBackup();
    const { toast } = useToast();

    const [dialogOpen, setDialogOpen] = useState(false);
    const [editingSource, setEditingSource] = useState<Source | undefined>();
    const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
    const [deletingSource, setDeletingSource] = useState<Source | undefined>();

    if (isLoading) return <LoadingState />;
    if (error) return <ErrorState message="Failed to load sources" />;

    const handleEdit = (source: Source) => {
        setEditingSource(source);
        setDialogOpen(true);
    };

    const handleDelete = (source: Source) => {
        setDeletingSource(source);
        setDeleteDialogOpen(true);
    };

    const confirmDelete = async () => {
        if (!deletingSource) return;
        try {
            await deleteSource.mutateAsync(deletingSource.id);
            toast({
                title: 'Source deleted',
                description: `${deletingSource.name} has been deleted`,
            });
        } catch (error) {
            toast({
                title: 'Error',
                description: 'Failed to delete source',
                variant: 'destructive',
            });
        } finally {
            setDeleteDialogOpen(false);
            setDeletingSource(undefined);
        }
    };

    const handleRunBackup = async (sourceId: number, sourceName: string) => {
        try {
            await runBackup.mutateAsync(sourceId);
            toast({
                title: 'Backup started',
                description: `Backup for ${sourceName} has been queued`,
            });
        } catch (error) {
            toast({
                title: 'Error',
                description: 'Failed to start backup',
                variant: 'destructive',
            });
        }
    };

    const getTargetName = (targetId: number | null) => {
        if (!targetId) return 'No target';
        const target = targets?.find((t) => t.id === targetId);
        return target?.name || 'Unknown';
    };

    return (
        <div className="space-y-6">
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-3xl font-bold text-foreground">Sources</h1>
                    <p className="text-muted-foreground mt-1">Manage your backup sources</p>
                </div>
                <Button onClick={() => setDialogOpen(true)}>
                    <Plus className="mr-2 h-4 w-4" />
                    Add Source
                </Button>
            </div>

            {sources?.length === 0 ? (
                <Card>
                    <CardContent className="flex flex-col items-center justify-center py-12">
                        <FolderOpen className="h-12 w-12 text-muted-foreground mb-4" />
                        <p className="text-muted-foreground text-center">
                            No sources configured yet.
                            <br />
                            Add a backup source to get started.
                        </p>
                        <Button onClick={() => setDialogOpen(true)} className="mt-4">
                            <Plus className="mr-2 h-4 w-4" />
                            Add Source
                        </Button>
                    </CardContent>
                </Card>
            ) : (
                <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                    {sources?.map((source) => (
                        <Card key={source.id} className="flex flex-col">
                            <CardHeader className="flex flex-row items-start justify-between space-y-0 pb-2">
                                <div className="flex items-center gap-3 overflow-hidden">
                                    <div className="p-2 bg-primary/10 rounded-md flex-shrink-0">
                                        <FolderOpen className="h-5 w-5 text-primary" />
                                    </div>
                                    <div className="min-w-0">
                                        <CardTitle className="text-base font-semibold truncate" title={source.name}>
                                            {source.name}
                                        </CardTitle>
                                        <CardDescription className="text-xs truncate" title={source.path}>
                                            {source.path}
                                        </CardDescription>
                                    </div>
                                </div>
                                <DropdownMenu>
                                    <DropdownMenuTrigger asChild>
                                        <Button variant="ghost" className="h-8 w-8 p-0 flex-shrink-0">
                                            <MoreHorizontal className="h-4 w-4" />
                                            <span className="sr-only">Open menu</span>
                                        </Button>
                                    </DropdownMenuTrigger>
                                    <DropdownMenuContent align="end">
                                        <DropdownMenuItem onClick={() => handleEdit(source)}>
                                            <Pencil className="mr-2 h-4 w-4" />
                                            Edit
                                        </DropdownMenuItem>
                                        <DropdownMenuItem onClick={() => handleDelete(source)} className="text-destructive focus:text-destructive">
                                            <Trash2 className="mr-2 h-4 w-4" />
                                            Delete
                                        </DropdownMenuItem>
                                    </DropdownMenuContent>
                                </DropdownMenu>
                            </CardHeader>
                            <CardContent className="flex-1 py-2">
                                <div className="grid gap-2 text-sm">
                                    <div className="flex items-center justify-between">
                                        <span className="text-muted-foreground">Target</span>
                                        <Badge variant="secondary" className="font-normal">
                                            {getTargetName(source.target_id)}
                                        </Badge>
                                    </div>
                                    <div className="flex items-center justify-between">
                                        <span className="text-muted-foreground">Exclusions</span>
                                        <span className="font-medium">{source.exclusions.length}</span>
                                    </div>
                                </div>
                            </CardContent>
                            <CardFooter className="grid grid-cols-2 gap-2 pt-2">
                                <Link to={`/sources/${source.id}`} className="w-full">
                                    <Button variant="outline" className="w-full h-9">
                                        View
                                    </Button>
                                </Link>
                                <Button
                                    onClick={() => handleRunBackup(source.id, source.name)}
                                    disabled={!source.target_id || runBackup.isPending}
                                    className="w-full h-9"
                                >
                                    <Play className="mr-2 h-4 w-4" />
                                    Run
                                </Button>
                            </CardFooter>
                        </Card>
                    ))}
                </div>
            )}

            <SourceDialog
                open={dialogOpen}
                onOpenChange={(open) => {
                    setDialogOpen(open);
                    if (!open) setEditingSource(undefined);
                }}
                source={editingSource}
            />

            <AlertDialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
                <AlertDialogContent>
                    <AlertDialogHeader>
                        <AlertDialogTitle>Delete Source</AlertDialogTitle>
                        <AlertDialogDescription>
                            Are you sure you want to delete "{deletingSource?.name}"? This action cannot be undone.
                        </AlertDialogDescription>
                    </AlertDialogHeader>
                    <AlertDialogFooter>
                        <AlertDialogCancel>Cancel</AlertDialogCancel>
                        <AlertDialogAction onClick={confirmDelete} className="bg-destructive text-destructive-foreground hover:bg-destructive/90">
                            Delete
                        </AlertDialogAction>
                    </AlertDialogFooter>
                </AlertDialogContent>
            </AlertDialog>
        </div>
    );
}
