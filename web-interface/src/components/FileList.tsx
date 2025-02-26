import React, { useRef } from 'react';
import { ProgressBarRef } from '../types'; 
import "react-perfect-scrollbar/dist/css/styles.css";
import PerfectScrollbar from 'react-perfect-scrollbar';
import ProgressBar from './ProgressBar.tsx'
import { useFileContext } from '../contexts/FileContext.tsx';
import UploadItem from './UploadItem.tsx';
import { LuTrash2 } from "react-icons/lu";

const FileBox : React.FC<{fileId: string, refCallback: (el: ProgressBarRef | null, fileId: string) => void}> = ({fileId, refCallback}) => {
    const { files, deleteFile } = useFileContext();

    return (
    <div className="flex flex-col justify-center items-center w-full h-[4.5rem] rounded-md bg-gray-100 px-2 mb-1 border-1 border-gray-200" key={fileId}>
        <div className="relative flex flex-row justify-start items-center w-full p-2">
            <div className="w-8 h-8">
                <UploadItem fileId={fileId}/>
            </div>
            <div className="flex flex-col w-full px-3">
                <span className="font-linik text-sm text-gray-700 font-bold">
                    {files && files[fileId].fileMeta.fileName}
                </span>
                <span className="text-xs text-gray-600">
                    {files && files[fileId].fileMeta.fileExtension}
                </span>
            </div>
            <div 
                className="absolute flex justify-center items-center rounded-full bg-white w-8 h-8 top-2 right-0 border border-gray-200 cursor-pointer"
                onClick={() => deleteFile(fileId)}
                >
                <LuTrash2 className='text-gray-600'/>
            </div>
        </div>
        <div className="flex justify-center items-center w-full">
            <ProgressBar 
                size='sm'
                ref={(el) => refCallback(el, fileId)} 
                className='max-w-[80%]'/>
        </div>
    </div>)
}

const FileList : React.FC = () => {
    const { progressBarRefs, files } = useFileContext();
    const refCallback = (el: ProgressBarRef | null, fileId: string) => {
        progressBarRefs.current[fileId] = el;
    }

    return (
        <div className="flex flex-col items-center justify-center w-full gap-2 px-2 max-h-[12rem]">
            {files && Object.keys(files).length > 0 && (
                <PerfectScrollbar
                    className="w-full h-full"
                    options={{
                        wheelSpeed: 1,
                        suppressScrollX: true,
                        minScrollbarLength: 30,
                    }}
                >
                    <div className="flex flex-col items-start w-full">
                        {Object.entries(files).map(([fileId, _], index) => (
                            <FileBox
                                key={fileId}
                                fileId={fileId}
                                refCallback={refCallback}
                            />
                        ))}
                    </div>
                </PerfectScrollbar>
            )}
        </div>
    );
}

export default FileList;