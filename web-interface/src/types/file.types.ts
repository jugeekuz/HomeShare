import { ProgressBarRefs } from "./progress.types";

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

export interface FileContextType {
    files: FileStore | null;
    setFiles: React.Dispatch<React.SetStateAction<FileStore | null>>;
    addFile: AddFile;
    uploadFiles: () => Promise<void>;
    addMd5Hash: (fileId: string, md5Hash: string) => void; 
    progressBarRefs: React.RefObject<ProgressBarRefs>;
}