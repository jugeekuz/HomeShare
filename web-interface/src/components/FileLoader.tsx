import React, { useRef } from 'react';
import { FileIcon, defaultStyles } from 'react-file-icon';
import { LuUpload } from "react-icons/lu";

import { useFileContext } from '../contexts/FileContext.tsx';


const FileUploadIcon : React.FC<{className: string | null, active?: boolean}> = ({className, active = false}) => (
    <div className={`relative w-12 h-12 ${className}`}>
        <div className="w-3/4 h-1/4">
            <FileIcon color='#f0f4f7' extension="" {...defaultStyles}/>
        </div>
        <div className={`absolute flex ${active ? "bg-blue-700" : "bg-[#051d41]"} rounded-full w-[50%] h-[50%] right-0 bottom-0 p-[5px] justify-center items-center`}>
            <LuUpload className='text-gray-50 w-full h-full'/>
        </div>
    </div>
)

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
                    className="flex flex-col w-full h-full justify-center items-center rounded-lg border-dashed border-2 border-gray-300 p-10 cursor-pointer"
                    onClick={() =>  fileInputRef.current?.click()}
                    >
                    <FileUploadIcon className="w-12 h-12"/>

                    <span className="font-linik  text-gray-800 text-sm mt-2">
                        Drag & Drop or <span className="text-blue-500">Choose files</span> to upload
                    </span>
                    <span className="font-linik  text-gray-400 text-xs ">
                        Maximum supported file size 4MB
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