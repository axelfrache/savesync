import { Menu } from 'lucide-react';
import { useUIStore } from '@/store/ui';
import { Button } from '@/components/ui/button';

export default function Header() {
    const { toggleSidebar } = useUIStore();

    return (
        <header className="flex h-16 items-center justify-between border-b border-border bg-background px-6">
            <div className="flex items-center gap-4">
                <Button
                    variant="ghost"
                    size="icon"
                    onClick={toggleSidebar}
                    className="lg:hidden"
                >
                    <Menu className="h-5 w-5" />
                </Button>
            </div>

            <div className="flex items-center gap-4">
                {/* Future: Add user menu, notifications, etc. */}
            </div>
        </header>
    );
}
