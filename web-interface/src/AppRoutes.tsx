import React from 'react';
import { Route, Routes } from 'react-router-dom';

import HomePage from './pages/HomePage';
import UploadPage from './pages/UploadPage';
import SharingPage from './pages/SharingPage';

const AppRoutes: React.FC = () => {
  return (
    <Routes>
      <Route path="/" element={<HomePage />} />
      <Route path="/upload" element={<UploadPage />} />
      <Route path="/sharing" element={<SharingPage />} />
      <Route path="*" element={<UploadPage />} />
    </Routes>
  );
};

export default AppRoutes;