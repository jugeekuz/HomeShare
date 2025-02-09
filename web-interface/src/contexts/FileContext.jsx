import React, { createContext, useState } from 'react';

export const FileContext = createContext();

export const FileProvider = ({ children }) => {
    const [file, setFile] = useState(null);
    const [progress, setProgress] = useState(0)
    return (
        <FileContext.Provider value={{ file, setFile, progress, setProgress }}>
        {children}
        </FileContext.Provider>
    );
};