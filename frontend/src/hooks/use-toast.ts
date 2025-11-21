import { useState, useCallback } from 'react';

export interface Toast {
    id: string;
    title: string;
    description?: string;
    variant?: 'default' | 'destructive';
}

interface ToastState {
    toasts: Toast[];
}

const toastState: ToastState = {
    toasts: [],
};

const listeners: Array<(state: ToastState) => void> = [];

function notify() {
    listeners.forEach((listener) => listener(toastState));
}

export function useToast() {
    const [, setToasts] = useState<Toast[]>([]);

    const subscribe = useCallback((listener: (state: ToastState) => void) => {
        listeners.push(listener);
        return () => {
            const index = listeners.indexOf(listener);
            if (index > -1) {
                listeners.splice(index, 1);
            }
        };
    }, []);

    const toast = useCallback(
        ({ title, description, variant = 'default' }: Omit<Toast, 'id'>) => {
            const id = Math.random().toString(36).substring(7);
            const newToast: Toast = { id, title, description, variant };

            toastState.toasts.push(newToast);
            notify();

            // Auto dismiss after 3 seconds
            setTimeout(() => {
                toastState.toasts = toastState.toasts.filter((t) => t.id !== id);
                notify();
            }, 3000);

            return id;
        },
        []
    );

    // Subscribe to toast state changes
    useState(() => {
        const unsubscribe = subscribe((state) => {
            setToasts([...state.toasts]);
        });
        return unsubscribe;
    });

    return { toast, toasts: toastState.toasts };
}
