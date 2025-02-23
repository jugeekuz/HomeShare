import React, { useEffect, useState } from 'react';
import config from '../configs/config'
import { FileIcon, defaultStyles, DefaultExtensionType } from 'react-file-icon';
import { FileIconType } from '../types';
import SparkMD5 from 'spark-md5';
import { useFileContext } from '../contexts/FileContext.tsx';

const UploadItem : React.FC<{fileId: string}> = ({ fileId }) => {
    const MAX_FILE_SIZE_MB = config.MAX_FILE_SIZE_MB;
    const { files, addMd5Hash } = useFileContext();
    const [file, setFile] = useState<File | null>(null);
    const [fileStyle, setFileStyle] = useState<FileIconType>({fileExtension: '', fileStyle: defaultStyles});
    const [fileMd5, setFileMd5] = useState<string | null>(null);

    useEffect(() => {
        if (!files) return;
        if (!(fileId in files)) {
            console.error(`FileId ${fileId} does not exist in File Store.`)
            return;
        }
        setFile(files[fileId].file);
    }, [])

    useEffect(() => {
        if(!files || !file) return;
        const fileExt = files[fileId].fileMeta.fileExtension.replace(/^\./, "").toLowerCase() as DefaultExtensionType;


        setFileStyle({
            fileExtension: fileExt,
            fileStyle: (defaultStyles[fileExt] || defaultStyles)
        })
        
    },[file])

    useEffect(() => {
        if (!file) return;

        if (file.size > MAX_FILE_SIZE_MB) {
            console.error(`File size is ${file.size} above maximum ${MAX_FILE_SIZE_MB}`);
            return;
        }

        const reader = new FileReader();
        let isCancelled = false;

        reader.onload = (e : ProgressEvent<FileReader>) => {
            if (isCancelled || !e.target?.result || typeof(e.target.result) == 'string') return;
            const buffer = e.target.result;
            try {
                const hash = SparkMD5.ArrayBuffer.hash(buffer);
                setFileMd5(hash);
            } catch (error) {
                console.error('Error generating MD5 hash:', error);
            }
        };

        reader.onerror = (error) => {
            if (isCancelled) return;
            console.error('Error reading file:', error);
        };

        reader.readAsArrayBuffer(file);

        // TODOOOO
        return () => {
            isCancelled = true;
            if (reader.readyState === FileReader.LOADING) {
                reader.abort();
                console.log('Aborted ongoing file read.');
            }
        };
    }, [file]);

    useEffect(() => {
        if (!fileMd5) return;

        addMd5Hash(fileId, fileMd5);
    }, [file, fileMd5])

    return (
        <FileIcon color='#fcba03' extension={fileStyle.fileExtension} {...fileStyle.fileStyle}  />
    )
}

export default UploadItem;