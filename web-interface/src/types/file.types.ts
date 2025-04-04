import { ProgressBarRefs } from "./progress.types";
import { DefaultExtensionType, FileIconProps } from "react-file-icon";

// Upload File Types
export interface FileMeta {
    fileId:         string;
    fileName:       string;
    fileExtension:  string;
    md5Hash:        string;
}

export interface FileItem {
    file:       File;
    fileMeta:   FileMeta;
}

export interface FileStore {
    [key: string]: FileItem
}

export type AddFile = (file: File) => void

export interface FileUploadContextType {
    files: FileStore | null;
    setFiles: React.Dispatch<React.SetStateAction<FileStore | null>>;
    addFile: AddFile;
    deleteFile: (fileId: string) => void;
    uploadFiles: (folderId ?: string) => Promise<void>;
    addMd5Hash: (fileId: string, md5Hash: string) => void;
    filesReady: boolean;
    filesUploading: boolean;
    progress: number,
    progressBarRefs: React.RefObject<ProgressBarRefs>;
}

export interface FileIconType {
    fileExtension: string;
    fileStyle: Record<DefaultExtensionType, Partial<FileIconProps>> | Partial<FileIconProps>
}
// Download Files Types
export interface FileDownloadItem {
    fileNameWoExt:      string;
    fileExtension:      string;
    fileSize:           number;
}

export interface FileDownloadItemStore {
    [key: string]: FileDownloadItem
}

export interface FileDownloadParams {
    file: string;
    folder_id: string;
}

export interface GetSharingFilesParams {
    folder_id : string
}

export interface FileDownloadContextType {
    files: FileDownloadItemStore | null;
    addFile: (file: FileDownloadItem, fileName: string) => void;
    setFiles:  React.Dispatch<React.SetStateAction<FileDownloadItemStore | null>>;
    downloadFile: (fileName: string, folderId: string) => void;
    downloadZip: (folderId: string) => void;
}