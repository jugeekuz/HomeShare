import React, { createContext, useContext, useState, ReactNode, } from 'react';

import { FileDownloadContextType, FileDownloadItemStore, FileDownloadItem, FileDownloadParams } from '../types';
import { useNotificationContext } from './NotificationContext';
import config from '../configs/config';
import api from '../api/api';
export const FileDownloadContext = createContext<FileDownloadContextType | undefined>(undefined);

export const FileDownloadProvider : React.FC<{children : ReactNode}> = ({children}) => {
    const { notifyError, notifyInfo } = useNotificationContext();

    const [files, setFiles] = useState<FileDownloadItemStore | null>(null);

    const downloadFile = async (fileName: string, folderId: string) => {
        if (!files) return;
        
        try {
            const params : FileDownloadParams = {
                folder_id: folderId,
                file: fileName
            }
            const response = await api.get(config.GET_DOWNLOAD_FILE_AVAILABLE_URL, {
                params: params,
            });

            if (response.status !== 200) {
                notifyError("Download Error", "Zip file is currently being processed, try again in a few moments")
                return
            }

            const url = `${config.DOWNLOAD_URL}?file=${encodeURIComponent(fileName)}&folder_id=${encodeURIComponent(folderId)}`;
            const link = document.createElement('a');
            link.href = url;
            link.setAttribute('download', fileName);
            document.body.appendChild(link);
            link.click();
            document.body.removeChild(link);
    
            notifyInfo("File Download", `${fileName} has successfully started downloading`)
        } catch (error) {
            notifyError("Download Error", "Zip file is currently being processed, try again in a few moments")
            return
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