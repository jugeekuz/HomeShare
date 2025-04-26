import React, { useEffect } from 'react';
import { Route, Routes, Navigate, useNavigate } from 'react-router-dom';
import { useAuth } from './contexts/AuthContext';

import NavBar from './components/NavBar';
import HomePage from './pages/HomePage';
import LoginPage from './pages/LoginPage';
import SharingGateway from './pages/SharingGateway';
import SharingPage from './pages/SharingPage';
import LoadingPage from './pages/LoadingPage';

interface PrivateRouteProps {
    children: React.ReactNode;
}

const PrivateSharingRoute: React.FC<PrivateRouteProps> = ({ children }) => {
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

const PrivateAdminRoute: React.FC<PrivateRouteProps> = ({ children }) => {
    const { isAuthenticated, isAdmin, refreshLoading } = useAuth();
    const navigate = useNavigate();
    useEffect(() => {
        if (refreshLoading) return;
        if (!isAdmin) {
            navigate(-1);
        }
    },[isAdmin, refreshLoading])
    if (refreshLoading) {
      return <LoadingPage />;
    }
    if (!isAuthenticated) {
      return <Navigate to="/login" replace />;
    }
    if (!isAdmin) {
      // declarative redirect
      return <Navigate to=".." replace />;   // or to="/some‑fallback"
    }
    return <>{children}</>;
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
            <PrivateAdminRoute>
                <>
                <NavBar />
                <HomePage />
                </>
            </PrivateAdminRoute>
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
            path="/sg-"
            element={
                <PublicRoute>
                    <>
                    <SharingGateway />
                    </>
                </PublicRoute>
                }
        />
        <Route
            path="/sharing"
            element={
            <PrivateSharingRoute>
                <>
                <NavBar />
                <SharingPage />
                </>
            </PrivateSharingRoute>
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
            <PrivateAdminRoute>
                <>
                <NavBar />
                <HomePage />
                </>
            </PrivateAdminRoute>
            }
        />
        </Routes>
  );
};

export default AppRoutes;
