import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { sourcesApi } from '@/lib/api';
import type { CreateSourceRequest, UpdateSourceRequest } from '@/types/api';

export const useSources = () => {
    return useQuery({
        queryKey: ['sources'],
        queryFn: sourcesApi.list,
    });
};

export const useSource = (id: number) => {
    return useQuery({
        queryKey: ['sources', id],
        queryFn: () => sourcesApi.get(id),
        enabled: !!id,
    });
};

export const useCreateSource = () => {
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: (data: CreateSourceRequest) => sourcesApi.create(data),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['sources'] });
        },
    });
};

export const useUpdateSource = () => {
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: ({ id, data }: { id: number; data: UpdateSourceRequest }) =>
            sourcesApi.update(id, data),
        onSuccess: (_, variables) => {
            queryClient.invalidateQueries({ queryKey: ['sources'] });
            queryClient.invalidateQueries({ queryKey: ['sources', variables.id] });
        },
    });
};

export const useDeleteSource = () => {
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: (id: number) => sourcesApi.delete(id),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['sources'] });
        },
    });
};

export const useRunBackup = () => {
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: (id: number) => sourcesApi.run(id),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['jobs'] });
            queryClient.invalidateQueries({ queryKey: ['snapshots'] });
        },
    });
};
