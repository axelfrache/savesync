import { useState } from 'react';
import { ChevronRight, ChevronDown, Folder, FolderOpen, File } from 'lucide-react';
import type { FileNode } from '@/hooks/useSnapshotFiles';

interface SnapshotFileTreeProps {
    root: FileNode;
}

interface TreeNodeProps {
    node: FileNode;
    level: number;
}

function TreeNode({ node, level }: TreeNodeProps) {
    const [isExpanded, setIsExpanded] = useState(level === 0);

    const handleToggle = () => {
        if (node.is_dir && node.children && node.children.length > 0) {
            setIsExpanded(!isExpanded);
        }
    };

    const formatSize = (bytes: number) => {
        const units = ['B', 'KB', 'MB', 'GB'];
        let size = bytes;
        let unitIndex = 0;
        while (size >= 1024 && unitIndex < units.length - 1) {
            size /= 1024;
            unitIndex++;
        }
        return `${size.toFixed(size < 10 && unitIndex > 0 ? 1 : 0)} ${units[unitIndex]}`;
    };

    const formatDate = (dateStr: string) => {
        try {
            const date = new Date(dateStr);
            return date.toLocaleString();
        } catch {
            return dateStr;
        }
    };

    return (
        <div>
            <div
                className="flex items-center space-x-2 py-1.5 px-2 hover:bg-accent rounded-md cursor-pointer group"
                onClick={handleToggle}
                style={{ paddingLeft: `${level * 1.5 + 0.5}rem` }}
            >
                {/* Expand/collapse icon */}
                {node.is_dir && node.children && node.children.length > 0 ? (
                    isExpanded ? (
                        <ChevronDown className="h-4 w-4 text-muted-foreground shrink-0" />
                    ) : (
                        <ChevronRight className="h-4 w-4 text-muted-foreground shrink-0" />
                    )
                ) : (
                    <span className="w-4" />
                )}

                {/* Icon */}
                {node.is_dir ? (
                    isExpanded ? (
                        <FolderOpen className="h-4 w-4 text-primary shrink-0" />
                    ) : (
                        <Folder className="h-4 w-4 text-primary shrink-0" />
                    )
                ) : (
                    <File className="h-4 w-4 text-muted-foreground shrink-0" />
                )}

                {/* Name */}
                <span className="flex-1 truncate text-sm">{node.name}</span>

                {/* Metadata */}
                {!node.is_dir && node.size !== undefined && (
                    <span className="text-xs text-muted-foreground opacity-0 group-hover:opacity-100 transition-opacity">
                        {formatSize(node.size)}
                    </span>
                )}
                {!node.is_dir && node.mod_time && (
                    <span className="text-xs text-muted-foreground opacity-0 group-hover:opacity-100 transition-opacity hidden md:inline">
                        {formatDate(node.mod_time)}
                    </span>
                )}
            </div>

            {/* Children */}
            {node.is_dir && isExpanded && node.children && (
                <div>
                    {node.children.map((child, index) => (
                        <TreeNode key={`${child.path}-${index}`} node={child} level={level + 1} />
                    ))}
                </div>
            )}
        </div>
    );
}

export default function SnapshotFileTree({ root }: SnapshotFileTreeProps) {
    return (
        <div className="space-y-1">
            <TreeNode node={root} level={0} />
        </div>
    );
}
