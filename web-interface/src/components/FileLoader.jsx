import { useRef, useState, useContext } from 'react';
import { FileContext } from '../contexts/FileContext';


const ArrowIcon = () => (
    <svg width="45px" height="45px" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path d="M8 8L12 4M12 4L16 8M12 4V16M4 20H20" stroke="#000000" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
    </svg>
)

const FileLoader = () => {

    const { files, setFiles, progress, setProgress } = useContext(FileContext);
    const [isUploading, setIsUploading] = useState(false);
    const fileInputRef = useRef(null);

    const handleFileChange = (e) => {
        const selectedFiles = Array.from(e.target.files);
        setFiles(selectedFiles);
    };

    return (
    <div className="flex items-center justify-center">
        <div 
            className="flex rounded-full items-center justify-center w-[5.5rem] h-[5.5rem] bg-blue-50 shadow-lg cursor-pointer"

            onClick={() =>  fileInputRef.current?.click()}
        >
            <ArrowIcon/>
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
};

export default FileLoader;
