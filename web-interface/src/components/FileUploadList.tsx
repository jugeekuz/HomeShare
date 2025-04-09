import React from 'react';

import PerfectScrollbar from 'react-perfect-scrollbar';
import { HiXMark } from "react-icons/hi2";
import "react-perfect-scrollbar/dist/css/styles.css";
import { ScrollShadow } from '@heroui/react';

import ProgressBar from './ProgressBar.tsx'
import { ProgressBarRef } from '../types/index.ts'; 
import { useFileUploadContext } from '../contexts/FileUploadContext.tsx';
import UploadItem from './UploadItem.tsx';


const FileBox : React.FC<{fileId: string, refCallback: (el: ProgressBarRef | null, fileId: string) => void}> = ({fileId, refCallback}) => {
    const { files, deleteFile } = useFileUploadContext();

    const convertBytes = (bytes: number) : string => {
        if (bytes === 0) return "0 B";

        const mapping = ["B", "kB", "MB", "GB", "TB"];
        const n = Math.floor(Math.log(bytes) / Math.log(1024));
        
        if (n >= mapping.length) throw new Error("Byte string is too large to convert.");

        const value = bytes / (1024 ** n);
        
        const formattedValue = parseFloat(value.toFixed(2));

        return `${formattedValue} ${mapping[n]}`;
    }

    return (
    <div className="flex flex-col justify-center items-center w-full h-[4.5rem] rounded-lg bg-contrast px-2 mb-1 border-1 border-gray-200" key={fileId}>
        <div className="relative flex flex-row justify-start items-center w-full p-2">
            <div className="flex justify-center items-center w-7 h-7">
                <UploadItem fileId={fileId}/>
            </div>
            <div className="flex flex-col w-full px-3 justify-center overflow-hidden text-nowrap">
                <span className="font-linik text-[13px] text-gray-700 font-bold">
                    {files && files[fileId].fileMeta.fileName}{files && files[fileId].fileMeta.fileExtension}
                </span>
                <div className="flex flex-row justify-start items-center">
                    <span className="text-xs text-gray-600">
                        {files && convertBytes(files[fileId].file.size)}
                    </span>
                </div>
            </div>
            <div 
                className="absolute flex justify-center items-center rounded-full w-5 h-5 top-0 right-0 -mr-1  cursor-pointer"
                onClick={() => deleteFile(fileId)}
                >
                <HiXMark className='text-gray-600'/>
            </div>
        </div>
        <div className="flex justify-center items-center w-full">
            <ProgressBar 
                size='sm'
                ref={(el) => refCallback(el, fileId)} 
                className='max-w-[90%]'/>
        </div>
    </div>)
}

const FileUploadList : React.FC = () => {
    const { progressBarRefs, files } = useFileUploadContext();
    const refCallback = (el: ProgressBarRef | null, fileId: string) => {
        progressBarRefs.current[fileId] = el;
    }

    return (
        <div className="flex flex-col items-center justify-center w-full gap-2 max-h-[10.5rem]">
            {files && Object.keys(files).length > 0 && (
                <PerfectScrollbar
                    className="w-full h-full"
                    options={{
                        wheelSpeed: 1,
                        suppressScrollX: true,
                        minScrollbarLength: 30,
                    }}
                >   
                    <ScrollShadow 
                        className='w-full h-full'
                        visibility='bottom'
                        size={12}
                    >
                        <div className="flex flex-col items-start w-full gap-[2px]">
                            {Object.entries(files).map(([fileId, _], index) => (
                                <FileBox
                                    key={`${index}${fileId}`}
                                    fileId={fileId}
                                    refCallback={refCallback}
                                />
                            ))}
                        </div>
                    </ScrollShadow>
                </PerfectScrollbar>
            )}
        </div>
    );
}

export default FileUploadList;