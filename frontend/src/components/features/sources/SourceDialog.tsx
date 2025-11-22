import { useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { useCreateSource, useUpdateSource } from '@/hooks/useSources';
import { useTargets } from '@/hooks/useTargets';
import type { Source, CreateSourceRequest } from '@/types/api';
import { FilePicker } from '@/components/ui/file-picker';
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select';
import { useToast } from '@/hooks/use-toast';

interface SourceDialogProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    source?: Source;
}

export default function SourceDialog({ open, onOpenChange, source }: SourceDialogProps) {
    const { data: targets } = useTargets();
    const createSource = useCreateSource();
    const updateSource = useUpdateSource();
    const { toast } = useToast();

    const {
        register,
        handleSubmit,
        reset,
        setValue,
        watch,
        formState: { errors },
    } = useForm<CreateSourceRequest>({
        defaultValues: {
            name: '',
            path: '',
            exclusions: [],
            target_id: null,
        },
    });

    const exclusionsValue = watch('exclusions');

    useEffect(() => {
        if (source) {
            reset({
                name: source.name,
                path: source.path,
                exclusions: source.exclusions,
                target_id: source.target_id,
            });
        } else {
            reset({
                name: '',
                path: '',
                exclusions: [],
                target_id: null,
            });
        }
    }, [source, reset]);

    const onSubmit = async (data: CreateSourceRequest) => {
        try {
            if (source) {
                await updateSource.mutateAsync({ id: source.id, data });
                toast({
                    title: 'Source updated',
                    description: `${data.name} has been updated`,
                });
            } else {
                await createSource.mutateAsync(data);
                toast({
                    title: 'Source created',
                    description: `${data.name} has been created`,
                });
            }
            onOpenChange(false);
        } catch (error) {
            toast({
                title: 'Error',
                description: source ? 'Failed to update source' : 'Failed to create source',
                variant: 'destructive',
            });
        }
    };

    const handleExclusionsChange = (value: string) => {
        const exclusions = value.split(',').map((s) => s.trim()).filter(Boolean);
        setValue('exclusions', exclusions);
    };

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="bg-card border-border text-foreground sm:max-w-[500px]">
                <DialogHeader>
                    <DialogTitle>{source ? 'Edit Source' : 'Add Source'}</DialogTitle>
                    <DialogDescription className="text-muted-foreground">
                        {source ? 'Update the source configuration' : 'Configure a new backup source'}
                    </DialogDescription>
                </DialogHeader>

                <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
                    <div className="space-y-2">
                        <Label htmlFor="name">Name</Label>
                        <Input
                            id="name"
                            {...register('name', { required: 'Name is required' })}
                            placeholder="My Documents"
                            className="bg-input border-input"
                        />
                        {errors.name && (
                            <p className="text-sm text-red-400">{errors.name.message}</p>
                        )}
                    </div>

                    <div className="space-y-2">
                        <Label htmlFor="path">Path</Label>
                        <FilePicker
                            value={watch('path')}
                            onChange={(value) => setValue('path', value, { shouldValidate: true })}
                            placeholder="/home/user/documents"
                        />
                        {errors.path && (
                            <p className="text-sm text-red-400">{errors.path.message}</p>
                        )}
                    </div>

                    <div className="space-y-2">
                        <Label htmlFor="target_id">Target</Label>
                        <Select
                            value={watch('target_id')?.toString() || ''}
                            onValueChange={(value) => {
                                const id = parseInt(value);
                                setValue('target_id', isNaN(id) ? null : id);
                            }}
                        >
                            <SelectTrigger className="bg-input border-input">
                                <SelectValue placeholder="Select a target" />
                            </SelectTrigger>
                            <SelectContent className="bg-popover border-border">
                                {targets?.map((target) => (
                                    <SelectItem key={target.id} value={target.id.toString()}>
                                        {target.name} ({target.type})
                                    </SelectItem>
                                ))}
                            </SelectContent>
                        </Select>
                    </div>

                    <div className="space-y-2">
                        <Label htmlFor="exclusions">Exclusions (comma-separated)</Label>
                        <Input
                            id="exclusions"
                            value={exclusionsValue.join(', ')}
                            onChange={(e) => handleExclusionsChange(e.target.value)}
                            placeholder="*.tmp, *.log, node_modules"
                            className="bg-input border-input"
                        />
                        <p className="text-xs text-slate-500">
                            Glob patterns to exclude from backup
                        </p>
                    </div>

                    <DialogFooter>
                        <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
                            Cancel
                        </Button>
                        <Button
                            type="submit"
                            disabled={createSource.isPending || updateSource.isPending}
                        >
                            {source ? 'Update' : 'Create'}
                        </Button>
                    </DialogFooter>
                </form>
            </DialogContent>
        </Dialog>
    );
}
