import React, { createContext, useContext, useState, ReactNode, } from 'react';
import { v4 as uuidv4 } from 'uuid';

import { FileDownloadContextType, FileDownloadItemStore, FileDownloadItem, FileDownloadParams } from '../types';
import { useNotificationContext } from './NotificationContext';
import config from '../configs/config';
import api from '../api/api';
export const FileDownloadContext = createContext<FileDownloadContextType | undefined>(undefined);

export const FileDownloadProvider : React.FC<{children : ReactNode}> = ({children}) => {
    const { notifyError, notifyInfo, notifySuccess } = useNotificationContext();

    const [files, setFiles] = useState<FileDownloadItemStore | null>(null);

    const downloadFile = async (fileName: string, folderId: string) => {
        if (!files) return;
        
        try {
            const params: FileDownloadParams = {
                file: fileName,
                folder_id: folderId
            };
            const response = await api.get(config.DOWNLOAD_URL, {
                responseType: 'blob',
                params: params,
            });
            
            const blob = new Blob([response.data]);
            const downloadUrl = window.URL.createObjectURL(blob);
            
            const link = document.createElement('a');
            link.href = downloadUrl;
            link.setAttribute('download', fileName);
            
            document.body.appendChild(link);
            link.click();
            document.body.removeChild(link);

            window.URL.revokeObjectURL(downloadUrl);
            
            return true;
        } catch (error) {
            console.error('Download failed:', error);
            notifyError("Download Error", `Encountered unexpected error when attempting to download file ${fileName}`)
            return false;
        }
          
    }

    const downloadZip = async (folderId: string) => {
        const fileName = `${folderId}.zip`;

        await downloadFile(fileName, folderId);
    }

    const addFile = (fileItem: FileDownloadItem, fileName: string) => {
        setFiles((prev) => ({
            ...prev,
            [fileName]: fileItem
        }))   
    }    

    return (
        <FileDownloadContext.Provider value={{ files, setFiles, addFile, downloadFile, downloadZip }}>
        {children}
        </FileDownloadContext.Provider>
    );
}

export const useFileDownloadContext = () => {
    const context = useContext(FileDownloadContext);
    if (!context) {
      throw new Error("useFileDownloadContext must be used within a FileDownloadProvider");
    }
    return context;
};