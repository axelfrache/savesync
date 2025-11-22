import { create } from 'zustand';
import { authApi } from '@/lib/api';

interface User {
    id: number;
    email: string;
    created_at: string;
    updated_at: string;
}

interface AuthState {
    user: User | null;
    token: string | null;
    isAuthenticated: boolean;
    isLoading: boolean;
    login: (email: string, password: string) => Promise<void>;
    register: (email: string, password: string) => Promise<void>;
    logout: () => void;
    loadUser: () => Promise<void>;
}

export const useAuthStore = create<AuthState>((set) => ({
    user: null,
    token: localStorage.getItem('auth_token'),
    isAuthenticated: !!localStorage.getItem('auth_token'),
    isLoading: false,

    login: async (email: string, password: string) => {
        set({ isLoading: true });
        try {
            const response = await authApi.login(email, password);
            const { token, user } = response;

            localStorage.setItem('auth_token', token);
            set({ token, user, isAuthenticated: true, isLoading: false });
        } catch (error) {
            set({ isLoading: false });
            throw error;
        }
    },

    register: async (email: string, password: string) => {
        set({ isLoading: true });
        try {
            const response = await authApi.register(email, password);
            const { token, user } = response;

            localStorage.setItem('auth_token', token);
            set({ token, user, isAuthenticated: true, isLoading: false });
        } catch (error) {
            set({ isLoading: false });
            throw error;
        }
    },

    logout: () => {
        localStorage.removeItem('auth_token');
        set({ user: null, token: null, isAuthenticated: false });
    },

    loadUser: async () => {
        const token = localStorage.getItem('auth_token');
        if (!token) {
            set({ isAuthenticated: false, user: null });
            return;
        }

        set({ isLoading: true });
        try {
            const user = await authApi.me();
            set({ user, isAuthenticated: true, isLoading: false });
        } catch (error) {
            localStorage.removeItem('auth_token');
            set({ user: null, token: null, isAuthenticated: false, isLoading: false });
        }
    },
}));
