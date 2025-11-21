import { useQuery } from '@tanstack/react-query';
import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

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
            const response = await axios.get(`${API_BASE_URL}/snapshots/${snapshotId}/files`);
            return response.data.data as FileNode;
        },
        enabled: !!snapshotId,
    });
}
