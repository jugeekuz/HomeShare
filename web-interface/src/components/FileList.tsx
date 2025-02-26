import React, { useRef } from 'react';
import { ProgressBarRef } from '../types'; 
import ProgressBar from './ProgressBar.tsx'
import { useFileContext } from '../contexts/FileContext.tsx';
import UploadItem from './UploadItem.tsx';


const FileBox : React.FC<{fileId: string, refCallback: (el: ProgressBarRef | null, fileId: string) => void}> = ({fileId, refCallback}) => {
    const { files } = useFileContext();

    return (
    <div className="flex flex-col justify-center items-center w-full h-[4.5rem] rounded-md bg-gray-100 px-2 mb-1 border-1 border-gray-200" key={fileId}>
        <div className="flex flex-row justify-start items-center w-full p-2">
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
        <div className="flex flex-col items-center justify-center w-full gap-2 px-2 max-h-[14rem]">
            {
                !files || Object.keys(files).length === 0 ? <></>
                :
                <div className="flex flex-col items-start w-full h-full  overflow-y-scroll">
                    {
                        Object.entries(files).map(([fileId, _], index) => (
                            <FileBox fileId={fileId} refCallback={refCallback}/>
                        ))
                    }
                </div>
            }
        </div>
    );
}

export default FileList;