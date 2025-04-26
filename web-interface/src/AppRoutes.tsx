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
    return <>{
        refreshLoading ?
            <LoadingPage />
        : (isAuthenticated && isAdmin?
            <>{children}</>
            :<Navigate to="/login" />)
    }</>
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
