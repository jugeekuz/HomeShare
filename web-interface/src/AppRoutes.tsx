import React from 'react';
import { Route, Routes, Navigate } from 'react-router-dom';
import { useAuth } from './contexts/AuthContext';

import NavBar from './components/NavBar';
import HomePage from './pages/HomePage';
import LoginPage from './pages/LoginPage';
import SharingPage from './pages/SharingPage';
import LoadingPage from './pages/LoadingPage';

interface PrivateRouteProps {
    children: React.ReactNode;
}

const PrivateRoute: React.FC<PrivateRouteProps> = ({ children }) => {
    const { isAuthenticated, refreshLoading } = useAuth();

    return<> 
        {
        refreshLoading ? <LoadingPage/>
            : isAuthenticated 
                ? <>{children}</>
                : <Navigate to="/login" />
        }
    </>
};

interface PublicRouteProps {
    children: React.ReactNode;
}

const PublicRoute: React.FC<PublicRouteProps> = ({ children }) => {
    return <>{children}</>;
};

const AppRoutes: React.FC = () => {
    return (
        <Routes>
        <Route
            path="/"
            element={
            <PrivateRoute>
                <>
                <NavBar />
                <HomePage />
                </>
            </PrivateRoute>
            }
        />
        <Route
            path="/login"
            element={
            <PublicRoute>
                <LoginPage />
            </PublicRoute>
            }
        />
        <Route
            path="/sharing"
            element={
            <PrivateRoute>
                <>
                <NavBar />
                <SharingPage />
                </>
            </PrivateRoute>
            }
        />
        <Route
            path="/loading"
            element={
            <PublicRoute>
                <>
                <LoadingPage/>
                </>
            </PublicRoute>
            }
        />
        <Route
            path="*"
            element={
            <PrivateRoute>
                <>
                <NavBar />
                <HomePage />
                </>
            </PrivateRoute>
            }
        />
        </Routes>
  );
};

export default AppRoutes;
