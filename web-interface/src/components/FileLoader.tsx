import React, { useRef } from 'react';
import { ProgressBarRef } from '../types'; 
import ProgressBar from './ProgressBar.tsx'
import { useFileContext } from '../contexts/FileContext.tsx';
import UploadItem from './UploadItem.tsx';
const ArrowIcon : React.FC = () => (
    <svg width="45px" height="45px" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path d="M8 8L12 4M12 4L16 8M12 4V16M4 20H20" stroke="#000000" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
    </svg>
)

const FileLoader : React.FC = () => {
    const { progressBarRefs } = useFileContext();
    const { files, addFile } = useFileContext();
    const fileInputRef = useRef<HTMLInputElement | null>(null);

    const handleFileChange = (e : React.ChangeEvent<HTMLInputElement>) => {
        if (e.target.files == null ) return;
        const selectedFiles : File[] = Array.from(e.target.files);

        for (const file of selectedFiles) {
            addFile(file);
        }
    }

    const refCallback = (el: ProgressBarRef | null, fileId: string) => {
        progressBarRefs.current[fileId] = el;
    }

    return (
        <div className="flex items-center justify-center w-full h-full">
            {
                !files || Object.keys(files).length === 0 ?
                    <div 
                        className="flex rounded-full items-center justify-center w-[5.5rem] h-[5.5rem] bg-blue-50 shadow-lg cursor-pointer"

                        onClick={() =>  fileInputRef.current?.click()}
                    >
                        <ArrowIcon/>
                    </div>
                :   
                    <div className="flex flex-row flex-wrap overflow-y-scroll items-start w-full h-full gap-x-7 gap-y-4 p-8 px-10">
                        {
                            Object.entries(files).map(([fileId, _], index) => (
                                <div className="flex flex-col max-w-12" key={fileId}>
                                    <UploadItem 
                                        fileId={fileId}/>
                                    <ProgressBar 
                                        ref={(el) => refCallback(el, fileId)} 
                                        className='w-12 mt-1'/>
                                </div>
                            ))
                        }
                    </div>
            }
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