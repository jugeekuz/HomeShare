import React, { useCallback } from 'react';
import { FileUploadIcon } from './FileIcons.tsx';
import { useDropzone } from 'react-dropzone';
import { useFileUploadContext } from '../contexts/FileUploadContext';

const FileLoader : React.FC = () => {
    const { addFile } = useFileUploadContext();

    const onDrop = useCallback((acceptedFiles: File[]) => {
        acceptedFiles.forEach((file) => {
        addFile(file);
        });
    }, [addFile]);

    const { getRootProps, getInputProps, isDragActive } = useDropzone({
        onDrop,
        multiple: true,
    });


    return (
        <div className="flex flex-col items-center justify-center w-full h-full" {...getRootProps()}>
            <input
                {...getInputProps()}
                className="sr-only"
                aria-label="Select files to upload"
            />
            <div className="w-full h-full p-2">
                <div className={`relative flex flex-col w-full h-full justify-center items-center rounded-lg border-dashed border-2  p-10 cursor-pointer ${isDragActive ? "border-blue-500" : "border-gray-300"} bg-wsecondary`}>
                    {/* Overlay */}
                    {isDragActive && (
                        <div className="absolute inset-0 bg-primary-gradient opacity-30 rounded-lg pointer-events-none z-20" />
                    )}

                    {/* Content */}
                    <div className={`flex flex-col items-center justify-center z-10 ${isDragActive ? "opacity-30" : ""} `}>
                        <FileUploadIcon className="w-12 h-12" />
                        <span className="font-linik text-gray-800 text-sm mt-2">
                            Drag & Drop or{' '}
                            <span className="bg-secondary-gradient bg-clip-text text-transparent font-medium">
                                Choose files
                            </span>{' '}
                            to upload
                        </span>
                        <span className="font-linik text-gray-400 text-xs">
                            Maximum supported file size 5GB
                        </span>
                    </div>
                </div>
            </div>
        </div>
    );
}

export default FileLoader;