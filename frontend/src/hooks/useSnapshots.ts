import { useQuery } from '@tanstack/react-query';
import { snapshotsApi } from '@/lib/api';

export const useSnapshots = () => {
    return useQuery({
        queryKey: ['snapshots'],
        queryFn: snapshotsApi.list,
    });
};

export const useSnapshot = (id: number) => {
    return useQuery({
        queryKey: ['snapshots', id],
        queryFn: () => snapshotsApi.get(id),
        enabled: !!id,
    });
};
