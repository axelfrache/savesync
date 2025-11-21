import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { targetsApi } from '@/lib/api';
import type { CreateTargetRequest } from '@/types/api';

export const useTargets = () => {
    return useQuery({
        queryKey: ['targets'],
        queryFn: targetsApi.list,
    });
};

export const useTarget = (id: number) => {
    return useQuery({
        queryKey: ['targets', id],
        queryFn: () => targetsApi.get(id),
        enabled: !!id,
    });
};

export const useCreateTarget = () => {
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: (data: CreateTargetRequest) => targetsApi.create(data),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['targets'] });
        },
    });
};

export const useUpdateTarget = () => {
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: ({ id, data }: { id: number; data: Partial<CreateTargetRequest> }) =>
            targetsApi.update(id, data),
        onSuccess: (_, variables) => {
            queryClient.invalidateQueries({ queryKey: ['targets'] });
            queryClient.invalidateQueries({ queryKey: ['targets', variables.id] });
        },
    });
};

export const useDeleteTarget = () => {
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: (id: number) => targetsApi.delete(id),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['targets'] });
        },
    });
};
