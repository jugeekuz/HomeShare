import React, { createContext, useState, useEffect } from 'react';

export const FileContext = createContext();

export const FileProvider = ({ children }) => {
    const [files, setFiles] = useState([]);
    const [filesMeta, setFilesMeta] = useState([]);
    const [progress, setProgress] = useState(0);

    const appendFilesMeta = (data) => {
        if (!data?.fileId || !data?.fileName || !data?.fileExtension || !data?.md5Hash) return;
        setFilesMeta((prevFilesMeta) => [...prevFilesMeta, data]);
    }

    return (
        <FileContext.Provider value={{ files, setFiles, filesMeta, setFilesMeta, appendFilesMeta, progress, setProgress }}>
        {children}
        </FileContext.Provider>
    );
};