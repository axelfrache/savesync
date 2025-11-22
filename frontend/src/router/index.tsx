import { createBrowserRouter, Navigate } from 'react-router-dom';
import AppLayout from '@/components/layout/AppLayout';
import ProtectedRoute from '@/components/shared/ProtectedRoute';
import Dashboard from '@/pages/Dashboard';
import Sources from '@/pages/Sources';
import SourceDetail from '@/pages/SourceDetail';
import Snapshots from '@/pages/Snapshots';
import SnapshotDetail from '@/pages/SnapshotDetail';
import Targets from '@/pages/Targets';
import Jobs from '@/pages/Jobs';
import Settings from '@/pages/Settings';
import LoginPage from '@/pages/Login';
import RegisterPage from '@/pages/Register';

export const router = createBrowserRouter([
    {
        path: '/login',
        element: <LoginPage />,
    },
    {
        path: '/register',
        element: <RegisterPage />,
    },
    {
        path: '/',
        element: (
            <ProtectedRoute>
                <AppLayout />
            </ProtectedRoute>
        ),
        children: [
            {
                index: true,
                element: <Navigate to="/dashboard" replace />,
            },
            {
                path: 'dashboard',
                element: <Dashboard />,
            },
            {
                path: 'sources',
                element: <Sources />,
            },
            {
                path: 'sources/:id',
                element: <SourceDetail />,
            },
            {
                path: 'snapshots',
                element: <Snapshots />,
            },
            {
                path: 'snapshots/:id',
                element: <SnapshotDetail />,
            },
            {
                path: 'targets',
                element: <Targets />,
            },
            {
                path: 'jobs',
                element: <Jobs />,
            },
            {
                path: 'settings',
                element: <Settings />,
            },
        ],
    },
]);
