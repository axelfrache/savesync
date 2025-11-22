import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { adminApi } from '@/lib/api';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { Trash2, Shield, ShieldOff, UserPlus } from 'lucide-react';
import { useToast } from '@/hooks/use-toast';
import { useAuthStore } from '@/store/authStore';
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';

export default function UserManagement() {
    const { toast } = useToast();
    const queryClient = useQueryClient();
    const { user: currentUser } = useAuthStore();
    const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
    const [newUserEmail, setNewUserEmail] = useState('');
    const [newUserPassword, setNewUserPassword] = useState('');

    const { data: users, isLoading } = useQuery({
        queryKey: ['users'],
        queryFn: adminApi.listUsers,
    });

    const toggleAdminMutation = useMutation({
        mutationFn: ({ id, isAdmin }: { id: number; isAdmin: boolean }) =>
            adminApi.toggleAdmin(id, isAdmin),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['users'] });
            toast({ title: 'User updated successfully' });
        },
        onError: (error: any) => {
            toast({
                title: 'Failed to update user',
                description: error.response?.data?.error?.message || 'Unknown error',
                variant: 'destructive',
            });
        },
    });

    const deleteUserMutation = useMutation({
        mutationFn: adminApi.deleteUser,
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['users'] });
            toast({ title: 'User deleted successfully' });
        },
        onError: (error: any) => {
            toast({
                title: 'Failed to delete user',
                description: error.response?.data?.error?.message || 'Unknown error',
                variant: 'destructive',
            });
        },
    });

    const createUserMutation = useMutation({
        mutationFn: (data: any) => adminApi.createUser(data),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['users'] });
            setIsCreateDialogOpen(false);
            setNewUserEmail('');
            setNewUserPassword('');
            toast({ title: 'User created successfully' });
        },
        onError: (error: any) => {
            toast({
                title: 'Failed to create user',
                description: error.response?.data?.error?.message || 'Unknown error',
                variant: 'destructive',
            });
        },
    });

    const handleCreateUser = (e: React.FormEvent) => {
        e.preventDefault();
        createUserMutation.mutate({ email: newUserEmail, password: newUserPassword });
    };

    if (isLoading) {
        return <div>Loading users...</div>;
    }

    return (
        <Card>
            <CardHeader className="flex flex-row items-center justify-between">
                <div>
                    <CardTitle>User Management</CardTitle>
                    <CardDescription>Manage users and their roles</CardDescription>
                </div>
                <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
                    <DialogTrigger asChild>
                        <Button className="flex items-center gap-2">
                            <UserPlus className="h-4 w-4" />
                            Add User
                        </Button>
                    </DialogTrigger>
                    <DialogContent>
                        <DialogHeader>
                            <DialogTitle>Add New User</DialogTitle>
                            <DialogDescription>
                                Create a new user account manually.
                            </DialogDescription>
                        </DialogHeader>
                        <form onSubmit={handleCreateUser} className="space-y-4">
                            <div className="space-y-2">
                                <Label htmlFor="email">Email</Label>
                                <Input
                                    id="email"
                                    type="email"
                                    value={newUserEmail}
                                    onChange={(e) => setNewUserEmail(e.target.value)}
                                    required
                                />
                            </div>
                            <div className="space-y-2">
                                <Label htmlFor="password">Password</Label>
                                <Input
                                    id="password"
                                    type="password"
                                    value={newUserPassword}
                                    onChange={(e) => setNewUserPassword(e.target.value)}
                                    required
                                />
                            </div>
                            <DialogFooter>
                                <Button type="submit" disabled={createUserMutation.isPending}>
                                    {createUserMutation.isPending ? 'Creating...' : 'Create User'}
                                </Button>
                            </DialogFooter>
                        </form>
                    </DialogContent>
                </Dialog>
            </CardHeader>
            <CardContent>
                <Table>
                    <TableHeader>
                        <TableRow>
                            <TableHead>ID</TableHead>
                            <TableHead>Email</TableHead>
                            <TableHead>Role</TableHead>
                            <TableHead className="text-right">Actions</TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {users?.map((user) => (
                            <TableRow key={user.id}>
                                <TableCell>{user.id}</TableCell>
                                <TableCell>{user.email}</TableCell>
                                <TableCell>
                                    {user.is_admin ? (
                                        <Badge variant="default">Admin</Badge>
                                    ) : (
                                        <Badge variant="secondary">User</Badge>
                                    )}
                                </TableCell>
                                <TableCell className="text-right space-x-2">
                                    {user.id !== currentUser?.id && (
                                        <>
                                            <Button
                                                variant="ghost"
                                                size="icon"
                                                onClick={() =>
                                                    toggleAdminMutation.mutate({
                                                        id: user.id,
                                                        isAdmin: !user.is_admin,
                                                    })
                                                }
                                                title={user.is_admin ? 'Remove Admin' : 'Make Admin'}
                                            >
                                                {user.is_admin ? (
                                                    <ShieldOff className="h-4 w-4 text-orange-500" />
                                                ) : (
                                                    <Shield className="h-4 w-4 text-green-500" />
                                                )}
                                            </Button>
                                            <Button
                                                variant="ghost"
                                                size="icon"
                                                onClick={() => {
                                                    if (
                                                        confirm(
                                                            'Are you sure you want to delete this user?'
                                                        )
                                                    ) {
                                                        deleteUserMutation.mutate(user.id);
                                                    }
                                                }}
                                                title="Delete User"
                                            >
                                                <Trash2 className="h-4 w-4 text-destructive" />
                                            </Button>
                                        </>
                                    )}
                                </TableCell>
                            </TableRow>
                        ))}
                    </TableBody>
                </Table>
            </CardContent>
        </Card>
    );
}
