import { useQuery } from '@tanstack/react-query';
import { jobsApi } from '@/lib/api';

export const useJobs = () => {
    return useQuery({
        queryKey: ['jobs'],
        queryFn: jobsApi.list,
    });
};

export const useJob = (id: number) => {
    return useQuery({
        queryKey: ['jobs', id],
        queryFn: () => jobsApi.get(id),
        enabled: !!id,
    });
};
