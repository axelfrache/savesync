import { useState, useEffect } from 'react';
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Folder, File, ChevronUp, Loader2 } from 'lucide-react';
import axios from 'axios';

interface FileEntry {
    name: string;
    path: string;
    is_dir: boolean;
}

interface FileExplorerProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    onSelect: (path: string) => void;
    initialPath?: string;
}

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

export function FileExplorer({ open, onOpenChange, onSelect, initialPath = '/' }: FileExplorerProps) {
    const [currentPath, setCurrentPath] = useState(initialPath);
    const [entries, setEntries] = useState<FileEntry[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        if (open) {
            loadDirectory(currentPath);
        }
    }, [open, currentPath]);

    const loadDirectory = async (path: string) => {
        setLoading(true);
        setError(null);
        try {
            const response = await axios.get(`${API_BASE_URL}/system/files`, {
                params: { path },
            });
            setEntries(response.data.data.entries);
            setCurrentPath(response.data.data.current_path);
        } catch (err) {
            setError('Failed to load directory');
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    const handleEntryClick = (entry: FileEntry) => {
        if (entry.is_dir) {
            loadDirectory(entry.path);
        }
    };

    const handleSelectCurrent = () => {
        onSelect(currentPath);
        onOpenChange(false);
    };

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="bg-card border-border text-foreground sm:max-w-[600px] h-[500px] flex flex-col">
                <DialogHeader>
                    <DialogTitle>Select Directory</DialogTitle>
                </DialogHeader>

                <div className="flex items-center space-x-2 mb-2">
                    <Input
                        value={currentPath}
                        onChange={(e) => setCurrentPath(e.target.value)}
                        onKeyDown={(e) => e.key === 'Enter' && loadDirectory(currentPath)}
                        className="bg-input border-input"
                    />
                    <Button variant="outline" size="icon" onClick={() => loadDirectory(currentPath)}>
                        <Loader2 className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
                    </Button>
                </div>

                {error && <div className="text-destructive text-sm mb-2">{error}</div>}

                <ScrollArea className="flex-1 overflow-y-auto border border-border rounded-md p-2 bg-muted/50 h-[350px]">
                    <div className="space-y-1">
                        {entries.map((entry) => (
                            <div
                                key={entry.path}
                                className="flex items-center space-x-2 p-2 hover:bg-accent hover:text-accent-foreground rounded-md cursor-pointer"
                                onClick={() => handleEntryClick(entry)}
                            >
                                {entry.name === '..' ? (
                                    <ChevronUp className="h-4 w-4 text-muted-foreground" />
                                ) : entry.is_dir ? (
                                    <Folder className="h-4 w-4 text-primary" />
                                ) : (
                                    <File className="h-4 w-4 text-muted-foreground" />
                                )}
                                <span className={entry.name === '..' ? 'text-muted-foreground' : ''}>
                                    {entry.name}
                                </span>
                            </div>
                        ))}
                        {entries.length === 0 && !loading && (
                            <div className="text-muted-foreground text-center py-4">Empty directory</div>
                        )}
                    </div>
                </ScrollArea>

                <div className="flex justify-end space-x-2 mt-4">
                    <Button variant="outline" onClick={() => onOpenChange(false)}>
                        Cancel
                    </Button>
                    <Button onClick={handleSelectCurrent}>
                        Select Current Directory
                    </Button>
                </div>
            </DialogContent>
        </Dialog>
    );
}
