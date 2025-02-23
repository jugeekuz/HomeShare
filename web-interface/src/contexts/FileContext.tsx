import React, { createContext, useContext, useState, useRef, ReactNode } from 'react';
import { v4 as uuidv4 } from 'uuid';

import { FileMeta, FileItem, FileContextType, FileStore } from '../types';
import { chunkAndUpload } from '../services/chunkFile';
import { ProgressBarRefs } from '../types';
export const FileContext = createContext<FileContextType | undefined>(undefined);

export const FileProvider : React.FC<{children : ReactNode}> = ({children}) => {
    const [files, setFiles] = useState<FileStore | null>(null);
    const progressBarRefs = useRef<ProgressBarRefs>({});

    const addFile = (file: File) => {
        const fileName = file.name;
        const parts = fileName.split(/(\..+)$/);

        const fileMeta : FileMeta = {
            fileId:         uuidv4(),
            fileName:       parts[0],
            fileExtension:  parts[1] || '',
            md5Hash:        ""
        }

        const fileItem : FileItem = {
            file: file,
            fileMeta: fileMeta
        }

        setFiles((prev) => ({
            ...prev,
            [fileMeta.fileId]: fileItem
        }))
    }

    const uploadFiles = async () => {
        for (const fileId in files) {
            const fileItem = files[fileId];

            const setFileProgress = (progress: number) => {
                progressBarRefs.current[fileId]?.updateProgress(progress);
            }
            
            await chunkAndUpload(setFileProgress, fileItem.fileMeta, fileItem.file);
        }
    }

    const addMd5Hash = (fileId: string, md5Hash: string) => {
        if (!files) return;
        
        setFiles((prev) => {
            if (!prev) return prev;

            const fileItem = prev[fileId];
            if (!fileItem) return prev;

            return ({
                ...prev,
                [fileId]: {
                    ...fileItem,
                    fileMeta: {
                        ...fileItem.fileMeta,
                        md5Hash: md5Hash
                    }
                }
            })
        })
    }

    return (
        <FileContext.Provider value={{ files, setFiles, addFile, uploadFiles, progressBarRefs, addMd5Hash }}>
        {children}
        </FileContext.Provider>
    );
}

export const useFileContext = () => {
    const context = useContext(FileContext);
    if (!context) {
      throw new Error("useFileContext must be used within a FileProvider");
    }
    return context;
};