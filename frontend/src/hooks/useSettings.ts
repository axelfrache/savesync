import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { settingsApi } from '@/lib/api';
import { useAuthStore } from '@/store/authStore';

export function useSettings() {
    const queryClient = useQueryClient();
    const { user } = useAuthStore();

    // Only fetch settings if user is admin
    const { data: settings, isLoading } = useQuery({
        queryKey: ['settings'],
        queryFn: settingsApi.getAll,
        enabled: !!user?.is_admin,
    });

    const updateSetting = useMutation({
        mutationFn: ({ key, value }: { key: string; value: string }) =>
            settingsApi.update(key, value),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['settings'] });
        },
    });

    return {
        settings,
        isLoading,
        updateSetting,
        isRegistrationEnabled: settings?.['registration_enabled'] === 'true',
    };
}
