import React from 'react';
import { FileIcon, defaultStyles } from 'react-file-icon';
import { LuUpload } from "react-icons/lu";

import { LuDownload } from "react-icons/lu";


export const FileUploadIcon : React.FC<{className: string | null, active?: boolean}> = ({className, active = false}) => (
    <div className={`relative w-12 h-12 ${className}`}>
        <div className="w-3/4 h-1/4">
            <FileIcon color='#ebeff2' extension="" {...defaultStyles}/>
        </div>
        <div className={`absolute flex ${active ? "bg-blue-700" : "bg-[#051d41]"} rounded-full w-[50%] h-[50%] right-0 bottom-0 p-[5px] justify-center items-center`}>
            <LuUpload className='text-gray-50 w-full h-full '/>
        </div>
    </div>
)

export const FileDownloadIcon : React.FC<{className: string | null, active?: boolean}> = ({className, active = false}) => (
    <div className={`relative w-12 h-12 ${className}`}>
        <div className="w-3/4 h-1/4">
            <FileIcon color='#f0f4f7' extension="" {...defaultStyles}/>
        </div>
        <div className={`absolute flex ${active ? "bg-blue-700" : "bg-[#051d41]"} rounded-full w-[50%] h-[50%] right-0 bottom-0 p-[5px] justify-center items-center`}>
            <LuDownload className='text-gray-50 w-full h-full'/>
        </div>
    </div>
)
