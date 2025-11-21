import { useState } from 'react';
import { useTargets, useDeleteTarget } from '@/hooks/useTargets';
import LoadingState from '@/components/shared/LoadingState';
import ErrorState from '@/components/shared/ErrorState';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Plus, Pencil, Trash2, Target as TargetIcon } from 'lucide-react';
import TargetDialog from '@/components/features/targets/TargetDialog';
import type { Target } from '@/types/api';
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

export default function Targets() {
    const { data: targets, isLoading, error } = useTargets();
    const deleteTarget = useDeleteTarget();
    const { toast } = useToast();

    const [dialogOpen, setDialogOpen] = useState(false);
    const [editingTarget, setEditingTarget] = useState<Target | undefined>();
    const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
    const [deletingTarget, setDeletingTarget] = useState<Target | undefined>();

    if (isLoading) return <LoadingState />;
    if (error) return <ErrorState message="Failed to load targets" />;

    const handleEdit = (target: Target) => {
        setEditingTarget(target);
        setDialogOpen(true);
    };

    const handleDelete = (target: Target) => {
        setDeletingTarget(target);
        setDeleteDialogOpen(true);
    };

    const confirmDelete = async () => {
        if (!deletingTarget) return;
        try {
            await deleteTarget.mutateAsync(deletingTarget.id);
            toast({
                title: 'Target deleted',
                description: `${deletingTarget.name} has been deleted`,
            });
        } catch (error) {
            toast({
                title: 'Error',
                description: 'Failed to delete target',
                variant: 'destructive',
            });
        } finally {
            setDeleteDialogOpen(false);
            setDeletingTarget(undefined);
        }
    };

    const getTypeLabel = (type: string) => {
        const labels: Record<string, string> = {
            local: 'Local Filesystem',
            s3: 'S3 Compatible',
            sftp: 'SFTP',
        };
        return labels[type] || type;
    };

    return (
        <div className="space-y-6">
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-3xl font-bold text-foreground">Targets</h1>
                    <p className="text-muted-foreground mt-1">Manage backup storage targets</p>
                </div>
                <Button onClick={() => setDialogOpen(true)}>
                    <Plus className="mr-2 h-4 w-4" />
                    Add Target
                </Button>
            </div>

            {targets?.length === 0 ? (
                <Card>
                    <CardContent className="flex flex-col items-center justify-center py-12">
                        <TargetIcon className="h-12 w-12 text-muted-foreground mb-4" />
                        <p className="text-muted-foreground text-center">
                            No targets configured yet.
                            <br />
                            Add a storage target to start backing up.
                        </p>
                        <Button onClick={() => setDialogOpen(true)} className="mt-4">
                            <Plus className="mr-2 h-4 w-4" />
                            Add Target
                        </Button>
                    </CardContent>
                </Card>
            ) : (
                <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                    {targets?.map((target) => (
                        <Card key={target.id}>
                            <CardHeader>
                                <CardTitle className="text-foreground flex items-center justify-between">
                                    <span>{target.name}</span>
                                    <TargetIcon className="h-5 w-5 text-muted-foreground" />
                                </CardTitle>
                                <CardDescription>
                                    {getTypeLabel(target.type)}
                                </CardDescription>
                            </CardHeader>
                            <CardContent className="space-y-3">
                                <div className="text-sm text-muted-foreground">
                                    {target.type === 'local' && (
                                        <p>Path: <span className="text-foreground">{target.config.path}</span></p>
                                    )}
                                    {target.type === 's3' && (
                                        <>
                                            <p>Bucket: <span className="text-foreground">{target.config.bucket}</span></p>
                                            <p>Region: <span className="text-foreground">{target.config.region}</span></p>
                                        </>
                                    )}
                                    {target.type === 'sftp' && (
                                        <>
                                            <p>Host: <span className="text-foreground">{target.config.host}</span></p>
                                            <p>Path: <span className="text-foreground">{target.config.path}</span></p>
                                        </>
                                    )}
                                </div>
                                <div className="flex gap-2">
                                    <Button
                                        size="sm"
                                        variant="outline"
                                        className="flex-1"
                                        onClick={() => handleEdit(target)}
                                    >
                                        <Pencil className="mr-1 h-3 w-3" />
                                        Edit
                                    </Button>
                                    <Button
                                        size="sm"
                                        variant="ghost"
                                        onClick={() => handleDelete(target)}
                                    >
                                        <Trash2 className="h-3 w-3" />
                                    </Button>
                                </div>
                            </CardContent>
                        </Card>
                    ))}
                </div>
            )}

            <TargetDialog
                open={dialogOpen}
                onOpenChange={(open) => {
                    setDialogOpen(open);
                    if (!open) setEditingTarget(undefined);
                }}
                target={editingTarget}
            />

            <AlertDialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
                <AlertDialogContent>
                    <AlertDialogHeader>
                        <AlertDialogTitle>Delete Target</AlertDialogTitle>
                        <AlertDialogDescription>
                            Are you sure you want to delete "{deletingTarget?.name}"? This action cannot be undone.
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
