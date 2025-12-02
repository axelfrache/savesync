import { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import { useCreateTarget, useUpdateTarget } from '@/hooks/useTargets';
import type { Target, CreateTargetRequest, TargetType } from '@/types/api';
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
import { Checkbox } from '@/components/ui/checkbox';
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

const AWS_REGIONS = [
    { value: 'us-east-1', label: 'US East (N. Virginia)' },
    { value: 'us-east-2', label: 'US East (Ohio)' },
    { value: 'us-west-1', label: 'US West (N. California)' },
    { value: 'us-west-2', label: 'US West (Oregon)' },
    { value: 'eu-west-1', label: 'EU (Ireland)' },
    { value: 'eu-central-1', label: 'EU (Frankfurt)' },
    { value: 'ap-southeast-1', label: 'Asia Pacific (Singapore)' },
    { value: 'ap-northeast-1', label: 'Asia Pacific (Tokyo)' },
];

export default function TargetDialog({ open, onOpenChange, target }: TargetDialogProps) {
    const createTarget = useCreateTarget();
    const updateTarget = useUpdateTarget();
    const { toast } = useToast();
    const [targetType, setTargetType] = useState<TargetType>('s3_generic');
    const [pathStyle, setPathStyle] = useState(true);
    const [useTLS, setUseTLS] = useState(true);

    const {
        register,
        handleSubmit,
        reset,
        setValue,
        watch,
        formState: { errors },
    } = useForm<CreateTargetRequest>({
        defaultValues: {
            name: '',
            type: 's3_generic',
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
            if (target.type === 's3_generic' && typeof target.config === 'object' && target.config !== null && 'path_style' in target.config) {
                setPathStyle(target.config.path_style ?? true);
                setUseTLS(target.config.use_tls ?? true);
            }
        } else {
            reset({
                name: '',
                type: 's3_generic',
                config: {},
            });
            setTargetType('s3_generic');
            setPathStyle(true);
            setUseTLS(true);
        }
    }, [target, reset]);

    const onSubmit = async (data: CreateTargetRequest) => {
        try {
            // Build config based on type
            let config: any = {};

            if (targetType === 's3_generic') {
                config = {
                    endpoint: (data.config as any).endpoint,
                    bucket: (data.config as any).bucket,
                    region: (data.config as any).region || 'us-east-1',
                    access_key: (data.config as any).access_key,
                    secret_key: (data.config as any).secret_key,
                    path_style: pathStyle,
                    use_tls: useTLS,
                };
            } else if (targetType === 's3_aws') {
                config = {
                    bucket: (data.config as any).bucket,
                    region: (data.config as any).region,
                    access_key: (data.config as any).access_key,
                    secret_key: (data.config as any).secret_key,
                };
            } else if (targetType === 'local') {
                config = {
                    path: (data.config as any).path,
                };
            } else if (targetType === 'sftp') {
                config = {
                    host: (data.config as any).host,
                    user: (data.config as any).user,
                    password: (data.config as any).password,
                    path: (data.config as any).path,
                };
            }

            const payload = {
                name: data.name,
                type: targetType,
                config,
            };

            if (target) {
                await updateTarget.mutateAsync({ id: target.id, data: payload });
                toast({
                    title: 'Target updated',
                    description: `${data.name} has been updated`,
                });
            } else {
                await createTarget.mutateAsync(payload);
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

    const handleTypeChange = (type: TargetType) => {
        setTargetType(type);
        setValue('type', type);
        setValue('config', {});
        if (type === 's3_generic') {
            setPathStyle(true);
            setUseTLS(true);
        }
    };

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="bg-card border-border text-foreground sm:max-w-[500px]">
                <DialogHeader>
                    <DialogTitle>{target ? 'Edit Target' : 'Add Target'}</DialogTitle>
                    <DialogDescription className="text-muted-foreground">
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
                            className="bg-input border-input"
                        />
                        {errors.name && (
                            <p className="text-sm text-red-400">{errors.name.message}</p>
                        )}
                    </div>

                    <div className="space-y-2">
                        <Label htmlFor="type">Type</Label>
                        <Select value={targetType} onValueChange={handleTypeChange}>
                            <SelectTrigger className="bg-input border-input">
                                <SelectValue />
                            </SelectTrigger>
                            <SelectContent className="bg-popover border-border">
                                <SelectItem value="local">Local Filesystem</SelectItem>
                                <SelectItem value="s3_generic">
                                    S3 Compatible (MinIO, Garage, Ceph, R2…)
                                </SelectItem>
                                <SelectItem value="s3_aws">AWS S3</SelectItem>
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
                                className="bg-input border-input"
                            />
                        </div>
                    )}

                    {targetType === 's3_generic' && (
                        <>
                            <div className="space-y-2">
                                <Label htmlFor="endpoint">Endpoint</Label>
                                <Input
                                    id="endpoint"
                                    {...register('config.endpoint', { required: 'Endpoint is required' })}
                                    placeholder="https://minio.example.com"
                                    className="bg-input border-input"
                                />
                            </div>
                            <div className="space-y-2">
                                <Label htmlFor="bucket">Bucket</Label>
                                <Input
                                    id="bucket"
                                    {...register('config.bucket', { required: 'Bucket is required' })}
                                    placeholder="savesync"
                                    className="bg-input border-input"
                                />
                            </div>
                            <div className="space-y-2">
                                <Label htmlFor="region">Region (optional)</Label>
                                <Input
                                    id="region"
                                    {...register('config.region')}
                                    placeholder="us-east-1"
                                    className="bg-input border-input"
                                />
                            </div>
                            <div className="space-y-2">
                                <Label htmlFor="access_key">Access Key</Label>
                                <Input
                                    id="access_key"
                                    {...register('config.access_key', { required: 'Access key is required' })}
                                    placeholder="YOUR_ACCESS_KEY"
                                    className="bg-input border-input"
                                />
                            </div>
                            <div className="space-y-2">
                                <Label htmlFor="secret_key">Secret Key</Label>
                                <Input
                                    id="secret_key"
                                    type="password"
                                    {...register('config.secret_Key', { required: 'Secret key is required' })}
                                    placeholder="YOUR_SECRET_KEY"
                                    className="bg-input border-input"
                                />
                            </div>
                            <div className="flex items-center space-x-2">
                                <Checkbox
                                    id="path_style"
                                    checked={pathStyle}
                                    onCheckedChange={(checked) => setPathStyle(checked as boolean)}
                                />
                                <Label htmlFor="path_style" className="font-normal">
                                    Use path-style URLs (recommended for MinIO/Garage)
                                </Label>
                            </div>
                            <div className="flex items-center space-x-2">
                                <Checkbox
                                    id="use_tls"
                                    checked={useTLS}
                                    onCheckedChange={(checked) => setUseTLS(checked as boolean)}
                                />
                                <Label htmlFor="use_tls" className="font-normal">
                                    Use TLS/HTTPS
                                </Label>
                            </div>
                        </>
                    )}

                    {targetType === 's3_aws' && (
                        <>
                            <div className="space-y-2">
                                <Label htmlFor="aws_bucket">Bucket</Label>
                                <Input
                                    id="aws_bucket"
                                    {...register('config.bucket', { required: 'Bucket is required' })}
                                    placeholder="my-backup-bucket"
                                    className="bg-input border-input"
                                />
                            </div>
                            <div className="space-y-2">
                                <Label htmlFor="aws_region">Region</Label>
                                <Select
                                    value={watch('config.region') as string}
                                    onValueChange={(value) => setValue('config.region', value)}
                                >
                                    <SelectTrigger className="bg-input border-input">
                                        <SelectValue placeholder="Select region" />
                                    </SelectTrigger>
                                    <SelectContent className="bg-popover border-border">
                                        {AWS_REGIONS.map((region) => (
                                            <SelectItem key={region.value} value={region.value}>
                                                {region.label}
                                            </SelectItem>
                                        ))}
                                    </SelectContent>
                                </Select>
                            </div>
                            <div className="space-y-2">
                                <Label htmlFor="aws_access_key">Access Key</Label>
                                <Input
                                    id="aws_access_key"
                                    {...register('config.access_key', { required: 'Access key is required' })}
                                    placeholder="YOUR_ACCESS_KEY"
                                    className="bg-input border-input"
                                />
                            </div>
                            <div className="space-y-2">
                                <Label htmlFor="aws_secret_key">Secret Key</Label>
                                <Input
                                    id="aws_secret_key"
                                    type="password"
                                    {...register('config.secret_key', { required: 'Secret key is required' })}
                                    placeholder="YOUR_SECRET_KEY"
                                    className="bg-input border-input"
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
                                    className="bg-input border-input"
                                />
                            </div>
                            <div className="space-y-2">
                                <Label htmlFor="user">User</Label>
                                <Input
                                    id="user"
                                    {...register('config.user', { required: 'User is required' })}
                                    placeholder="backup-user"
                                    className="bg-input border-input"
                                />
                            </div>
                            <div className="space-y-2">
                                <Label htmlFor="password">Password</Label>
                                <Input
                                    id="password"
                                    type="password"
                                    {...register('config.password')}
                                    placeholder="password"
                                    className="bg-input border-input"
                                />
                            </div>
                            <div className="space-y-2">
                                <Label htmlFor="sftp_path">Path</Label>
                                <Input
                                    id="sftp_path"
                                    {...register('config.path', { required: 'Path is required' })}
                                    placeholder="/backups"
                                    className="bg-input border-input"
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
