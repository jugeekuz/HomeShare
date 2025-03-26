import React, { useRef } from 'react';
import { FileUploadIcon } from './FileIcons.tsx';

import { useFileContext } from '../contexts/FileContext.tsx';

const FileLoader : React.FC = () => {
    const { addFile } = useFileContext();
    const fileInputRef = useRef<HTMLInputElement | null>(null);

    const handleFileChange = (e : React.ChangeEvent<HTMLInputElement>) => {
        if (e.target.files == null ) return;
        const selectedFiles : File[] = Array.from(e.target.files);

        for (const file of selectedFiles) {
            addFile(file);
        }
    }
    return (
        <div className="flex flex-col items-center justify-center w-full h-full">
            <div className="w-full h-full p-2">
                <div 
                    className="flex flex-col w-full h-full justify-center items-center rounded-lg border-dashed border-2 border-gray-300 p-10 cursor-pointer bg-wsecondary"
                    onClick={() =>  fileInputRef.current?.click()}
                    >
                    <FileUploadIcon className="w-12 h-12"/>

                    <span className="font-linik  text-gray-800 text-sm mt-2">
                        Drag & Drop or <span className="bg-secondary-gradient bg-clip-text text-transparent font-medium">Choose files</span> to upload
                    </span>
                    <span className="font-linik  text-gray-400 text-xs ">
                        Maximum supported file size 2GB
                    </span>
                </div>
            </div>
            <input
                ref={fileInputRef}
                type="file"
                className="sr-only"
                onChange={handleFileChange}
                accept="*"
                multiple
                aria-label="Select files to upload"
            />
        </div>
    );
}

export default FileLoader;