import { useState } from 'react';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { FolderOpen } from 'lucide-react';
import { FileExplorer } from './file-explorer';

interface FilePickerProps {
    value?: string;
    onChange: (value: string) => void;
    placeholder?: string;
    className?: string;
}

export function FilePicker({ value, onChange, placeholder, className }: FilePickerProps) {
    const [explorerOpen, setExplorerOpen] = useState(false);

    return (
        <div className={`flex w-full items-center space-x-2 ${className}`}>
            <Input
                type="text"
                value={value}
                onChange={(e) => onChange(e.target.value)}
                placeholder={placeholder}
                className="bg-input border-input text-foreground"
            />
            <Button
                type="button"
                variant="outline"
                size="icon"
                className="shrink-0"
                onClick={() => setExplorerOpen(true)}
            >
                <FolderOpen className="h-4 w-4" />
            </Button>

            <FileExplorer
                open={explorerOpen}
                onOpenChange={setExplorerOpen}
                onSelect={onChange}
                initialPath={value || '/'}
            />
        </div>
    );
}
