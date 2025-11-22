import { useUIStore } from '@/store/ui';
import { useAuthStore } from '@/store/authStore';
import { useSettings } from '@/hooks/useSettings';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Label } from '@/components/ui/label';
import { Switch } from '@/components/ui/switch';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Shield, Settings as SettingsIcon, Users } from 'lucide-react';
import UserManagement from '@/components/features/settings/UserManagement';

export default function Settings() {
    const { theme, setTheme } = useUIStore();
    const { user } = useAuthStore();
    const { isRegistrationEnabled, updateSetting } = useSettings();

    const handleToggleRegistration = (checked: boolean) => {
        updateSetting.mutate({
            key: 'registration_enabled',
            value: checked ? 'true' : 'false',
        });
    };

    return (
        <div className="space-y-6">
            <div>
                <h1 className="text-3xl font-bold text-foreground">Settings</h1>
                <p className="text-muted-foreground mt-1">Manage application settings</p>
            </div>

            <Tabs defaultValue="general" className="space-y-4">
                <TabsList>
                    <TabsTrigger value="general" className="flex items-center gap-2">
                        <SettingsIcon className="h-4 w-4" />
                        General
                    </TabsTrigger>
                    {user?.is_admin && (
                        <>
                            <TabsTrigger value="security" className="flex items-center gap-2">
                                <Shield className="h-4 w-4" />
                                Security
                            </TabsTrigger>
                            <TabsTrigger value="users" className="flex items-center gap-2">
                                <Users className="h-4 w-4" />
                                Users
                            </TabsTrigger>
                        </>
                    )}
                </TabsList>

                <TabsContent value="general" className="space-y-4">
                    <Card>
                        <CardHeader>
                            <CardTitle>Appearance</CardTitle>
                            <CardDescription>
                                Customize the look and feel of the application
                            </CardDescription>
                        </CardHeader>
                        <CardContent className="space-y-4">
                            <div className="flex items-center justify-between">
                                <div className="space-y-0.5">
                                    <Label htmlFor="dark-mode">Dark Mode</Label>
                                    <p className="text-sm text-muted-foreground">
                                        Use dark theme for the interface
                                    </p>
                                </div>
                                <Switch
                                    id="dark-mode"
                                    checked={theme === 'dark'}
                                    onCheckedChange={(checked) => setTheme(checked ? 'dark' : 'light')}
                                />
                            </div>
                        </CardContent>
                    </Card>

                    <Card>
                        <CardHeader>
                            <CardTitle>About</CardTitle>
                            <CardDescription>
                                Application information
                            </CardDescription>
                        </CardHeader>
                        <CardContent className="space-y-2 text-sm">
                            <div>
                                <span className="text-muted-foreground">Version:</span>
                                <span className="ml-2 text-foreground">0.1.0</span>
                            </div>
                            <div>
                                <span className="text-muted-foreground">API Endpoint:</span>
                                <span className="ml-2 text-foreground">
                                    {import.meta.env.VITE_API_URL || 'http://localhost:8080'}
                                </span>
                            </div>
                        </CardContent>
                    </Card>
                </TabsContent>

                {user?.is_admin && (
                    <>
                        <TabsContent value="security" className="space-y-4">
                            <Card>
                                <CardHeader>
                                    <CardTitle>Security Settings</CardTitle>
                                    <CardDescription>Control access and registration</CardDescription>
                                </CardHeader>
                                <CardContent className="space-y-4">
                                    <div className="flex items-center justify-between">
                                        <div className="space-y-0.5">
                                            <Label htmlFor="registration">Enable User Registration</Label>
                                            <p className="text-sm text-muted-foreground">
                                                Allow new users to create accounts
                                            </p>
                                        </div>
                                        <Switch
                                            id="registration"
                                            checked={isRegistrationEnabled}
                                            onCheckedChange={handleToggleRegistration}
                                        />
                                    </div>
                                </CardContent>
                            </Card>
                        </TabsContent>

                        <TabsContent value="users" className="space-y-4">
                            <UserManagement />
                        </TabsContent>
                    </>
                )}
            </Tabs>
        </div>
    );
}
