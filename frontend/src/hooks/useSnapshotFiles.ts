import { useQuery } from '@tanstack/react-query';
import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export interface FileNode {
    name: string;
    path: string;
    is_dir: boolean;
    size?: number;
    mod_time?: string;
    children?: FileNode[];
}

export function useSnapshotFiles(snapshotId: number | undefined) {
    return useQuery({
        queryKey: ['snapshots', snapshotId, 'files'],
        queryFn: async () => {
            const token = localStorage.getItem('auth_token');
            const response = await axios.get(`${API_BASE_URL}/api/snapshots/${snapshotId}/files`, {
                headers: token ? { Authorization: `Bearer ${token}` } : {},
            });
            return response.data.data as FileNode;
        },
        enabled: !!snapshotId,
    });
}
