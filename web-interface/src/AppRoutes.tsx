import React from 'react';
import { Route, Routes } from 'react-router-dom';

import NavBar from './components/NavBar';
import HomePage from './pages/HomePage';
import LoginPage from './pages/LoginPage';
import SharingPage from './pages/SharingPage';

const AppRoutes: React.FC = () => {
    return (
        <Routes>
            <Route path="/" element={
                <>
                <NavBar />
                <HomePage />
                </>
            } />
            <Route path="/login" element={<LoginPage />} />
            <Route path="/sharing" element={
                <>
                <NavBar />
                <SharingPage />
                </>
            } />
            <Route path="*" element={
                <>
                <NavBar />
                <HomePage />
                </>
            } />
        </Routes>
    );
};

export default AppRoutes;