import { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import { useCreateTarget, useUpdateTarget } from '@/hooks/useTargets';
import type { Target, CreateTargetRequest } from '@/types/api';
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

interface TargetDialogProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    target?: Target;
}

export default function TargetDialog({ open, onOpenChange, target }: TargetDialogProps) {
    const createTarget = useCreateTarget();
    const updateTarget = useUpdateTarget();
    const { toast } = useToast();
    const [targetType, setTargetType] = useState<'local' | 's3' | 'sftp'>('local');

    const {
        register,
        handleSubmit,
        reset,
        setValue,
        formState: { errors },
    } = useForm<CreateTargetRequest>({
        defaultValues: {
            name: '',
            type: 'local',
            config: {},
        },
    });

    useEffect(() => {
        if (target) {
            reset({
                name: target.name,
                type: target.type,
                config: target.config,
            });
            setTargetType(target.type);
        } else {
            reset({
                name: '',
                type: 'local',
                config: {},
            });
            setTargetType('local');
        }
    }, [target, reset]);

    const onSubmit = async (data: CreateTargetRequest) => {
        try {
            if (target) {
                await updateTarget.mutateAsync({ id: target.id, data });
                toast({
                    title: 'Target updated',
                    description: `${data.name} has been updated`,
                });
            } else {
                await createTarget.mutateAsync(data);
                toast({
                    title: 'Target created',
                    description: `${data.name} has been created`,
                });
            }
            onOpenChange(false);
        } catch (error) {
            toast({
                title: 'Error',
                description: target ? 'Failed to update target' : 'Failed to create target',
                variant: 'destructive',
            });
        }
    };

    const handleTypeChange = (type: 'local' | 's3' | 'sftp') => {
        setTargetType(type);
        setValue('type', type);
        setValue('config', {});
    };

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="bg-slate-900 border-slate-800 text-white sm:max-w-[500px]">
                <DialogHeader>
                    <DialogTitle>{target ? 'Edit Target' : 'Add Target'}</DialogTitle>
                    <DialogDescription className="text-slate-400">
                        {target ? 'Update the target configuration' : 'Configure a new storage target'}
                    </DialogDescription>
                </DialogHeader>

                <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
                    <div className="space-y-2">
                        <Label htmlFor="name">Name</Label>
                        <Input
                            id="name"
                            {...register('name', { required: 'Name is required' })}
                            placeholder="My Backup Storage"
                            className="bg-slate-800 border-slate-700"
                        />
                        {errors.name && (
                            <p className="text-sm text-red-400">{errors.name.message}</p>
                        )}
                    </div>

                    <div className="space-y-2">
                        <Label htmlFor="type">Type</Label>
                        <Select value={targetType} onValueChange={handleTypeChange}>
                            <SelectTrigger className="bg-slate-800 border-slate-700">
                                <SelectValue />
                            </SelectTrigger>
                            <SelectContent className="bg-slate-800 border-slate-700">
                                <SelectItem value="local">Local Filesystem</SelectItem>
                                <SelectItem value="s3">S3 Compatible</SelectItem>
                                <SelectItem value="sftp">SFTP</SelectItem>
                            </SelectContent>
                        </Select>
                    </div>

                    {targetType === 'local' && (
                        <div className="space-y-2">
                            <Label htmlFor="path">Path</Label>
                            <Input
                                id="path"
                                {...register('config.path', { required: 'Path is required' })}
                                placeholder="/backups"
                                className="bg-slate-800 border-slate-700"
                            />
                        </div>
                    )}

                    {targetType === 's3' && (
                        <>
                            <div className="space-y-2">
                                <Label htmlFor="bucket">Bucket</Label>
                                <Input
                                    id="bucket"
                                    {...register('config.bucket', { required: 'Bucket is required' })}
                                    placeholder="my-backups"
                                    className="bg-slate-800 border-slate-700"
                                />
                            </div>
                            <div className="space-y-2">
                                <Label htmlFor="region">Region</Label>
                                <Input
                                    id="region"
                                    {...register('config.region')}
                                    placeholder="us-east-1"
                                    className="bg-slate-800 border-slate-700"
                                />
                            </div>
                            <div className="space-y-2">
                                <Label htmlFor="access_key">Access Key</Label>
                                <Input
                                    id="access_key"
                                    {...register('config.access_key')}
                                    placeholder="YOUR_ACCESS_KEY"
                                    className="bg-slate-800 border-slate-700"
                                />
                            </div>
                            <div className="space-y-2">
                                <Label htmlFor="secret_key">Secret Key</Label>
                                <Input
                                    id="secret_key"
                                    type="password"
                                    {...register('config.secret_key')}
                                    placeholder="YOUR_SECRET_KEY"
                                    className="bg-slate-800 border-slate-700"
                                />
                            </div>
                        </>
                    )}

                    {targetType === 'sftp' && (
                        <>
                            <div className="space-y-2">
                                <Label htmlFor="host">Host</Label>
                                <Input
                                    id="host"
                                    {...register('config.host', { required: 'Host is required' })}
                                    placeholder="backup.example.com"
                                    className="bg-slate-800 border-slate-700"
                                />
                            </div>
                            <div className="space-y-2">
                                <Label htmlFor="user">User</Label>
                                <Input
                                    id="user"
                                    {...register('config.user', { required: 'User is required' })}
                                    placeholder="backup-user"
                                    className="bg-slate-800 border-slate-700"
                                />
                            </div>
                            <div className="space-y-2">
                                <Label htmlFor="password">Password</Label>
                                <Input
                                    id="password"
                                    type="password"
                                    {...register('config.password')}
                                    placeholder="password"
                                    className="bg-slate-800 border-slate-700"
                                />
                            </div>
                            <div className="space-y-2">
                                <Label htmlFor="sftp_path">Path</Label>
                                <Input
                                    id="sftp_path"
                                    {...register('config.path', { required: 'Path is required' })}
                                    placeholder="/backups"
                                    className="bg-slate-800 border-slate-700"
                                />
                            </div>
                        </>
                    )}

                    <DialogFooter>
                        <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
                            Cancel
                        </Button>
                        <Button
                            type="submit"
                            disabled={createTarget.isPending || updateTarget.isPending}
                        >
                            {target ? 'Update' : 'Create'}
                        </Button>
                    </DialogFooter>
                </form>
            </DialogContent>
        </Dialog>
    );
}
