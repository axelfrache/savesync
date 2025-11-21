import { useUIStore } from '@/store/ui';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Label } from '@/components/ui/label';
import { Switch } from '@/components/ui/switch';

export default function Settings() {
    const { theme, setTheme } = useUIStore();

    return (
        <div className="space-y-6">
            <div>
                <h1 className="text-3xl font-bold text-foreground">Settings</h1>
                <p className="text-muted-foreground mt-1">Manage application settings</p>
            </div>

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
                            {import.meta.env.VITE_API_URL || 'http://localhost:8080/api'}
                        </span>
                    </div>
                </CardContent>
            </Card>

            <Card>
                <CardHeader>
                    <CardTitle>Notifications</CardTitle>
                    <CardDescription>
                        Configure notification preferences
                    </CardDescription>
                </CardHeader>
                <CardContent>
                    <p className="text-sm text-muted-foreground">
                        Notification settings coming soon
                    </p>
                </CardContent>
            </Card>
        </div>
    );
}
