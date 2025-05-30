import React, { createContext, useContext, useState, useRef, ReactNode, useEffect } from 'react';
import { v4 as uuidv4 } from 'uuid';

import { FileMeta, FileItem, FileUploadContextType, FileStore } from '../types';
import { chunkAndUpload } from '../services/chunkFile';
import { ProgressBarRefs } from '../types';
import { useNotificationContext } from './NotificationContext';
export const FileUploadContext = createContext<FileUploadContextType | undefined>(undefined);

export const FileUploadProvider : React.FC<{children : ReactNode}> = ({children}) => {
    const { notifyError, notifyInfo, notifySuccess } = useNotificationContext();
    const [files, setFiles] = useState<FileStore | null>(null);
    const [filesReady, setFilesReady] = useState<boolean>(false);
    const [filesUploading, setFilesUploading] = useState<boolean>(false);
    const totalFileSize = useRef<number>(0);
    const totalFileSizeSent = useRef<number>(0);
    const [progress, setProgress] = useState<number>(0);
    const progressBarRefs = useRef<ProgressBarRefs>({});

    // Calculate total size
    useEffect(() => {
        if (!files) return;
        
        totalFileSize.current = 0;
        totalFileSizeSent.current = 0;
        
        for (const fileItem of Object.values(files)) {
            if (!fileItem?.file) continue;
            totalFileSize.current += fileItem.file.size;
        }
    }, [files]);

    // See if all hashes are calculated
    useEffect(() => {
        if (!files) return;
        setFilesReady(Object.values(files).every(fileItem => !!fileItem?.fileMeta.md5Hash))
    },[files])

    const addFile = (file: File) => {
        const fileName = file.name;
        const lastDotIndex = fileName.lastIndexOf('.');
        const baseName = lastDotIndex !== -1 ? fileName.slice(0, lastDotIndex) : fileName;
        const extension = lastDotIndex !== -1 ? fileName.slice(lastDotIndex) : '';

        const fileMeta : FileMeta = {
            fileId:         uuidv4(),
            fileName:       baseName,
            fileExtension:  extension,
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

    const deleteFile = (fileId: string) => {
        if (!fileId || !files || !files[fileId]) return;
        
        const newFiles = {...files}

        delete newFiles[fileId];

        setFiles(newFiles);
    }

    const uploadFiles = async (folderId ?: string) => {
        for (const fileId in files) {
            setFilesUploading(true);

            const fileItem = files[fileId];
            const fileSize = fileItem.file.size;

            const setFileProgress = (progress: number) => {
                progressBarRefs.current[fileId]?.updateProgress(Math.ceil(progress));
            }
            try {
                const uploadResponse = await chunkAndUpload(setFileProgress, fileItem.fileMeta, fileItem.file, folderId);
                
                if (uploadResponse.success) {
                    totalFileSizeSent.current += fileSize;
                    setProgress((_) => (totalFileSizeSent.current/totalFileSize.current)*100);
                } else {
                    notifyError("Upload Error", `File \`${fileItem.fileMeta.fileName}${fileItem.fileMeta.fileExtension}\` failed to upload`)
                }
            } catch (error) {
                notifyError("Upload Error", `File \`${fileItem.fileMeta.fileName}${fileItem.fileMeta.fileExtension}\` failed to upload`)
            }
        }
        if ((totalFileSizeSent.current/totalFileSize.current)*100 === 100) { // dont use state as it is async
            notifySuccess("Upload Success", "Files finished uploaded successfully")
        } else {
            notifyInfo("Upload Error", "Files finished uploading. Some files failed to upload")
        }
        setFilesUploading(false);
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
        <FileUploadContext.Provider value={{ files, setFiles, addFile, deleteFile, filesReady, uploadFiles, progressBarRefs, progress, filesUploading, addMd5Hash }}>
        {children}
        </FileUploadContext.Provider>
    );
}

export const useFileUploadContext = () => {
    const context = useContext(FileUploadContext);
    if (!context) {
      throw new Error("useFileUploadContext must be used within a FileUploadProvider");
    }
    return context;
};